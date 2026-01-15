// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// GetBoundaryByID retrieves a single boundary from Armis using the boundary's ID.
func (c *Client) GetBoundaryByID(ctx context.Context, boundaryID string) (*BoundarySettings, error) {
	if boundaryID == "" {
		return nil, fmt.Errorf("%w", ErrBoundaryID)
	}

	// URL encode the boundary ID
	encodedBoundaryID := url.QueryEscape(boundaryID)

	// Create a new request
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/boundaries/%s/", c.apiVersion, encodedBoundaryID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetBoundary: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch boundary: %w", err)
	}

	// Parse the response
	var response GetBoundaryByID
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse boundary response: %w", err)
	}

	return &response.Data, nil
}

// GetBoundaries retrieves all boundaries from Armis.
func (c *Client) GetBoundaries(ctx context.Context) ([]BoundarySettings, error) {
	// Create a new request
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/boundaries/", c.apiVersion), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetBoundaries: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch boundaries: %w", err)
	}

	// Parse the response
	var response GetBoundaries
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse boundaries response: %w", err)
	}

	return response.Data.Boundaries, nil
}
