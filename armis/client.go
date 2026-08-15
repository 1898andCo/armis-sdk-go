// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

// Package armis provides a Go client for interacting with the Armis Centrix API.
// It is designed to be safe for concurrent use, idiomatic, and easy to extend.
package armis

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Defaults used when a caller omits configuration.
const (
	defaultAPIVersion = "v1"
)

// Config holds the inputs required to build a Client. Use functional options
// (With* helpers) to set values instead of mutating the struct directly.
//
// Example:
//
//	client, err := armis.NewClient("my-api-key", "https://staging-api.armis.com",
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
func NewClient(apiKey string, apiURL string, opts ...Option) (*Client, error) {
	cfg := &Config{
		APIKey:     apiKey,
		APIURL:     apiURL,
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

	if err := c.validateClient(); err != nil {
		return nil, err
	}

	if err := c.authenticate(context.Background()); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAuthFailed, err)
	}

	return c, nil
}

// validateClient checks that the Client is in a usable state.
func (c *Client) validateClient() error {
	if c.apiKey == "" {
		return ErrNoAPIKey
	}

	if c.apiURL == "" {
		return ErrNoAPIURL
	}

	if c.httpClient == nil {
		return ErrNoHTTPClient
	}
	return nil
}

// validateRequest checks that the request parameters are valid.
func (c *Client) validateRequest(ctx context.Context, method, path string) error {
	if ctx == nil {
		return ErrNilContext
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if method == "" {
		return ErrEmptyMethod
	}
	if path == "" {
		return ErrEmptyPath
	}
	return nil
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
	if err := c.validateRequest(ctx, method, path); err != nil {
		return nil, err
	}
	if err := c.validateClient(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, c.apiURL+path, body)
	if err != nil {
		return nil, err
	}

	// For long running processes, tokens will expire with subsequent API calls.
	// The token needs to be validated before each request.
	c.mu.RLock()
	expired := c.accessTokenExpires.Before(time.Now())
	c.mu.RUnlock()

	if expired {
		if err := c.authenticate(ctx); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrAuthFailed, err)
		}
	}

	c.mu.RLock()
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
// codes or an *APIError otherwise. It automatically retries once on 401
// Unauthorized by refreshing the access token.
func (c *Client) doRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	return c.doRequestWithRetry(ctx, req, true)
}

// doRequestWithRetry is the internal implementation that tracks whether a
// retry is allowed. This prevents infinite retry loops on persistent 401s.
func (c *Client) doRequestWithRetry(ctx context.Context, req *http.Request, canRetry bool) ([]byte, error) {
	// Buffer the request body so we can replay it on retry. This is necessary
	// because the body reader is consumed during the first attempt.
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	res, err := c.httpClient.Do(req) //nolint:gosec // G704: URL is constructed from caller-configured apiURL, not user input
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized && canRetry {
		// The server can reject a token before its local expiry: Armis keeps a
		// single active token per secret key, so another consumer
		// authenticating with the same key invalidates ours. Drop the cached
		// token first so authenticate fetches a fresh one instead of
		// early-returning on the unexpired cache, then retry once.
		c.invalidateToken(req.Header.Get("Authorization"))
		if err := c.authenticate(ctx); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrAuthFailed, err)
		}
		c.mu.RLock()
		req.Header.Set("Authorization", c.accessToken)
		c.mu.RUnlock()
		// Reset the request body for the retry attempt
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
		return c.doRequestWithRetry(ctx, req, false)
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return body, nil
	}

	return nil, &APIError{StatusCode: res.StatusCode, Body: body}
}
