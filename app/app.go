package app

import (
	"context"
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

// ensures the driver can establish a connection to the database
// and creates the migrations table if not exists
func (a *App) EnsureDriver() error {

	if err := a.Driver.Connect(a.Config.Url); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := a.Driver.CreateMigrationsTable(ctx); err != nil {
		return err
	}

	return nil
}
