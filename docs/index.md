# incident.io Plugin for Steampipe

Use SQL to query incidents, alerts, escalations, users, and more from [incident.io](https://incident.io).

## Installation

Clone and build the plugin:

```sh
git clone https://github.com/vthiery/steampipe-plugin-incidentio.git
cd steampipe-plugin-incidentio
mkdir -p ~/.steampipe/plugins/local/incidentio
go build -o ~/.steampipe/plugins/local/incidentio/steampipe-plugin-incidentio.plugin .
```

## Configuration

Copy the sample config:

```sh
cp config/incidentio.spc ~/.steampipe/config/incidentio.spc
```

Edit `~/.steampipe/config/incidentio.spc`:

```hcl
connection "incidentio" {
  plugin  = "local/incidentio"

  # API key from your incident.io dashboard → Settings → API keys.
  # Required scopes depend on the tables you query.
  api_key = "YOUR_API_KEY"
}
```

## Tables

| Table | Description |
|-------|-------------|
| [incidentio_incident](tables/incidentio_incident.md) | List and inspect incidents. |
| [incidentio_action](tables/incidentio_action.md) | List action items tracked against incidents. |
| [incidentio_severity](tables/incidentio_severity.md) | List severity levels configured in your account. |
| [incidentio_incident_type](tables/incidentio_incident_type.md) | List incident types configured in your account. |
| [incidentio_followups](tables/incidentio_followups.md) | List follow-up items tracked against incidents. |
| [incidentio_incident_updates](tables/incidentio_incident_updates.md) | List status updates posted on incidents. |
| [incidentio_users](tables/incidentio_users.md) | List users in your incident.io account. |
| [incidentio_alerts](tables/incidentio_alerts.md) | List alerts ingested from connected alert sources. |
| [incidentio_incident_roles](tables/incidentio_incident_roles.md) | List incident role definitions. |
| [incidentio_incident_statuses](tables/incidentio_incident_statuses.md) | List incident status definitions. |
| [incidentio_custom_fields](tables/incidentio_custom_fields.md) | List custom field definitions. |
| [incidentio_escalations](tables/incidentio_escalations.md) | List escalations triggered in your account. |
