// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestGetRoles(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/roles/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": []map[string]any{{
					"roleId": 1,
					"name":   "Example",
				}},
			})
		},
	})
	defer cleanup()

	roles, err := client.GetRoles(context.Background())
	if err != nil {
		t.Fatalf("get roles: %v", err)
	}
	if len(roles) != 1 || roles[0].Name != "Example" {
		t.Fatalf("unexpected roles: %+v", roles)
	}
}

func TestGetRoleByName(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/roles/": func(w http.ResponseWriter, r *http.Request) {
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": []map[string]any{{
					"roleId": 2,
					"name":   "Target",
				}},
			})
		},
	})
	defer cleanup()

	role, err := client.GetRoleByName(context.Background(), "Target")
	if err != nil {
		t.Fatalf("get role by name: %v", err)
	}
	if role.Name != "Target" {
		t.Fatalf("unexpected role: %+v", role)
	}
}

func TestGetRoleByNameRequiresName(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.GetRoleByName(context.Background(), ""); !errors.Is(err, ErrRoleName) {
		t.Fatalf("expected ErrRoleName, got %v", err)
	}
}

func TestGetRoleByID(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/roles/": func(w http.ResponseWriter, r *http.Request) {
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": []map[string]any{{
					"roleId": 3,
					"name":   "Example",
				}},
			})
		},
	})
	defer cleanup()

	role, err := client.GetRoleByID(context.Background(), "3")
	if err != nil {
		t.Fatalf("get role by id: %v", err)
	}
	if role.ID != 3 {
		t.Fatalf("unexpected role: %+v", role)
	}
}

func TestGetRoleByIDRequiresID(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.GetRoleByID(context.Background(), ""); !errors.Is(err, ErrRoleID) {
		t.Fatalf("expected ErrRoleID, got %v", err)
	}
}

func TestCreateRole(t *testing.T) {
	t.Parallel()

	role := RoleSettings{Name: "Example"}
	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/roles/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", r.Method)
			}
			var got RoleSettings
			if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if got.Name != role.Name {
				t.Fatalf("unexpected role name: %q", got.Name)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{"success": true})
		},
	})
	defer cleanup()

	ok, err := client.CreateRole(context.Background(), role)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}
	if !ok {
		t.Fatalf("expected create success")
	}
}

func TestCreateRoleRequiresName(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.CreateRole(context.Background(), RoleSettings{}); !errors.Is(err, ErrRoleName) {
		t.Fatalf("expected ErrRoleName, got %v", err)
	}
}

func TestUpdateRole(t *testing.T) {
	t.Parallel()

	role := RoleSettings{Name: "Updated"}
	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/roles/1/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodPatch {
				t.Fatalf("expected PATCH, got %s", r.Method)
			}
			var got RoleSettings
			if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if got.Name != "Updated" {
				t.Fatalf("unexpected body: %+v", got)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"name":   "Updated",
				"roleId": 1,
			})
		},
	})
	defer cleanup()

	res, err := client.UpdateRole(context.Background(), role, "1")
	if err != nil {
		t.Fatalf("update role: %v", err)
	}
	if res.Name != "Updated" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestUpdateRoleRequiresID(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.UpdateRole(context.Background(), RoleSettings{}, ""); !errors.Is(err, ErrRoleID) {
		t.Fatalf("expected ErrRoleID, got %v", err)
	}
}

func TestDeleteRole(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/roles/1/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodDelete {
				t.Fatalf("expected DELETE, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{"success": true})
		},
	})
	defer cleanup()

	ok, err := client.DeleteRole(context.Background(), "1")
	if err != nil {
		t.Fatalf("delete role: %v", err)
	}
	if !ok {
		t.Fatalf("expected delete success")
	}
}

func TestDeleteRoleRequiresID(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.DeleteRole(context.Background(), ""); !errors.Is(err, ErrRoleID) {
		t.Fatalf("expected ErrRoleID, got %v", err)
	}
}
