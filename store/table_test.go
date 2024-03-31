package store

import (
	"os"
	"testing"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	model "github.com/orangeseeds/blitzbase/models"
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
	base.SetID(uuid.NewString())

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

	col := model.NewCollection(uuid.NewString(), "test_collection", model.BASE)
	col.Schema.AddField(&model.Field{
		Id:      uuid.NewString(),
		Name:    "field_one",
		Type:    model.FieldTypeText,
		Options: nil,
	})
	col.Schema.AddField(&model.Field{
		Id:      uuid.NewString(),
		Name:    "field_two",
		Type:    model.FieldTypeNumber,
		Options: nil,
	})

	err = store.CreateCollectionTable(col)
	if err != nil {
		t.Error(err)
	}
}
