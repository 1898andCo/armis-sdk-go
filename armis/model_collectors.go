// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

// CollectorAPIResponse struct to match the entire API response.
type CollectorAPIResponse struct {
	Data    Collectors `json:"data"`
	Success bool       `json:"success,omitempty"`
	Count   *int       `json:"count,omitempty"`
	Next    *string    `json:"next,omitempty"`
	Prev    *string    `json:"prev,omitempty"`
	Total   *int       `json:"total,omitempty"`
}

// Collectors struct for the "data" key.
type Collectors struct {
	Collectors []CollectorSettings `json:"collectors,omitempty"`
}

// CollectorSettings struct for individual collector settings.
type CollectorSettings struct {
	BootTime           string `json:"bootTime,omitempty"`
	ClusterID          int    `json:"clusterId,omitempty"`
	CollectorNumber    int    `json:"collectorNumber,omitempty"`
	DefaultGateway     string `json:"defaultGateway,omitempty"`
	HTTPSProxyRedacted string `json:"httpsProxyRedacted,omitempty"`
	IPAddress          string `json:"ipAddress,omitempty"`
	LastSeen           string `json:"lastSeen,omitempty"`
	MacAddress         string `json:"macAddress,omitempty"`
	Name               string `json:"name,omitempty"`
	Namespace          string `json:"namespace,omitempty"`
	ProductSerial      string `json:"productSerial,omitempty"`
	Status             string `json:"status,omitempty"`
	Subnet             string `json:"subnet,omitempty"`
	SystemVendor       string `json:"systemVendor,omitempty"`
	Type               string `json:"type,omitempty"`
}

type CreateCollectorSettings struct {
	DeploymentType string `json:"deploymentType"`
	Name           string `json:"name"`
}

type UpdateCollectorSettings struct {
	DeploymentType string `json:"deploymentType"`
	Name           string `json:"name"`
}

// CreateCollectorAPIResponse struct for creating a new collector.
type CreateCollectorAPIResponse struct {
	Data    NewCollectorSettings `json:"data"`
	Success bool                 `json:"success,omitempty"`
}

// GetCollectorAPIResponse struct for getting collector details.
type GetCollectorAPIResponse struct {
	Data    CollectorSettings `json:"data"`
	Success bool              `json:"success,omitempty"`
}

// UpdateCollectorAPIResponse struct for updating collector.
type UpdateCollectorAPIResponse struct {
	Data    CollectorSettings `json:"data"`
	Success bool              `json:"success"`
}

// DeleteCollectorAPIResponse struct for deleting collector.
type DeleteCollectorAPIResponse struct {
	Success bool `json:"success"`
}

// NewCollectorSettings struct for creating a new collector response.
type NewCollectorSettings struct {
	CollectorID int    `json:"collectorId,omitempty"`
	LicenseKey  string `json:"licenseKey,omitempty"`
	Password    string `json:"password,omitempty"`
	User        string `json:"user,omitempty"`
}
