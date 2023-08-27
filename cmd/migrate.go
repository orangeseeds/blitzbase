package main

import (
	"fmt"
	"log"

	"github.com/orangeseeds/blitzbase/core"
	"github.com/orangeseeds/blitzbase/utils/migrations"
	"github.com/spf13/cobra"
)

func NewMigrateCommand(app *core.App) *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "run migrations"}

	migrateCmd.AddCommand(migrateUpCommand(app))
	migrateCmd.AddCommand(migrateDownCommand(app))
	migrateCmd.AddCommand(migrateNewCommand(app))
	return migrateCmd
}

func migrateUpCommand(app *core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "up",
		Short: "migrate up",
		RunE: func(command *cobra.Command, args []string) error {
			app.Store.Connect()

			err := migrations.MigrateUp(app.Store)
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
			err := migrations.MigrateDown(app.Store)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("ran migrations down")
			return nil
		},
	}
	return command
}

func migrateNewCommand(app *core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "new",
		Short: "migrate new",
		RunE: func(command *cobra.Command, args []string) error {
			app.Store.Connect()
			err := migrations.CreateInitTable(app.Store)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("ran migrations new")
			return nil
		},
	}
	return command
}
