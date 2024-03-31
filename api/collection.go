package api

import (
	"log"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/orangeseeds/blitzbase/core"
	model "github.com/orangeseeds/blitzbase/models"
	"github.com/orangeseeds/blitzbase/request"
	"github.com/orangeseeds/blitzbase/store"
)

type CollectionAPI struct {
	app core.App
}

func (a *CollectionAPI) index(c echo.Context) error {
	var colOne model.Collection
	var col []model.Collection
	err := a.app.Store().DB().Select().From(colOne.TableName()).All(&col)
	if err != nil {
		return NewApiError(500, "", err)
	}

	a.app.OnCollectionIndex().Trigger(&core.CollectionEvent{
		Type:    core.CreateEvent,
		Request: &c,
	})

	return c.JSON(200, col)
}

func (a *CollectionAPI) detail(c echo.Context) error {
	name := c.Param("collection")
	exec := store.Wrap(a.app.Store().DB())
	col, err := a.app.Store().FindCollectionByNameorId(exec, name)
	if err != nil {
		return NewNotFoundError("", err)
	}

	a.app.OnCollectionDetail().Trigger(&core.CollectionEvent{
		Type:       core.DetailEvent,
		Collection: col,
		Request:    &c,
	})

	return c.JSON(200, map[string]any{
		"collection": col,
	})
}

func (a *CollectionAPI) save(c echo.Context) error {
	req, err := request.JsonValidate[model.Collection, request.CollectionSaveRequest](c)
	if err != nil {
		return NewBadRequestError("", err)
	}
	col := req.Model()

	if col.IsAuth() {
		for _, v := range model.AuthRecordFields() {
			if !col.Schema.HasField(v) {
				log.Println(v)
				f := model.Field{
					Id:      uuid.NewString(),
					Name:    v,
					Type:    model.FieldTypeText,
					Options: nil,
				}
				col.Schema.AddField(&f)
			}
		}
	}
	col.SetID(uuid.NewString())
	exec := store.Wrap(a.app.Store().DB())
	err = a.app.Store().SaveCollection(exec, &col)
	if err != nil {
		return NewBadRequestError("Failed to save collection", err)
	}

	if !a.app.Store().TableExists(exec, col.GetName()) {
		err = a.app.Store().CreateCollectionTable(exec, &col)
		if err != nil {
			return NewBadRequestError("Failed to create collection table", err)
		}
	}

	a.app.OnCollectionDetail().Trigger(&core.CollectionEvent{
		Type:       core.CreateEvent,
		Collection: &col,
		Request:    &c,
	})

	return c.JSON(200, map[string]any{
		"message":    "saved successfully",
		"collection": col,
	})
}

func (a *CollectionAPI) delete(c echo.Context) error {
	var col model.Collection
	col.SetID(c.Param("collection"))
	exec := store.Wrap(a.app.Store().DB())
	err := a.app.Store().DeleteCollection(exec, &col)
	if err != nil {
		return NewBadRequestError("", map[string]any{"data": err.Error()})
	}

	a.app.OnCollectionDelete().Trigger(&core.CollectionEvent{
		Type:       core.DeleteEvent,
		Collection: &col,
		Request:    &c,
	})

	return c.JSON(200, map[string]any{
		"message":    "deleted successfully",
		"collection": col,
	})
}
