// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// GetSearch executes an AQL search request. The `aql` parameter must be a valid
// Armis Query Language string. The caller can control inclusion of sample
// results and totals via the boolean flags.
func (c *Client) GetSearch(ctx context.Context, aql string, includeSample, includeTotal bool) (SearchData, error) {
	if strings.TrimSpace(aql) == "" {
		return SearchData{}, fmt.Errorf("%w", ErrSearchAQL)
	}

	params := url.Values{}
	params.Set("aql", aql)
	params.Set("includeSample", strconv.FormatBool(includeSample))
	params.Set("includeTotal", strconv.FormatBool(includeTotal))

	encodedParams := params.Encode()

	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/search/?%s", c.apiVersion, encodedParams), nil)
	if err != nil {
		return SearchData{}, fmt.Errorf("failed to create request for GetSearch: %w", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		return SearchData{}, fmt.Errorf("failed to execute search: %w", err)
	}

	var response SearchAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return SearchData{}, fmt.Errorf("failed to parse search response: %w", err)
	}

	if !response.Success {
		return SearchData{}, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return response.Data, nil
}
