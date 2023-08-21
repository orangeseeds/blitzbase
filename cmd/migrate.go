package main

import (
	"fmt"
	"log"

	"github.com/orangeseeds/blitzbase/core"
	"github.com/spf13/cobra"
)

func NewMigrateCommand(app *core.App) *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "run migrations"}

	migrateCmd.AddCommand(migrateUpCommand(app))
	migrateCmd.AddCommand(migrateDownCommand(app))
	return migrateCmd
}

func migrateUpCommand(app *core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "up",
		Short: "migrate up",
		RunE: func(command *cobra.Command, args []string) error {
			app.Store.Connect()
			err := core.MigrateUp(app.Store.Path, app.Store.DB.DB())
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println("ran migrations up")
			return nil
		},
	}
	return command
}

func migrateDownCommand(app *core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "down",
		Short: "migrate down",
		RunE: func(command *cobra.Command, args []string) error {
			app.Store.Connect()
			err := core.MigrateDown(app.Store.Path, app.Store.DB.DB())
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("ran migrations down")
			return nil
		},
	}
	return command
}
