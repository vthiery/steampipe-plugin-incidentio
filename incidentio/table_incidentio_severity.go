package incidentio

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// severitiesResponse is the envelope returned by GET /v1/severities.
type severitiesResponse struct {
	Severities []Severity `json:"severities"`
}

//// TABLE DEFINITION

func tableIncidentioSeverity() *plugin.Table {
	return &plugin.Table{
		Name:        "incidentio_severity",
		Description: "List severity levels configured in incident.io.",
		List: &plugin.ListConfig{
			Hydrate: listSeverities,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the severity."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the severity level."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the severity level."},
			{Name: "rank", Type: proto.ColumnType_INT, Description: "Rank of the severity (lower = more severe)."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the severity was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the severity was last updated."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listSeverities(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result severitiesResponse
	if err := client.get(ctx, "/v1/severities", nil, &result); err != nil {
		return nil, fmt.Errorf("listing severities: %w", err)
	}

	for _, severity := range result.Severities {
		d.StreamListItem(ctx, severity)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
