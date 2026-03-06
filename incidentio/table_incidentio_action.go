package incidentio

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// Action represents an action item linked to an incident.
type Action struct {
	Assignee    *User  `json:"assignee"`
	CompletedAt string `json:"completed_at"`
	CreatedAt   string `json:"created_at"`
	Description string `json:"description"`
	ID          string `json:"id"`
	IncidentID  string `json:"incident_id"`
	Status      string `json:"status"`
	UpdatedAt   string `json:"updated_at"`
}

// actionsResponse is the envelope returned by GET /v2/actions.
type actionsResponse struct {
	Actions []Action `json:"actions"`
}

//// TABLE DEFINITION

func tableIncidentioAction() *plugin.Table {
	return &plugin.Table{
		Name:        "incidentio_action",
		Description: "List action items tracked against incidents in incident.io.",
		List: &plugin.ListConfig{
			Hydrate: listActions,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "incident_id", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the action."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the action."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Current status of the action (outstanding, completed, deleted, not_doing)."},
			{Name: "incident_id", Type: proto.ColumnType_STRING, Description: "ID of the incident this action belongs to."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the action was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the action was last updated."},
			{Name: "completed_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the action was completed (if applicable)."},
			{Name: "assignee", Type: proto.ColumnType_JSON, Description: "User assigned to this action."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listActions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{}
	if id := d.EqualsQualString("incident_id"); id != "" {
		params["incident_id"] = id
	}

	var result actionsResponse
	if err := client.get(ctx, "/v2/actions", params, &result); err != nil {
		return nil, fmt.Errorf("listing actions: %w", err)
	}

	for _, action := range result.Actions {
		d.StreamListItem(ctx, action)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
