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

func TestPolicyValidateSuccess(t *testing.T) {
	t.Parallel()

	policy := validPolicy()
	if err := policy.Validate(); err != nil {
		t.Fatalf("expected validation success, got %v", err)
	}
}

func TestPolicyValidateErrors(t *testing.T) {
	t.Parallel()

	policy := PolicySettings{}
	if err := policy.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestValidateRejectsWhitespaceOnlyName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		policyName string
		wantErr    bool
	}{
		{"empty string", "", true},
		{"whitespace only", "   ", true},
		{"tabs only", "\t\t", true},
		{"mixed whitespace", " \t \n ", true},
		{"valid name", "My Policy", false},
		{"name with surrounding whitespace", "  My Policy  ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			policy := validPolicy()
			policy.Name = tt.policyName

			err := policy.Validate()
			if tt.wantErr && !errors.Is(err, ErrPolicyName) {
				t.Errorf("expected ErrPolicyName for name %q, got %v", tt.policyName, err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error for name %q, got %v", tt.policyName, err)
			}
		})
	}
}

func TestValidateRejectsLongDescription(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		wantErr     bool
	}{
		{"empty description", "", false},
		{"short description", "A short description", false},
		{"exactly 500 chars", string(make([]byte, 500)), false},
		{"501 chars", string(make([]byte, 501)), true},
		{"very long description", string(make([]byte, 1000)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			policy := validPolicy()
			// Fill description with 'a' characters for readable test
			if len(tt.description) > 0 && tt.description[0] == 0 {
				desc := make([]byte, len(tt.description))
				for i := range desc {
					desc[i] = 'a'
				}
				policy.Description = string(desc)
			} else {
				policy.Description = tt.description
			}

			err := policy.Validate()
			if tt.wantErr && !errors.Is(err, ErrPolicyDescription) {
				t.Errorf("expected ErrPolicyDescription for length %d, got %v", len(policy.Description), err)
			}
			if !tt.wantErr && errors.Is(err, ErrPolicyDescription) {
				t.Errorf("expected no description error for length %d, got %v", len(policy.Description), err)
			}
		})
	}
}

func TestValidateRejectsInvalidRuleType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		ruleType string
		wantErr  bool
	}{
		{"ACTIVITY", "ACTIVITY", false},
		{"IP_CONNECTION", "IP_CONNECTION", false},
		{"DEVICE", "DEVICE", false},
		{"VULNERABILITY", "VULNERABILITY", false},
		{"empty string", "", true},
		{"lowercase activity", "activity", true},
		{"invalid type", "INVALID", true},
		{"partial match", "ACTIV", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			policy := validPolicy()
			policy.RuleType = tt.ruleType

			err := policy.Validate()
			if tt.wantErr && !errors.Is(err, ErrPolicyRuleType) {
				t.Errorf("expected ErrPolicyRuleType for ruleType %q, got %v", tt.ruleType, err)
			}
			if !tt.wantErr && errors.Is(err, ErrPolicyRuleType) {
				t.Errorf("expected no ruleType error for %q, got %v", tt.ruleType, err)
			}
		})
	}
}

func TestCreatePolicy(t *testing.T) {
	t.Parallel()

	policy := validPolicy()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/policies/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", r.Method)
			}
			var got PolicySettings
			if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if got.Name != policy.Name {
				t.Fatalf("unexpected policy name: %q", got.Name)
			}
			respondJSON(t, w, http.StatusCreated, map[string]any{
				"success": true,
				"data":    map[string]any{"id": 123},
			})
		},
	})
	defer cleanup()

	id, err := client.CreatePolicy(context.Background(), policy)
	if err != nil {
		t.Fatalf("create policy: %v", err)
	}
	if id.ID != 123 {
		t.Fatalf("unexpected policy id: %+v", id)
	}
}

func TestCreatePolicyValidationFailure(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.CreatePolicy(context.Background(), PolicySettings{}); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestGetPolicy(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/policies/1/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"name":     "Example",
					"ruleType": "ACTIVITY",
					"rules": map[string]any{
						"and": []any{"foo"},
					},
				},
			})
		},
	})
	defer cleanup()

	policy, err := client.GetPolicy(context.Background(), "1")
	if err != nil {
		t.Fatalf("get policy: %v", err)
	}
	if policy.Name != "Example" {
		t.Fatalf("unexpected policy: %+v", policy)
	}
}

func TestGetPolicyRequiresID(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.GetPolicy(context.Background(), ""); !errors.Is(err, ErrPolicyID) {
		t.Fatalf("expected ErrPolicyID, got %v", err)
	}
}

func TestGetAllPolicies(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/policies/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"count": 1,
					"next":  nil,
					"prev":  0,
					"total": 1,
					"policies": []map[string]any{{
						"id":       "1",
						"name":     "Example",
						"ruleType": "ACTIVITY",
						"rules": map[string]any{
							"and": []any{"foo"},
						},
					}},
				},
			})
		},
	})
	defer cleanup()

	policies, err := client.GetAllPolicies(context.Background())
	if err != nil {
		t.Fatalf("get all policies: %v", err)
	}
	if len(policies) != 1 || policies[0].Name != "Example" {
		t.Fatalf("unexpected policies: %+v", policies)
	}
}

func TestUpdatePolicy(t *testing.T) {
	t.Parallel()

	policy := validPolicy()
	policy.Name = "Updated"

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/policies/1/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodPatch {
				t.Fatalf("expected PATCH, got %s", r.Method)
			}
			var got PolicySettings
			if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if got.Name != "Updated" {
				t.Fatalf("unexpected policy name: %q", got.Name)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data":    map[string]any{"name": "Updated"},
			})
		},
	})
	defer cleanup()

	res, err := client.UpdatePolicy(context.Background(), policy, "1")
	if err != nil {
		t.Fatalf("update policy: %v", err)
	}
	if res.Name != "Updated" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestUpdatePolicyRequiresID(t *testing.T) {
	t.Parallel()

	policy := validPolicy()
	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.UpdatePolicy(context.Background(), policy, ""); !errors.Is(err, ErrPolicyID) {
		t.Fatalf("expected ErrPolicyID, got %v", err)
	}
}

func TestUpdatePolicyRequiresName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		policyName string
		wantErr    bool
	}{
		{"empty name", "", true},
		{"whitespace only", "   ", true},
		{"valid name", "My Policy", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
				"/api/v1/policies/123/": func(w http.ResponseWriter, r *http.Request) {
					respondJSON(t, w, http.StatusOK, map[string]any{
						"success": true,
						"data":    map[string]any{"name": tt.policyName},
					})
				},
			})
			defer cleanup()

			policy := validPolicy()
			policy.Name = tt.policyName

			_, err := client.UpdatePolicy(context.Background(), policy, "123")
			if tt.wantErr && !errors.Is(err, ErrPolicyName) {
				t.Errorf("expected ErrPolicyName for name %q, got %v", tt.policyName, err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error for name %q, got %v", tt.policyName, err)
			}
		})
	}
}

func TestUpdatePolicyValidatesIDWhitespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		policyID string
		wantErr  bool
	}{
		{"empty ID", "", true},
		{"whitespace only", "   ", true},
		{"tabs only", "\t\t", true},
		{"valid ID", "123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
				"/api/v1/policies/123/": func(w http.ResponseWriter, r *http.Request) {
					respondJSON(t, w, http.StatusOK, map[string]any{
						"success": true,
						"data":    map[string]any{"name": "Updated"},
					})
				},
			})
			defer cleanup()

			policy := validPolicy()
			_, err := client.UpdatePolicy(context.Background(), policy, tt.policyID)
			if tt.wantErr && !errors.Is(err, ErrPolicyID) {
				t.Errorf("expected ErrPolicyID for ID %q, got %v", tt.policyID, err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error for ID %q, got %v", tt.policyID, err)
			}
		})
	}
}

func TestDeletePolicy(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/policies/1/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodDelete {
				t.Fatalf("expected DELETE, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{"success": true})
		},
	})
	defer cleanup()

	ok, err := client.DeletePolicy(context.Background(), "1")
	if err != nil {
		t.Fatalf("delete policy: %v", err)
	}
	if !ok {
		t.Fatalf("expected delete success")
	}
}

func TestDeletePolicyRequiresID(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.DeletePolicy(context.Background(), ""); !errors.Is(err, ErrPolicyID) {
		t.Fatalf("expected ErrPolicyID, got %v", err)
	}
}

func validPolicy() PolicySettings {
	return PolicySettings{
		Name:        "Example Policy",
		Description: "short description",
		RuleType:    "ACTIVITY",
		Rules: Rules{
			And: []any{"protocol:HTTP"},
		},
	}
}
