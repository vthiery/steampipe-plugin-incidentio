# Table: incidentio_action

List action items tracked against incidents in your incident.io account.

## Example queries

**List all outstanding actions:**

```sql
select id, description, incident_id, created_at
from incidentio_action
where status = 'outstanding';
```

**List outstanding actions for a specific incident:**

```sql
select id, description, created_at
from incidentio_action
where incident_id = '01FDAG4SAP5TYPT98WGR2N7W91'
  and status = 'outstanding';
```
