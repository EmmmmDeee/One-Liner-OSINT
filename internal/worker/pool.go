package worker

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Task represents a work item
type Task func(ctx context.Context) error

// Pool manages a pool of workers
type Pool struct {
	workers    int
	tasks      chan Task
	results    chan error
	limiter    *rate.Limiter
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	maxRetries int
	retryDelay time.Duration
}

// NewPool creates a new worker pool
func NewPool(workers int, rps int, burst int, maxRetries int, retryDelay time.Duration) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	return &Pool{
		workers:    workers,
		tasks:      make(chan Task, workers*2),
		results:    make(chan error, workers*2),
		limiter:    rate.NewLimiter(rate.Limit(rps), burst),
		ctx:        ctx,
		cancel:     cancel,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

// Start starts the worker pool
func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// worker processes tasks from the task channel
func (p *Pool) worker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.tasks:
			if !ok {
				return
			}

			// Rate limiting
			if err := p.limiter.Wait(p.ctx); err != nil {
				p.results <- err
				continue
			}

			// Execute task with retries
			var err error
			for attempt := 0; attempt <= p.maxRetries; attempt++ {
				if attempt > 0 {
					time.Sleep(p.retryDelay)
				}

				err = task(p.ctx)
				if err == nil {
					break
				}
			}

			p.results <- err
		}
	}
}

// Submit submits a task to the pool
func (p *Pool) Submit(task Task) {
	p.tasks <- task
}

// Wait waits for all tasks to complete and returns errors
func (p *Pool) Wait() []error {
	close(p.tasks)
	p.wg.Wait()
	close(p.results)

	var errors []error
	for err := range p.results {
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// Stop stops the worker pool
func (p *Pool) Stop() {
	p.cancel()
	p.wg.Wait()
}

// Results returns the results channel
func (p *Pool) Results() <-chan error {
	return p.results
}
