// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"bytes"
	"encoding/json"
)

// SearchAPIResponse represents the response returned by the search endpoint.
type SearchAPIResponse struct {
	Data    SearchData `json:"data"`
	Success bool       `json:"success"`
}

// SearchData contains pagination metadata and the list of search results.
type SearchData struct {
	Count   int            `json:"count"`
	Next    *int           `json:"next"`
	Prev    *int           `json:"prev"`
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
}

// SearchResult represents a single search hit returned by Armis.
type SearchResult struct {
	ActivityUUIDs        []string               `json:"activityUUIDs,omitempty"`
	AffectedDevicesCount int                    `json:"affectedDevicesCount,omitempty"`
	AlertID              int                    `json:"alertId,omitempty"`
	Category             string                 `json:"category,omitempty"`
	Classification       string                 `json:"classification,omitempty"`
	ConnectionIDs        []string               `json:"connectionIds,omitempty"`
	Description          string                 `json:"description,omitempty"`
	DestinationEndpoints []SearchEndpoint       `json:"destinationEndpoints,omitempty"`
	Devices              int                    `json:"devices,omitempty"`
	DeviceIDs            []int                  `json:"deviceIds,omitempty"`
	Evidence             *RiskFactorEvidence    `json:"evidence,omitempty"`
	FirstSeen            string                 `json:"firstSeen,omitempty"`
	Group                string                 `json:"group,omitempty"`
	LastSeen             string                 `json:"lastSeen,omitempty"`
	LastAlertUpdateTime  string                 `json:"lastAlertUpdateTime,omitempty"`
	MitreAttackLabels    []string               `json:"mitreAttackLabels,omitempty"`
	Policy               any                    `json:"policy,omitempty"`
	PolicyID             string                 `json:"policyId,omitempty"`
	PolicyLabels         []string               `json:"policyLabels,omitempty"`
	PolicyTitle          string                 `json:"policyTitle,omitempty"`
	Remediation          *RiskFactorRemediation `json:"remediation,omitempty"`
	Score                string                 `json:"score,omitempty"`
	Severity             string                 `json:"severity,omitempty"`
	Source               string                 `json:"source,omitempty"`
	SourceEndpoints      []SearchEndpoint       `json:"sourceEndpoints,omitempty"`
	Status               string                 `json:"status,omitempty"`
	StatusChangeTime     string                 `json:"statusChangeTime,omitempty"`
	Time                 string                 `json:"time,omitempty"`
	Title                string                 `json:"title,omitempty"`
	Type                 string                 `json:"type,omitempty"`
	Action               string                 `json:"action,omitempty"`
	AdditionalInfo       *AuditAdditionalInfo   `json:"additionalInfo,omitempty"`
	ID                   int                    `json:"id,omitempty"`
	TimeUtc              string                 `json:"timeUtc,omitempty"`
	Trigger              string                 `json:"trigger,omitempty"`
	User                 string                 `json:"user,omitempty"`
	UserIP               string                 `json:"userIp,omitempty"`
	Band                 string                 `json:"band,omitempty"`
	BSSID                string                 `json:"bssid,omitempty"`
	Channel              string                 `json:"channel,omitempty"`
	Duration             int                    `json:"duration,omitempty"`
	EndTimestamp         string                 `json:"endTimestamp,omitempty"`
	InboundTraffic       int                    `json:"inboundTraffic,omitempty"`
	OutboundTraffic      int                    `json:"outboundTraffic,omitempty"`
	Protocol             string                 `json:"protocol,omitempty"`
	Risk                 string                 `json:"risk,omitempty"`
	RSSI                 string                 `json:"rssi,omitempty"`
	SNR                  string                 `json:"snr,omitempty"`
	SourceID             int                    `json:"sourceId,omitempty"`
	SSID                 string                 `json:"ssid,omitempty"`
	StartTimestamp       string                 `json:"startTimestamp,omitempty"`
	TargetID             int                    `json:"targetId,omitempty"`
	Traffic              int                    `json:"traffic,omitempty"`
	Sensor               Sensor                 `json:"sensor"`
	Site                 SingleSite             `json:"site"`
	Sites                []SingleSite           `json:"sites"`
}

type SearchEndpointID string

func (id *SearchEndpointID) UnmarshalJSON(b []byte) error {
	// Trim leading whitespace to detect quoted vs numeric IDs.
	i := 0
	for ; i < len(b) && (b[i] == ' ' || b[i] == '\n' || b[i] == '\r' || b[i] == '\t'); i++ {
	}

	if i < len(b) && b[i] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*id = SearchEndpointID(s)
		return nil
	}

	// Otherwise, treat it as a number.
	var n json.Number
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	*id = SearchEndpointID(n.String())
	return nil
}

type SearchEndpointIPs []string

func (ips *SearchEndpointIPs) UnmarshalJSON(b []byte) error {
	trimmed := bytes.TrimSpace(b)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		*ips = nil
		return nil
	}

	if trimmed[0] == '[' {
		var list []string
		if err := json.Unmarshal(trimmed, &list); err != nil {
			return err
		}
		*ips = list
		return nil
	}

	var single string
	if err := json.Unmarshal(trimmed, &single); err != nil {
		return err
	}
	*ips = []string{single}
	return nil
}

// SearchEndpoint represents an endpoint referenced in a search result.
type SearchEndpoint struct {
	ID             SearchEndpointID  `json:"id,omitempty"`
	IP             SearchEndpointIPs `json:"ip,omitempty"`
	Name           string            `json:"name,omitempty"`
	Risk           int               `json:"risk,omitempty"`
	Type           string            `json:"type,omitempty"`
	MacAddress     []string          `json:"macAddress,omitempty"`
	DataSources    []string          `json:"dataSources,omitempty"`
	BusinessImpact string            `json:"businessImpact,omitempty"`
	RiskLevel      int               `json:"riskLevel,omitempty"`
}

// AuditAdditionalInfo represents additional information in audit log entries.
type AuditAdditionalInfo struct {
	Data string `json:"data,omitempty"`
	Type string `json:"type,omitempty"`
}

// RiskFactorEvidence represents evidence details for a risk factor.
type RiskFactorEvidence struct {
	AQL          string `json:"AQL,omitempty"`
	WhatHappened string `json:"whatHappened,omitempty"`
}

// RecommendedAction represents a single recommended action in remediation.
type RecommendedAction struct {
	Description string `json:"description,omitempty"`
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Type        string `json:"type,omitempty"`
}

// RiskFactorRemediation represents remediation details for a risk factor.
type RiskFactorRemediation struct {
	Category           string              `json:"category,omitempty"`
	Description        string              `json:"description,omitempty"`
	RecommendedActions []RecommendedAction `json:"recommendedActions,omitempty"`
	Type               string              `json:"type,omitempty"`
}

// Sensor represents the name and type of a sensor in a connections search.
type Sensor struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// SingleSite represents a single site in a connections search.
type SingleSite struct {
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`
}
