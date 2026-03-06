// Package incidentio provides a Steampipe plugin for querying incident.io resources using SQL.
package incidentio

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Plugin returns the definition of the incident.io Steampipe plugin.
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-incidentio",
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			"incidentio_incident":          tableIncidentioIncident(),
			"incidentio_action":            tableIncidentioAction(),
			"incidentio_severity":          tableIncidentioSeverity(),
			"incidentio_incident_type":     tableIncidentioIncidentType(),
			"incidentio_followups":         tableIncidentioFollowups(),
			"incidentio_incident_updates":  tableIncidentioIncidentUpdates(),
			"incidentio_users":             tableIncidentioUsers(),
			"incidentio_alerts":            tableIncidentioAlerts(),
			"incidentio_incident_roles":    tableIncidentioIncidentRoles(),
			"incidentio_incident_statuses": tableIncidentioIncidentStatuses(),
			"incidentio_custom_fields":     tableIncidentioCustomFields(),
			"incidentio_escalations":       tableIncidentioEscalations(),
		},
	}
	return p
}
