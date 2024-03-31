package store

import (
	"testing"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/mattn/go-sqlite3"
	model "github.com/orangeseeds/blitzbase/models"
	"github.com/orangeseeds/blitzbase/utils"
)

func TestSQliteStore(t *testing.T) {
	db, err := dbx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	store := NewSQliteStore(db)

	col := model.NewCollection(utils.RandStr(10), "test_co", model.BASE)
	err = store.SaveCollection(&DBWrapper{store.DB()}, col)

	gotCol, err := store.FindCollectionByNameorId(&DBWrapper{store.DB()}, col.GetID())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", gotCol)

}
