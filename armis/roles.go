// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// GetRoles fetches all roles from the API.
func (c *Client) GetRoles(ctx context.Context) ([]RoleSettings, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/roles/", c.apiVersion), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetRoles: %w", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}

	var response RolesAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse roles response: %w", err)
	}

	return response.Data, nil
}

// GetRoleByName fetches a specific role by name.
func (c *Client) GetRoleByName(ctx context.Context, name string) (*RoleSettings, error) {
	if name == "" {
		return nil, fmt.Errorf("%w", ErrRoleName)
	}

	roles, err := c.GetRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles to find role %q: %w", name, err)
	}

	for _, r := range roles {
		if r.Name == name {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("%w: %q", ErrRoleNotFound, name)
}

// GetRoleByID fetches a specific role by ID.
func (c *Client) GetRoleByID(ctx context.Context, id string) (*RoleSettings, error) {
	if id == "" {
		return nil, fmt.Errorf("%w", ErrRoleID)
	}

	roles, err := c.GetRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles to find role ID %q: %w", id, err)
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid role ID format: %w", err)
	}

	for _, r := range roles {
		if r.ID == idInt {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("%w: %q", ErrRoleNotFound, id)
}

// CreateRole creates a new role in the API.
func (c *Client) CreateRole(ctx context.Context, role RoleSettings) (bool, error) {
	roleData, err := json.Marshal(role)
	if err != nil {
		return false, fmt.Errorf("failed to marshal role data: %w", err)
	}

	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/api/%s/roles/", c.apiVersion), bytes.NewReader(roleData))
	if err != nil {
		return false, fmt.Errorf("failed to create request for CreateRole: %w", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		return false, fmt.Errorf("failed to create role: %w", err)
	}

	var response CreateRoleAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return false, fmt.Errorf("failed to parse role creation response: %w", err)
	}

	return response.Success, nil
}

// UpdateRole updates an existing role in the API.
func (c *Client) UpdateRole(ctx context.Context, role RoleSettings, id string) (*RoleSettings, error) {
	if id == "" {
		return nil, fmt.Errorf("%w", ErrRoleID)
	}

	roleData, err := json.Marshal(role)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal role data: %w", err)
	}

	encodedRoleID := url.QueryEscape(id)

	req, err := c.newRequest(ctx, "PATCH", fmt.Sprintf("/api/%s/roles/%s/", c.apiVersion, encodedRoleID), bytes.NewReader(roleData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for UpdateRole: %w", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	var response RoleSettings
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse role update response: %w", err)
	}

	return &response, nil
}

// DeleteRole deletes a role by ID.
func (c *Client) DeleteRole(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("%w", ErrRoleID)
	}

	encodedRoleID := url.QueryEscape(id)

	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("/api/%s/roles/%s/", c.apiVersion, encodedRoleID), nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request for DeleteRole: %w", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		return false, fmt.Errorf("failed to delete role: %w", err)
	}

	var response DeleteRoleAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return false, fmt.Errorf("failed to parse role deletion response: %w", err)
	}

	return response.Success, nil
}
