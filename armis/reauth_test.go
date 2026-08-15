// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// newRotatingTokenServer returns a test server whose auth endpoint issues a
// new token on every call ("token-1", "token-2", ...) and routes every other
// path to handler. authCalls reports how many times authentication happened.
func newRotatingTokenServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *atomic.Int32) {
	t.Helper()

	var authCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == authPath {
			n := authCalls.Add(1)
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"access_token":   fmt.Sprintf("token-%d", n),
					"expiration_utc": testExpiry,
				},
			})
			return
		}
		handler(w, r)
	}))
	return server, &authCalls
}

// TestDoRequest_ReauthenticatesWhenTokenInvalidatedServerSide reproduces the
// production failure behind MSSCI-19407: the cached token is rejected by the
// server while its local expiry is still in the future (Armis invalidates the
// previous token whenever a new one is issued for the same secret key). The
// 401 retry must fetch a fresh token rather than replaying the dead one.
func TestDoRequest_ReauthenticatesWhenTokenInvalidatedServerSide(t *testing.T) {
	t.Parallel()

	server, authCalls := newRotatingTokenServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("Authorization") {
		case "token-1":
			// token-1 was invalidated server-side.
			respondJSON(t, w, http.StatusUnauthorized, map[string]any{"success": false})
		case "token-2":
			respondJSON(t, w, http.StatusOK, map[string]any{"success": true})
		default:
			t.Errorf("unexpected Authorization header: %q", r.Header.Get("Authorization"))
			respondJSON(t, w, http.StatusUnauthorized, map[string]any{"success": false})
		}
	})
	defer server.Close()

	client, err := NewClient(testAPIKey, server.URL, WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	// The client holds token-1 and believes it valid until 2099.
	client.mu.RLock()
	cached := client.accessToken
	client.mu.RUnlock()
	if cached != "token-1" {
		t.Fatalf("expected cached token-1, got %q", cached)
	}

	req, err := client.newRequest(context.Background(), http.MethodGet, "/api/v1/probe/", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if _, err := client.doRequest(context.Background(), req); err != nil {
		t.Fatalf("expected request to succeed after forced re-auth, got %v", err)
	}

	if got := authCalls.Load(); got != 2 {
		t.Fatalf("expected 2 auth calls (initial + forced refresh), got %d", got)
	}
	client.mu.RLock()
	final := client.accessToken
	client.mu.RUnlock()
	if final != "token-2" {
		t.Fatalf("expected cached token-2 after re-auth, got %q", final)
	}
}

// TestDoRequest_PersistentUnauthorizedFailsAfterOneRetry ensures the forced
// re-auth still performs exactly one retry: a persistently rejecting server
// must produce an APIError, not an infinite auth/retry loop.
func TestDoRequest_PersistentUnauthorizedFailsAfterOneRetry(t *testing.T) {
	t.Parallel()

	server, authCalls := newRotatingTokenServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(t, w, http.StatusUnauthorized, map[string]any{"success": false})
	})
	defer server.Close()

	client, err := NewClient(testAPIKey, server.URL, WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	req, err := client.newRequest(context.Background(), http.MethodGet, "/api/v1/probe/", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	_, err = client.doRequest(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for persistent 401")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected APIError with status 401, got %v", err)
	}
	if got := authCalls.Load(); got != 2 {
		t.Fatalf("expected 2 auth calls (initial + one forced refresh), got %d", got)
	}
}

// TestInvalidateToken_KeepsNewerToken verifies the stale-token guard: a 401
// handler must not discard a token that another goroutine already refreshed.
func TestInvalidateToken_KeepsNewerToken(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	client.mu.Lock()
	client.accessToken = "token-new"
	client.accessTokenExpires = time.Now().Add(time.Hour)
	client.mu.Unlock()

	// Invalidation keyed to a token that is no longer cached is a no-op.
	client.invalidateToken("token-old")
	client.mu.RLock()
	kept := client.accessToken
	client.mu.RUnlock()
	if kept != "token-new" {
		t.Fatalf("expected newer token to survive, got %q", kept)
	}

	// Invalidation keyed to the cached token clears it.
	client.invalidateToken("token-new")
	client.mu.RLock()
	cleared := client.accessToken
	expires := client.accessTokenExpires
	client.mu.RUnlock()
	if cleared != "" || !expires.IsZero() {
		t.Fatalf("expected token cleared, got %q (expires %v)", cleared, expires)
	}
}
