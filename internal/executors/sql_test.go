package executors

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLExecutor(t *testing.T) {
	t.Run("NewSQLExecutor", func(t *testing.T) {
		executor := NewSQLExecutor()
		assert.NotNil(t, executor)
		assert.Equal(t, "sql", executor.Name())
	})

	t.Run("Validate configuration", func(t *testing.T) {
		tests := []struct {
			name    string
			config  map[string]interface{}
			wantErr bool
			errMsg  string
		}{
			{
				name: "valid configuration",
				config: map[string]interface{}{
					"driver":   "postgres",
					"host":     "localhost",
					"port":     5432,
					"database": "testdb",
					"username": "testuser",
					"password": "testpass",
					"sslmode":  "disable",
				},
				wantErr: false,
			},
			{
				name: "missing driver",
				config: map[string]interface{}{
					"host":     "localhost",
					"database": "testdb",
					"username": "testuser",
				},
				wantErr: true,
				errMsg:  "driver is required",
			},
			{
				name: "unsupported driver",
				config: map[string]interface{}{
					"driver":   "mysql",
					"host":     "localhost",
					"database": "testdb",
					"username": "testuser",
				},
				wantErr: true,
				errMsg:  "unsupported driver: mysql",
			},
			{
				name: "missing host",
				config: map[string]interface{}{
					"driver":   "postgres",
					"database": "testdb",
					"username": "testuser",
				},
				wantErr: true,
				errMsg:  "host is required",
			},
			{
				name: "missing database",
				config: map[string]interface{}{
					"driver":   "postgres",
					"host":     "localhost",
					"username": "testuser",
				},
				wantErr: true,
				errMsg:  "database is required",
			},
			{
				name: "missing username",
				config: map[string]interface{}{
					"driver":   "postgres",
					"host":     "localhost",
					"database": "testdb",
				},
				wantErr: true,
				errMsg:  "username is required",
			},
			{
				name: "default port and sslmode",
				config: map[string]interface{}{
					"driver":   "postgres",
					"host":     "localhost",
					"database": "testdb",
					"username": "testuser",
					"password": "testpass",
				},
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				executor := NewSQLExecutor()
				err := executor.Validate(tt.config)
				if tt.wantErr {
					assert.Error(t, err)
					if tt.errMsg != "" {
						assert.Contains(t, err.Error(), tt.errMsg)
					}
				} else {
					assert.NoError(t, err)
					if tt.config["port"] == nil {
						assert.Equal(t, 5432, executor.config.Port)
					}
					if tt.config["sslmode"] == nil {
						assert.Equal(t, "disable", executor.config.SSLMode)
					}
				}
			})
		}
	})

	t.Run("buildDSN", func(t *testing.T) {
		executor := &SQLExecutor{
			config: SQLConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "testuser",
				Password: "testpass",
				Database: "testdb",
				SSLMode:  "disable",
			},
		}

		dsn := executor.buildDSN()
		expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
		assert.Equal(t, expected, dsn)
	})

	t.Run("buildDSN with environment variable", func(t *testing.T) {
		os.Setenv("TEST_DB_PASSWORD", "secret123")
		defer os.Unsetenv("TEST_DB_PASSWORD")

		executor := &SQLExecutor{
			config: SQLConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "testuser",
				Password: "${TEST_DB_PASSWORD}",
				Database: "testdb",
				SSLMode:  "require",
			},
		}

		dsn := executor.buildDSN()
		expected := "host=localhost port=5432 user=testuser password=secret123 dbname=testdb sslmode=require"
		assert.Equal(t, expected, dsn)
	})

	t.Run("splitSQLStatements", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected []string
		}{
			{
				name:  "single statement",
				input: "CREATE TABLE users (id INT PRIMARY KEY)",
				expected: []string{
					"CREATE TABLE users (id INT PRIMARY KEY)",
				},
			},
			{
				name:  "multiple statements",
				input: "CREATE TABLE users (id INT PRIMARY KEY); INSERT INTO users VALUES (1); SELECT * FROM users;",
				expected: []string{
					"CREATE TABLE users (id INT PRIMARY KEY)",
					"INSERT INTO users VALUES (1)",
					"SELECT * FROM users",
				},
			},
			{
				name:  "statements with newlines",
				input: "CREATE TABLE users (\n  id INT PRIMARY KEY\n);\nINSERT INTO users VALUES (1);",
				expected: []string{
					"CREATE TABLE users (\n  id INT PRIMARY KEY\n)",
					"INSERT INTO users VALUES (1)",
				},
			},
			{
				name:     "empty input",
				input:    "",
				expected: []string{},
			},
			{
				name:     "only semicolons",
				input:    ";;;",
				expected: []string{},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := splitSQLStatements(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestSQLExecutorWithMock(t *testing.T) {
	t.Run("Execute with successful query", func(t *testing.T) {
		// Create mock database
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// Create executor with mock db
		executor := &SQLExecutor{
			config: SQLConfig{
				Driver:   "postgres",
				Host:     "localhost",
				Port:     5432,
				Database: "testdb",
				Username: "testuser",
			},
			db: db,
		}

		// Create test SQL file
		sqlContent := "INSERT INTO users (name) VALUES ('test')"
		tmpFile, err := os.CreateTemp("", "test*.sql")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(sqlContent)
		require.NoError(t, err)
		tmpFile.Close()

		// Set up mock expectations
		mock.ExpectExec("INSERT INTO users").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Execute
		result, err := executor.Execute(context.Background(), ExecutionFile{
			Path: tmpFile.Name(),
		})

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Contains(t, result.Output, "1 rows affected")
		assert.GreaterOrEqual(t, result.Duration, int64(0))

		// Ensure all expectations were met
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Execute with transaction", func(t *testing.T) {
		// Create mock database
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// Create executor with mock db
		executor := &SQLExecutor{
			config: SQLConfig{
				Driver:   "postgres",
				Host:     "localhost",
				Port:     5432,
				Database: "testdb",
				Username: "testuser",
			},
			db: db,
		}

		// Create test SQL file with multiple statements
		sqlContent := "INSERT INTO users (name) VALUES ('test1'); INSERT INTO users (name) VALUES ('test2')"
		tmpFile, err := os.CreateTemp("", "test*.sql")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(sqlContent)
		require.NoError(t, err)
		tmpFile.Close()

		// TODO: Add transaction support to Execute method
		// For now, test without transaction
		mock.ExpectExec("INSERT INTO users").
			WithArgs().
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO users").
			WithArgs().
			WillReturnResult(sqlmock.NewResult(2, 1))

		// Execute
		result, err := executor.Execute(context.Background(), ExecutionFile{
			Path: tmpFile.Name(),
		})

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Contains(t, result.Output, "Statement 1: 1 rows affected")
		assert.Contains(t, result.Output, "Statement 2: 1 rows affected")

		// Ensure all expectations were met
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Execute with SQL error", func(t *testing.T) {
		// Create mock database
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// Create executor with mock db
		executor := &SQLExecutor{
			config: SQLConfig{
				Driver:   "postgres",
				Host:     "localhost",
				Port:     5432,
				Database: "testdb",
				Username: "testuser",
			},
			db: db,
		}

		// Create test SQL file
		sqlContent := "INSERT INTO nonexistent_table (name) VALUES ('test')"
		tmpFile, err := os.CreateTemp("", "test*.sql")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(sqlContent)
		require.NoError(t, err)
		tmpFile.Close()

		// Set up mock expectations
		mock.ExpectExec("INSERT INTO nonexistent_table").
			WillReturnError(fmt.Errorf("table does not exist"))

		// Execute
		result, err := executor.Execute(context.Background(), ExecutionFile{
			Path: tmpFile.Name(),
		})

		// Verify
		assert.NoError(t, err) // Execute returns error in result, not as error
		assert.NotNil(t, result)
		assert.False(t, result.Success)
		assert.NotNil(t, result.Error)
		assert.Contains(t, result.Error.Error(), "table does not exist")

		// Ensure all expectations were met
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Execute with nonexistent file", func(t *testing.T) {
		// Create executor with valid config to avoid connection error
		executor := &SQLExecutor{
			config: SQLConfig{
				Driver:   "postgres",
				Host:     "localhost",
				Port:     5432,
				Database: "testdb",
				Username: "testuser",
				Password: "testpass",
				SSLMode:  "disable",
			},
		}

		// Execute with nonexistent file
		result, err := executor.Execute(context.Background(), ExecutionFile{
			Path: "/nonexistent/file.sql",
		})

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.Success)
		assert.NotNil(t, result.Error)
		// Should fail to connect or read file
		assert.True(t,
			strings.Contains(result.Error.Error(), "failed to read SQL file") ||
				strings.Contains(result.Error.Error(), "failed to connect to database"),
			"Error should mention file read or connection failure, got: %s", result.Error.Error())
	})
}

func TestSQLExecutorTransactions(t *testing.T) {
	t.Run("executeInTransaction success", func(t *testing.T) {
		// Create mock database
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		executor := &SQLExecutor{db: db}

		// Set up expectations
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("UPDATE users").
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// Execute
		sql := "INSERT INTO users (name) VALUES ('test'); UPDATE users SET active = true WHERE id = 1"
		output, err := executor.executeInTransaction(context.Background(), sql)

		// Verify
		assert.NoError(t, err)
		assert.Contains(t, output, "Statement 1: 1 rows affected")
		assert.Contains(t, output, "Statement 2: 1 rows affected")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("executeInTransaction rollback on error", func(t *testing.T) {
		// Create mock database
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		executor := &SQLExecutor{db: db}

		// Set up expectations
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO invalid_table").
			WillReturnError(fmt.Errorf("table does not exist"))
		mock.ExpectRollback()

		// Execute
		sql := "INSERT INTO users (name) VALUES ('test'); INSERT INTO invalid_table VALUES (1)"
		output, err := executor.executeInTransaction(context.Background(), sql)

		// Verify
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "statement 2 failed")
		assert.Contains(t, output, "Statement 1: 1 rows affected")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

// Helper function to test database connection
func TestSQLExecutorConnection(t *testing.T) {
	t.Run("connect with invalid credentials", func(t *testing.T) {
		executor := &SQLExecutor{
			config: SQLConfig{
				Driver:   "postgres",
				Host:     "nonexistent-host",
				Port:     5432,
				Database: "testdb",
				Username: "testuser",
				Password: "testpass",
				SSLMode:  "disable",
			},
		}

		err := executor.connect()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect to database")
	})
}
