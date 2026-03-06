# Table: incidentio_incident_roles

List incident role definitions configured in your incident.io account.

## Example queries

**List all incident roles:**

```sql
select id, name, role_type, shortform, required
from incidentio_incident_roles
order by role_type, name;
```

**List custom roles only:**

```sql
select id, name, shortform, instructions
from incidentio_incident_roles
where role_type = 'custom';
```
