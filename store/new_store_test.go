package store

import (
	"testing"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	model "github.com/orangeseeds/blitzbase/models"
)

func TestSQliteStore(t *testing.T) {
	db, err := dbx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	store := NewSQliteStore(db)

	col := model.NewCollection(uuid.NewString(), "test_co", model.BASE)
	err = store.SaveCollection(&DBWrapper{store.DB()}, col)

	gotCol, err := store.FindCollectionByNameorId(&DBWrapper{store.DB()}, col.GetID())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", gotCol)

}
