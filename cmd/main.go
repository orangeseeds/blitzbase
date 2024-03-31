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
	log.Printf("Serving on %s", ":9900")
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
		ServerAddr: ":9900",
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
