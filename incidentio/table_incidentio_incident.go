package incidentio

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// IncidentStatus represents the status of an incident.
type IncidentStatus struct {
	Category  string `json:"category"`
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Rank      int    `json:"rank"`
	UpdatedAt string `json:"updated_at"`
}

// Severity represents an incident severity level.
type Severity struct {
	CreatedAt   string `json:"created_at"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Rank        int    `json:"rank"`
	UpdatedAt   string `json:"updated_at"`
}

// IncidentType represents the type of incident.
type IncidentType struct {
	CreateInTriage       string `json:"create_in_triage"`
	CreatedAt            string `json:"created_at"`
	Description          string `json:"description"`
	ID                   string `json:"id"`
	IsDefault            bool   `json:"is_default"`
	Name                 string `json:"name"`
	PrivateIncidentsOnly bool   `json:"private_incidents_only"`
	UpdatedAt            string `json:"updated_at"`
}

// ExternalIssueReference represents a linked issue in an external tracker.
type ExternalIssueReference struct {
	IssueName      string `json:"issue_name"`
	IssuePermalink string `json:"issue_permalink"`
	Provider       string `json:"provider"`
}

// Incident represents an incident.io incident.
type Incident struct {
	CallURL                 string                 `json:"call_url"`
	CreatedAt               string                 `json:"created_at"`
	CustomFieldEntries      []interface{}          `json:"custom_field_entries"`
	DurationMetrics         []interface{}          `json:"duration_metrics"`
	ExternalIssueReference  ExternalIssueReference `json:"external_issue_reference"`
	HasDebrief              bool                   `json:"has_debrief"`
	ID                      string                 `json:"id"`
	IncidentRoleAssignments []interface{}          `json:"incident_role_assignments"`
	IncidentStatus          IncidentStatus         `json:"incident_status"`
	IncidentTimestampValues []interface{}          `json:"incident_timestamp_values"`
	IncidentType            *IncidentType          `json:"incident_type"`
	Mode                    string                 `json:"mode"`
	Name                    string                 `json:"name"`
	Permalink               string                 `json:"permalink"`
	PostmortemDocumentIDs   []string               `json:"postmortem_document_ids"`
	PostmortemDocumentURL   string                 `json:"postmortem_document_url"`
	Reference               string                 `json:"reference"`
	Severity                *Severity              `json:"severity"`
	SlackChannelID          string                 `json:"slack_channel_id"`
	SlackChannelName        string                 `json:"slack_channel_name"`
	SlackTeamID             string                 `json:"slack_team_id"`
	Summary                 string                 `json:"summary"`
	UpdatedAt               string                 `json:"updated_at"`
	Visibility              string                 `json:"visibility"`
	WorkloadMinutesLate     float64                `json:"workload_minutes_late"`
	WorkloadMinutesSleeping float64                `json:"workload_minutes_sleeping"`
	WorkloadMinutesTotal    float64                `json:"workload_minutes_total"`
	WorkloadMinutesWorking  float64                `json:"workload_minutes_working"`
}

// incidentsResponse is the envelope returned by GET /v2/incidents.
type incidentsResponse struct {
	Incidents      []Incident     `json:"incidents"`
	PaginationMeta PaginationMeta `json:"pagination_meta"`
}

// incidentResponse is the envelope returned by GET /v2/incidents/{id}.
type incidentResponse struct {
	Incident Incident `json:"incident"`
}

//// TABLE DEFINITION

func tableIncidentioIncident() *plugin.Table {
	return &plugin.Table{
		Name:        "incidentio_incident",
		Description: "List and inspect incidents from incident.io.",
		List: &plugin.ListConfig{
			Hydrate: listIncidents,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getIncident,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the incident."},
			{Name: "reference", Type: proto.ColumnType_STRING, Description: "Human-readable reference (e.g. INC-123)."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the incident."},
			{Name: "summary", Type: proto.ColumnType_STRING, Description: "Summary of the incident."},
			{Name: "mode", Type: proto.ColumnType_STRING, Description: "Mode of the incident (standard, retrospective, test, tutorial)."},
			{Name: "visibility", Type: proto.ColumnType_STRING, Description: "Visibility of the incident (public or private)."},
			{Name: "permalink", Type: proto.ColumnType_STRING, Description: "URL of the incident in the incident.io app."},
			{Name: "call_url", Type: proto.ColumnType_STRING, Description: "URL of the call associated with the incident."},
			{Name: "has_debrief", Type: proto.ColumnType_BOOL, Description: "Whether the incident has a debrief."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the incident was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the incident was last updated."},

			// Status
			{Name: "status_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("IncidentStatus.ID"), Description: "ID of the current incident status."},
			{Name: "status_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("IncidentStatus.Name"), Description: "Name of the current incident status."},
			{Name: "status_category", Type: proto.ColumnType_STRING, Transform: transform.FromField("IncidentStatus.Category"), Description: "Category of the current status (triage, live, learning, closed, etc)."},

			// Severity
			{Name: "severity_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Severity.ID"), Description: "ID of the incident severity."},
			{Name: "severity_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Severity.Name"), Description: "Name of the incident severity."},
			{Name: "severity_rank", Type: proto.ColumnType_INT, Transform: transform.FromField("Severity.Rank"), Description: "Rank of the incident severity (lower = more severe)."},

			// Incident type
			{Name: "incident_type_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("IncidentType.ID"), Description: "ID of the incident type."},
			{Name: "incident_type_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("IncidentType.Name"), Description: "Name of the incident type."},

			// Slack
			{Name: "slack_channel_id", Type: proto.ColumnType_STRING, Description: "ID of the Slack channel for this incident."},
			{Name: "slack_channel_name", Type: proto.ColumnType_STRING, Description: "Name of the Slack channel for this incident."},
			{Name: "slack_team_id", Type: proto.ColumnType_STRING, Description: "ID of the Slack team."},

			// Postmortem
			{Name: "postmortem_document_url", Type: proto.ColumnType_STRING, Description: "URL of the postmortem document."},
			{Name: "postmortem_document_ids", Type: proto.ColumnType_JSON, Description: "IDs of the postmortem documents."},

			// External issue
			{Name: "external_issue_reference", Type: proto.ColumnType_JSON, Description: "Linked issue in an external tracker (e.g. Linear, Jira)."},

			// Workload
			{Name: "workload_minutes_total", Type: proto.ColumnType_DOUBLE, Description: "Total workload in minutes spent on this incident."},
			{Name: "workload_minutes_working", Type: proto.ColumnType_DOUBLE, Description: "Working-hours workload in minutes."},
			{Name: "workload_minutes_sleeping", Type: proto.ColumnType_DOUBLE, Description: "Sleeping-hours workload in minutes."},
			{Name: "workload_minutes_late", Type: proto.ColumnType_DOUBLE, Description: "Late-hours workload in minutes."},

			// Full JSON blobs
			{Name: "incident_role_assignments", Type: proto.ColumnType_JSON, Description: "Role assignments for the incident."},
			{Name: "custom_field_entries", Type: proto.ColumnType_JSON, Description: "Custom field values set on the incident."},
			{Name: "incident_timestamp_values", Type: proto.ColumnType_JSON, Description: "Timestamps recorded for this incident."},
			{Name: "duration_metrics", Type: proto.ColumnType_JSON, Description: "Duration metrics for this incident."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listIncidents(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

		var result incidentsResponse
		if err := client.get(ctx, "/v2/incidents", params, &result); err != nil {
			return nil, fmt.Errorf("listing incidents: %w", err)
		}

		for _, incident := range result.Incidents {
			d.StreamListItem(ctx, incident)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.PaginationMeta.After == "" || len(result.Incidents) == 0 {
			break
		}
		after = result.PaginationMeta.After
	}

	return nil, nil
}

func getIncident(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result incidentResponse
	if err := client.get(ctx, fmt.Sprintf("/v2/incidents/%s", id), nil, &result); err != nil {
		return nil, fmt.Errorf("getting incident %s: %w", id, err)
	}

	return result.Incident, nil
}
