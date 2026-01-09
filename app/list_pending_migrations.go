package app

import (
	"context"
	"fmt"
	"time"

	"github.com/iandaly/migrator/filesystem"
)

func (a *App) ListPendingMigrations() error {
	// get last migration that was applied to the database
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// last applied migration in the database
	lastApplied, isFresh, err := a.Driver.LastAppliedMigration(ctx)

	// if no migrations are applied
	if isFresh {
		folderContents, err := filesystem.GetMigrationsFolderContents(a.Config.Folder)
		if err != nil {
			return err
		}

		printResult(folderContents)
		return nil
	}

	if err != nil {
		return err
	}

	// if a migration has been applied before we need to get the list of items created
	// after the last applied migration
	pending, err := filesystem.GetPendingMigrations(lastApplied, a.Config.Folder)
	if err != nil {
		return err
	}

	if len(pending) == 0 {
		fmt.Println("No migrations pending")
		return nil
	}

	printResult(pending)
	return nil
}

// prints the pending result in bullet point list
func printResult(migrations []string) {
	fmt.Println("***************** Pending migrations *****************")

	for _, item := range migrations {
		fmt.Printf("• %v\n", item)
	}

	fmt.Println("******************************************************")
}
