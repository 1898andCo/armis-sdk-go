// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

// GetReportsResponse represents the entire API response for retrieving all reports.
type GetReportsResponse struct {
	Data    ReportsList `json:"data"`
	Success bool        `json:"success,omitempty"`
}

// GetReportByIDResponse represents the API response for retrieving a single report.
type GetReportByIDResponse struct {
	Data    Report `json:"data"`
	Success bool   `json:"success,omitempty"`
}

// ReportsList represents a list of reports with pagination info.
type ReportsList struct {
	Reports []Report `json:"items"`
	Total   int      `json:"total"`
}

// Report represents a single report.
type Report struct {
	ID           int      `json:"id,omitempty"`
	ReportName   string   `json:"reportName,omitempty"`
	ReportType   string   `json:"reportType,omitempty"`
	Asq          string   `json:"asq,omitempty"`
	Schedule     Schedule `json:"schedule"`
	CreationTime string   `json:"creationTime,omitempty"`
	IsScheduled  bool     `json:"isScheduled,omitempty"`
}

// Schedule represents a report schedule.
type Schedule struct {
	Email []string `json:"email,omitempty"`
	// RepeatAmount is a float64 to support decimal interval values (e.g., 0.5 for half-day intervals).
	// The Armis API returns this as a numeric value that can include decimals.
	RepeatAmount     float64  `json:"repeatAmount,omitempty"`
	RepeatUnit       string   `json:"repeatUnit,omitempty"`
	ReportFileFormat string   `json:"reportFileFormat,omitempty"`
	TimeOfDay        string   `json:"timeOfDay,omitempty"`
	Timezone         string   `json:"timezone,omitempty"`
	Weekdays         []string `json:"weekdays,omitempty"`
}

// CreateReportRequest represents the request payload for creating a new report.
type CreateReportRequest struct {
	ASQ                 string              `json:"asq"`
	EmailSubject        string              `json:"emailSubject,omitempty"`
	ExportConfiguration ExportConfiguration `json:"exportConfiguration"`
	ReportName          string              `json:"reportName"`
	Schedule            CreateSchedule      `json:"schedule"`
}

// ExportConfiguration represents the export configuration for a report.
type ExportConfiguration struct {
	Columns ExportColumns `json:"columns"`
}

// ExportColumns represents the columns to export for different entity types.
type ExportColumns struct {
	Devices         []string `json:"devices,omitempty"`
	Vulnerabilities []string `json:"vulnerabilities,omitempty"`
	Activities      []string `json:"activities,omitempty"`
}

// CreateSchedule represents the schedule configuration when creating a report.
// RepeatAmount is a string in the request payload (unlike Schedule which uses float64 for responses).
type CreateSchedule struct {
	Email            []string `json:"email,omitempty"`
	RepeatAmount     string   `json:"repeatAmount,omitempty"`
	RepeatUnit       string   `json:"repeatUnit,omitempty"`
	ReportFileFormat string   `json:"reportFileFormat,omitempty"`
	TimeOfDay        string   `json:"timeOfDay,omitempty"`
	Timezone         string   `json:"timezone,omitempty"`
	Weekdays         []string `json:"weekdays,omitempty"`
}

// CreateReportResponse represents the API response for creating a report.
type CreateReportResponse struct {
	Data    Report `json:"data"`
	Success bool   `json:"success,omitempty"`
}

// DeleteReportResponse represents the API response for deleting a report.
type DeleteReportResponse struct {
	Success bool `json:"success"`
}

// UpdateReportRequest represents the request payload for updating a report.
// All fields are optional since PATCH allows partial updates.
// ExportConfiguration and Schedule use pointer types to enable true partial updates -
// nil values are omitted from JSON, allowing updates to only the specified fields.
type UpdateReportRequest struct {
	ASQ                 string               `json:"asq,omitempty"`
	EmailSubject        string               `json:"emailSubject,omitempty"`
	ExportConfiguration *ExportConfiguration `json:"exportConfiguration,omitempty"`
	ReportName          string               `json:"reportName,omitempty"`
	Schedule            *CreateSchedule      `json:"schedule,omitempty"`
}

// UpdateReportResponse represents the API response for updating a report.
type UpdateReportResponse struct {
	Data    Report `json:"data"`
	Success bool   `json:"success"`
}
