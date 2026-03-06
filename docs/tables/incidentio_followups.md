# Table: incidentio_followups

List follow-up items tracked against incidents in your incident.io account.

## Example queries

**List all outstanding follow-ups:**

```sql
select id, title, status, priority_name, incident_id
from incidentio_followups
where status = 'outstanding';
```

**Count follow-ups by priority:**

```sql
select priority_name, count(*) as total
from incidentio_followups
group by priority_name
order by total desc;
```
