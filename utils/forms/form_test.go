package forms

import (
	"log"
	"testing"

	"github.com/orangeseeds/blitzbase/store"
)

var s *store.Storage

func init() {
	s = store.NewStorage("../../test.db", "../../migrations")
	s.Connect()

	err := s.DB.DB().Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetCollectionData(t *testing.T) {

	q := s.DB.Select("*").From("_collections")
	data, err := q.Rows()
	check(t, err)

	type collections struct {
		Id     int
		Name   string
		Schema string
		Ctype  string
	}
	var (
		col  collections
		cols []collections
	)

	for data.Next() {
		err := data.ScanStruct(&col)
		cols = append(cols, col)

		check(t, err)
	}
    t.Log(cols)
}
