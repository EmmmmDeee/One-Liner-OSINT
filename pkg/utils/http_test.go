package utils

import (
	"testing"
	"context"
	"time"
)

func TestNewHTTPClient(t *testing.T) {
	timeout := 30 * time.Second
	userAgent := "Test-Agent/1.0"

	client := NewHTTPClient(timeout, userAgent)

	if client == nil {
		t.Fatal("Expected client to be created")
	}

	if client.userAgent != userAgent {
		t.Errorf("Expected user agent %s, got %s", userAgent, client.userAgent)
	}
}

func TestHTTPClientSetHeader(t *testing.T) {
	client := NewHTTPClient(30*time.Second, "Test-Agent")

	client.SetHeader("X-Custom-Header", "test-value")

	if client.headers["X-Custom-Header"] != "test-value" {
		t.Errorf("Expected header value 'test-value', got '%s'", client.headers["X-Custom-Header"])
	}
}

func TestHTTPClientGet(t *testing.T) {
	client := NewHTTPClient(30*time.Second, "Test-Agent")
	ctx := context.Background()

	// Test with a reliable endpoint
	resp, err := client.Get(ctx, "https://httpbin.org/user-agent")
	if err != nil {
		t.Skipf("Skipping test due to network error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}
