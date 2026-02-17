// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// GetReportByID returns the specified report based on the ID presented.
func (c *Client) GetReportByID(ctx context.Context, reportID string) (*Report, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if strings.TrimSpace(reportID) == "" {
		return nil, ErrReportID
	}

	// URL encode the report ID
	encodedReportID := url.QueryEscape(reportID)

	// Create a new request
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/reports/%s/", c.apiVersion, encodedReportID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetReport: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch report: %w", err)
	}

	// Parse the response
	var response GetReportByIDResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse report response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return &response.Data, nil
}

// GetReports returns all reports from the Armis API.
func (c *Client) GetReports(ctx context.Context) ([]Report, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Create a new request
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/reports/", c.apiVersion), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetReports: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reports: %w", err)
	}

	// Parse the response
	var response GetReportsResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse reports response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return response.Data.Reports, nil
}

// CreateReport creates a new report in Armis.
func (c *Client) CreateReport(ctx context.Context, report CreateReportRequest) (*Report, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if strings.TrimSpace(report.ReportName) == "" {
		return nil, ErrReportName
	}

	if strings.TrimSpace(report.ASQ) == "" {
		return nil, ErrReportASQ
	}

	reportData, err := json.Marshal(report)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal report data: %w", err)
	}

	// Create a new request
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/api/%s/reports/", c.apiVersion), bytes.NewReader(reportData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for CreateReport: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	// Parse the response
	var response CreateReportResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse create report response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return &response.Data, nil
}

// UpdateReport updates an existing report by its ID.
func (c *Client) UpdateReport(ctx context.Context, reportID string, report UpdateReportRequest) (*Report, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if strings.TrimSpace(reportID) == "" {
		return nil, ErrReportID
	}

	reportData, err := json.Marshal(report)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal report data: %w", err)
	}

	// URL encode the report ID
	encodedReportID := url.QueryEscape(reportID)

	// Create a new request
	req, err := c.newRequest(ctx, "PATCH", fmt.Sprintf("/api/%s/reports/%s/", c.apiVersion, encodedReportID), bytes.NewReader(reportData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for UpdateReport: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update report: %w", err)
	}

	// Parse the response
	var response UpdateReportResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse update report response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return &response.Data, nil
}

// DeleteReport deletes a report by its ID.
func (c *Client) DeleteReport(ctx context.Context, reportID string) (bool, error) {
	if ctx == nil {
		return false, ErrNilContext
	}

	if err := ctx.Err(); err != nil {
		return false, err
	}

	if strings.TrimSpace(reportID) == "" {
		return false, ErrReportID
	}

	// URL encode the report ID
	encodedReportID := url.QueryEscape(reportID)

	// Create a new request
	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("/api/%s/reports/%s/", c.apiVersion, encodedReportID), nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request for DeleteReport: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(ctx, req)
	if err != nil {
		return false, fmt.Errorf("failed to delete report: %w", err)
	}

	// Parse the response
	var response DeleteReportResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return false, fmt.Errorf("failed to parse delete report response: %w", err)
	}

	if !response.Success {
		return false, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return response.Success, nil
}
