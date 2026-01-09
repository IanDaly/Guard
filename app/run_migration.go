package app

import (
	"context"
	"fmt"
	"time"

	"github.com/iandaly/migrator/filesystem"
)

func (a *App) handleRunMigrations() error {

	fmt.Println("Running migrations....")

	// get last migration that was applied to the database
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// last applied migration in the database
	lastApplied, isFresh, err := a.Driver.LastAppliedMigration(ctx)

	// if no migrations applied
	if isFresh {
		folderContents, err := filesystem.GetMigrationsFolderContents(a.Config.Folder)
		if err != nil {
			return err
		}

		for _, file := range folderContents {
			a.runMigration(file)
		}

		return nil
	}

	if err != nil {
		return err
	}

	pending, err := filesystem.GetPendingMigrations(lastApplied, a.Config.Folder)
	if err != nil {
		return err
	}

	for _, file := range pending {
		a.runMigration(file)
	}

	return nil
}

func (a *App) runMigration(file string) error {

}
