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

// GetCollectorByID retrieves a single collector from Armis using the collector's ID.
func (c *Client) GetCollectorByID(ctx context.Context, collectorID string) (*CollectorSettings, error) {
	if collectorID == "" {
		return nil, fmt.Errorf("%w", ErrCollectorID)
	}

	// URL encode the collector ID
	encodedCollectorID := url.QueryEscape(collectorID)

	// Create a new request
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/collectors/%s/", c.apiVersion, encodedCollectorID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetCollector: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collector: %w", err)
	}

	// Parse the response
	var response GetCollectorAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse collector response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return &response.Data, nil
}

// GetCollectors retrieves all collectors from Armis.
func (c *Client) GetCollectors(ctx context.Context) ([]CollectorSettings, error) {
	// Create a new request
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/collectors/", c.apiVersion), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetCollectors: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collectors: %w", err)
	}

	// Parse the response
	var response CollectorAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse collectors response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return response.Data.Collectors, nil
}

// CreateCollector creates a new collector in Armis.
func (c *Client) CreateCollector(ctx context.Context, collector CreateCollectorSettings) (*NewCollectorSettings, error) {
	if collector.Name == "" {
		return nil, fmt.Errorf("%w", ErrCollectorName)
	}

	if collector.DeploymentType == "" {
		return nil, fmt.Errorf("%w", ErrCollectorType)
	}

	collectorData, err := json.Marshal(collector)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal collector data: %w", err)
	}

	// Create a new request
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/api/%s/collectors/", c.apiVersion), bytes.NewReader(collectorData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for CreateCollector: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create collector: %w", err)
	}

	// Parse the response
	var response CreateCollectorAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse create collector response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return &response.Data, nil
}

// UpdateCollector updates an existing collector in Armis.
func (c *Client) UpdateCollector(ctx context.Context, collectorID string, collector UpdateCollectorSettings) (*CollectorSettings, error) {
	if collectorID == "" {
		return nil, fmt.Errorf("%w", ErrCollectorID)
	}

	if collector.Name == "" {
		return nil, fmt.Errorf("%w", ErrCollectorName)
	}

	if collector.DeploymentType == "" {
		return nil, fmt.Errorf("%w", ErrCollectorType)
	}

	collectorData, err := json.Marshal(collector)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal collector data: %w", err)
	}

	// URL encode the collector ID
	encodedCollectorID := url.QueryEscape(collectorID)

	// Create a new request
	req, err := c.newRequest(ctx, "PATCH", fmt.Sprintf("/api/%s/collectors/%s/", c.apiVersion, encodedCollectorID), bytes.NewBuffer(collectorData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for UpdateCollector: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update collector: %w", err)
	}

	// Parse the response
	var response UpdateCollectorAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse update collector response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return &response.Data, nil
}

// DeleteCollector removes a collector from Armis.
func (c *Client) DeleteCollector(ctx context.Context, collectorID string) (bool, error) {
	if collectorID == "" {
		return false, fmt.Errorf("%w", ErrCollectorID)
	}

	// URL encode the collector ID
	encodedCollectorID := url.QueryEscape(collectorID)

	// Create a new request
	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("/api/%s/collectors/%s/", c.apiVersion, encodedCollectorID), nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request for DeleteCollector: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return false, fmt.Errorf("failed to delete collector: %w", err)
	}

	// Parse the response
	var response DeleteCollectorAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return false, fmt.Errorf("failed to parse collector response: %w", err)
	}

	if !response.Success {
		return false, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return response.Success, nil
}
