package incidentio

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// FollowUpPriority represents the priority of a follow-up.
type FollowUpPriority struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Rank        int    `json:"rank"`
}

// FollowUp represents a follow-up item linked to an incident.
type FollowUp struct {
	Assignee               *User                  `json:"assignee"`
	CompletedAt            string                 `json:"completed_at"`
	CreatedAt              string                 `json:"created_at"`
	Description            string                 `json:"description"`
	ExternalIssueReference ExternalIssueReference `json:"external_issue_reference"`
	ID                     string                 `json:"id"`
	IncidentID             string                 `json:"incident_id"`
	Labels                 []string               `json:"labels"`
	Priority               *FollowUpPriority      `json:"priority"`
	Status                 string                 `json:"status"`
	Title                  string                 `json:"title"`
	UpdatedAt              string                 `json:"updated_at"`
}

// followUpsResponse is the envelope returned by GET /v2/follow_ups.
type followUpsResponse struct {
	FollowUps []FollowUp `json:"follow_ups"`
}

//// TABLE DEFINITION

func tableIncidentioFollowups() *plugin.Table {
	return &plugin.Table{
		Name:        "incidentio_followups",
		Description: "List follow-up items tracked against incidents in incident.io.",
		List: &plugin.ListConfig{
			Hydrate: listFollowUps,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "incident_id", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the follow-up."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the follow-up."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the follow-up."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Current status of the follow-up (outstanding, completed, deleted, not_doing)."},
			{Name: "incident_id", Type: proto.ColumnType_STRING, Description: "ID of the incident this follow-up belongs to."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Labels attached to the follow-up."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the follow-up was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the follow-up was last updated."},
			{Name: "completed_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the follow-up was completed (if applicable)."},
			{Name: "assignee", Type: proto.ColumnType_JSON, Description: "User assigned to this follow-up."},
			{Name: "priority_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Priority.ID"), Description: "ID of the follow-up priority."},
			{Name: "priority_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Priority.Name"), Description: "Name of the follow-up priority."},
			{Name: "priority_rank", Type: proto.ColumnType_INT, Transform: transform.FromField("Priority.Rank"), Description: "Rank of the follow-up priority."},
			{Name: "external_issue_reference", Type: proto.ColumnType_JSON, Description: "Linked issue in an external tracker (e.g. Linear, Jira)."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listFollowUps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{}
	if id := d.EqualsQualString("incident_id"); id != "" {
		params["incident_id"] = id
	}

	var result followUpsResponse
	if err := client.get(ctx, "/v2/follow_ups", params, &result); err != nil {
		return nil, fmt.Errorf("listing follow-ups: %w", err)
	}

	for _, f := range result.FollowUps {
		d.StreamListItem(ctx, f)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
