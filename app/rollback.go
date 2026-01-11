package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

func (a *App) Rollback() error {

	// get the last migration applied
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, isFresh, err := a.Driver.LastAppliedMigration(ctx)
	if err != nil {
		if isFresh {
			fmt.Println("No migrations applied")
			return nil
		}
		return err
	}

	var steps int = 1

	// user has passed a rollback value
	if len(os.Args) == 3 {
		stepsNum := os.Args[2]
		steps, err = strconv.Atoi(stepsNum)
		if err != nil {
			return errors.New("Invalid integer for steps")
		}
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := a.Driver.RollbackMigrations(ctx, a.Config.Folder, steps); err != nil {
		return err
	}

	return nil
}
