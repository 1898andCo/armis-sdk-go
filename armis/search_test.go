// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestGetAlertSearch(t *testing.T) {
	t.Parallel()

	const (
		aql           = "in:alerts status:Open"
		includeSample = true
		includeTotal  = true
	)

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/search/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}

			values := r.URL.Query()
			if got := values.Get("aql"); got != aql {
				t.Fatalf("unexpected aql: %q", got)
			}
			if got := values.Get("includeSample"); got != "true" {
				t.Fatalf("unexpected includeSample: %q", got)
			}
			if got := values.Get("includeTotal"); got != "true" {
				t.Fatalf("unexpected includeTotal: %q", got)
			}

			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"count": 1,
					"next":  nil,
					"prev":  nil,
					"total": 1,
					"results": []map[string]any{{
						"title":  "Example Alert",
						"status": "Open",
					}},
				},
			})
		},
	})
	defer cleanup()

	res, err := client.GetSearch(context.Background(), aql, includeSample, includeTotal)
	if err != nil {
		t.Fatalf("get search: %v", err)
	}
	if res.Total != 1 || len(res.Results) != 1 || res.Results[0].Title != "Example Alert" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestGetActivitySearch(t *testing.T) {
	t.Parallel()

	const (
		aql           = "in:activity status:Open"
		includeSample = true
		includeTotal  = true
	)

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/search/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}

			values := r.URL.Query()
			if got := values.Get("aql"); got != aql {
				t.Fatalf("unexpected aql: %q", got)
			}
			if got := values.Get("includeSample"); got != "true" {
				t.Fatalf("unexpected includeSample: %q", got)
			}
			if got := values.Get("includeTotal"); got != "true" {
				t.Fatalf("unexpected includeTotal: %q", got)
			}

			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"count": 1,
					"next":  nil,
					"prev":  nil,
					"total": 1,
					"results": []map[string]any{{
						"title":         "Example Activity",
						"activityUUIDs": []string{"11111111-2222-3333-4444-555555555555"},
						"time":          "2025-10-23T04:48:10.804405Z",
						"sourceEndpoints": []map[string]any{
							{
								"id":   97,
								"ip":   []string{"192.168.2.2"},
								"name": "192.168.2.2",
							},
						},
					}},
				},
			})
		},
	})
	defer cleanup()

	res, err := client.GetSearch(context.Background(), aql, includeSample, includeTotal)
	if err != nil {
		t.Fatalf("get search: %v", err)
	}
	if res.Total != 1 || len(res.Results) != 1 || res.Results[0].Title != "Example Activity" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestSearchEndpointIDUnmarshal(t *testing.T) {
	payload := []byte(`{"sourceEndpoints":[{"id":123},{"id":"456"}]}`)
	var res struct {
		Source []SearchEndpoint `json:"sourceEndpoints"`
	}
	if err := json.Unmarshal(payload, &res); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := res.Source[0].ID; got != SearchEndpointID("123") {
		t.Fatalf("expected 123, got %q", got)
	}
	if got := res.Source[1].ID; got != SearchEndpointID("456") {
		t.Fatalf("expected 456, got %q", got)
	}
}

func TestSearchEndpointIPsUnmarshal(t *testing.T) {
	t.Parallel()

	t.Run("slice input", func(t *testing.T) {
		t.Parallel()
		payload := []byte(`{"sourceEndpoints":[{"ip":["10.0.0.1","fe80::1"]}]}`)
		var res struct {
			Source []SearchEndpoint `json:"sourceEndpoints"`
		}
		if err := json.Unmarshal(payload, &res); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := []string(res.Source[0].IP)
		if len(got) != 2 || got[0] != "10.0.0.1" || got[1] != "fe80::1" {
			t.Fatalf("unexpected ips: %#v", got)
		}
	})

	t.Run("string input", func(t *testing.T) {
		t.Parallel()
		payload := []byte(`{"sourceEndpoints":[{"ip":"10.0.0.1"}]}`)
		var res struct {
			Source []SearchEndpoint `json:"sourceEndpoints"`
		}
		if err := json.Unmarshal(payload, &res); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := []string(res.Source[0].IP)
		if len(got) != 1 || got[0] != "10.0.0.1" {
			t.Fatalf("unexpected ips: %#v", got)
		}
	})
}

func TestGetAuditLogSearch(t *testing.T) {
	t.Parallel()

	const (
		aql           = "in:auditLogs"
		includeSample = false
		includeTotal  = false
	)

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/search/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}

			values := r.URL.Query()
			if got := values.Get("aql"); got != aql {
				t.Fatalf("unexpected aql: %q", got)
			}
			if got := values.Get("includeSample"); got != "false" {
				t.Fatalf("unexpected includeSample: %q", got)
			}
			if got := values.Get("includeTotal"); got != "false" {
				t.Fatalf("unexpected includeTotal: %q", got)
			}

			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"count": 1,
					"next":  10,
					"prev":  nil,
					"results": []map[string]any{{
						"action": "API Call",
						"additionalInfo": map[string]any{
							"data": "Endpoint: search/, Status: 200",
							"type": "TEXT",
						},
						"id":      3505970941,
						"time":    "2025-11-06T17:24:01.314669+00:00",
						"timeUtc": "2025-11-06T17:24:01.314669+00:00",
						"trigger": "User Action",
						"user":    "Michael Rosenfeld (michael.rosenfeld@1898andco.io)",
						"userIp":  "None",
					}},
				},
			})
		},
	})
	defer cleanup()

	res, err := client.GetSearch(context.Background(), aql, includeSample, includeTotal)
	if err != nil {
		t.Fatalf("get search: %v", err)
	}
	if len(res.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res.Results))
	}

	auditLog := res.Results[0]
	if auditLog.Action != "API Call" {
		t.Errorf("expected action 'API Call', got %q", auditLog.Action)
	}
	if auditLog.ID != 3505970941 {
		t.Errorf("expected id 3505970941, got %d", auditLog.ID)
	}
	if auditLog.Trigger != "User Action" {
		t.Errorf("expected trigger 'User Action', got %q", auditLog.Trigger)
	}
	if auditLog.User != "Michael Rosenfeld (michael.rosenfeld@1898andco.io)" {
		t.Errorf("expected user 'Michael Rosenfeld (michael.rosenfeld@1898andco.io)', got %q", auditLog.User)
	}
	if auditLog.UserIP != "None" {
		t.Errorf("expected userIp 'None', got %q", auditLog.UserIP)
	}
	if auditLog.Time != "2025-11-06T17:24:01.314669+00:00" {
		t.Errorf("expected time '2025-11-06T17:24:01.314669+00:00', got %q", auditLog.Time)
	}
	if auditLog.TimeUtc != "2025-11-06T17:24:01.314669+00:00" {
		t.Errorf("expected timeUtc '2025-11-06T17:24:01.314669+00:00', got %q", auditLog.TimeUtc)
	}
	if auditLog.AdditionalInfo == nil {
		t.Fatal("expected additionalInfo to be non-nil")
	}
	if auditLog.AdditionalInfo.Data != "Endpoint: search/, Status: 200" {
		t.Errorf("expected additionalInfo.data 'Endpoint: search/, Status: 200', got %q", auditLog.AdditionalInfo.Data)
	}
	if auditLog.AdditionalInfo.Type != "TEXT" {
		t.Errorf("expected additionalInfo.type 'TEXT', got %q", auditLog.AdditionalInfo.Type)
	}
}

func TestGetRiskFactorSearch(t *testing.T) {
	t.Parallel()

	const (
		aql           = "in:riskFactors"
		includeSample = false
		includeTotal  = false
	)

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/search/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}

			values := r.URL.Query()
			if got := values.Get("aql"); got != aql {
				t.Fatalf("unexpected aql: %q", got)
			}
			if got := values.Get("includeSample"); got != "false" {
				t.Fatalf("unexpected includeSample: %q", got)
			}
			if got := values.Get("includeTotal"); got != "false" {
				t.Fatalf("unexpected includeTotal: %q", got)
			}

			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"count": 2,
					"next":  nil,
					"prev":  nil,
					"total": 2,
					"results": []map[string]any{
						{
							"category":    "Profile",
							"description": "Windows 10, build 19045, will reach end-of-support on October 14, 2025",
							"devices":     1,
							"evidence": map[string]any{
								"AQL":          "",
								"whatHappened": "This device has been observed running an operating system that will reach its end-of-support date on October 14, 2025.",
							},
							"firstSeen": "2025-06-19T19:11:24.291901+00:00",
							"group":     "Deprecated SW/HW (Impending)",
							"lastSeen":  "2025-08-21T15:42:31.220047+00:00",
							"policy":    nil,
							"remediation": map[string]any{
								"category":    "System and Software Management",
								"description": "Upgrade devices running an end-of-life operating system to a current, supported version.",
								"recommendedActions": []map[string]any{
									{
										"description": "Ensure that the latest security patches and updates are applied to the operating systems.",
										"id":          1,
										"title":       "Apply Security Patches",
										"type":        "Remediation",
									},
									{
										"description": "If immediate OS upgrade is not possible, isolate the devices running end-of-life operating systems from the main network.",
										"id":          2,
										"title":       "Isolate Legacy OSes",
										"type":        "Mitigation",
									},
								},
								"type": "Upgrade Operating System",
							},
							"score":  "Low",
							"source": "Armis",
							"status": "Open",
							"type":   "Impending End-of-Support Operating System",
						},
						{
							"category":    "Behavioural",
							"description": "Policy Violation",
							"devices":     107,
							"evidence": map[string]any{
								"AQL":          "in:alerts classification:Security",
								"whatHappened": "These devices have been observed violating one or more existing policies.",
							},
							"firstSeen": "2025-02-19T16:27:44+00:00",
							"group":     "Policies",
							"lastSeen":  "2025-11-11T00:25:36+00:00",
							"policy":    nil,
							"remediation": map[string]any{
								"category":    "Compliance",
								"description": "Ensure all devices adhere to existing policies by updating configurations and conducting regular audits.",
								"type":        "Enforce Policy Compliance",
							},
							"score":  "Low",
							"source": "Armis",
							"status": "Open",
							"type":   "Policy Violations",
						},
					},
				},
			})
		},
	})
	defer cleanup()

	res, err := client.GetSearch(context.Background(), aql, includeSample, includeTotal)
	if err != nil {
		t.Fatalf("get search: %v", err)
	}
	if res.Total != 2 || len(res.Results) != 2 {
		t.Fatalf("expected 2 results, got %d (total: %d)", len(res.Results), res.Total)
	}

	// Test first risk factor (Profile category)
	rf1 := res.Results[0]
	if rf1.Category != "Profile" {
		t.Errorf("expected category 'Profile', got %q", rf1.Category)
	}
	if rf1.Description != "Windows 10, build 19045, will reach end-of-support on October 14, 2025" {
		t.Errorf("unexpected description: %q", rf1.Description)
	}
	if rf1.Devices != 1 {
		t.Errorf("expected devices 1, got %d", rf1.Devices)
	}
	if rf1.FirstSeen != "2025-06-19T19:11:24.291901+00:00" {
		t.Errorf("expected firstSeen '2025-06-19T19:11:24.291901+00:00', got %q", rf1.FirstSeen)
	}
	if rf1.LastSeen != "2025-08-21T15:42:31.220047+00:00" {
		t.Errorf("expected lastSeen '2025-08-21T15:42:31.220047+00:00', got %q", rf1.LastSeen)
	}
	if rf1.Group != "Deprecated SW/HW (Impending)" {
		t.Errorf("expected group 'Deprecated SW/HW (Impending)', got %q", rf1.Group)
	}
	if rf1.Score != "Low" {
		t.Errorf("expected score 'Low', got %q", rf1.Score)
	}
	if rf1.Source != "Armis" {
		t.Errorf("expected source 'Armis', got %q", rf1.Source)
	}
	if rf1.Status != "Open" {
		t.Errorf("expected status 'Open', got %q", rf1.Status)
	}
	if rf1.Type != "Impending End-of-Support Operating System" {
		t.Errorf("expected type 'Impending End-of-Support Operating System', got %q", rf1.Type)
	}

	// Test evidence
	if rf1.Evidence == nil {
		t.Fatal("expected evidence to be non-nil")
	}
	if rf1.Evidence.AQL != "" {
		t.Errorf("expected evidence.AQL to be empty, got %q", rf1.Evidence.AQL)
	}
	if rf1.Evidence.WhatHappened != "This device has been observed running an operating system that will reach its end-of-support date on October 14, 2025." {
		t.Errorf("unexpected evidence.whatHappened: %q", rf1.Evidence.WhatHappened)
	}

	// Test remediation
	if rf1.Remediation == nil {
		t.Fatal("expected remediation to be non-nil")
	}
	if rf1.Remediation.Category != "System and Software Management" {
		t.Errorf("expected remediation.category 'System and Software Management', got %q", rf1.Remediation.Category)
	}
	if rf1.Remediation.Type != "Upgrade Operating System" {
		t.Errorf("expected remediation.type 'Upgrade Operating System', got %q", rf1.Remediation.Type)
	}
	if len(rf1.Remediation.RecommendedActions) != 2 {
		t.Fatalf("expected 2 recommended actions, got %d", len(rf1.Remediation.RecommendedActions))
	}
	if rf1.Remediation.RecommendedActions[0].ID != 1 {
		t.Errorf("expected first action id 1, got %d", rf1.Remediation.RecommendedActions[0].ID)
	}
	if rf1.Remediation.RecommendedActions[0].Title != "Apply Security Patches" {
		t.Errorf("expected first action title 'Apply Security Patches', got %q", rf1.Remediation.RecommendedActions[0].Title)
	}
	if rf1.Remediation.RecommendedActions[0].Type != "Remediation" {
		t.Errorf("expected first action type 'Remediation', got %q", rf1.Remediation.RecommendedActions[0].Type)
	}

	// Test second risk factor (Behavioural category)
	rf2 := res.Results[1]
	if rf2.Category != "Behavioural" {
		t.Errorf("expected category 'Behavioural', got %q", rf2.Category)
	}
	if rf2.Devices != 107 {
		t.Errorf("expected devices 107, got %d", rf2.Devices)
	}
	if rf2.Evidence == nil {
		t.Fatal("expected evidence to be non-nil")
	}
	if rf2.Evidence.AQL != "in:alerts classification:Security" {
		t.Errorf("expected evidence.AQL 'in:alerts classification:Security', got %q", rf2.Evidence.AQL)
	}
}

func TestGetConnectionsSearch(t *testing.T) {
	t.Parallel()

	const (
		aql           = "in:connections protocol:BMS"
		includeSample = false
		includeTotal  = true
	)

	client, cleanup := newTestClient(t, map[string]http.HandlerFunc{
		"/api/v1/search/": func(w http.ResponseWriter, r *http.Request) {
			assertAuthHeader(t, r)
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET, got %s", r.Method)
			}

			values := r.URL.Query()
			if got := values.Get("aql"); got != aql {
				t.Fatalf("unexpected aql: %q", got)
			}
			if got := values.Get("includeSample"); got != "false" {
				t.Fatalf("unexpected includeSample: %q", got)
			}
			if got := values.Get("includeTotal"); got != "true" {
				t.Fatalf("unexpected includeTotal: %q", got)
			}

			respondJSON(t, w, http.StatusOK, map[string]any{
				"success": true,
				"data": map[string]any{
					"count": 1,
					"next":  nil,
					"prev":  nil,
					"total": 1,
					"results": []map[string]any{{
						"band":            nil,
						"bssid":           nil,
						"channel":         nil,
						"duration":        135745097,
						"endTimestamp":    "2030-01-01T00:00:00+00:00",
						"id":              27,
						"inboundTraffic":  0,
						"outboundTraffic": 0,
						"protocol":        "BMS",
						"risk":            "Medium",
						"rssi":            nil,
						"sensor": map[string]any{
							"name": "SPAN 8153 ens20 (PLC)",
							"type": "SPAN",
						},
						"site": map[string]any{
							"location": "Houston",
							"name":     "Lab",
						},
						"sites": []map[string]any{
							{
								"location": "Houston",
								"name":     "Lab",
							},
						},
						"snr":            nil,
						"sourceId":       7,
						"ssid":           nil,
						"startTimestamp": "2025-09-12T21:01:42.738562+00:00",
						"targetId":       319,
						"title":          "Connection between modbus device - unit id: 1 and Workstation",
						"traffic":        0,
					}},
				},
			})
		},
	})
	defer cleanup()

	res, err := client.GetSearch(context.Background(), aql, includeSample, includeTotal)
	if err != nil {
		t.Fatalf("get search: %v", err)
	}
	if res.Total != 1 || len(res.Results) != 1 {
		t.Fatalf("expected 1 result, got %d (total: %d)", len(res.Results), res.Total)
	}

	conn := res.Results[0]
	if conn.ID != 27 {
		t.Errorf("expected id 27, got %d", conn.ID)
	}
	if conn.Protocol != "BMS" {
		t.Errorf("expected protocol 'BMS', got %q", conn.Protocol)
	}
	if conn.Risk != "Medium" {
		t.Errorf("expected risk 'Medium', got %q", conn.Risk)
	}
	if conn.Duration != 135745097 {
		t.Errorf("expected duration 135745097, got %d", conn.Duration)
	}
	if conn.InboundTraffic != 0 {
		t.Errorf("expected inboundTraffic 0, got %d", conn.InboundTraffic)
	}
	if conn.OutboundTraffic != 0 {
		t.Errorf("expected outboundTraffic 0, got %d", conn.OutboundTraffic)
	}
	if conn.Traffic != 0 {
		t.Errorf("expected traffic 0, got %d", conn.Traffic)
	}
	if conn.SourceID != 7 {
		t.Errorf("expected sourceId 7, got %d", conn.SourceID)
	}
	if conn.TargetID != 319 {
		t.Errorf("expected targetId 319, got %d", conn.TargetID)
	}
	if conn.StartTimestamp != "2025-09-12T21:01:42.738562+00:00" {
		t.Errorf("expected startTimestamp '2025-09-12T21:01:42.738562+00:00', got %q", conn.StartTimestamp)
	}
	if conn.EndTimestamp != "2030-01-01T00:00:00+00:00" {
		t.Errorf("expected endTimestamp '2030-01-01T00:00:00+00:00', got %q", conn.EndTimestamp)
	}
	if conn.Title != "Connection between modbus device - unit id: 1 and Workstation" {
		t.Errorf("expected title 'Connection between modbus device - unit id: 1 and Workstation', got %q", conn.Title)
	}

	if conn.Sensor.Name != "SPAN 8153 ens20 (PLC)" {
		t.Errorf("expected sensor.name 'SPAN 8153 ens20 (PLC)', got %q", conn.Sensor.Name)
	}
	if conn.Sensor.Type != "SPAN" {
		t.Errorf("expected sensor.type 'SPAN', got %q", conn.Sensor.Type)
	}

	if conn.Site.Location != "Houston" {
		t.Errorf("expected site.location 'Houston', got %q", conn.Site.Location)
	}
	if conn.Site.Name != "Lab" {
		t.Errorf("expected site.name 'Lab', got %q", conn.Site.Name)
	}

	if len(conn.Sites) != 1 {
		t.Fatalf("expected 1 site, got %d", len(conn.Sites))
	}
	if conn.Sites[0].Location != "Houston" {
		t.Errorf("expected sites[0].location 'Houston', got %q", conn.Sites[0].Location)
	}
	if conn.Sites[0].Name != "Lab" {
		t.Errorf("expected sites[0].name 'Lab', got %q", conn.Sites[0].Name)
	}
}

func TestGetSearchRequiresAQL(t *testing.T) {
	t.Parallel()

	client, cleanup := newTestClient(t, nil)
	defer cleanup()

	if _, err := client.GetSearch(context.Background(), " ", false, false); !errors.Is(err, ErrSearchAQL) {
		t.Fatalf("expected ErrSearchAQL, got %v", err)
	}
}
