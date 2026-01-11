package drivers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type PostgresDriver struct {
	// database connection
	db *pgx.Conn
}

// connect to the database
func (p *PostgresDriver) Connect(connection string) error {
	db, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	p.db = db
	return nil
}

// close the database connection
func (p *PostgresDriver) Close() error {
	if p.db != nil {
		return p.db.Close(context.Background())
	}
	return nil
}

// runs a migration at the given path
// migration param is the migration folder
// direction is either up / down
func (p *PostgresDriver) RunMigrationFile(ctx context.Context, baseFolder, migration, direction string) error {

	// Read the SQL file
	migrationFile := fmt.Sprintf("%v/%v/%v.sql", baseFolder, migration, direction)

	sqlFile, err := os.ReadFile(migrationFile)
	if err != nil {
		panic(err)
	}

	// split the migration file by semicolons
	statements := strings.Split(string(sqlFile), ";")

	// begin a database transaction
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	// rollback incase any failures
	defer tx.Rollback(context.Background())

	// run each statement
	for _, statement := range statements {
		stmt := strings.TrimSpace(statement)
		if stmt == "" {
			continue
		}

		if _, err := tx.Exec(ctx, stmt); err != nil {
			return err
		}
	}

	// if applying a migration we insert a record
	// if down we need to delete the migration record
	var query string

	if direction == "up" {
		query = "INSERT INTO migrations(name) VALUES($1)"
	} else {
		query = "DELETE FROM migrations WHERE name = $1"
	}

	if _, err := tx.Exec(ctx, query, migration); err != nil {
		return err
	}

	// commit if no failures present
	return tx.Commit(ctx)
}

// finds the last migration applied
func (p *PostgresDriver) LastAppliedMigration(ctx context.Context) (string, bool, error) {

	query := "SELECT name FROM migrations ORDER BY created_at DESC LIMIT 1"

	var lastApplied string
	// execute the query
	err := p.db.QueryRow(ctx, query).Scan(&lastApplied)
	if err != nil {
		// no runs
		if err.Error() == "no rows in result set" {
			return "", true, err
		}

		return "", false, err
	}

	return lastApplied, false, nil
}

// returns a list of database applied migrations
func (p *PostgresDriver) ListAppliedMigrations(ctx context.Context) ([]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := "SELECT name FROM migrations ORDER BY name ASC"

	migrations := make([]string, 0)

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m string
		if err := rows.Scan(&m); err != nil {
			return nil, err
		}

		migrations = append(migrations, m)
	}

	return migrations, nil
}

// Creates the migration table to keep track of applied migrations in the database
func (p *PostgresDriver) CreateMigrationsTable(ctx context.Context) error {

	query := `
		CREATE TABLE IF NOT EXISTS migrations
		(
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT now(),
			updated_at TIMESTAMP NOT NULL DEFAULT now()
		);
	`

	_, err := p.db.Exec(ctx, query)
	return err
}

func (p *PostgresDriver) GetRollbackMigrations(ctx context.Context, steps int) ([]string, error) {

	migrations := make([]string, 0)

	// get the list of migrations applied from the database
	// order reversed so we can rollback by number of steps
	query := "SELECT name FROM migrations ORDER BY name DESC LIMIT $1"

	rows, err := p.db.Query(ctx, query, steps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m string
		if err := rows.Scan(&m); err != nil {
			return nil, err
		}

		migrations = append(migrations, m)
	}

	return migrations, nil
}
