// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetCollectors(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/collectors/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"collectors": []map[string]any{{
						"collectorNumber": 1,
						"name":            "Primary",
					}},
				},
				"success": true,
			})
		},
	})
	defer cleanup()

	res, err := client.GetCollectors(context.Background())
	if err != nil {
		t.Fatalf("get collectors: %v", err)
	}
	if len(res) != 1 || res[0].Name != "Primary" {
		t.Fatalf("unexpected collectors: %+v", res)
	}
}

func TestGetCollectorByID(t *testing.T) {
	t.Parallel()

	id := "collector-1"
	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/collectors/collector-1/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{"name": "Primary"},
			})
		},
	})
	defer cleanup()

	res, err := client.GetCollectorByID(context.Background(), id)
	if err != nil {
		t.Fatalf("get collector by id: %v", err)
	}
	if res.Name != "Primary" {
		t.Fatalf("unexpected collector: %+v", res)
	}
}

func TestCreateCollector(t *testing.T) {
	t.Parallel()

	payload := CreateCollectorSettings{Name: "New", DeploymentType: "OVA"}
	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/collectors/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", r.Method)
			}
			var got CreateCollectorSettings
			if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if got != payload {
				t.Fatalf("unexpected payload: %+v", got)
			}
			respondJSON(t, w, http.StatusCreated, map[string]any{
				"data": map[string]any{
					"collectorId": 99,
					"licenseKey":  "abc",
				},
				"success": true,
			})
		},
	})
	defer cleanup()

	res, err := client.CreateCollector(context.Background(), payload)
	if err != nil {
		t.Fatalf("create collector: %v", err)
	}
	if res.CollectorID != 99 {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestCreateCollectorValidation(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.CreateCollector(context.Background(), CreateCollectorSettings{}); err == nil {
		t.Fatalf("expected error for missing name and type")
	}
}

func TestUpdateCollector(t *testing.T) {
	t.Parallel()

	payload := UpdateCollectorSettings{Name: "Updated", DeploymentType: "OVA"}
	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/collectors/collector-1/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodPatch {
				t.Fatalf("expected PATCH, got %s", r.Method)
			}
			var got UpdateCollectorSettings
			if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if got != payload {
				t.Fatalf("unexpected payload: %+v", got)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{"name": "Updated"},
			})
		},
	})
	defer cleanup()

	res, err := client.UpdateCollector(context.Background(), "collector-1", payload)
	if err != nil {
		t.Fatalf("update collector: %v", err)
	}
	if res.Name != "Updated" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestDeleteCollector(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/collectors/collector-1/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodDelete {
				t.Fatalf("expected DELETE, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{"success": true})
		},
	})
	defer cleanup()

	ok, err := client.DeleteCollector(context.Background(), "collector-1")
	if err != nil {
		t.Fatalf("delete collector: %v", err)
	}
	if !ok {
		t.Fatalf("expected success deleting collector")
	}
}
