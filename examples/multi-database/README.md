# Multi-Database Example

This example demonstrates how to use Plexr to manage multiple PostgreSQL databases in a single execution plan.

## Overview

The configuration sets up three separate databases:

1. **Main Database** (`myapp_main`) - Core application data
   - User management
   - User profiles
   - Core business entities

2. **Analytics Database** (`myapp_analytics`) - Event tracking and analytics
   - Event tracking
   - Page views
   - Aggregated summaries

3. **Logs Database** (`myapp_logs`) - Application and audit logging
   - Application logs with different severity levels
   - Audit logs for compliance
   - Log retention policies

## Prerequisites

- PostgreSQL server running locally or accessible via network
- Environment variables set:
  - `DB_HOST` (default: localhost)
  - `DB_PORT` (default: 5432)
  - `DB_USER` (default: postgres)
  - `DB_PASSWORD` (required)
  - `DB_SSLMODE` (default: disable)

## Database Setup

Before running this example, ensure the databases exist:

```bash
# Create databases (run as PostgreSQL superuser)
createdb myapp_main
createdb myapp_analytics
createdb myapp_logs
```

Or using SQL:

```sql
CREATE DATABASE myapp_main;
CREATE DATABASE myapp_analytics;
CREATE DATABASE myapp_logs;
```

## Running the Example

```bash
# Set required environment variables
export DB_PASSWORD="your_password"

# Execute the plan
plexr execute

# Check status
plexr status
```

## Configuration Details

Each database is configured as a separate executor in the plan:

```yaml
executors:
  main_db:
    type: sql
    database:
      driver: postgres
      database: myapp_main
      # ... connection details
      
  analytics_db:
    type: sql
    database:
      driver: postgres
      database: myapp_analytics
      # ... connection details
      
  log_db:
    type: sql
    database:
      driver: postgres
      database: myapp_logs
      # ... connection details
```

## SQL Files Organization

The SQL files are organized by database:

```
sql/
├── main/           # Main database scripts
│   ├── 001_users.sql
│   └── 002_seed_users.sql
├── analytics/      # Analytics database scripts
│   ├── 001_events.sql
│   └── 002_aggregates.sql
└── logs/           # Logs database scripts
    ├── 001_logs.sql
    └── 002_retention_policy.sql
```

## Use Cases

This multi-database setup is useful for:

1. **Separation of Concerns**: Keep different types of data in separate databases
2. **Performance**: Isolate high-volume logging/analytics from main application
3. **Security**: Apply different access controls to different databases
4. **Scaling**: Scale databases independently based on load
5. **Backup Strategy**: Different backup schedules for different data types

## Extending the Example

To add more databases:

1. Add a new executor configuration in `plan.yml`
2. Create a new directory under `sql/` for your database scripts
3. Add steps that reference the new executor

Example:

```yaml
executors:
  cache_db:
    type: sql
    database:
      driver: postgres
      database: myapp_cache
      # ... connection details

steps:
  - name: "Setup cache database"
    executor: cache_db
    sql_file: sql/cache/001_cache_tables.sql
```

## Troubleshooting

If you encounter connection errors:

1. Verify PostgreSQL is running: `pg_isready`
2. Check connection parameters: `psql -h $DB_HOST -p $DB_PORT -U $DB_USER -l`
3. Ensure databases exist: `psql -l`
4. Check PostgreSQL logs for authentication errors
5. Verify network connectivity if using remote database

## Best Practices

1. **Use transactions** for related changes within a database
2. **Order your steps** carefully when there are dependencies
3. **Use environment variables** for sensitive information
4. **Test in development** before running in production
5. **Keep SQL files idempotent** using IF NOT EXISTS clauses