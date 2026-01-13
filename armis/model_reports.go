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
