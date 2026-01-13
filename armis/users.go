// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// GetUsers fetches all users from Armis.
func (c *Client) GetUsers(ctx context.Context) ([]UserSettings, error) {
	// Create a new request
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/users/", c.apiVersion), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetUsers: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	// Parse the response
	var response UserAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse users response: %w", err)
	}

	return response.Data.Users, nil
}

// GetUser fetches a user from Armis using the user's ID or email.
func (c *Client) GetUser(ctx context.Context, userID string) (*UserSettings, error) {
	if userID == "" {
		return nil, fmt.Errorf("%w", ErrUserID)
	}

	// URL encode the user ID
	encodedUserID := url.QueryEscape(userID)

	// Create a new request
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/users/%s/", c.apiVersion, encodedUserID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Get User: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Parse the response
	var response GetUserAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return &response.Data, nil
}

// CreateUser creates a new user in Armis.
func (c *Client) CreateUser(ctx context.Context, user UserSettings) (*UserSettings, error) {
	if user.Name == "" {
		return nil, fmt.Errorf("%w", ErrUserName)
	}

	if user.Email == "" {
		return nil, fmt.Errorf("%w", ErrUserEmail)
	}

	userData, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/api/%s/users/", c.apiVersion), bytes.NewReader(userData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for CreateUser: %w", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user %q: %w", user.Name, err)
	}

	var apiResponse CreateUserAPIResponse
	if err := json.Unmarshal(res, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("%w: %+v", ErrHTTPResponse, apiResponse)
	}

	// Return the parsed user settings directly
	return &apiResponse.Data, nil
}

// UpdateUser updates a user in Armis.
func (c *Client) UpdateUser(ctx context.Context, user UserSettings, userID string) (*UserSettings, error) {
	if user.Name == "" {
		return nil, fmt.Errorf("%w", ErrUserName)
	}

	if user.Email == "" {
		return nil, fmt.Errorf("%w", ErrUserEmail)
	}

	if userID == "" {
		return nil, fmt.Errorf("%w", ErrUserID)
	}

	userData, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	// URL encode the user ID
	encodedUserID := url.QueryEscape(userID)

	req, err := c.newRequest(ctx, "PATCH", fmt.Sprintf("/api/%s/users/%s/", c.apiVersion, encodedUserID), bytes.NewReader(userData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for UpdateUser: %w", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user %q: %w", user.Name, err)
	}

	var apiResponse UpdateUserAPIResponse
	if err := json.Unmarshal(res, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("%w: %+v", ErrHTTPResponse, apiResponse)
	}

	// Return the parsed user settings directly
	return &apiResponse.Data, nil
}

// DeleteUser deletes a user from Armis.
func (c *Client) DeleteUser(ctx context.Context, userID string) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("%w", ErrUserID)
	}

	// URL encode the user ID
	encodedUserID := url.QueryEscape(userID)

	// Create a new request
	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("/api/%s/users/%s/", c.apiVersion, encodedUserID), nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request for DeleteUser: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return false, fmt.Errorf("failed to delete user: %w", err)
	}

	// Parse the response
	var response DeleteUserAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return false, fmt.Errorf("failed to parse user response: %w", err)
	}

	return response.Success, nil
}
