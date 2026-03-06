# Table: incidentio_users

List users in your incident.io account.

## Example queries

**List all users with their base role:**

```sql
select id, name, email, role, base_role_name
from incidentio_users
order by name;
```

**Find a user by email:**

```sql
select id, name, email, base_role_name
from incidentio_users
where email = 'alice@example.com';
```
