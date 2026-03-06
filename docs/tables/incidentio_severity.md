# Table: incidentio_severity

List severity levels configured in your incident.io account.

## Example queries

**List all severity levels ordered by rank:**

```sql
select id, name, description, rank
from incidentio_severity
order by rank;
```
