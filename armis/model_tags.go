// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

// TagsAPIResponse represents the response structure for retrieving tags.
type TagsAPIResponse struct {
	Data    TagsData `json:"data"`
	Success bool     `json:"success,omitempty"`
}

// TagsData represents the data field in the tags API response.
type TagsData struct {
	Count int      `json:"count,omitempty"`
	Next  *int     `json:"next"`
	Prev  *int     `json:"prev"`
	Tags  []string `json:"tags,omitempty"`
}
