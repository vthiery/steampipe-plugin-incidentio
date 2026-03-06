package incidentio

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// incidentTypesResponse is the envelope returned by GET /v1/incident_types.
type incidentTypesResponse struct {
	IncidentTypes []IncidentType `json:"incident_types"`
}

//// TABLE DEFINITION

func tableIncidentioIncidentType() *plugin.Table {
	return &plugin.Table{
		Name:        "incidentio_incident_type",
		Description: "List incident types configured in incident.io.",
		List: &plugin.ListConfig{
			Hydrate: listIncidentTypes,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the incident type."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the incident type."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the incident type."},
			{Name: "is_default", Type: proto.ColumnType_BOOL, Description: "Whether this is the default incident type."},
			{Name: "private_incidents_only", Type: proto.ColumnType_BOOL, Description: "Whether incidents of this type are always private."},
			{Name: "create_in_triage", Type: proto.ColumnType_STRING, Description: "Whether incidents of this type are created in triage (always, optional)."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the incident type was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the incident type was last updated."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listIncidentTypes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result incidentTypesResponse
	if err := client.get(ctx, "/v1/incident_types", nil, &result); err != nil {
		return nil, fmt.Errorf("listing incident types: %w", err)
	}

	for _, t := range result.IncidentTypes {
		d.StreamListItem(ctx, t)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
