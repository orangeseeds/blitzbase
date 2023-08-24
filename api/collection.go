package api

import (
	"encoding/json"
	"net/http"
	"strings"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/orangeseeds/blitzbase/store"
	"github.com/orangeseeds/blitzbase/utils/migrations"
	"github.com/orangeseeds/blitzbase/utils/schema"
)

func (api *collectionServer) handleCollList(w http.ResponseWriter, r *http.Request) {}

func (api *collectionServer) handleCollCreate(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Name   string
		Type   string
		Schema []struct {
			Name     string
			Type     string
			Required bool
			Options  map[string]any
		}
		CreateRule string
		ListRule   string
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	sch := []schema.FieldSchema{}
	for _, s := range req.Schema {
		switch strings.ToLower(s.Type) {
		case "text":
			sch = append(sch, schema.FieldSchema{
				Name:     s.Name,
				Type:     s.Type,
				Options:  schema.TextFieldOptions{},
				Required: s.Required,
			})
		case "integer":
			{
				sch = append(sch, schema.FieldSchema{
					Name:     s.Name,
					Type:     s.Type,
					Options:  schema.IntFieldOptions{},
					Required: s.Required,
				})

			}
		default:
			continue
		}
	}

	err = migrations.CreateNewTable(api.app.Store, store.BaseColletion{
		Name:   req.Name,
		Type:   req.Type,
		Schema: sch,
	})
	if err != nil {
		return
	}

	schemaJson, _ := json.Marshal(req.Schema)

	q := api.app.Store.DB.Insert("_collections", dbx.Params{
		"name":   req.Name,
		"type":   req.Type,
		"schema": schemaJson,
	})
	_, err = q.Execute()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
