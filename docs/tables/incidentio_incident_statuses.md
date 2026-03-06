# Table: incidentio_incident_statuses

List incident status definitions configured in your incident.io account.

## Example queries

**List all statuses ordered by category and rank:**

```sql
select id, name, category, rank
from incidentio_incident_statuses
order by category, rank;
```

**List only live-phase statuses:**

```sql
select id, name, rank
from incidentio_incident_statuses
where category = 'live'
order by rank;
```
