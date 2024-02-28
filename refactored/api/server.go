package api

import (
	"github.com/labstack/echo"
	"github.com/orangeseeds/blitzbase/refactored/core"
)

func SetupServer(app core.App) *echo.Echo {
	e := echo.New()
	LoadCollectionAPI(app, e)
	return e
}
