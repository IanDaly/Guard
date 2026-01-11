package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

const (
	// The name of the migration config file
	ConfigFilename = "migrator.yaml"
)

type MigratorConfig struct {
	// database url
	Url string `yaml:"url"`

	// path to migrations folder
	Folder string `yaml:"folder"`

	// database driver e.g postgres, mysql
	Driver string `yaml:"driver"`
}

// checks if the config file exists
func Exists() (bool, error) {
	path, err := os.Getwd()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(fmt.Sprintf("%v/%v", path, ConfigFilename))
	return err == nil, nil
}

// creates a config file
func Create() error {
	// empty config file body
	emptyConfig := []byte("---\nurl:\nfolder:\ndriver:")

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	// write the migration
	file := fmt.Sprintf("%v/%v", path, ConfigFilename)

	if err := os.WriteFile(file, emptyConfig, 0755); err != nil {
		return err
	}

	return nil
}

// loads the config file
func Load() (*MigratorConfig, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(fmt.Sprintf("%v/%v", path, ConfigFilename))
	if err != nil {
		return nil, err
	}

	config := new(MigratorConfig)

	// read the file
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
