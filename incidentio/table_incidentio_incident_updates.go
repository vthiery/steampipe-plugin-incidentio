package incidentio

import (
	"context"
	"errors"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// IncidentUpdate represents a status update posted on an incident.
type IncidentUpdate struct {
	CreatedAt            string          `json:"created_at"`
	ID                   string          `json:"id"`
	IncidentID           string          `json:"incident_id"`
	MergedIntoIncidentID string          `json:"merged_into_incident_id"`
	Message              string          `json:"message"`
	NewIncidentStatus    *IncidentStatus `json:"new_incident_status"`
	NewSeverity          *Severity       `json:"new_severity"`
	Updater              interface{}     `json:"updater"`
}

// incidentUpdatesResponse is the envelope returned by GET /v2/incident_updates.
type incidentUpdatesResponse struct {
	IncidentUpdates []IncidentUpdate `json:"incident_updates"`
	PaginationMeta  PaginationMeta   `json:"pagination_meta"`
}

//// TABLE DEFINITION

func tableIncidentioIncidentUpdates() *plugin.Table {
	return &plugin.Table{
		Name:        "incidentio_incident_updates",
		Description: "List status updates posted on incidents in incident.io.",
		List: &plugin.ListConfig{
			Hydrate: listIncidentUpdates,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "incident_id", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the incident update."},
			{Name: "incident_id", Type: proto.ColumnType_STRING, Description: "ID of the incident this update belongs to."},
			{Name: "message", Type: proto.ColumnType_STRING, Description: "Message content of the update."},
			{Name: "merged_into_incident_id", Type: proto.ColumnType_STRING, Description: "If set, the incident was merged into this incident ID at this update."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the update was created."},
			// New status fields
			{Name: "new_status_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NewIncidentStatus.ID"), Description: "ID of the status set by this update."},
			{Name: "new_status_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("NewIncidentStatus.Name"), Description: "Name of the status set by this update."},
			{Name: "new_status_category", Type: proto.ColumnType_STRING, Transform: transform.FromField("NewIncidentStatus.Category"), Description: "Category of the status set by this update."},
			// New severity fields
			{Name: "new_severity_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NewSeverity.ID"), Description: "ID of the severity set by this update."},
			{Name: "new_severity_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("NewSeverity.Name"), Description: "Name of the severity set by this update."},
			{Name: "new_severity_rank", Type: proto.ColumnType_INT, Transform: transform.FromField("NewSeverity.Rank"), Description: "Rank of the severity set by this update."},
			// Updater
			{Name: "updater", Type: proto.ColumnType_JSON, Description: "The user, API key, or workflow that created this update."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listIncidentUpdates(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	pageSize := "250"
	after := ""

	for {
		params := map[string]string{"page_size": pageSize}
		if after != "" {
			params["after"] = after
		}
		if id := d.EqualsQualString("incident_id"); id != "" {
			params["incident_id"] = id
		}

		var result incidentUpdatesResponse
		if err := client.get(ctx, "/v2/incident_updates", params, &result); err != nil {
			if errors.Is(err, ErrNotFound) {
				return nil, nil
			}
			return nil, fmt.Errorf("listing incident updates: %w", err)
		}

		for _, u := range result.IncidentUpdates {
			d.StreamListItem(ctx, u)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.PaginationMeta.After == "" || len(result.IncidentUpdates) == 0 {
			break
		}
		after = result.PaginationMeta.After
	}

	return nil, nil
}
