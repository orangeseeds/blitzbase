package api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/orangeseeds/blitzbase/core"
)

type RealtimeAPI struct {
	app core.App
}

// TODO: implement subscribe for certain actions
// 1. collection
// 2. record_id
// 3. specific icons : CREATE, VIEW,
func (api *RealtimeAPI) setSubscription(c echo.Context) error {
	// var req request.SetSubscriptionRequest
	// if err := c.Bind(&req); err != nil {
	//
	// 	return err
	// }
	//
	// err := c.Validate(req)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (api *RealtimeAPI) subscribe(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)

	bus := make(chan map[string]any, 1)

	createId := api.app.OnRecordCreate().Add(func(e *core.RecordEvent) error {
		bus <- map[string]any{
			"event":      "create",
			"collection": e.Record.Collection().GetName(),
			"record":     e.Record,
		}
		return nil
	})

	updateId := api.app.OnRecordUpdate().Add(func(e *core.RecordEvent) error {
		bus <- map[string]any{
			"event":      "update",
			"collection": e.Record.Collection().GetName(),
			"record":     e.Record,
		}
		return nil
	})

	deleteId := api.app.OnRecordDelete().Add(func(e *core.RecordEvent) error {
		bus <- map[string]any{
			"event":      "delete",
			"collection": e.Record.Collection().GetName(),
			"record":     e.Record,
		}
		return nil
	})

	enc := json.NewEncoder(c.Response())
	for l := range bus {
		if err := enc.Encode(l); err != nil {
			return err
		}
		c.Response().Flush()
	}

	api.app.OnRecordCreate().Remove(createId)
	api.app.OnRecordCreate().Remove(updateId)
	api.app.OnRecordCreate().Remove(deleteId)
	return nil
}
