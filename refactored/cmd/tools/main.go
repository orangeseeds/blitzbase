package main

import (
	"log"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/mattn/go-sqlite3"
	model "github.com/orangeseeds/blitzbase/refactored/models"
	"github.com/orangeseeds/blitzbase/refactored/store"
	"github.com/orangeseeds/blitzbase/utils"
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

	col := model.NewCollection(utils.RandStr(10), "test_collection", model.BASE)
	col.Schema.AddField(&model.Field{"1", "FieldName", model.FieldTypeText, nil})
	col.Schema.AddField(&model.Field{"2", "FieldName2", model.FieldTypeText, nil})

	store.SaveCollection(store.DB(), col)
	store.CreateCollectionTable(col)

	admin := model.Admin{
		BaseModel: model.BaseModel{
			Id:        utils.RandStr(10),
			CreatedAt: time.Now().String(),
			UpdatedAt: time.Now().String(),
		},
		Email: "admin@mail.com",
	}
	admin.SetPassword("1234567890")
	admin.RefreshToken()
	err = db.Model(&admin).Insert()
	if err != nil {
		log.Println(err)
	}
}
