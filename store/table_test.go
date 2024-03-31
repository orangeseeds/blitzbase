package store

import (
	"os"
	"testing"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/mattn/go-sqlite3"
	model "github.com/orangeseeds/blitzbase/models"
	"github.com/orangeseeds/blitzbase/utils"
)

func TestCreateCollectionsTable(t *testing.T) {
	os.Remove("./test.db")
	db, err := dbx.Open("sqlite3", "./test.db")
	if err != nil {
		t.Error(err)
	}
	store := BaseStore{db: db}

	err = store.CreateCollectionMetaTable()
	if err != nil {
		t.Error(err)
	}
}

func TestCreateAdminTable(t *testing.T) {
	db, err := dbx.Open("sqlite3", "./test.db")
	if err != nil {
		t.Error(err)
	}

	store := BaseStore{db: db}

	var base model.BaseModel
	base.SetID(utils.RandStr(10))

	err = store.CreateAdminTable()
	if err != nil {
		t.Error(err)
	}
}

func TestCreateCollectionTable(t *testing.T) {
	db, err := dbx.Open("sqlite3", "./test.db")
	if err != nil {
		t.Error(err)
	}
	store := BaseStore{db: db}

	col := model.NewCollection(utils.RandStr(10), "test_collection", model.BASE)
	col.Schema.AddField(&model.Field{
		Id:      utils.RandStr(10),
		Name:    "field_one",
		Type:    model.FieldTypeText,
		Options: nil,
	})
	col.Schema.AddField(&model.Field{
		Id:      utils.RandStr(10),
		Name:    "field_two",
		Type:    model.FieldTypeNumber,
		Options: nil,
	})

	err = store.CreateCollectionTable(col)
	if err != nil {
		t.Error(err)
	}
}
