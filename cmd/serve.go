package main

import (
	"github.com/orangeseeds/rtdatabase/api"
	"github.com/orangeseeds/rtdatabase/core"
	"github.com/spf13/cobra"
)

func NewServerCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "serve",
		Short: "starts the server defaults to (128.0.0.1:3300)",
		Run: func(command *cobra.Command, args []string) {
            publisher := core.NewPublisher()
			store := core.NewStorage(publisher)
			app := core.NewApp(store, publisher)

            api.Serve(*app, ":3300")
		},
	}
	return command
}
