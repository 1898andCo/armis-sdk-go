// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
)

func TestGetTags(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/tags/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"count": 3,
					"next":  nil,
					"prev":  0,
					"tags":  []string{"ARP Proxy", "Asset at risk", "Critical asset at risk"},
				},
				"success": true,
			})
		},
	})
	defer cleanup()

	tags, err := client.GetTags(context.Background())
	if err != nil {
		t.Fatalf("get tags: %v", err)
	}
	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(tags))
	}
	if tags[0] != "ARP Proxy" {
		t.Fatalf("unexpected first tag: %q", tags[0])
	}
}

func TestGetTags_EmptyList(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/tags/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"count": 0,
					"next":  nil,
					"prev":  0,
					"tags":  []string{},
				},
				"success": true,
			})
		},
	})
	defer cleanup()

	tags, err := client.GetTags(context.Background())
	if err != nil {
		t.Fatalf("get tags: %v", err)
	}
	if len(tags) != 0 {
		t.Fatalf("expected 0 tags, got %d", len(tags))
	}
}

func TestGetTags_Pagination(t *testing.T) {
	t.Parallel()

	var call atomic.Int32
	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/tags/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}

			n := call.Add(1)
			switch n {
			case 1:
				from := r.URL.Query().Get("from")
				if from != "" {
					t.Fatalf("expected no 'from' on first call, got %q", from)
				}
				next := 3
				respondJSON(t, w, http.StatusOK, map[string]any{
					"data": map[string]any{
						"count": 3,
						"next":  next,
						"prev":  0,
						"tags":  []string{"ARP Proxy", "Asset at risk", "Critical asset at risk"},
					},
					"success": true,
				})
			case 2:
				from := r.URL.Query().Get("from")
				if from != "3" {
					t.Fatalf("expected from=3 on second call, got %q", from)
				}
				respondJSON(t, w, http.StatusOK, map[string]any{
					"data": map[string]any{
						"count": 2,
						"next":  nil,
						"prev":  0,
						"tags":  []string{"DNS Server", "Nested Device"},
					},
					"success": true,
				})
			default:
				t.Fatalf("unexpected call number: %d", n)
			}
		},
	})
	defer cleanup()

	tags, err := client.GetTags(context.Background())
	if err != nil {
		t.Fatalf("get tags: %v", err)
	}
	if len(tags) != 5 {
		t.Fatalf("expected 5 tags, got %d", len(tags))
	}
	if tags[3] != "DNS Server" {
		t.Fatalf("unexpected fourth tag: %q", tags[3])
	}
}

func TestGetTags_NilContext(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	//nolint:staticcheck // SA1012 testing nil context handling
	if _, err := client.GetTags(nil); !errors.Is(err, ErrNilContext) {
		t.Fatalf("expected ErrNilContext, got %v", err)
	}
}

func TestGetTags_APIFailure(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/tags/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data":    map[string]any{},
				"success": false,
			})
		},
	})
	defer cleanup()

	_, err := client.GetTags(context.Background())
	if err == nil {
		t.Fatal("expected error for unsuccessful response")
	}
	if !errors.Is(err, ErrHTTPResponse) {
		t.Fatalf("expected ErrHTTPResponse, got %v", err)
	}
}
