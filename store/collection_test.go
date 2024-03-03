package store

import (
	"testing"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/mattn/go-sqlite3"
	model "github.com/orangeseeds/blitzbase/refactored/models"
)

func TestAllCollection(t *testing.T) {

	db, err := dbx.Open("sqlite3", "./test.db")
	if err != nil {
		t.Error(err)
	}
	var col []model.Collection
    err =db.Select().From("_collection").All(&col)
    if err != nil {
        t.Error(err)
    }
    t.Log(col)


}

func TestFindCollectionByName(t *testing.T) {
	db, err := dbx.Open("sqlite3", "./test.db")
	if err != nil {
		t.Error(err)
	}
	store := NewBaseStore(db)

	col, err := store.FindCollectionByNameorId(store.DB(), "test_collection")
	if err != nil {
		t.Error(err)
	}
	t.Log(col)

}

func TestDeleteCollection(t *testing.T) {
	db, err := dbx.Open("sqlite3", "./test.db")
	if err != nil {
		t.Error(err)
	}
	store := NewBaseStore(db)

	err = store.DeleteCollection(store.DB(), model.NewCollection("f760a436eef7c0db659e", "test_collection", model.BASE))
	if err != nil {
		t.Error(err)
	}
}
