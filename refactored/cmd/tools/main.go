package main

import (
	"log"

	dbx "github.com/go-ozzo/ozzo-dbx"
	model "github.com/orangeseeds/blitzbase/refactored/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/orangeseeds/blitzbase/refactored/store"
)

func main() {

	dbPath := "./data.db"
	db, err := dbx.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	store := store.NewBaseStore(db)

	store.CreateAdminTable()
	store.CreateCollectionMetaTable()

	col := model.NewCollection("id", "test_collection", model.BASE)
	col.Schema.AddField(&model.Field{"1", "FieldName", model.Text, nil})
	col.Schema.AddField(&model.Field{"2", "FieldName2", model.Text, nil})

	store.CreateCollectionTable(col)
}
