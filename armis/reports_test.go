// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestGetReports(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/reports/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"items": []map[string]any{
						{
							"id":           1,
							"reportName":   "Test Report",
							"reportType":   "DEVICE",
							"asq":          "in:devices",
							"creationTime": "2024-01-15T10:00:00Z",
							"isScheduled":  true,
							"schedule": map[string]any{
								"email":            []string{"test@example.com"},
								"repeatAmount":     1,
								"repeatUnit":       "WEEK",
								"reportFileFormat": "PDF",
								"timeOfDay":        "09:00",
								"timezone":         "UTC",
								"weekdays":         []string{"MONDAY"},
							},
						},
						{
							"id":           2,
							"reportName":   "Second Report",
							"reportType":   "VULNERABILITY",
							"asq":          "in:vulnerabilities",
							"creationTime": "2024-01-16T11:00:00Z",
							"isScheduled":  false,
						},
					},
					"total": 2,
				},
				"success": true,
			})
		},
	})
	defer cleanup()

	res, err := client.GetReports(context.Background())
	if err != nil {
		t.Fatalf("get reports: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(res))
	}
	if res[0].ReportName != "Test Report" {
		t.Fatalf("unexpected first report name: %s", res[0].ReportName)
	}
	if res[1].ReportName != "Second Report" {
		t.Fatalf("unexpected second report name: %s", res[1].ReportName)
	}
}

func TestGetReportByID(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/reports/123/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"id":           123,
					"reportName":   "Specific Report",
					"reportType":   "DEVICE",
					"asq":          "in:devices",
					"creationTime": "2024-01-15T10:00:00Z",
					"isScheduled":  true,
					"schedule": map[string]any{
						"email":            []string{"admin@example.com"},
						"repeatAmount":     2,
						"repeatUnit":       "DAY",
						"reportFileFormat": "CSV",
						"timeOfDay":        "08:00",
						"timezone":         "America/New_York",
						"weekdays":         []string{},
					},
				},
				"success": true,
			})
		},
	})
	defer cleanup()

	res, err := client.GetReportByID(context.Background(), "123")
	if err != nil {
		t.Fatalf("get report by id: %v", err)
	}
	if res.ReportName != "Specific Report" {
		t.Fatalf("unexpected report name: %s", res.ReportName)
	}
	if res.ID != 123 {
		t.Fatalf("unexpected report ID: %d", res.ID)
	}
	if res.ReportType != "DEVICE" {
		t.Fatalf("unexpected report type: %s", res.ReportType)
	}
	if !res.IsScheduled {
		t.Fatal("expected report to be scheduled")
	}
	if len(res.Schedule.Email) != 1 || res.Schedule.Email[0] != "admin@example.com" {
		t.Fatalf("unexpected schedule email: %v", res.Schedule.Email)
	}
}

func TestGetReportByID_EmptyID(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	_, err := client.GetReportByID(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty report ID")
	}
	if !errors.Is(err, ErrReportID) {
		t.Fatalf("expected ErrReportID, got: %v", err)
	}
}

func TestGetReportByID_URLEncoding(t *testing.T) {
	t.Parallel()

	// Test that report IDs with special characters are properly URL-encoded
	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/reports/report%2Fwith%2Fslashes/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"id":         999,
					"reportName": "Encoded Report",
					"reportType": "DEVICE",
				},
				"success": true,
			})
		},
	})
	defer cleanup()

	res, err := client.GetReportByID(context.Background(), "report/with/slashes")
	if err != nil {
		t.Fatalf("get report with special chars: %v", err)
	}
	if res.ReportName != "Encoded Report" {
		t.Fatalf("unexpected report name: %s", res.ReportName)
	}
}

func TestGetReports_EmptyList(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/reports/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"items": []map[string]any{},
					"total": 0,
				},
				"success": true,
			})
		},
	})
	defer cleanup()

	res, err := client.GetReports(context.Background())
	if err != nil {
		t.Fatalf("get reports: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("expected 0 reports, got %d", len(res))
	}
}
