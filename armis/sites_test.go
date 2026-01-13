// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"net/http"
	"testing"
)

func TestGetSites(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/sites/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"sites": []map[string]any{{"name": "Example", "siteId": 42}},
				},
			})
		},
	})
	defer cleanup()

	res, err := client.GetSites(context.Background())
	if err != nil {
		t.Fatalf("get sites: %v", err)
	}
	if len(res) != 1 || res[0].Name != "Example" {
		t.Fatalf("unexpected response: %+v", res)
	}
}
