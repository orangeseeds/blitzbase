package main

import (
	"log"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/orangeseeds/blitzbase/api"
	"github.com/orangeseeds/blitzbase/core"
	"github.com/orangeseeds/blitzbase/store"
)

func main() {
	port := ":5173"
	log.Printf("Serving on %s", port)
	dbPath := "data.db"
	db, err := dbx.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	st := store.NewSQliteStore(db)

	exec := store.Wrap(db)
	st.CreateAdminTable(exec)
	st.CreateCollectionMetaTable(exec)

	app := core.NewDBApp(core.DBAppConfig{
		DbPath:     dbPath,
		ServerAddr: port,
	}, st)

	go func() {
		app.OnAdminAuth().Add(func(e *core.AdminEvent) error {
			log.Println("Admin logged in", e.Admin)
			return nil
		})
	}()

	e := api.SetupServer(app)
	e.Start(app.Addr())
}
