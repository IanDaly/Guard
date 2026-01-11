package app

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/iandaly/migrator/filesystem"
)

// creates a new migration inside the migrations folder
func (a *App) CreateMigration() error {
	if err := ensureAllArgs(); err != nil {
		return err
	}

	// check if the migrations folder exists,
	// if it does not we need to create it
	if filesystem.MigrationsFolderExists(a.Config.Folder) == false {
		os.Mkdir(a.Config.Folder, os.ModePerm)
	}

	// we want to get the last arguments to build the migration filename
	// for example the user will enter
	// migrator make:migration create users table
	// so we make the filename create_users_table
	// we also replace any slashes incase a user
	// adds a slash, this will mess up paths
	userInput := strings.ReplaceAll(strings.Join(os.Args[2:], "_"), "/", "_")

	// prefix the filename with the current unix timestamp to ensure uniqueness
	now := strconv.FormatInt(time.Now().Unix(), 10)

	filename := fmt.Sprintf("%v_%v", now, userInput)

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

	fmt.Println("Migration created successfully")

	return nil
}

// ensures all args are passed
func ensureAllArgs() error {
	if len(os.Args) < 3 {
		return errors.New("Not enough arguments supplied")
	}
	return nil
}
