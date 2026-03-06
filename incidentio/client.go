package incidentio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const baseURL = "https://api.incident.io"

// Client wraps the incident.io REST API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// newClient creates a new Client using the given API key.
func newClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// getClient retrieves a configured Client from the plugin connection config.
func getClient(ctx context.Context, d *plugin.QueryData) (*Client, error) {
	cfg := GetConfig(d.Connection)
	if cfg.APIKey == nil {
		return nil, fmt.Errorf("api_key must be configured for the incidentio plugin")
	}
	return newClient(*cfg.APIKey), nil
}

// get performs an authenticated GET request and unmarshals the response body into result.
func (c *Client) get(ctx context.Context, path string, params map[string]string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	if params != nil {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("incident.io API error: status=%d body=%s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("unmarshalling response: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Shared sub-types used across multiple tables
// ---------------------------------------------------------------------------

// User represents of an incident.io user reference.
type User struct {
	Email       string `json:"email"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	SlackUserID string `json:"slack_user_id"`
}

// PaginationMeta holds cursor-based pagination info returned by list endpoints.
type PaginationMeta struct {
	After            string `json:"after"`
	PageSize         int    `json:"page_size"`
	TotalRecordCount int    `json:"total_record_count"`
}
