package main

import (
	"errors"

	"github.com/orangeseeds/blitzbase/core"
	"github.com/spf13/cobra"
)

func NewAdminCommand(app *core.App) *cobra.Command {

	adminCmd := &cobra.Command{
		Use:   "admin",
		Short: "admin commands",
	}

	adminCmd.AddCommand(createAdminCommand(app))
	adminCmd.AddCommand(listAdminCommand(app))
	return adminCmd
}

func createAdminCommand(app *core.App) *cobra.Command {
	command := &cobra.Command{
		Use:           "create",
		Short:         "create new admin",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(command *cobra.Command, args []string) error {

			app.Store.Connect()

			if len(args) != 2 {
				return errors.New("command needs an email and a password")
			}

			err := app.CreateNewAdmin(args[0], args[1])
			if err != nil {
				return err
			}

			return nil
		},
	}
	return command
}

func listAdminCommand(app *core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "l",
		Short: "list all admins",
		RunE: func(command *cobra.Command, args []string) error {
			app.Store.Connect()
			err := app.ListAdmins()
			if err != nil {
				return err
			}
			return nil
		},
	}
	return command
}
