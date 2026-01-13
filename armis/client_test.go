// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/ping": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			respondJSON(t, w, http.StatusOK, map[string]string{"message": "pong"})
		},
	})
	defer cleanup()

	req, err := client.newRequest(context.Background(), http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatalf("newRequest failed: %v", err)
	}

	if got := req.Header.Get("Authorization"); got != testToken {
		t.Fatalf("expected Authorization header %q, got %q", testToken, got)
	}
}

func TestNewClientRequiresAPIKey(t *testing.T) {
	t.Parallel()

	c, err := NewClient("")
	if !errors.Is(err, ErrNoAPIKey) {
		t.Fatalf("expected ErrNoAPIKey, got %v", err)
	}
	if c != nil {
		t.Fatalf("expected nil client when no API key provided")
	}
}
