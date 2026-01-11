package drivers

import (
	"context"
	"errors"
)

const (
	DriverPostgres   = "postgres"
	DriverMySQL      = "mysql"
	DriverClickhouse = "clickhouse"
	DriverSQLite     = "sqlite"
	DriverTurso      = "turso"
)

// basic interface to ensure all driver methods are implemented
type Driver interface {
	// connection to the database driver
	Connect(connection string) error

	// close the database connection
	Close() error

	// runs a specific migration
	// we pass the context so migration does not run forever
	RunMigrationFile(ctx context.Context, baseFolder, migration, direction string) error

	// retrieves the last migration applied
	// if boolean is true, it means no migrations applied
	LastAppliedMigration(ctx context.Context) (string, bool, error)

	// gets a list of migrations applied to the database
	ListAppliedMigrations(ctx context.Context) ([]string, error)

	// creates the migration table to keep track of migrations applied
	CreateMigrationsTable(ctx context.Context) error

	// rollback a migration based on steps provided
	RollbackMigrations(ctx context.Context, baseFolder string, steps int) error
}

// returns a new driver
// if no driver is found return an error
func NewDriver(driver string) (Driver, error) {
	drivers := map[string]Driver{
		DriverPostgres: &PostgresDriver{db: nil},
	}

	d, ok := drivers[driver]
	if ok == false {
		return nil, errors.New("Invalid driver")
	}

	return d, nil
}
