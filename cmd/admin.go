package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/orangeseeds/blitzbase/core"
	"github.com/spf13/cobra"
	"golang.org/x/term"
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
		Example:       "blitzbase admin create example@mail.com",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(command *cobra.Command, args []string) error {

			if len(args) != 1 {
				return errors.New("[admin create] command needs an email")
			}

			fmt.Print("Enter password:")
			pass, err := term.ReadPassword(0)
			fmt.Println("")
			if err != nil {
				return err
			}

			app.Store.Connect()

			_, err = app.CreateNewAdmin(args[0], string(pass))
			if err != nil {
				return err
			}

			log.Printf("New admin %s added\n", args[0])

			return nil
		},
	}
	return command
}

func listAdminCommand(app *core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "lists all admins",
		RunE: func(command *cobra.Command, args []string) error {
			app.Store.Connect()
			admins, err := app.ListAdmins()
			if err != nil {
				return err
			}

			for _, admin := range admins {
				fmt.Println(admin.Email)
			}

			return nil
		},
	}
	return command
}
