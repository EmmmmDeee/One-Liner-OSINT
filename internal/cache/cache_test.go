package cache

import (
	"os"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	dir := "/tmp/test_cache"
	defer os.RemoveAll(dir)

	cache, err := NewCache(dir, 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	if cache.dir != dir {
		t.Errorf("Expected cache dir %s, got %s", dir, cache.dir)
	}
}

func TestCacheSetGet(t *testing.T) {
	dir := "/tmp/test_cache_setget"
	defer os.RemoveAll(dir)

	cache, err := NewCache(dir, 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Set a value
	key := "test-key"
	value := map[string]string{"foo": "bar"}

	err = cache.Set(key, value)
	if err != nil {
		t.Fatalf("Failed to set cache value: %v", err)
	}

	// Get the value
	var result map[string]string
	found, err := cache.Get(key, &result)
	if err != nil {
		t.Fatalf("Failed to get cache value: %v", err)
	}

	if !found {
		t.Fatal("Expected to find cached value")
	}

	if result["foo"] != "bar" {
		t.Errorf("Expected cached value 'bar', got '%s'", result["foo"])
	}
}

func TestCacheExpiration(t *testing.T) {
	dir := "/tmp/test_cache_expiration"
	defer os.RemoveAll(dir)

	// Create cache with 1 second TTL
	cache, err := NewCache(dir, 1*time.Second)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	key := "test-key"
	value := "test-value"

	err = cache.Set(key, value)
	if err != nil {
		t.Fatalf("Failed to set cache value: %v", err)
	}

	// Wait for expiration
	time.Sleep(2 * time.Second)

	var result string
	found, err := cache.Get(key, &result)
	if err != nil {
		t.Fatalf("Failed to get cache value: %v", err)
	}

	if found {
		t.Error("Expected cached value to be expired")
	}
}

func TestCacheClear(t *testing.T) {
	dir := "/tmp/test_cache_clear"
	defer os.RemoveAll(dir)

	cache, err := NewCache(dir, 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Set multiple values
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// Clear cache
	err = cache.Clear()
	if err != nil {
		t.Fatalf("Failed to clear cache: %v", err)
	}

	// Verify cache is empty
	var result string
	found, _ := cache.Get("key1", &result)
	if found {
		t.Error("Expected cache to be cleared")
	}
}
