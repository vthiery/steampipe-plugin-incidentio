# Table: incidentio_alerts

List alerts ingested from connected alert sources in your incident.io account.

## Example queries

**List all currently firing alerts:**

```sql
select id, title, alert_source_id, deduplication_key, created_at
from incidentio_alerts
where status = 'firing'
order by created_at desc;
```

**Count alerts by status:**

```sql
select status, count(*) as total
from incidentio_alerts
group by status;
```
