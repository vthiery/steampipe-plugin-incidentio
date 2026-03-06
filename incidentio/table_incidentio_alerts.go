package incidentio

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// Alert represents an alert from incident.io.
type Alert struct {
	AlertSourceID    string        `json:"alert_source_id"`
	Attributes       []interface{} `json:"attributes"`
	CreatedAt        string        `json:"created_at"`
	DeduplicationKey string        `json:"deduplication_key"`
	Description      string        `json:"description"`
	ID               string        `json:"id"`
	ResolvedAt       string        `json:"resolved_at"`
	SourceURL        string        `json:"source_url"`
	Status           string        `json:"status"`
	Title            string        `json:"title"`
	UpdatedAt        string        `json:"updated_at"`
}

// alertsResponse is the envelope returned by GET /v2/alerts.
type alertsResponse struct {
	Alerts         []Alert        `json:"alerts"`
	PaginationMeta PaginationMeta `json:"pagination_meta"`
}

//// TABLE DEFINITION

func tableIncidentioAlerts() *plugin.Table {
	return &plugin.Table{
		Name:        "incidentio_alerts",
		Description: "List alerts in your incident.io account.",
		List: &plugin.ListConfig{
			Hydrate: listAlerts,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the alert."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the alert."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the alert."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Current status of the alert (firing, resolved)."},
			{Name: "alert_source_id", Type: proto.ColumnType_STRING, Description: "ID of the alert source that produced this alert."},
			{Name: "deduplication_key", Type: proto.ColumnType_STRING, Description: "Deduplication key used to group related alerts."},
			{Name: "source_url", Type: proto.ColumnType_STRING, Description: "URL to the alert in the originating system."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the alert was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the alert was last updated."},
			{Name: "resolved_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the alert was resolved (if applicable)."},
			{Name: "attributes", Type: proto.ColumnType_JSON, Description: "Attributes attached to the alert."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listAlerts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	after := ""
	for {
		params := map[string]string{"page_size": "50"} // max for this endpoint
		if after != "" {
			params["after"] = after
		}

		var result alertsResponse
		if err := client.get(ctx, "/v2/alerts", params, &result); err != nil {
			return nil, fmt.Errorf("listing alerts: %w", err)
		}

		for _, a := range result.Alerts {
			d.StreamListItem(ctx, a)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.PaginationMeta.After == "" || len(result.Alerts) == 0 {
			break
		}
		after = result.PaginationMeta.After
	}

	return nil, nil
}
