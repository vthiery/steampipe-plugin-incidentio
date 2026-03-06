# Table: incidentio_incident_type

List incident types configured in your incident.io account.

## Example queries

**List all incident types:**

```sql
select id, name, description, is_default, create_in_triage
from incidentio_incident_type;
```
