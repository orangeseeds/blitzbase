package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/orangeseeds/blitzbase/refactored/core"
)

func SetupServer(app core.App) *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// e.Pre(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
	// 	Skipper: func(c echo.Context) bool {
	// 		// enable by default only for the API routes
	// 		return !strings.HasPrefix(c.Request().URL.Path, "/api/")
	// 	},
	// }))
	LoadAdminAPI(app, e)
	LoadCollectionAPI(app, e)
	LoadRecordAPI(app, e)
    LoadAuthRecordAPI(app,e)

	return e
}
