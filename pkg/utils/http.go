package utils

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

// HTTPClient is a wrapper around http.Client with additional features
type HTTPClient struct {
	client    *http.Client
	userAgent string
	headers   map[string]string
	proxy     *url.URL
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(timeout time.Duration, userAgent string) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
		},
		userAgent: userAgent,
		headers:   make(map[string]string),
	}
}

// SetProxy sets the proxy for the HTTP client
func (c *HTTPClient) SetProxy(proxyURL string) error {
	if proxyURL == "" {
		return nil
	}

	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}

	c.proxy = proxy
	c.client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxy)
	return nil
}

// SetHeader sets a custom header
func (c *HTTPClient) SetHeader(key, value string) {
	c.headers[key] = value
}

// Get performs a GET request
func (c *HTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set user agent
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	// Set custom headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	return c.client.Do(req)
}

// Post performs a POST request
func (c *HTTPClient) Post(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	// Set user agent
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	// Set content type
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// Set custom headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	return c.client.Do(req)
}

// Do performs a custom HTTP request
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Set user agent if not already set
	if req.Header.Get("User-Agent") == "" && c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	// Set custom headers if not already set
	for k, v := range c.headers {
		if req.Header.Get(k) == "" {
			req.Header.Set(k, v)
		}
	}

	return c.client.Do(req)
}

// Close closes idle connections
func (c *HTTPClient) Close() {
	c.client.CloseIdleConnections()
}
