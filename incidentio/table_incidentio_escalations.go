package incidentio

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// EscalationPriority represents the priority of an escalation.
type EscalationPriority struct {
	Name string `json:"name"`
}

// Escalation represents an escalation event.
type Escalation struct {
	CreatedAt        string              `json:"created_at"`
	Creator          interface{}         `json:"creator"`
	EscalationPathID string              `json:"escalation_path_id"`
	Events           []interface{}       `json:"events"`
	ID               string              `json:"id"`
	Priority         *EscalationPriority `json:"priority"`
	RelatedAlerts    []interface{}       `json:"related_alerts"`
	RelatedIncidents []interface{}       `json:"related_incidents"`
	Status           string              `json:"status"`
	Title            string              `json:"title"`
	UpdatedAt        string              `json:"updated_at"`
}

// escalationsResponse is the envelope returned by GET /v2/escalations.
type escalationsResponse struct {
	Escalations    []Escalation   `json:"escalations"`
	PaginationMeta PaginationMeta `json:"pagination_meta"`
}

//// TABLE DEFINITION

func tableIncidentioEscalations() *plugin.Table {
	return &plugin.Table{
		Name:        "incidentio_escalations",
		Description: "List escalations in your incident.io account.",
		List: &plugin.ListConfig{
			Hydrate: listEscalations,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the escalation."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the escalation."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Current status (pending, triggered, acked, resolved, expired, cancelled)."},
			{Name: "escalation_path_id", Type: proto.ColumnType_STRING, Description: "ID of the escalation path used."},
			{Name: "priority_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Priority.Name"), Description: "Priority name of the escalation (e.g. P1)."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the escalation was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the escalation was last updated."},
			{Name: "creator", Type: proto.ColumnType_JSON, Description: "The user, alert, or workflow that triggered the escalation."},
			{Name: "events", Type: proto.ColumnType_JSON, Description: "Timeline of events for this escalation."},
			{Name: "related_alerts", Type: proto.ColumnType_JSON, Description: "Alerts related to this escalation."},
			{Name: "related_incidents", Type: proto.ColumnType_JSON, Description: "Incidents related to this escalation."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listEscalations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

		var result escalationsResponse
		if err := client.get(ctx, "/v2/escalations", params, &result); err != nil {
			return nil, fmt.Errorf("listing escalations: %w", err)
		}

		for _, e := range result.Escalations {
			d.StreamListItem(ctx, e)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.PaginationMeta.After == "" || len(result.Escalations) == 0 {
			break
		}
		after = result.PaginationMeta.After
	}

	return nil, nil
}
