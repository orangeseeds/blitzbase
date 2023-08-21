package main

import (
	"fmt"
	"os"

	"github.com/orangeseeds/blitzbase/core"
	"github.com/orangeseeds/blitzbase/store"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "blitzbase",
	Short: "Blitzbase is a realtime database.",
}

func main() {

	dbPath := "./test.db"
	store := store.NewStorage(dbPath)
	app := core.NewApp(store)

	rootCmd.AddCommand(NewMigrateCommand(app))
	rootCmd.AddCommand(NewAdminCommand(app))
	rootCmd.AddCommand(NewServerCommand(app))
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
