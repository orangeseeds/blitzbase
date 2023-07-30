package main

import (
	"github.com/orangeseeds/blitzbase/api"
	"github.com/orangeseeds/blitzbase/core"
	"github.com/spf13/cobra"
)

func NewServerCommand(app *core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "serve",
		Short: "starts the server defaults to [:3300]",
		Run: func(command *cobra.Command, args []string) {
			app.Store.Connect()
			api.Serve(*app, ":3300")
		},
	}
	return command
}
