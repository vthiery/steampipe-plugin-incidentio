package incidentio

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// incidentStatusesResponse is the envelope returned by GET /v1/incident_statuses.
type incidentStatusesResponse struct {
	IncidentStatuses []IncidentStatus `json:"incident_statuses"`
}

//// TABLE DEFINITION

func tableIncidentioIncidentStatuses() *plugin.Table {
	return &plugin.Table{
		Name:        "incidentio_incident_statuses",
		Description: "List incident status definitions configured in incident.io.",
		List: &plugin.ListConfig{
			Hydrate: listIncidentStatuses,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the status."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the status (e.g. Closed)."},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "Category of the status (triage, live, learning, closed, etc)."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the status."},
			{Name: "rank", Type: proto.ColumnType_INT, Description: "Ordering rank of the status."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the status was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the status was last updated."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listIncidentStatuses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result incidentStatusesResponse
	if err := client.get(ctx, "/v1/incident_statuses", nil, &result); err != nil {
		return nil, fmt.Errorf("listing incident statuses: %w", err)
	}

	for _, s := range result.IncidentStatuses {
		d.StreamListItem(ctx, s)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
