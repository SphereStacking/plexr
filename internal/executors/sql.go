package executors

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/mitchellh/mapstructure"
)

// SQLExecutor implements the Executor interface for SQL scripts
type SQLExecutor struct {
	config SQLConfig
	db     *sql.DB
}

// SQLConfig represents the configuration for SQL executor
type SQLConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`
}

// SQLStepConfig represents step-specific configuration
type SQLStepConfig struct {
	Transaction bool   `mapstructure:"transaction"`
	SkipIf      string `mapstructure:"skip_if"`
}

// NewSQLExecutor creates a new SQL executor
func NewSQLExecutor() *SQLExecutor {
	return &SQLExecutor{}
}

// Name returns the name of this executor
func (e *SQLExecutor) Name() string {
	return "sql"
}

// Validate validates the configuration
func (e *SQLExecutor) Validate(config map[string]interface{}) error {
	// Parse configuration using mapstructure
	var sqlConfig SQLConfig
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &sqlConfig,
		WeaklyTypedInput: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(config); err != nil {
		return fmt.Errorf("invalid SQL configuration: %w", err)
	}

	// Validate driver
	if sqlConfig.Driver == "" {
		return fmt.Errorf("driver is required")
	}
	if sqlConfig.Driver != "postgres" {
		return fmt.Errorf("unsupported driver: %s (only 'postgres' is supported)", sqlConfig.Driver)
	}

	// Validate connection parameters
	if sqlConfig.Host == "" {
		return fmt.Errorf("host is required")
	}
	if sqlConfig.Port == 0 {
		sqlConfig.Port = 5432 // Default PostgreSQL port
	}
	if sqlConfig.Database == "" {
		return fmt.Errorf("database is required")
	}
	if sqlConfig.Username == "" {
		return fmt.Errorf("username is required")
	}
	if sqlConfig.SSLMode == "" {
		sqlConfig.SSLMode = "disable"
	}

	e.config = sqlConfig
	return nil
}

// Execute executes the SQL file
func (e *SQLExecutor) Execute(ctx context.Context, file ExecutionFile) (*ExecutionResult, error) {
	start := time.Now()

	// Connect to database if not connected
	if e.db == nil {
		if err := e.connect(); err != nil {
			return &ExecutionResult{
				Success:  false,
				Error:    err,
				Duration: time.Since(start).Milliseconds(),
			}, nil
		}
	}

	// Determine transaction mode
	useTransaction := false
	switch file.TransactionMode {
	case "all":
		useTransaction = true
	case "each":
		// Each statement in its own transaction (default behavior)
		useTransaction = false
	case "none", "":
		// No transaction
		useTransaction = false
	}

	// Read SQL file
	content, err := os.ReadFile(file.Path)
	if err != nil {
		return &ExecutionResult{
			Success:  false,
			Error:    fmt.Errorf("failed to read SQL file: %w", err),
			Duration: time.Since(start).Milliseconds(),
		}, nil
	}

	// Expand environment variables
	sqlContent := os.ExpandEnv(string(content))

	// Execute SQL
	var output string
	var execErr error

	if useTransaction {
		output, execErr = e.executeInTransaction(ctx, sqlContent)
	} else {
		output, execErr = e.executeDirect(ctx, sqlContent)
	}

	if execErr != nil {
		return &ExecutionResult{
			Success:  false,
			Error:    execErr,
			Output:   output,
			Duration: time.Since(start).Milliseconds(),
		}, nil
	}

	return &ExecutionResult{
		Success:  true,
		Output:   output,
		Duration: time.Since(start).Milliseconds(),
	}, nil
}

// connect establishes a database connection
func (e *SQLExecutor) connect() error {
	dsn := e.buildDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	e.db = db
	return nil
}

// buildDSN builds PostgreSQL connection string
func (e *SQLExecutor) buildDSN() string {
	// Expand environment variables in password
	password := os.ExpandEnv(e.config.Password)

	// Build DSN
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		e.config.Host,
		e.config.Port,
		e.config.Username,
		password,
		e.config.Database,
		e.config.SSLMode,
	)

	return dsn
}

// executeDirect executes SQL without transaction
func (e *SQLExecutor) executeDirect(ctx context.Context, sqlContent string) (string, error) {
	// Split SQL statements by semicolon
	statements := splitSQLStatements(sqlContent)

	var results []string
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		result, err := e.db.ExecContext(ctx, stmt)
		if err != nil {
			return strings.Join(results, "\n"), fmt.Errorf("statement %d failed: %w", i+1, err)
		}

		rowsAffected, _ := result.RowsAffected()
		results = append(results, fmt.Sprintf("Statement %d: %d rows affected", i+1, rowsAffected))
	}

	return strings.Join(results, "\n"), nil
}

// executeInTransaction executes SQL within a transaction
func (e *SQLExecutor) executeInTransaction(ctx context.Context, sqlContent string) (string, error) {
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback() // Will be no-op if committed
	}()

	// Split SQL statements
	statements := splitSQLStatements(sqlContent)

	var results []string
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		result, err := tx.ExecContext(ctx, stmt)
		if err != nil {
			return strings.Join(results, "\n"), fmt.Errorf("statement %d failed: %w", i+1, err)
		}

		rowsAffected, _ := result.RowsAffected()
		results = append(results, fmt.Sprintf("Statement %d: %d rows affected", i+1, rowsAffected))
	}

	if err := tx.Commit(); err != nil {
		return strings.Join(results, "\n"), fmt.Errorf("failed to commit transaction: %w", err)
	}

	return strings.Join(results, "\n"), nil
}

// splitSQLStatements splits SQL content into individual statements
func splitSQLStatements(sqlContent string) []string {
	// Simple implementation - split by semicolon
	// TODO: Handle semicolons within strings, comments, etc.
	statements := strings.Split(sqlContent, ";")

	// Filter out empty statements and trim whitespace
	var filtered []string
	for _, stmt := range statements {
		trimmed := strings.TrimSpace(stmt)
		if trimmed != "" {
			filtered = append(filtered, trimmed)
		}
	}

	// Return empty slice instead of nil
	if len(filtered) == 0 {
		return []string{}
	}

	return filtered
}

// Close closes the database connection
func (e *SQLExecutor) Close() error {
	if e.db != nil {
		return e.db.Close()
	}
	return nil
}

// Clone creates a new instance of SQLExecutor with the same configuration
func (e *SQLExecutor) Clone() *SQLExecutor {
	return &SQLExecutor{
		config: e.config,
		db:     nil, // Each instance gets its own connection
	}
}
