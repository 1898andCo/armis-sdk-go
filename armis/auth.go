// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package armis

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// authResponse mirrors the JSON returned by /access_token/. All field names
// follow Go's initialism rules while the struct tags preserve Armis' schema.
type authResponse struct {
	Success bool `json:"success"`
	Data    struct {
		AccessToken   string `json:"access_token"`
		ExpirationUTC string `json:"expiration_utc"`
		UserID        int    `json:"user_id"`
	} `json:"data"`
}

// authenticate exchanges the API key for a short-lived bearer token. It is
// concurrency-safe and returns early if the cached token is still valid.
func (c *Client) authenticate(ctx context.Context) error {
	// This uses double-checked locking to prevent a TOCTOU race: an RLock fast
	// path checks the token, then a Lock with a second check ensures only one
	// goroutine authenticates if multiple detect an expired token simultaneously.
	c.mu.RLock()
	if c.accessToken != "" && time.Now().Before(c.accessTokenExpires) {
		c.mu.RUnlock()
		return nil // cached token still good
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.accessToken != "" && time.Now().Before(c.accessTokenExpires) {
		return nil // another goroutine already authenticated
	}

	form := url.Values{}
	form.Set("secret_key", c.apiKey)

	endpoint := fmt.Sprintf("%s/api/%s/access_token/", c.apiURL, c.apiVersion)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("armis: build auth request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return &APIError{StatusCode: res.StatusCode, Body: b}
	}

	var ar authResponse
	if err := json.NewDecoder(res.Body).Decode(&ar); err != nil {
		return fmt.Errorf("armis: decode auth response: %w", err)
	}
	if !ar.Success {
		return ErrAuthFailed
	}

	expiry, err := time.Parse(time.RFC3339Nano, ar.Data.ExpirationUTC)
	if err != nil {
		return fmt.Errorf("armis: parse expiry: %w", err)
	}
	expiry = expiry.Add(-5 * time.Minute) // expire early for safety

	c.accessToken = ar.Data.AccessToken
	c.accessTokenExpires = expiry
	c.userID = ar.Data.UserID

	return nil
}
