# SQL Setup Example

This example demonstrates how to use the SQL executor to set up a PostgreSQL database.

## Prerequisites

- PostgreSQL server running and accessible
- Database credentials

## Configuration

Set the following environment variables:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=myapp_dev
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_SSLMODE=disable
```

## Running the Example

```bash
# Validate the plan
plexr validate plan.yml

# Execute the plan
plexr execute plan.yml
```

## What it does

1. **create_schema**: Creates the database tables (users, products, orders)
2. **create_indexes**: Creates indexes for better query performance
3. **seed_data**: Inserts sample data into the tables

## Testing

You can verify the setup by connecting to the database:

```bash
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "\dt"
```

This should show the three tables created by the setup.