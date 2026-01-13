// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// GetReportByID returns the specified report based on the ID presented.
func (c *Client) GetReportByID(ctx context.Context, reportID string) (*Report, error) {
	if reportID == "" {
		return nil, fmt.Errorf("%w", ErrReportID)
	}

	// URL encode the report ID
	encodedReportID := url.QueryEscape(reportID)

	// Create a new request
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/reports/%s/", c.apiVersion, encodedReportID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetReport: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch report: %w", err)
	}

	// Parse the response
	var response GetReportByIDResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse report response: %w", err)
	}

	return &response.Data, nil
}

// GetReports returns all reports from the Armis API.
func (c *Client) GetReports(ctx context.Context) ([]Report, error) {
	// Create a new request
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/reports/", c.apiVersion), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetReports: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reports: %w", err)
	}

	// Parse the response
	var response GetReportsResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse reports response: %w", err)
	}

	return response.Data.Reports, nil
}
