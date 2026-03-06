# Table: incidentio_escalations

List escalations triggered in your incident.io account.

## Example queries

**List recent escalations:**

```sql
select id, title, status, priority_name, created_at
from incidentio_escalations
order by created_at desc
limit 20;
```

**Count escalations by status:**

```sql
select status, count(*) as total
from incidentio_escalations
group by status
order by total desc;
```
