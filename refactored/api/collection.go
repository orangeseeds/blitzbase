package api

import (
	"github.com/labstack/echo"
	"github.com/orangeseeds/blitzbase/refactored/core"
	model "github.com/orangeseeds/blitzbase/refactored/models"
)

func LoadCollectionAPI(app core.App, e *echo.Echo) {
	api := CollectionAPI{app: app}

	grp := e.Group("/collections")

	grp.GET("", api.index)
	grp.GET("/:collection", api.detail)
	grp.POST("", api.save)
	grp.DELETE("/:collection", api.delete)
}

type CollectionAPI struct {
	app core.App
}

func (a *CollectionAPI) index(c echo.Context) error {
	var col []model.Collection
	err := a.app.Store().DB().Select().From("_collection").All(&col)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	return c.JSON(200, col)
}

func (a *CollectionAPI) detail(c echo.Context) error {
    name := c.Param("collection")
	col, err := a.app.Store().FindCollectionByNameorId(a.app.Store().DB(), name)
	if err != nil {
		return err
	}
	return c.JSON(200, map[string]any{
		"collection": col,
	})
}

func (a *CollectionAPI) save(c echo.Context) error {
	var col model.Collection
	err := a.app.Store().SaveCollection(a.app.Store().DB(), &col)
	if err != nil {
		return err
	}
	return c.JSON(200, map[string]any{
		"message": "saved successfully",
	})
}

func (a *CollectionAPI) delete(c echo.Context) error {
	var col model.Collection
	err := a.app.Store().DeleteCollection(a.app.Store().DB(), &col)
	if err != nil {
		return err
	}
	return c.JSON(200, map[string]any{
		"message": "deleted successfully",
	})
}
