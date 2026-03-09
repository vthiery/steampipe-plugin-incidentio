## v0.0.2 [2026-03-09]

_Bug fixes_

- `incidentio_incident` — querying by `id` for a non-existent incident now returns an empty result set instead of a 404 error ([#incidentio_incident](docs/tables/incidentio_incident.md))
- `incidentio_incident_updates` — filtering by `incident_id` for a non-existent incident now returns an empty result set instead of a 404 error ([#incidentio_incident_updates](docs/tables/incidentio_incident_updates.md))

## v0.0.1 [2026-03-06]

_What's new?_

- Initial release of the plugin with 12 tables:
  - [incidentio_action](docs/tables/incidentio_action.md) — list action items tracked against incidents
  - [incidentio_alerts](docs/tables/incidentio_alerts.md) — list alerts in your incident.io account
  - [incidentio_custom_fields](docs/tables/incidentio_custom_fields.md) — list custom field definitions
  - [incidentio_escalations](docs/tables/incidentio_escalations.md) — list escalations in your incident.io account
  - [incidentio_followups](docs/tables/incidentio_followups.md) — list follow-up items tracked against incidents
  - [incidentio_incident](docs/tables/incidentio_incident.md) — list and inspect incidents
  - [incidentio_incident_roles](docs/tables/incidentio_incident_roles.md) — list incident role definitions
  - [incidentio_incident_statuses](docs/tables/incidentio_incident_statuses.md) — list incident status definitions
  - [incidentio_incident_type](docs/tables/incidentio_incident_type.md) — list incident types
  - [incidentio_incident_updates](docs/tables/incidentio_incident_updates.md) — list status updates posted on incidents
  - [incidentio_severity](docs/tables/incidentio_severity.md) — list severity levels
  - [incidentio_users](docs/tables/incidentio_users.md) — list users in your incident.io account
