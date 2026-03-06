# incident.io Plugin for Steampipe

Use SQL to query incidents, actions, alerts, escalations, schedules, custom fields, and more from [incident.io](https://incident.io).

- **[Get started →](https://github.com/vthiery/steampipe-plugin-incidentio)**
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/vthiery/steampipe-plugin-incidentio/issues)

## Quick start

### Install

Install the plugin with [Steampipe](https://steampipe.io):

```sh
steampipe plugin install ghcr.io/vthiery/incidentio
```

### Configure

Copy the sample config and set your API key:

```sh
cp config/incidentio.spc ~/.steampipe/config/incidentio.spc
```

Edit `~/.steampipe/config/incidentio.spc`:

```hcl
connection "incidentio" {
  plugin  = "ghcr.io/vthiery/incidentio"

  # API key from your incident.io dashboard → Settings → API keys.
  # See https://docs.incident.io/api-reference/introduction for details.
  api_key = "YOUR_API_KEY"
}
```

### Run a query

```shell
steampipe query
```

List all live incidents:

```sql
select
  id,
  reference,
  name,
  status_category,
  severity_name
from
  incidentio_incident
where
  status_category = 'live';
```

```
+-------------------------+----------+-------------------------------+-----------------+---------------+
| id                      | reference | name                         | status_category | severity_name |
+-------------------------+----------+-------------------------------+-----------------+---------------+
| 01FDAG4SAP5TYPT98WGR2N7 | INC-42   | Payments service degraded    | live            | SEV-2         |
| 01FDAG4SAP5TYPT98WGR2N8 | INC-43   | Login failures in EU region  | live            | SEV-1         |
+-------------------------+----------+-------------------------------+-----------------+---------------+
```

## Tables

| Table | Description |
|-------|-------------|
| [incidentio_incident](incidentio/table_incidentio_incident.go) | List and inspect incidents. |
| [incidentio_action](incidentio/table_incidentio_action.go) | List action items tracked against incidents. |
| [incidentio_severity](incidentio/table_incidentio_severity.go) | List severity levels configured in your account. |
| [incidentio_incident_type](incidentio/table_incidentio_incident_type.go) | List incident types configured in your account. |
| [incidentio_followups](incidentio/table_incidentio_followups.go) | List follow-up items tracked against incidents. |
| [incidentio_incident_updates](incidentio/table_incidentio_incident_updates.go) | List status updates posted on incidents. |
| [incidentio_users](incidentio/table_incidentio_users.go) | List users in your incident.io account. |
| [incidentio_alerts](incidentio/table_incidentio_alerts.go) | List alerts ingested from connected alert sources. |
| [incidentio_incident_roles](incidentio/table_incidentio_incident_roles.go) | List incident role definitions (e.g. Incident Lead). |
| [incidentio_incident_statuses](incidentio/table_incidentio_incident_statuses.go) | List incident status definitions. |
| [incidentio_custom_fields](incidentio/table_incidentio_custom_fields.go) | List custom field definitions configured in your account. |
| [incidentio_escalations](incidentio/table_incidentio_escalations.go) | List escalations triggered in your account. |

## Development

### Prerequisites

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

### Build and Install

```sh
make install
```

Configure the plugin:

```sh
cp config/incidentio.spc ~/.steampipe/config/incidentio.spc
vi ~/.steampipe/config/incidentio.spc
```

## Testing

Run a smoke query against every table:

```sh
make test
```

The test script ([scripts/test_tables.sh](scripts/test_tables.sh)) builds the plugin, queries each table, and reports pass/fail/skip (scope-restricted tables are skipped rather than failed).

### Further reading

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [incident.io API reference](https://docs.incident.io/api-reference/introduction)
