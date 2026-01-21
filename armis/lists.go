// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetLists retrieves all lists from Armis.
func (c *Client) GetLists(ctx context.Context) ([]ListSettings, error) {
	// Create a new request
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/lists/", c.apiVersion), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetLists: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Lists: %w", err)
	}

	// Parse the response
	var response ListsAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse Lists response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return response.Data.Lists, nil
}
