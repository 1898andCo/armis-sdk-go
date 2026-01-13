// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testAPIKey = "test-api-key"
	testToken  = "Bearer test-token"
	testExpiry = "2099-01-01T00:00:00Z"
	authPath   = "/api/v1/access_token/"
)

func newTestClient(t *testing.T, handlers map[string]http.HandlerFunc) (*Client, func()) {
	t.Helper()

	if handlers == nil {
		handlers = make(map[string]http.HandlerFunc)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == authPath {
			if err := r.ParseForm(); err != nil {
				t.Fatalf("parse auth form: %v", err)
			}
			if r.Method != http.MethodPost {
				t.Fatalf("unexpected auth method: %s", r.Method)
			}
			if got := r.Form.Get("secret_key"); got != testAPIKey {
				t.Fatalf("unexpected secret key: %q", got)
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

		path := r.URL.EscapedPath()
		handler, ok := handlers[path]
		if !ok {
			t.Fatalf("unexpected path: %s", path)
		}
		handler(w, r)
	}))

	client, err := NewClient(testAPIKey, WithAPIURL(server.URL), WithHTTPClient(server.Client()))
	if err != nil {
		server.Close()
		t.Fatalf("new client: %v", err)
	}

	return client, server.Close
}

func respondJSON(t *testing.T, w http.ResponseWriter, status int, body any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}

func assertAuthHeader(t *testing.T, r *http.Request) {
	t.Helper()
	if got := r.Header.Get("Authorization"); got != testToken {
		t.Fatalf("unexpected Authorization header: %q", got)
	}
}
