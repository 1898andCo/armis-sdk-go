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

func TestCreateReport(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/reports/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"id":           42,
					"reportName":   "My Report",
					"reportType":   "DEVICE",
					"asq":          "in:devices timeFrame:\"1 Day\"",
					"creationTime": "2024-01-15T15:00:00Z",
					"isScheduled":  true,
					"schedule": map[string]any{
						"email":            []string{"test@armis.com"},
						"repeatAmount":     2,
						"repeatUnit":       "Days",
						"reportFileFormat": "csv",
						"timeOfDay":        "15:00",
						"timezone":         "Asia/Jerusalem",
						"weekdays":         []string{"Monday"},
					},
				},
				"success": true,
			})
		},
	})
	defer cleanup()

	req := CreateReportRequest{
		ASQ:          "in:devices timeFrame:\"1 Day\"",
		EmailSubject: "Subject",
		ExportConfiguration: ExportConfiguration{
			Columns: ExportColumns{
				Devices: []string{"Site", "Type", "Category"},
			},
		},
		ReportName: "My Report",
		Schedule: CreateSchedule{
			Email:            []string{"test@armis.com"},
			RepeatAmount:     "2",
			RepeatUnit:       "Days",
			ReportFileFormat: "csv",
			TimeOfDay:        "15:00",
			Timezone:         "Asia/Jerusalem",
			Weekdays:         []string{"Monday"},
		},
	}

	res, err := client.CreateReport(context.Background(), req)
	if err != nil {
		t.Fatalf("create report: %v", err)
	}
	if res.ID != 42 {
		t.Fatalf("expected report ID 42, got %d", res.ID)
	}
	if res.ReportName != "My Report" {
		t.Fatalf("unexpected report name: %s", res.ReportName)
	}
	if res.Asq != "in:devices timeFrame:\"1 Day\"" {
		t.Fatalf("unexpected ASQ: %s", res.Asq)
	}
	if !res.IsScheduled {
		t.Fatal("expected report to be scheduled")
	}
}

func TestCreateReport_EmptyName(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	req := CreateReportRequest{
		ASQ:        "in:devices",
		ReportName: "",
	}

	_, err := client.CreateReport(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty report name")
	}
	if !errors.Is(err, ErrReportName) {
		t.Fatalf("expected ErrReportName, got: %v", err)
	}
}

func TestCreateReport_EmptyASQ(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	req := CreateReportRequest{
		ASQ:        "",
		ReportName: "My Report",
	}

	_, err := client.CreateReport(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty ASQ")
	}
	if !errors.Is(err, ErrReportASQ) {
		t.Fatalf("expected ErrReportASQ, got: %v", err)
	}
}

func TestDeleteReport(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/reports/42/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodDelete {
				t.Fatalf("expected DELETE, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
			})
		},
	})
	defer cleanup()

	success, err := client.DeleteReport(context.Background(), "42")
	if err != nil {
		t.Fatalf("delete report: %v", err)
	}
	if !success {
		t.Fatal("expected success to be true")
	}
}

func TestDeleteReport_EmptyID(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	_, err := client.DeleteReport(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty report ID")
	}
	if !errors.Is(err, ErrReportID) {
		t.Fatalf("expected ErrReportID, got: %v", err)
	}
}

func TestDeleteReport_URLEncoding(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/reports/report%2Fwith%2Fslashes/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodDelete {
				t.Fatalf("expected DELETE, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
			})
		},
	})
	defer cleanup()

	success, err := client.DeleteReport(context.Background(), "report/with/slashes")
	if err != nil {
		t.Fatalf("delete report with special chars: %v", err)
	}
	if !success {
		t.Fatal("expected success to be true")
	}
}

func TestUpdateReport(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/reports/42/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodPatch {
				t.Fatalf("expected PATCH, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"id":           42,
					"reportName":   "Updated Report",
					"reportType":   "DEVICE",
					"asq":          "in:devices timeFrame:\"1 Day\"",
					"creationTime": "2024-01-15T15:00:00Z",
					"isScheduled":  true,
					"schedule": map[string]any{
						"email":            []string{"updated@armis.com"},
						"repeatAmount":     2,
						"repeatUnit":       "Days",
						"reportFileFormat": "csv",
						"timeOfDay":        "15:00",
						"timezone":         "Asia/Jerusalem",
						"weekdays":         []string{"Monday"},
					},
				},
				"success": true,
			})
		},
	})
	defer cleanup()

	req := UpdateReportRequest{
		ASQ:          "in:devices timeFrame:\"1 Day\"",
		EmailSubject: "Updated Subject",
		ExportConfiguration: ExportConfiguration{
			Columns: ExportColumns{
				Devices: []string{"Site", "Type", "Category"},
			},
		},
		ReportName: "Updated Report",
		Schedule: CreateSchedule{
			Email:            []string{"updated@armis.com"},
			RepeatAmount:     "2",
			RepeatUnit:       "Days",
			ReportFileFormat: "csv",
			TimeOfDay:        "15:00",
			Timezone:         "Asia/Jerusalem",
			Weekdays:         []string{"Monday"},
		},
	}

	res, err := client.UpdateReport(context.Background(), "42", req)
	if err != nil {
		t.Fatalf("update report: %v", err)
	}
	if res.ID != 42 {
		t.Fatalf("expected report ID 42, got %d", res.ID)
	}
	if res.ReportName != "Updated Report" {
		t.Fatalf("unexpected report name: %s", res.ReportName)
	}
	if res.Asq != "in:devices timeFrame:\"1 Day\"" {
		t.Fatalf("unexpected ASQ: %s", res.Asq)
	}
	if !res.IsScheduled {
		t.Fatal("expected report to be scheduled")
	}
}

func TestUpdateReport_PartialUpdate(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/reports/42/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodPatch {
				t.Fatalf("expected PATCH, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"id":           42,
					"reportName":   "Only Name Updated",
					"reportType":   "DEVICE",
					"asq":          "in:devices",
					"creationTime": "2024-01-15T15:00:00Z",
					"isScheduled":  false,
				},
				"success": true,
			})
		},
	})
	defer cleanup()

	// Only updating the report name
	req := UpdateReportRequest{
		ReportName: "Only Name Updated",
	}

	res, err := client.UpdateReport(context.Background(), "42", req)
	if err != nil {
		t.Fatalf("update report: %v", err)
	}
	if res.ReportName != "Only Name Updated" {
		t.Fatalf("unexpected report name: %s", res.ReportName)
	}
}

func TestUpdateReport_EmptyID(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	req := UpdateReportRequest{
		ReportName: "Updated Report",
	}

	_, err := client.UpdateReport(context.Background(), "", req)
	if err == nil {
		t.Fatal("expected error for empty report ID")
	}
	if !errors.Is(err, ErrReportID) {
		t.Fatalf("expected ErrReportID, got: %v", err)
	}
}

func TestUpdateReport_URLEncoding(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/reports/report%2Fwith%2Fslashes/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodPatch {
				t.Fatalf("expected PATCH, got %s", r.Method)
			}
			respondJSON(t, w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"id":         999,
					"reportName": "Encoded Report Updated",
					"reportType": "DEVICE",
				},
				"success": true,
			})
		},
	})
	defer cleanup()

	req := UpdateReportRequest{
		ReportName: "Encoded Report Updated",
	}

	res, err := client.UpdateReport(context.Background(), "report/with/slashes", req)
	if err != nil {
		t.Fatalf("update report with special chars: %v", err)
	}
	if res.ReportName != "Encoded Report Updated" {
		t.Fatalf("unexpected report name: %s", res.ReportName)
	}
}
