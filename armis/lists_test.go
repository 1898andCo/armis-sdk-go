// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"net/http"
	"testing"
)

func TestGetLists(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/lists/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"lists": []map[string]any{
						{
							"list_id":          1,
							"list_name":        "Test List",
							"list_type":        "IP",
							"description":      "A test list",
							"created_by":       "admin@example.com",
							"creation_time":    "2024-01-01T00:00:00Z",
							"last_updated_by":  "admin@example.com",
							"last_update_time": "2024-01-02T00:00:00Z",
						},
						{
							"list_id":          2,
							"list_name":        "Another List",
							"list_type":        "MAC",
							"description":      "Another test list",
							"created_by":       "user@example.com",
							"creation_time":    "2024-01-03T00:00:00Z",
							"last_updated_by":  "user@example.com",
							"last_update_time": "2024-01-04T00:00:00Z",
						},
					},
				},
			})
		},
	})
	defer cleanup()

	res, err := client.GetLists(context.Background())
	if err != nil {
		t.Fatalf("get lists: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 lists, got %d", len(res))
	}
	if res[0].ListName != "Test List" {
		t.Fatalf("unexpected first list name: %s", res[0].ListName)
	}
	if res[0].ListID != 1 {
		t.Fatalf("unexpected first list ID: %d", res[0].ListID)
	}
	if res[1].ListName != "Another List" {
		t.Fatalf("unexpected second list name: %s", res[1].ListName)
	}
}

func TestGetLists_EmptyList(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/lists/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"lists": []map[string]any{},
				},
			})
		},
	})
	defer cleanup()

	res, err := client.GetLists(context.Background())
	if err != nil {
		t.Fatalf("get lists: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("expected 0 lists, got %d", len(res))
	}
}

func TestGetLists_APIError(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/lists/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			respondJSON(t, w, http.StatusInternalServerError, map[string]any{
				"success": false,
				"error":   "internal server error",
			})
		},
	})
	defer cleanup()

	_, err := client.GetLists(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
