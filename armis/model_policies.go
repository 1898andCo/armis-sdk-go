// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

// CreatePolicyAPIResponse represents the response structure for creating a policy.
type CreatePolicyAPIResponse struct {
	Data    PolicyID `json:"data"`
	Success bool     `json:"success,omitempty"`
}

// GetPolicyAPIResponse and UpdatePolicyAPIResponse represent the response structures for retrieving and updating policies, respectively.
type GetPolicyAPIResponse struct {
	Data    GetPolicySettings `json:"data"`
	Success bool              `json:"success,omitempty"`
}

type GetAllPoliciesAPIResponse struct {
	Data    GetAllPolicySettings `json:"data"`
	Success bool                 `json:"success,omitempty"`
}

type UpdatePolicyAPIResponse struct {
	Data    UpdatePolicySettings `json:"data"`
	Success bool                 `json:"success,omitempty"`
}

// DeletePolicyAPIResponse represents the response structure for deleting a policy.
type DeletePolicyAPIResponse struct {
	Success bool `json:"success"`
}

type PolicyID struct {
	ID int `json:"id"`
}

// Consolidation represents the consolidation parameters.
type Consolidation struct {
	Amount int    `json:"amount,omitempty"`
	Unit   string `json:"unit,omitempty"`
}

// Params represents the parameters of an action.
type Params struct {
	Consolidation Consolidation `json:"consolidation"`
	Severity      string        `json:"severity,omitempty"`
	Title         string        `json:"title,omitempty"`
	Type          string        `json:"type,omitempty"`
	Endpoint      string        `json:"endpoint,omitempty"`
	Tags          []string      `json:"tags,omitempty"`
}

// Action represents an individual action.
type Action struct {
	Params Params `json:"params"`
	Type   string `json:"type,omitempty"`
}

// Rules represents the rules configuration for the policy.
type Rules struct {
	And []any `json:"and,omitempty"`
	Or  []any `json:"or,omitempty"`
}

// PolicySettings represents the main JSON structure.
type PolicySettings struct {
	Actions           []Action `json:"actions,omitempty"`
	Description       string   `json:"description,omitempty"`
	IsEnabled         bool     `json:"isEnabled,omitempty"`
	Labels            []string `json:"labels,omitempty"`
	MitreAttackLabels []string `json:"mitreAttackLabels,omitempty"`
	Name              string   `json:"name,omitempty"`
	RuleType          string   `json:"ruleType,omitempty"`
	Rules             Rules    `json:"rules"`
}

// UpdatePolicySettings represents the structure for updating policies.
type UpdatePolicySettings struct {
	Actions           []Action           `json:"actions,omitempty"`
	Description       string             `json:"description,omitempty"`
	IsEnabled         bool               `json:"isEnabled,omitempty"`
	Labels            []string           `json:"labels,omitempty"`
	MitreAttackLabels []MitreAttackLabel `json:"mitreAttackLabels,omitempty"`
	Name              string             `json:"name,omitempty"`
	RuleType          string             `json:"ruleType,omitempty"`
	Rules             Rules              `json:"rules"`
}

// GetPolicySettings represents the structure for retrieving policies.
type GetPolicySettings struct {
	Actions           []Action           `json:"actions,omitempty"`
	Description       string             `json:"description,omitempty"`
	IsEnabled         bool               `json:"isEnabled,omitempty"`
	Labels            []string           `json:"labels,omitempty"`
	MitreAttackLabels []MitreAttackLabel `json:"mitreAttackLabels,omitempty"`
	Name              string             `json:"name,omitempty"`
	RuleType          string             `json:"ruleType,omitempty"`
	Rules             Rules              `json:"rules"`
}

type MitreAttackLabel struct {
	Matrix       string `json:"matrix"`
	SubTechnique string `json:"subTechnique"`
	Tactic       string `json:"tactic"`
	Technique    string `json:"technique"`
}

// GetAllPolicySettings represents the structure for retrieving multiple policies.
type GetAllPolicySettings struct {
	Count    int            `json:"count"`
	Next     *int           `json:"next"`
	Prev     *int           `json:"prev"`
	Total    int            `json:"total"`
	Policies []SinglePolicy `json:"policies"`
}

// SinglePolicy represents an individual policy in the list of policies.
type SinglePolicy struct {
	Action            Action             `json:"action"`
	Actions           []Action           `json:"actions"`
	Description       string             `json:"description"`
	ID                string             `json:"id"`
	IsEnabled         bool               `json:"isEnabled"`
	Labels            []string           `json:"labels"`
	MitreAttackLabels []MitreAttackLabel `json:"mitreAttackLabels"`
	Name              string             `json:"name"`
	RiskFactorData    *any               `json:"riskFactorData"`
	RuleType          string             `json:"ruleType"`
	Rules             Rules              `json:"rules"`
}
