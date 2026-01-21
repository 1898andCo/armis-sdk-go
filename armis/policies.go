// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const descriptionLimit = 500

var allowedRuleTypes = map[string]struct{}{
	"ACTIVITY":      {},
	"IP_CONNECTION": {},
	"DEVICE":        {},
	"VULNERABILITY": {},
}

// CreatePolicy creates a new policy in Armis.
func (c *Client) CreatePolicy(ctx context.Context, policy PolicySettings) (PolicyID, error) {
	err := policy.Validate()
	if err != nil {
		return PolicyID{}, fmt.Errorf("failed to validate policy: %w", err)
	}

	policyData, err := json.Marshal(policy)
	if err != nil {
		return PolicyID{}, fmt.Errorf("failed to marshal policy data: %w", err)
	}

	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/api/%s/policies/", c.apiVersion), bytes.NewReader(policyData))
	if err != nil {
		return PolicyID{}, fmt.Errorf("failed to create request for CreatePolicy: %w", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		return PolicyID{}, fmt.Errorf("failed to create policy %q: %w", policy.Name, err)
	}

	var apiResponse CreatePolicyAPIResponse
	if err := json.Unmarshal(res, &apiResponse); err != nil {
		return PolicyID{}, fmt.Errorf("failed to parse policy response: %w", err)
	}

	if !apiResponse.Success {
		return PolicyID{}, fmt.Errorf("%w:%v", ErrHTTPResponse, apiResponse)
	}

	return apiResponse.Data, nil
}

// Validate performs policy validation to ensure that a name and rules exist,
// the description field is less than 500 characters,
// and the rule type must be ACTIVITY, IP CONNECTION, DEVICE, or VULNERABILITY.
func (p PolicySettings) Validate() error {
	var errs []error
	if p.Name == "" {
		errs = append(errs, ErrPolicyName)
	}

	if len(p.Rules.And) == 0 && len(p.Rules.Or) == 0 {
		errs = append(errs, ErrPolicyRules)
	}

	if len(p.Description) > descriptionLimit {
		errs = append(errs, ErrPolicyDescription)
	}

	if _, ok := allowedRuleTypes[p.RuleType]; !ok {
		errs = append(errs, ErrPolicyRuleType)
	}

	return errors.Join(errs...)
}

// GetPolicy fetches a policy from Armis using the policy's ID.
func (c *Client) GetPolicy(ctx context.Context, policyID string) (GetPolicySettings, error) {
	if policyID == "" {
		return GetPolicySettings{}, fmt.Errorf("%w", ErrPolicyID)
	}

	// URL encode the policy ID
	encodedPolicyID := url.QueryEscape(policyID)

	// Create a new request
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/policies/%s/", c.apiVersion, encodedPolicyID), nil)
	if err != nil {
		return GetPolicySettings{}, fmt.Errorf("failed to create request for Get Policy: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return GetPolicySettings{}, fmt.Errorf("failed to fetch policy: %w", err)
	}

	// Parse the response
	var response GetPolicyAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return GetPolicySettings{}, fmt.Errorf("failed to parse policy response: %w", err)
	}

	if !response.Success {
		return GetPolicySettings{}, fmt.Errorf("%w:%v", ErrHTTPResponse, response)
	}

	return response.Data, nil
}

// GetAllPolicies retrieves all policies from Armis with automatic pagination.
func (c *Client) GetAllPolicies(ctx context.Context) ([]SinglePolicy, error) {
	var allPolicies []SinglePolicy
	from := 0
	length := 100

	for {
		req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/%s/policies/?from=%d&length=%d", c.apiVersion, from, length), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request for Get All Policies: %w", err)
		}

		res, err := c.doRequest(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch policies page (from=%d): %w", from, err)
		}

		var response GetAllPoliciesAPIResponse
		if err := json.Unmarshal(res, &response); err != nil {
			return nil, fmt.Errorf("failed to parse policies response: %w", err)
		}

		if !response.Success {
			return nil, fmt.Errorf("%w:%v", ErrHTTPResponse, response)
		}

		if len(allPolicies) == 0 {
			allPolicies = make([]SinglePolicy, 0, response.Data.Total)
		}

		allPolicies = append(allPolicies, response.Data.Policies...)

		if response.Data.Next == nil {
			break
		}

		from = *response.Data.Next
	}

	return allPolicies, nil
}

// UpdatePolicy updates a policy in Armis.
func (c *Client) UpdatePolicy(ctx context.Context, policy PolicySettings, policyID string) (UpdatePolicySettings, error) {
	if err := validateUpdateInput(policy, policyID); err != nil {
		return UpdatePolicySettings{}, err
	}

	policyData, err := json.Marshal(policy)
	if err != nil {
		return UpdatePolicySettings{}, fmt.Errorf("failed to marshal policy data: %w", err)
	}

	// URL encode the policy ID
	encodedPolicyID := url.QueryEscape(policyID)

	req, err := c.newRequest(ctx, "PATCH", fmt.Sprintf("/api/%s/policies/%s/", c.apiVersion, encodedPolicyID), bytes.NewReader(policyData))
	if err != nil {
		return UpdatePolicySettings{}, fmt.Errorf("failed to create request for UpdatePolicy: %w", err)
	}

	res, err := c.doRequest(req)
	if err != nil {
		return UpdatePolicySettings{}, fmt.Errorf("failed to update policy %q: %w", policy.Name, err)
	}

	var apiResponse UpdatePolicyAPIResponse
	if err := json.Unmarshal(res, &apiResponse); err != nil {
		return UpdatePolicySettings{}, fmt.Errorf("failed to parse policy response: %w", err)
	}

	if !apiResponse.Success {
		return UpdatePolicySettings{}, fmt.Errorf("%w:%v", ErrHTTPResponse, apiResponse)
	}

	// Return the parsed policy settings directly
	return apiResponse.Data, nil
}

func validateUpdateInput(policy PolicySettings, id string) error {
	var errs []error

	if strings.TrimSpace(policy.Name) == "" {
		errs = append(errs, ErrPolicyName)
	}
	if strings.TrimSpace(id) == "" {
		errs = append(errs, ErrPolicyID)
	}
	return errors.Join(errs...)
}

// DeletePolicy deletes a policy from Armis.
func (c *Client) DeletePolicy(ctx context.Context, policyID string) (bool, error) {
	if policyID == "" {
		return false, fmt.Errorf("%w", ErrPolicyID)
	}

	// URL encode the policy ID
	encodedPolicyID := url.QueryEscape(policyID)

	// Create a new request
	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("/api/%s/policies/%s/", c.apiVersion, encodedPolicyID), nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request for DeletePolicy: %w", err)
	}

	// Perform the request
	res, err := c.doRequest(req)
	if err != nil {
		return false, fmt.Errorf("failed to delete policy: %w", err)
	}

	// Parse the response
	var response DeletePolicyAPIResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return false, fmt.Errorf("failed to parse policy response: %w", err)
	}

	if !response.Success {
		return false, fmt.Errorf("%w: %+v", ErrHTTPResponse, response)
	}

	return response.Success, nil
}

// Len returns the total number of policy rules (And + Or).
func (r Rules) Len() int {
	return len(r.And) + len(r.Or)
}
