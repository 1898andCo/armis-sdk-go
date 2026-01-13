// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import "errors"

// Sentinel errors shared across the Armis package.

var (
	ErrGetKey       = errors.New("failed to get API key")
	ErrGetURL       = errors.New("failed to get API URL")
	ErrNoAPIKey     = errors.New("missing API key")
	ErrAuthFailed   = errors.New("failed to authenticate")
	ErrHTTPResponse = errors.New("HTTP response error")

	// Collector errors.
	ErrCollectorID   = errors.New("collector ID cannot be empty")
	ErrCollectorName = errors.New("collector name cannot be empty")
	ErrCollectorType = errors.New("collector type cannot be empty")

	// Boundary errors.
	ErrBoundaryID = errors.New("boundary ID cannot be empty")

	// Policy errors.
	ErrPolicyID          = errors.New("policy ID cannot be empty")
	ErrPolicyName        = errors.New("policy name cannot be empty")
	ErrPolicyDescription = errors.New("policy description must be less than 500 characters")
	ErrPolicyRules       = errors.New("policy rules cannot be empty")
	ErrPolicyRuleType    = errors.New("policy rule type must be ACTIVITY, IP_CONNECTION, DEVICE, or VULNERABILITY")

	// User errors.
	ErrUserID    = errors.New("user ID cannot be empty")
	ErrUserName  = errors.New("user name cannot be empty")
	ErrUserEmail = errors.New("user email cannot be empty")

	// Role errors.
	ErrRoleID       = errors.New("role ID cannot be empty")
	ErrRoleName     = errors.New("role name cannot be empty")
	ErrRoleNotFound = errors.New("role not found")

	// Search errors.
	ErrSearchAQL = errors.New("search AQL cannot be empty")

	// Report errors.
	ErrReportID = errors.New("report ID cannot be empty")
)
