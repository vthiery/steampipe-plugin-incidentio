package incidentio

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// IncidentRole represents a role definition (e.g. Incident Lead).
type IncidentRole struct {
	CreatedAt    string `json:"created_at"`
	Description  string `json:"description"`
	ID           string `json:"id"`
	Instructions string `json:"instructions"`
	Name         string `json:"name"`
	Required     bool   `json:"required"`
	RoleType     string `json:"role_type"`
	Shortform    string `json:"shortform"`
	UpdatedAt    string `json:"updated_at"`
}

// incidentRolesResponse is the envelope returned by GET /v1/incident_roles.
type incidentRolesResponse struct {
	IncidentRoles []IncidentRole `json:"incident_roles"`
}

//// TABLE DEFINITION

func tableIncidentioIncidentRoles() *plugin.Table {
	return &plugin.Table{
		Name:        "incidentio_incident_roles",
		Description: "List incident role definitions configured in incident.io.",
		List: &plugin.ListConfig{
			Hydrate: listIncidentRoles,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the role."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the role (e.g. Incident Lead)."},
			{Name: "shortform", Type: proto.ColumnType_STRING, Description: "Short identifier for the role (e.g. lead)."},
			{Name: "role_type", Type: proto.ColumnType_STRING, Description: "Type of role (lead, reporter, custom)."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the role."},
			{Name: "instructions", Type: proto.ColumnType_STRING, Description: "Instructions shown to responders assigned this role."},
			{Name: "required", Type: proto.ColumnType_BOOL, Description: "Whether the role must be filled for an incident."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the role was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the role was last updated."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listIncidentRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result incidentRolesResponse
	if err := client.get(ctx, "/v1/incident_roles", nil, &result); err != nil {
		return nil, fmt.Errorf("listing incident roles: %w", err)
	}

	for _, r := range result.IncidentRoles {
		d.StreamListItem(ctx, r)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
