package app

import (
	"context"
	"fmt"
	"time"

	"github.com/iandaly/migrator/filesystem"
)

func (a *App) RunMigrations() error {

	// get last migration that was applied to the database
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	migrations := make([]string, 0)

	// last applied migration in the database
	lastApplied, isFresh, err := a.Driver.LastAppliedMigration(ctx)
	if err != nil {
		// if the error is no migrations run and is a fresh migration run
		if isFresh {
			migrations, err = filesystem.GetMigrationsFolderContents(a.Config.Folder)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// if not a fresh run we get the pending migrations
	if isFresh == false {
		migrations, err = filesystem.GetPendingMigrations(lastApplied, a.Config.Folder)
		if err != nil {
			return err
		}
	}

	// if no migrations
	if len(migrations) == 0 {
		fmt.Println("No migrations to run")
		return nil
	}

	fmt.Println("Running migrations....")

	// run each migration file
	for _, migration := range migrations {
		// allow 15 seconds per migration file running
		context, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		if err := a.Driver.RunMigrationFile(context, a.Config.Folder, migration, "up"); err != nil {
			return err
		}
	}

	return nil
}
