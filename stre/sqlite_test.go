package stre

import (
	"os"
	"testing"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	model "github.com/orangeseeds/blitzbase/models"
)

var (
	db          *dbx.DB
	sqliteStore *SQliteStore
)

func TestMain(m *testing.M) {
	var err error
	db, err = dbx.Open("sqlite3", ":memory:")
	defer db.Close()
	if err != nil {
		os.Exit(1)
	}
	sqliteStore = NewSQliteStore(db)

	// exec := Wrap(db)
	// sqliteStore.CreateAdminTable(exec)
	// sqliteStore.CreateCollectionMetaTable(exec)

	// col := model.NewCollection(uuid.NewString(), "test_collection", model.BASE)
	// col.Schema.AddField(&model.Field{"1", "FieldName", model.FieldTypeText, nil})
	// col.Schema.AddField(&model.Field{"2", "FieldName2", model.FieldTypeText, nil})

	// err = sqliteStore.SaveCollection(sqliteStore.DB(), col)
	// if err != nil {
	// 	log.Println(err)
	// }
	// err = sqliteStore.CreateCollectionTable(exec, col)
	// if err != nil {
	// 	log.Println(err)
	// }

	code := m.Run()
	os.Exit(code)
}

func TestFns(t *testing.T){
    // db.Select().Where().On
}

func TestCreateCollectionsTable(t *testing.T) {
	exec := Wrap(db)
	err := sqliteStore.CreateCollectionMetaTable(exec)
	if err != nil {
		t.Fatal(err)
	}

	if !sqliteStore.TableExists(exec, "_collection") {
		t.Fatal("_collection table not created")
	}
}

func TestCreateAdminTable(t *testing.T) {
	var base model.BaseModel
	base.SetID(uuid.NewString())

	exec := Wrap(db)
	err := sqliteStore.CreateAdminTable(exec)
	if err != nil {
		t.Fatal(err)
	}
	if !sqliteStore.TableExists(exec, "_collection") {
		t.Fatal("_collection table not created")
	}
}

func TestCreateCollectionTable(t *testing.T) {
	col := model.NewCollection(uuid.NewString(), "test_collection", model.BASE)
	col.Schema.AddField(&model.Field{"1", "FieldName", model.FieldTypeText, nil})
	col.Schema.AddField(&model.Field{"2", "FieldName2", model.FieldTypeText, nil})

	exec := Wrap(db)
	err := sqliteStore.SaveCollection(sqliteStore.DB(), col)
	if err != nil {
		t.Fatal(err)
	}
	err = sqliteStore.CreateCollectionTable(exec, col)
	if err != nil {
		t.Fatal(err)
	}

	if !sqliteStore.TableExists(exec, col.GetName()) {
		t.Fatalf("%s table not created", col.GetName())
	}
}

func TestAllCollection(t *testing.T) {
	var col []model.Collection
	err := sqliteStore.db.Select().From("_collection").All(&col)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindCollectionByName(t *testing.T) {
	_, err := sqliteStore.FindCollectionByNameorId(Wrap(db), "test_collection")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteCollection(t *testing.T) {
	err := sqliteStore.DeleteCollection(Wrap(db), model.NewCollection(uuid.NewString(), "test_collection", model.BASE))
	if err != nil {
		t.Fatal(err)
	}
}
