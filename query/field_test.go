package query

import (
	"testing"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/mattn/go-sqlite3"
	model "github.com/orangeseeds/blitzbase/models"
	"github.com/orangeseeds/blitzbase/store"
	"github.com/orangeseeds/blitzbase/utils"
)

func setup() (*store.BaseStore, error) {
	db, err := dbx.Open("sqlite3", "./test.db")
	if err != nil {
		return nil, err
	}
	store := store.NewBaseStore(db)
	return store, err

}

func TestCollection(t *testing.T) {
	store, err := setup()
	if err != nil {
		t.Error(err)
	}

	collection := model.NewCollection(utils.RandStr(10), "test_collection", model.BASE)

	fieldSpec := RecordFieldSpecifier{
		Collection:    collection,
		AllowedFields: []string{},
		Request:       nil,
	}

	filter := func(q *dbx.SelectQuery) error {
		expr, err := FilterRule(collection.IndexRule).BuildExpr(&fieldSpec)
		if err != nil {
			return err
		}
		if expr != nil {
			q.AndWhere(expr)
		}
		return nil
	}

	rec, err := store.FindRecordById(store.DB(), "id", collection.Name, filter)
	if err != nil {
		t.Error(err)
	}

	t.Log(rec)
}

func TestCollectionSave(t *testing.T) {
	store, err := setup()
	if err != nil {
		t.Error(err)
	}

	collection, err := store.FindCollectionByNameorId(store.DB(), "test_collection")
	rec := model.NewRecord(collection)
	rec.Load(map[string]any{
		"Id":        utils.RandStr(10),
		"field_one": "naice",
		"field_two": 10,
	})

	fieldSpec := RecordFieldSpecifier{
		Collection:    collection,
		AllowedFields: []string{},
		Request:       nil,
	}

	createFilter := func(q *dbx.SelectQuery) error {
		expr, err := FilterRule(collection.IndexRule).BuildExpr(&fieldSpec)
		if err != nil {
			return err
		}
		if expr != nil {
			q.AndWhere(expr)
		}
		return nil
	}

	err = store.SaveRecord(store.DB(), rec, createFilter)
	if err != nil {
		t.Error(err)
	}
}
