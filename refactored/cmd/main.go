package main

import (
	"log"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/orangeseeds/blitzbase/refactored/api"
	"github.com/orangeseeds/blitzbase/refactored/core"
	"github.com/orangeseeds/blitzbase/refactored/store"
)

func main() {
	dbPath := "data.db"
	db, err := dbx.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	store := store.NewBaseStore(db)

	store.CreateAdminTable()
	store.CreateCollectionMetaTable()

	app := core.NewDBApp(core.DBAppConfig{
		DbPath:     dbPath,
		ServerAddr: ":9900",
	}, store)

	go func() {
		app.OnAdminAuth().Add(func(e *core.AdminEvent) error {
			log.Println("Admin logged in", e.Admin)
			return nil
		})
	}()

	e := api.SetupServer(app)
	e.Start(app.Addr())
}
