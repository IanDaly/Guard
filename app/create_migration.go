package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/iandaly/migrator/filesystem"
)

// creates a new migration inside the migrations folder
func (a *App) handleCreateMigration() error {
	if err := ensureAllArgs(); err != nil {
		return err
	}

	// check if the migrations folder exists,
	// if it does not we need to create it
	if filesystem.MigrationsFolderExists(a.Config.Folder) == false {
		os.Mkdir(a.Config.Folder, os.ModePerm)
	}

	// get the folder name
	filename := filesystem.CreateMigrationFolderName()

	// filepath with migrations folder
	folderPath := fmt.Sprintf("%v/%v", a.Config.Folder, filename)

	// create the migrations folder
	if err := os.Mkdir(folderPath, os.ModePerm); err != nil {
		return errors.New("Error creating migration folder")
	}

	// inside the migration folder we need to create the up / down sql files
	if err := filesystem.CreateSqlFile(folderPath, "up"); err != nil {
		return err
	}

	if err := filesystem.CreateSqlFile(folderPath, "down"); err != nil {
		return err
	}

	return nil
}

// ensures all args are passed
func ensureAllArgs() error {
	if len(os.Args) < 3 {
		return errors.New("Not enough arguments supplied")
	}
	return nil
}
