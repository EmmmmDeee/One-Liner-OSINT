package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Cache provides file-based caching
type Cache struct {
	dir string
	ttl time.Duration
}

// CacheEntry represents a cached item
type CacheEntry struct {
	Key       string      `json:"key"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	TTL       int64       `json:"ttl"`
}

// NewCache creates a new cache instance
func NewCache(dir string, ttl time.Duration) (*Cache, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return &Cache{
		dir: dir,
		ttl: ttl,
	}, nil
}

// Get retrieves a value from cache
func (c *Cache) Get(key string, result interface{}) (bool, error) {
	path := c.getPath(key)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return false, err
	}

	// Check if cache entry is expired
	if time.Since(entry.Timestamp) > time.Duration(entry.TTL)*time.Second {
		os.Remove(path)
		return false, nil
	}

	// Unmarshal data into result
	dataBytes, err := json.Marshal(entry.Data)
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal(dataBytes, result); err != nil {
		return false, err
	}

	return true, nil
}

// Set stores a value in cache
func (c *Cache) Set(key string, value interface{}) error {
	entry := CacheEntry{
		Key:       key,
		Data:      value,
		Timestamp: time.Now(),
		TTL:       int64(c.ttl.Seconds()),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	path := c.getPath(key)
	return os.WriteFile(path, data, 0644)
}

// Delete removes a value from cache
func (c *Cache) Delete(key string) error {
	path := c.getPath(key)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Clear removes all cached items
func (c *Cache) Clear() error {
	return os.RemoveAll(c.dir)
}

// getPath returns the file path for a cache key
func (c *Cache) getPath(key string) string {
	hash := sha256.Sum256([]byte(key))
	filename := hex.EncodeToString(hash[:])
	return filepath.Join(c.dir, filename+".json")
}

// Cleanup removes expired cache entries
func (c *Cache) Cleanup() error {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(c.dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var cacheEntry CacheEntry
		if err := json.Unmarshal(data, &cacheEntry); err != nil {
			continue
		}

		// Remove expired entries
		if time.Since(cacheEntry.Timestamp) > time.Duration(cacheEntry.TTL)*time.Second {
			os.Remove(path)
		}
	}

	return nil
}

// GetOrSet retrieves a value from cache or sets it if not found
func (c *Cache) GetOrSet(key string, result interface{}, setter func() (interface{}, error)) error {
	// Try to get from cache
	found, err := c.Get(key, result)
	if err != nil {
		return err
	}

	if found {
		return nil
	}

	// Not in cache, call setter
	value, err := setter()
	if err != nil {
		return err
	}

	// Store in cache
	if err := c.Set(key, value); err != nil {
		return err
	}

	// Marshal and unmarshal to populate result
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, result)
}

// Stats returns cache statistics
func (c *Cache) Stats() (map[string]interface{}, error) {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_entries": len(entries),
		"cache_dir":     c.dir,
		"ttl_seconds":   int64(c.ttl.Seconds()),
	}

	var totalSize int64
	var expiredCount int
	now := time.Now()

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}
		totalSize += info.Size()

		// Check if expired
		path := filepath.Join(c.dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var cacheEntry CacheEntry
		if err := json.Unmarshal(data, &cacheEntry); err != nil {
			continue
		}

		if now.Sub(cacheEntry.Timestamp) > time.Duration(cacheEntry.TTL)*time.Second {
			expiredCount++
		}
	}

	stats["total_size_bytes"] = totalSize
	stats["expired_entries"] = expiredCount

	return stats, nil
}
