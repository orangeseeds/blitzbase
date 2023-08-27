package migrations

import (
	"fmt"
	"os"
	"testing"

	"github.com/orangeseeds/blitzbase/store"
	"github.com/orangeseeds/blitzbase/utils/schema"
)

func preCond(dir string) (*store.Storage, error) {
	_, err := os.Open(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	p := fmt.Sprintf("%s/migrations", dir)
	s := store.NewStorage(fmt.Sprintf("%s/test.db", dir), p)
	return s, nil
}

func TestCreateNewMigration(t *testing.T) {

	dir := "./test_tmp"
	s, err := preCond(dir)
	if err != nil {
		t.Error(err.Error())
	}
	s.Connect()

	fs := schema.FieldSchema{
		Name:     "test_column",
		Type:     schema.TEXT_FIELD,
		Options:  schema.TextFieldOptions{},
		Required: true,
	}

	c := store.BaseColletion{
		Name:   "test_collection2",
		Type:   "base",
		Schema: []schema.FieldSchema{fs},
	}

	if err := CreateNewTable(s, c); err != nil {
		t.Errorf("Error creating migration: %s", err.Error())
	}

	// err := os.RemoveAll(p)
	// if err != nil {
	// 	t.Error(err.Error())
	// }
	// err = os.Remove("./test.db")
	// if err != nil {
	// 	t.Error(err.Error())
	// }
}

func TestCreateInitTable(t *testing.T) {
	dir := "./test_tmp"
	s, err := preCond(dir)
	if err != nil {
		t.Error(err.Error())
	}
	s.Connect()

	err = CreateInitTable(s)
	if err != nil {
		t.Error(err.Error())
	}
}
