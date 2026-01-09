package filesystem

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// creates an sql file, used above to create the up / down files inside the migration folder
func CreateSqlFile(folderPath, name string) error {
	filepath := fmt.Sprintf("%v/%v.sql", folderPath, name)

	if _, err := os.OpenFile(filepath, os.O_CREATE, os.ModePerm); err != nil {
		return errors.New("Error creating " + name + " sql file")
	}

	return nil
}

// retrieves the folder name from the command line args
// and builds the string with the current unix timestamp
func CreateMigrationFolderName() string {
	// we want to get the last arguments to build the migration filename
	// for example the user will enter
	// migrator make:migration create users table
	// so we make the filename create_users_table
	// we also replace any slashes incase a user
	// adds a slash, this will mess up paths
	filename := strings.ReplaceAll(strings.Join(os.Args[2:], "_"), "/", "_")

	// prefix the filename with the current unix timestamp to ensure uniqueness
	now := strconv.FormatInt(time.Now().Unix(), 10)

	return now + "_" + filename
}

// check the migrations folder exists
func MigrationsFolderExists(folder string) bool {
	_, err := os.Stat(folder)
	return err == nil
}

// returns the contents inside the migrations folder
func GetMigrationsFolderContents(folder string) ([]string, error) {
	migrations := make([]string, 0)

	if MigrationsFolderExists(folder) == false {
		return nil, errors.New("Migrations folder does not exist")
	}

	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		migrations = append(migrations, entry.Name())
	}

	return migrations, nil
}

// returns an array of pending migrations
// last applied migration
// migrations folder
func GetPendingMigrations(lastApplied string, folder string) ([]string, error) {
	// since the migrations start with a unix timestamp we split the name to
	// get the unix timestamp
	// we do the same with each migration name to get migrations created
	// after the last applied

	// migrations are named unix_timestamp_name_of_migration
	lastAppliedUnix, err := getUnixFromString(strings.TrimSpace(strings.Split(lastApplied, "_")[0]))
	if err != nil {
		return nil, err
	}

	folderContents, err := GetMigrationsFolderContents(folder)
	if err != nil {
		return nil, err
	}

	pending := make([]string, 0)

	for _, value := range folderContents {
		// get the unix timestamp to compare
		migrationFileUnix, err := getUnixFromString(strings.TrimSpace(strings.Split(value, "_")[0]))
		if err != nil {
			return nil, err
		}

		// if the migration file is after the last applied
		// we append it to the pending migrations list
		if migrationFileUnix.After(*lastAppliedUnix) {
			pending = append(pending, value)
		}
	}

	return pending, nil
}

// helper method used to extract a time object from the timestamp
func getUnixFromString(unix string) (*time.Time, error) {
	// before converting the string to time,
	// we need to convert the string to int64 before
	unixParsed, err := strconv.ParseInt(unix, 10, 64)
	if err != nil {
		return nil, err
	}

	// Convert Unix timestamp to time.Time
	t := time.Unix(unixParsed, 0)

	return &t, nil
}
