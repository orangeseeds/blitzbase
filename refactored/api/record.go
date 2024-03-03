package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/orangeseeds/blitzbase/refactored/core"
	model "github.com/orangeseeds/blitzbase/refactored/models"
	"github.com/orangeseeds/blitzbase/utils"
)

func LoadRecordAPI(app core.App, e *echo.Echo) {
	api := RecordAPI{app: app}

	grp := e.Group("/collections/:collection", LoadCollectionContextFromPath(app))

	grp.GET("/records", api.index)
	grp.GET("/records/:record", api.detail)

	grp.POST("/records", api.save)
	grp.DELETE("/records/:record", api.delete, LoadJWT())
}

type RecordAPI struct {
	app core.App
}

func (a *RecordAPI) index(c echo.Context) error {
	collection, ok := c.Get(string(utils.JWTCollection)).(*model.Collection)
	if !ok {
		return c.JSON(500, fmt.Errorf("couldnt conver to collection").Error())
	}
	records, err := a.app.Store().FindRecordsAll(a.app.Store().DB(), collection.Name)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	return c.JSON(200, map[string]any{
		"records_co": records,
	})
}

func (a *RecordAPI) detail(c echo.Context) error {
	id := c.Param("record")
	col := c.Param("collection")
	record, err := a.app.Store().FindRecordById(a.app.Store().DB(), id, col)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, map[string]any{
		"record": record,
	})
}

func (a *RecordAPI) save(c echo.Context) error {
	col := c.Get(string(utils.JWTCollection)).(*model.Collection)

	record := model.NewRecord(col)

	err := c.Bind(record)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	if col.IsAuth() {
		record.SetPassword(record.GetString(model.FieldPassword.String()))
	}

	err = a.app.Store().SaveRecord(a.app.Store().DB(), record)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	return c.JSON(200, map[string]any{
		"message": "saved successfully",
		"record":  record,
	})
}

func (a *RecordAPI) delete(c echo.Context) error {
	col := c.Get(string(utils.JWTCollection)).(*model.Collection)
	record := model.NewRecord(col)
	err := a.app.Store().DeleteRecord(a.app.Store().DB(), record)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, map[string]any{
		"message": "deleted successfully",
	})
}
