package engine

import (
	"context"
	"sync"
	"time"

	"github.com/EmmmmDeee/One-Liner-OSINT/internal/cache"
	"github.com/EmmmmDeee/One-Liner-OSINT/internal/worker"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/config"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/logger"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/output"
)

// Engine is the core OSINT execution engine
type Engine struct {
	config  *config.Config
	logger  *logger.Logger
	cache   *cache.Cache
	pool    *worker.Pool
	modules map[string]Module
	mu      sync.RWMutex
}

// Module represents an OSINT module
type Module interface {
	Name() string
	Description() string
	Execute(ctx context.Context, target string) ([]output.Result, error)
}

// NewEngine creates a new OSINT engine
func NewEngine(cfg *config.Config) (*Engine, error) {
	log := logger.NewLogger()
	if cfg.Verbose {
		log.SetVerbose(true)
	}

	// Initialize cache
	cacheTTL := 1 * time.Hour
	cacheDir := cfg.CacheDir
	if cacheDir == "" {
		cacheDir = ".osint_cache"
	}

	c, err := cache.NewCache(cacheDir, cacheTTL)
	if err != nil {
		return nil, err
	}

	// Initialize worker pool
	rps := cfg.RateLimit.RequestsPerSecond
	if rps == 0 {
		rps = 10
	}

	burst := cfg.RateLimit.BurstSize
	if burst == 0 {
		burst = 20
	}

	retryAttempts := cfg.RateLimit.RetryAttempts
	if retryAttempts == 0 {
		retryAttempts = 3
	}

	retryDelay := cfg.RateLimit.RetryDelay
	if retryDelay == 0 {
		retryDelay = 2 * time.Second
	}

	pool := worker.NewPool(cfg.Workers, rps, burst, retryAttempts, retryDelay)

	return &Engine{
		config:  cfg,
		logger:  log,
		cache:   c,
		pool:    pool,
		modules: make(map[string]Module),
	}, nil
}

// RegisterModule registers an OSINT module
func (e *Engine) RegisterModule(module Module) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.modules[module.Name()] = module
	e.logger.Debugf("Registered module: %s", module.Name())
}

// GetModule retrieves a module by name
func (e *Engine) GetModule(name string) (Module, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	module, ok := e.modules[name]
	return module, ok
}

// ListModules returns all registered modules
func (e *Engine) ListModules() []Module {
	e.mu.RLock()
	defer e.mu.RUnlock()

	modules := make([]Module, 0, len(e.modules))
	for _, module := range e.modules {
		modules = append(modules, module)
	}
	return modules
}

// Execute executes a module on a target
func (e *Engine) Execute(ctx context.Context, moduleName string, target string) ([]output.Result, error) {
	module, ok := e.GetModule(moduleName)
	if !ok {
		e.logger.Errorf("Module not found: %s", moduleName)
		return nil, ErrModuleNotFound
	}

	e.logger.Infof("Executing module '%s' on target '%s'", moduleName, target)

	// Check cache first
	cacheKey := moduleName + ":" + target
	var cachedResults []output.Result
	found, err := e.cache.Get(cacheKey, &cachedResults)
	if err != nil {
		e.logger.Warnf("Cache error: %v", err)
	}

	if found {
		e.logger.Debugf("Cache hit for %s", cacheKey)
		return cachedResults, nil
	}

	// Execute module
	results, err := module.Execute(ctx, target)
	if err != nil {
		e.logger.Errorf("Module execution failed: %v", err)
		return nil, err
	}

	// Cache results
	if err := e.cache.Set(cacheKey, results); err != nil {
		e.logger.Warnf("Failed to cache results: %v", err)
	}

	e.logger.Infof("Module '%s' completed successfully with %d results", moduleName, len(results))
	return results, nil
}

// ExecuteMultiple executes multiple modules concurrently
func (e *Engine) ExecuteMultiple(ctx context.Context, moduleNames []string, target string) (map[string][]output.Result, error) {
	results := make(map[string][]output.Result)
	resultsMu := sync.Mutex{}
	var wg sync.WaitGroup

	e.pool.Start()
	defer e.pool.Stop()

	for _, moduleName := range moduleNames {
		moduleName := moduleName // capture loop variable

		wg.Add(1)
		e.pool.Submit(func(ctx context.Context) error {
			defer wg.Done()

			moduleResults, err := e.Execute(ctx, moduleName, target)
			if err != nil {
				e.logger.Errorf("Module '%s' failed: %v", moduleName, err)
				return err
			}

			resultsMu.Lock()
			results[moduleName] = moduleResults
			resultsMu.Unlock()

			return nil
		})
	}

	wg.Wait()
	return results, nil
}

// ExecuteAll executes all registered modules
func (e *Engine) ExecuteAll(ctx context.Context, target string) (map[string][]output.Result, error) {
	moduleNames := make([]string, 0, len(e.modules))
	for name := range e.modules {
		moduleNames = append(moduleNames, name)
	}

	return e.ExecuteMultiple(ctx, moduleNames, target)
}

// Cache returns the engine's cache
func (e *Engine) Cache() *cache.Cache {
	return e.cache
}

// Logger returns the engine's logger
func (e *Engine) Logger() *logger.Logger {
	return e.logger
}

// Config returns the engine's configuration
func (e *Engine) Config() *config.Config {
	return e.config
}

// Shutdown gracefully shuts down the engine
func (e *Engine) Shutdown() error {
	e.logger.Info("Shutting down engine...")
	e.pool.Stop()
	if err := e.cache.Cleanup(); err != nil {
		e.logger.Warnf("Cache cleanup failed: %v", err)
	}
	return nil
}

// ErrModuleNotFound is returned when a module is not found
var ErrModuleNotFound = NewEngineError("module not found")

// EngineError represents an engine error
type EngineError struct {
	message string
}

// NewEngineError creates a new engine error
func NewEngineError(message string) *EngineError {
	return &EngineError{message: message}
}

// Error returns the error message
func (e *EngineError) Error() string {
	return e.message
}
