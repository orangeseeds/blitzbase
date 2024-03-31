package main

import (
	"log"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	model "github.com/orangeseeds/blitzbase/models"
	"github.com/orangeseeds/blitzbase/store"
	"github.com/orangeseeds/blitzbase/utils"
)

func main() {

	dbPath := "./data.db"
	db, err := dbx.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	st := store.NewSQliteStore(db)

	exec := store.Wrap(db)
	st.CreateAdminTable(exec)
	st.CreateCollectionMetaTable(exec)

	col := model.NewCollection(uuid.NewString(), "test_collection", model.BASE)
	col.Schema.AddField(&model.Field{"1", "FieldName", model.FieldTypeText, nil})
	col.Schema.AddField(&model.Field{"2", "FieldName2", model.FieldTypeText, nil})

	err = st.SaveCollection(st.DB(), col)
	if err != nil {
		log.Println(err)
	}
	err = st.CreateCollectionTable(exec, col)
	if err != nil {
		log.Println(err)
	}

	admin := model.Admin{
		BaseModel: model.BaseModel{
			Id:        uuid.NewString(),
			CreatedAt: utils.NowDateTime(),
			UpdatedAt: utils.NowDateTime(),
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
