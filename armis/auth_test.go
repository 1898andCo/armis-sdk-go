// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestAuthenticate_Success(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	// Verify client was authenticated successfully
	client.mu.RLock()
	token := client.accessToken
	expires := client.accessTokenExpires
	client.mu.RUnlock()

	if token != testToken {
		t.Fatalf("expected token %q, got %q", testToken, token)
	}
	if expires.IsZero() {
		t.Fatal("expected non-zero expiration time")
	}
}

func TestAuthenticate_InvalidAPIKey(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == authPath {
			respondJSON(t, w, http.StatusUnauthorized, map[string]any{
				"success": false,
				"error":   "Invalid API key",
			})
			return
		}
	}))
	defer server.Close()

	_, err := NewClient("invalid-key", WithAPIURL(server.URL), WithHTTPClient(server.Client()))
	if err == nil {
		t.Fatal("expected error for invalid API key")
	}
	if !errors.Is(err, ErrAuthFailed) {
		t.Fatalf("expected ErrAuthFailed, got %v", err)
	}
}

func TestAuthenticate_MalformedResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == authPath {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"invalid json`))
			return
		}
	}))
	defer server.Close()

	_, err := NewClient(testAPIKey, WithAPIURL(server.URL), WithHTTPClient(server.Client()))
	if err == nil {
		t.Fatal("expected error for malformed response")
	}
	// NewClient wraps all auth errors with ErrAuthFailed
	if !errors.Is(err, ErrAuthFailed) {
		t.Fatalf("expected ErrAuthFailed wrapper, got %v", err)
	}
}

func TestAuthenticate_MissingAccessToken(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == authPath {
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"expiration_utc": testExpiry,
				},
			})
			return
		}
	}))
	defer server.Close()

	client, err := NewClient(testAPIKey, WithAPIURL(server.URL), WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Token should be empty string (zero value)
	client.mu.RLock()
	token := client.accessToken
	client.mu.RUnlock()

	if token != "" {
		t.Fatalf("expected empty token, got %q", token)
	}
}

func TestAuthenticate_InvalidExpiration(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == authPath {
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"access_token":   testToken,
					"expiration_utc": "invalid-date",
				},
			})
			return
		}
	}))
	defer server.Close()

	_, err := NewClient(testAPIKey, WithAPIURL(server.URL), WithHTTPClient(server.Client()))
	if err == nil {
		t.Fatal("expected error for invalid expiration")
	}
	// NewClient wraps all auth errors with ErrAuthFailed
	if !errors.Is(err, ErrAuthFailed) {
		t.Fatalf("expected ErrAuthFailed wrapper, got %v", err)
	}
}

func TestAuthenticate_HTTPError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		statusText string
	}{
		{"BadRequest", http.StatusBadRequest, "Bad Request"},
		{"Unauthorized", http.StatusUnauthorized, "Unauthorized"},
		{"Forbidden", http.StatusForbidden, "Forbidden"},
		{"InternalServerError", http.StatusInternalServerError, "Internal Server Error"},
		{"ServiceUnavailable", http.StatusServiceUnavailable, "Service Unavailable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.EscapedPath() == authPath {
					w.WriteHeader(tt.statusCode)
					_, _ = w.Write([]byte(`{"error": "auth failed"}`))
					return
				}
			}))
			defer server.Close()

			_, err := NewClient(testAPIKey, WithAPIURL(server.URL), WithHTTPClient(server.Client()))
			if err == nil {
				t.Fatalf("expected error for status %d", tt.statusCode)
			}

			// NewClient wraps errors, so check for ErrAuthFailed
			if !errors.Is(err, ErrAuthFailed) {
				t.Fatalf("expected ErrAuthFailed wrapper, got %v", err)
			}

			// Verify the error message contains the status code and text
			errMsg := err.Error()
			expectedMsg := fmt.Sprintf("API error %d", tt.statusCode)
			if !strings.Contains(errMsg, expectedMsg) {
				t.Fatalf("expected error message to contain %q, got %q", expectedMsg, errMsg)
			}
			if !strings.Contains(errMsg, tt.statusText) {
				t.Fatalf("expected error message to contain %q, got %q", tt.statusText, errMsg)
			}
		})
	}
}

func TestAuthenticate_TokenCaching(t *testing.T) {
	t.Parallel()

	authCallCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == authPath {
			authCallCount++
			if err := r.ParseForm(); err != nil {
				t.Fatalf("parse auth form: %v", err)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"access_token":   testToken,
					"expiration_utc": testExpiry,
				},
			})
			return
		}
	}))
	defer server.Close()

	client, err := NewClient(testAPIKey, WithAPIURL(server.URL), WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	if authCallCount != 1 {
		t.Fatalf("expected 1 auth call during NewClient, got %d", authCallCount)
	}

	// Call authenticate multiple times - should use cached token
	for i := range 5 {
		if err := client.authenticate(context.Background()); err != nil {
			t.Fatalf("authenticate call %d: %v", i, err)
		}
	}

	if authCallCount != 1 {
		t.Fatalf("expected 1 auth call total (cached), got %d", authCallCount)
	}
}

func TestAuthenticate_TokenExpiration(t *testing.T) {
	t.Parallel()

	authCallCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == authPath {
			authCallCount++
			if err := r.ParseForm(); err != nil {
				t.Fatalf("parse auth form: %v", err)
			}
			// Return a token that expires very soon
			expiry := time.Now().Add(1 * time.Millisecond).Format(time.RFC3339Nano)
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"access_token":   testToken,
					"expiration_utc": expiry,
				},
			})
			return
		}
	}))
	defer server.Close()

	client, err := NewClient(testAPIKey, WithAPIURL(server.URL), WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	if authCallCount != 1 {
		t.Fatalf("expected 1 auth call during NewClient, got %d", authCallCount)
	}

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Next authenticate call should refresh the token
	if err := client.authenticate(context.Background()); err != nil {
		t.Fatalf("authenticate after expiry: %v", err)
	}

	if authCallCount != 2 {
		t.Fatalf("expected 2 auth calls (initial + refresh), got %d", authCallCount)
	}
}

func TestAuthenticate_Concurrent(t *testing.T) {
	t.Parallel()

	authCallCount := 0
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == authPath {
			mu.Lock()
			authCallCount++
			mu.Unlock()

			if err := r.ParseForm(); err != nil {
				t.Fatalf("parse auth form: %v", err)
			}

			// Simulate slow auth endpoint
			time.Sleep(10 * time.Millisecond)

			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"access_token":   testToken,
					"expiration_utc": testExpiry,
				},
			})
			return
		}
	}))
	defer server.Close()

	client, err := NewClient(testAPIKey, WithAPIURL(server.URL), WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	mu.Lock()
	initialAuthCalls := authCallCount
	mu.Unlock()

	// Launch multiple concurrent authenticate calls
	const numGoroutines = 10
	var wg sync.WaitGroup
	errChan := make(chan error, numGoroutines)

	for range numGoroutines {
		wg.Go(func() {
			if err := client.authenticate(context.Background()); err != nil {
				errChan <- err
			}
		})
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		t.Fatalf("concurrent authenticate error: %v", err)
	}

	// Should still only have initial auth call (token cached)
	mu.Lock()
	finalAuthCalls := authCallCount
	mu.Unlock()

	if finalAuthCalls != initialAuthCalls {
		t.Fatalf("expected %d auth calls (cached), got %d", initialAuthCalls, finalAuthCalls)
	}
}

func TestAuthenticate_ContextCancellation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == authPath {
			// Simulate slow response
			time.Sleep(100 * time.Millisecond)
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"access_token":   testToken,
					"expiration_utc": testExpiry,
				},
			})
			return
		}
	}))
	defer server.Close()

	// Create client with custom HTTP client that has short timeout
	httpClient := &http.Client{Timeout: 10 * time.Millisecond}

	_, err := NewClient(testAPIKey, WithAPIURL(server.URL), WithHTTPClient(httpClient))
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestAuthenticate_SuccessFalse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == authPath {
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": false,
				"data": map[string]any{
					"access_token":   testToken,
					"expiration_utc": testExpiry,
				},
			})
			return
		}
	}))
	defer server.Close()

	_, err := NewClient(testAPIKey, WithAPIURL(server.URL), WithHTTPClient(server.Client()))
	if err == nil {
		t.Fatal("expected error when success is false")
	}
	if !errors.Is(err, ErrAuthFailed) {
		t.Fatalf("expected ErrAuthFailed, got %v", err)
	}
}

func TestAuthenticate_ExpiryWithSafetyMargin(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == authPath {
			if err := r.ParseForm(); err != nil {
				t.Fatalf("parse auth form: %v", err)
			}
			// Token expires in 10 minutes
			expiry := time.Now().Add(10 * time.Minute).Format(time.RFC3339Nano)
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"access_token":   testToken,
					"expiration_utc": expiry,
				},
			})
			return
		}
	}))
	defer server.Close()

	client, err := NewClient(testAPIKey, WithAPIURL(server.URL), WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	// Check that expiry has 5-minute safety margin applied
	client.mu.RLock()
	expires := client.accessTokenExpires
	client.mu.RUnlock()

	// Should expire in approximately 5 minutes (10 - 5 minute safety margin)
	timeUntilExpiry := time.Until(expires)
	if timeUntilExpiry < 4*time.Minute || timeUntilExpiry > 6*time.Minute {
		t.Fatalf("expected expiry in ~5 minutes, got %v", timeUntilExpiry)
	}
}
