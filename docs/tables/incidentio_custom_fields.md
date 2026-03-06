# Table: incidentio_custom_fields

List custom field definitions configured in your incident.io account.

## Example queries

**List all custom fields by type:**

```sql
select id, name, field_type, show_before_creation, show_before_closure
from incidentio_custom_fields
order by field_type, name;
```

**List required custom fields:**

```sql
select id, name, field_type, required_v2
from incidentio_custom_fields
where required_v2 = 'always';
```
