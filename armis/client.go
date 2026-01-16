// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

// Package armis provides a Go client for interacting with the Armis Centrix API.
// It is designed to be safe for concurrent use, idiomatic, and easy to extend.
package armis

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Defaults used when a caller omits configuration.
const (
	defaultAPIURL     = "https://api.armis.com"
	defaultAPIVersion = "v1"
)

// Config holds the inputs required to build a Client. Use functional options
// (With* helpers) to set values instead of mutating the struct directly.
//
// Example:
//
//	client, err := armis.NewClient("my-api-key",
//	    armis.WithAPIURL("https://staging-api.armis.com"),
//	    armis.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}))
type Config struct {
	APIKey     string
	APIURL     string
	apiVersion string
	HTTPClient *http.Client
}

// Option configures a Config. They are produced by With* helpers.
type Option func(*Config)

// WithAPIURL overrides the default API base URL.
func WithAPIURL(u string) Option { return func(c *Config) { c.APIURL = u } }

// WithAPIVersion overrides the default API version.
func WithAPIVersion(v string) Option { return func(c *Config) { c.apiVersion = v } }

// WithHTTPClient lets callers supply their own *http.Client (for custom timeouts,
// proxies, tracing, etc.).
func WithHTTPClient(h *http.Client) Option { return func(c *Config) { c.HTTPClient = h } }

// Client is a concurrency-safe Armis API client. Create it with NewClient.
// Do not instantiate it directly.
type Client struct {
	apiKey string

	apiURL     string
	apiVersion string

	httpClient *http.Client

	mu                 sync.RWMutex
	accessToken        string
	accessTokenExpires time.Time
	userID             int
}

// UserID returns the authenticated user's ID from the last successful
// authentication. Useful for external logging and audit trails.
func (c *Client) UserID() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.userID
}

// NewClient constructs a new Client. The first parameter (apiKey) is required.
// Optional parameters may be provided with functional options; see With* funcs.
//
// The function immediately performs authentication so the returned client is
// ready for use.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}

	cfg := &Config{
		APIKey:     apiKey,
		APIURL:     defaultAPIURL,
		apiVersion: defaultAPIVersion,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(cfg)
	}

	c := &Client{
		apiKey:     cfg.APIKey,
		apiURL:     cfg.APIURL,
		apiVersion: cfg.apiVersion,
		httpClient: cfg.HTTPClient,
	}

	if err := c.authenticate(context.Background()); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAuthFailed, err)
	}

	return c, nil
}

// APIError represents a non-2xx response from Armis. Code and Body are exposed
// so callers can inspect them programmatically.
type APIError struct {
	StatusCode int
	Body       []byte
}

// Error implements the error interface, returning a formatted error message
// containing the HTTP status code and its text description.
func (e *APIError) Error() string {
	return fmt.Sprintf("armis: API error %d: %s", e.StatusCode, http.StatusText(e.StatusCode))
}

// newRequest creates an *http.Request, applying authentication and common
// headers. The path should already include the API version prefix (e.g.
// "/v1/devices").
func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.apiURL+path, body)
	if err != nil {
		return nil, err
	}

	// For long running processes, tokens will expire with subsequent API calls.
	// The token needs to be validated before each request.
	c.mu.RLock()
	if c.accessTokenExpires.Before(time.Now()) {
		if err := c.authenticate(ctx); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrAuthFailed, err)
		}
	}
	token := c.accessToken
	c.mu.RUnlock()

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// doRequest executes the HTTP request and returns the response body for 2xx
// codes or an *APIError otherwise.
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return body, nil
	}

	return nil, &APIError{StatusCode: res.StatusCode, Body: body}
}
