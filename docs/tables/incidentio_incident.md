# Table: incidentio_incident

List and inspect incidents in your incident.io account.

## Example queries

**List all live incidents:**

```sql
select id, reference, name, status_category, severity_name
from incidentio_incident
where status_category = 'live';
```

**Top 10 most recently created incidents:**

```sql
select id, reference, name, status_name, severity_name, created_at
from incidentio_incident
order by created_at desc
limit 10;
```

**Count incidents by severity:**

```sql
select severity_name, count(*) as total
from incidentio_incident
group by severity_name
order by total desc;
```

**Get a specific incident by ID:**

```sql
select id, reference, name, status_name, severity_name, permalink
from incidentio_incident
where id = '01FDAG4SAP5TYPT98WGR2N7W91';
```
