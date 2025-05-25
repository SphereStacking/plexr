-- Setup foreign data wrapper to access main database
-- Note: This requires postgres_fdw extension and appropriate permissions

-- Create extension if not exists
CREATE EXTENSION IF NOT EXISTS postgres_fdw;

-- Create foreign server for main database
CREATE SERVER IF NOT EXISTS main_db_server
FOREIGN DATA WRAPPER postgres_fdw
OPTIONS (
    host '${MAIN_DB_HOST:-localhost}',
    port '${MAIN_DB_PORT:-5432}',
    dbname '${MAIN_DB_NAME:-myapp}'
);

-- Create user mapping
CREATE USER MAPPING IF NOT EXISTS FOR CURRENT_USER
SERVER main_db_server
OPTIONS (
    user '${MAIN_DB_USER:-postgres}',
    password '${MAIN_DB_PASSWORD}'
);

-- Import user tables from main database
CREATE SCHEMA IF NOT EXISTS main_db;

-- Import foreign tables
IMPORT FOREIGN SCHEMA public
LIMIT TO (users, user_profiles)
FROM SERVER main_db_server
INTO main_db;

-- Create a view joining local analytics with main database users
CREATE OR REPLACE VIEW user_event_details AS
SELECT 
    e.id,
    e.event_type,
    e.user_id,
    u.username,
    u.email,
    e.properties,
    e.created_at
FROM events e
LEFT JOIN main_db.users u ON e.user_id = u.id
ORDER BY e.created_at DESC;