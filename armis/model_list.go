// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

type ListSettings struct {
	CreatedBy      string `json:"created_by"`
	CreationTime   string `json:"creation_time"`
	Description    string `json:"description"`
	LastUpdateTime string `json:"last_update_time"`
	LastUpdatedBy  string `json:"last_updated_by"`
	ListID         int    `json:"list_id"`
	ListName       string `json:"list_name"`
	ListType       string `json:"list_type"`
}

// ListsAPIResponse matches the Armis lists endpoint, where data contains a nested lists array.
type ListsAPIResponse struct {
	Data struct {
		Lists []ListSettings `json:"lists"`
	} `json:"data"`
	Success bool `json:"success,omitempty"`
}
