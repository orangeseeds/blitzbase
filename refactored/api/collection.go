package api

import (
	"github.com/labstack/echo/v4"
	"github.com/orangeseeds/blitzbase/refactored/core"
	model "github.com/orangeseeds/blitzbase/refactored/models"
	"github.com/orangeseeds/blitzbase/refactored/store"
	"github.com/orangeseeds/blitzbase/utils"
)

func LoadCollectionAPI(app core.App, e *echo.Echo) {
	api := CollectionAPI{app: app}

	grp := e.Group("/collections", LoadJWT(), NeedsAdminAuth(app))

	grp.GET("", api.index)
	grp.GET("/:collection", api.detail)
	grp.POST("", api.save)
	grp.DELETE("/:collection", api.delete)
}

type CollectionAPI struct {
	app core.App
}

func (a *CollectionAPI) index(c echo.Context) error {
	var colOne model.Collection
	var col []model.Collection
	err := a.app.Store().DB().Select().From(colOne.TableName()).All(&col)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	return c.JSON(200, col)
}

func (a *CollectionAPI) detail(c echo.Context) error {
	name := c.Param("collection")
	col, err := a.app.Store().FindCollectionByNameorId(a.app.Store().DB(), name)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, map[string]any{
		"collection": col,
	})
}

func (a *CollectionAPI) save(c echo.Context) error {
	var col model.Collection

	err := c.Bind(&col)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	if col.IsAuth() {
		for _, v := range model.AuthRecordFields() {

			if !col.Schema.HasField(v) {
				f := model.Field{
					Id:      utils.RandStr(10),
					Name:    v,
					Type:    model.FieldTypeText,
					Options: nil,
				}
				col.Schema.AddField(&f)
			}
		}
	}

	err = a.app.Store().SaveCollection(a.app.Store().DB(), &col)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	if !a.app.Store().(*store.BaseStore).TableExists(col.GetName()) {
		err = a.app.Store().CreateCollectionTable(&col)
		if err != nil {
			return c.JSON(500, err.Error())
		}
	}

	return c.JSON(200, map[string]any{
		"message":    "saved successfully",
		"collection": col,
	})
}

func (a *CollectionAPI) delete(c echo.Context) error {
	var col model.Collection
	col.SetID(c.Param("collection"))
	err := a.app.Store().DeleteCollection(a.app.Store().DB(), &col)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, map[string]any{
		"message":    "deleted successfully",
		"collection": col,
	})
}
