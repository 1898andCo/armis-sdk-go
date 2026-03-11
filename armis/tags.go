// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// GetTags retrieves all tags from Armis with automatic pagination.
func (c *Client) GetTags(ctx context.Context) ([]string, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var allTags []string

	for {
		path := fmt.Sprintf("/api/%s/tags/", c.apiVersion)
		if len(allTags) > 0 {
			params := url.Values{}
			params.Set("from", strconv.Itoa(len(allTags)))
			path = fmt.Sprintf("/api/%s/tags/?%s", c.apiVersion, params.Encode())
		}

		req, err := c.newRequest(ctx, "GET", path, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request for GetTags: %w", err)
		}

		res, err := c.doRequest(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch tags: %w", err)
		}

		var response TagsAPIResponse
		if err := json.Unmarshal(res, &response); err != nil {
			return nil, fmt.Errorf("failed to parse tags response: %w", err)
		}

		if !response.Success {
			return nil, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
		}

		if len(response.Data.Tags) == 0 {
			break
		}

		allTags = append(allTags, response.Data.Tags...)

		if response.Data.Next == nil {
			break
		}
	}

	return allTags, nil
}
