// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetSites retrieves all sites from Armis.
func (c *Client) GetSites(ctx context.Context) ([]SiteSettings, error) {
	// Create a new request
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/sites/", c.apiVersion), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetSites: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sites: %w", err)
	}

	// Parse the response
	var response SiteApiResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse sites response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return response.Data.Sites, nil
}
