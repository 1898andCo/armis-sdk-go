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

func TestGetUsers(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/users/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"users": []map[string]any{{
						"email": "user@example.com",
						"name":  "Example",
					}},
				},
				"success": true,
			})
		},
	})
	defer cleanup()

	users, err := client.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("get users: %v", err)
	}
	if len(users) != 1 || users[0].Email != "user@example.com" {
		t.Fatalf("unexpected users: %+v", users)
	}
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	const email = "user@example.com"
	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/users/user%40example.com/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"email": email,
					"name":  "Example",
				},
			})
		},
	})
	defer cleanup()

	user, err := client.GetUser(context.Background(), email)
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if user.Email != email {
		t.Fatalf("unexpected user: %+v", user)
	}
}

func TestGetUserRequiresID(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.GetUser(context.Background(), ""); !errors.Is(err, ErrUserID) {
		t.Fatalf("expected ErrUserID, got %v", err)
	}
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	payload := UserSettings{
		Name:     "Example",
		Email:    "user@example.com",
		Username: "example",
	}

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/users/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", r.Method)
			}
			var got UserSettings
			if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if got.Email != payload.Email {
				t.Fatalf("unexpected payload: %+v", got)
			}
			respondJSON(t, w, http.StatusCreated, map[string]any{
				"success": true,
				"data": map[string]any{
					"email":    payload.Email,
					"username": payload.Username,
				},
			})
		},
	})
	defer cleanup()

	user, err := client.CreateUser(context.Background(), payload)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if user.Email != payload.Email {
		t.Fatalf("unexpected response: %+v", user)
	}
}

func TestCreateUserValidation(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.CreateUser(context.Background(), UserSettings{}); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	payload := UserSettings{
		Name:     "Example",
		Email:    "user@example.com",
		Username: "updated",
	}

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/users/user%40example.com/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodPatch {
				t.Fatalf("expected PATCH, got %s", r.Method)
			}
			var got UserSettings
			if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if got.Username != payload.Username {
				t.Fatalf("unexpected payload: %+v", got)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"email":    payload.Email,
					"username": payload.Username,
				},
			})
		},
	})
	defer cleanup()

	user, err := client.UpdateUser(context.Background(), payload, payload.Email)
	if err != nil {
		t.Fatalf("update user: %v", err)
	}
	if user.Username != payload.Username {
		t.Fatalf("unexpected response: %+v", user)
	}
}

func TestUpdateUserValidation(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	payload := UserSettings{Name: "", Email: ""}
	if _, err := client.UpdateUser(context.Background(), payload, ""); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/users/user%40example.com/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodDelete {
				t.Fatalf("expected DELETE, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{"success": true})
		},
	})
	defer cleanup()

	ok, err := client.DeleteUser(context.Background(), "user@example.com")
	if err != nil {
		t.Fatalf("delete user: %v", err)
	}
	if !ok {
		t.Fatalf("expected delete success")
	}
}

func TestDeleteUserRequiresID(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.DeleteUser(context.Background(), ""); !errors.Is(err, ErrUserID) {
		t.Fatalf("expected ErrUserID, got %v", err)
	}
}
