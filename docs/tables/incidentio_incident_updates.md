# Table: incidentio_incident_updates

List status updates posted on incidents in your incident.io account.

## Example queries

**List all updates for a specific incident:**

```sql
select id, created_at, message, new_status_name, new_severity_name
from incidentio_incident_updates
where incident_id = '01FDAG4SAP5TYPT98WGR2N7W91'
order by created_at desc;
```

**List the most recent updates across all incidents:**

```sql
select id, incident_id, created_at, new_status_name, new_severity_name
from incidentio_incident_updates
order by created_at desc
limit 20;
```
