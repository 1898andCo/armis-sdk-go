// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

// Struct to match the entire API response.
type SiteApiResponse struct {
	Data    Site `json:"data"`
	Success bool `json:"success,omitempty"`
}

// Struct for the "data" key.
type Site struct {
	Count int            `json:"count,omitempty"`
	Next  *int           `json:"next"` // Handle null as a pointer
	Prev  int            `json:"prev"`
	Sites []SiteSettings `json:"sites,omitempty"`
}

// Struct for individual site settings.
type SiteSettings struct {
	ID       string  `json:"id,omitempty"`
	Lat      float64 `json:"lat,omitempty"`
	Lng      float64 `json:"lng,omitempty"`
	Location string  `json:"location,omitempty"`
	Name     string  `json:"name,omitempty"`
	Tier     string  `json:"tier,omitempty"`
	User     string  `json:"user,omitempty"`
}
