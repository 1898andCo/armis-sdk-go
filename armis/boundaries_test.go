// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"net/http"
	"testing"
)

func TestGetBoundaries(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/boundaries/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"boundaries": []map[string]any{{"id": 1, "name": "Test"}},
				},
				"success": true,
			})
		},
	})
	defer cleanup()

	res, err := client.GetBoundaries(context.Background())
	if err != nil {
		t.Fatalf("get boundaries: %v", err)
	}
	if len(res) != 1 || res[0].Name != "Test" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestGetBoundaryByID(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/boundaries/1/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data":    map[string]any{"id": 1, "name": "Boundary"},
				"success": true,
			})
		},
	})
	defer cleanup()

	res, err := client.GetBoundaryByID(context.Background(), "1")
	if err != nil {
		t.Fatalf("get boundary by id: %v", err)
	}
	if res.Name != "Boundary" {
		t.Fatalf("unexpected response: %+v", res)
	}
}
