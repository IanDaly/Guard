package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/iandaly/migrator/config"
	"github.com/iandaly/migrator/drivers"
)

type App struct {
	Config *config.MigratorConfig
	Driver drivers.Driver
}

func New() *App {
	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	driver, err := drivers.NewDriver(config.Driver)
	if err != nil {
		panic(err)
	}

	return &App{
		Config: config,
		Driver: driver,
	}
}

func (a *App) Start() {

	// close db connection when function complete
	defer a.Driver.Close()

	first := os.Args[1]

	switch first {
	// if we are initializing a project
	case "init":
		if err := a.handleInit(); err != nil {
			panic(err)
		}

	// if we are creating a migration
	case "make:migration":
		if err := a.handleCreateMigration(); err != nil {
			panic(err)
		}

	case "migrate":
		if err := a.EnsureDriver(); err != nil {
			panic(err)
		}

		if err := a.handleRunMigrations(); err != nil {
			panic(err)
		}

	// prints a list of pending migrations
	case "pending":
		if err := a.EnsureDriver(); err != nil {
			panic(err)
		}

		if err := a.ListPendingMigrations(); err != nil {
			panic(err)
		}

	case "rollback":
		if err := a.EnsureDriver(); err != nil {
			panic(err)
		}
	}
}

func (a *App) handleInit() error {
	if config.Exists() {
		return errors.New("Config exists")
	}

	if err := config.Create(); err != nil {
		return err
	}

	fmt.Println("Initialize complete")

	return nil
}

// ensures the driver can establish a connection to the database
// and creates the migrations table if not exists
func (a *App) EnsureDriver() error {

	if err := a.Driver.Connect(a.Config.Url); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := a.Driver.CreateMigrationTable(ctx); err != nil {
		return err
	}

	return nil
}
