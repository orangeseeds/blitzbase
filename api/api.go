package api

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/orangeseeds/blitzbase/core"
	"github.com/orangeseeds/blitzbase/utils"
)

func LoadAllAPIRoutes(app core.App, e *echo.Echo) {
	adminAPI := AdminAPI{app: app}
	{
		grp := e.Group("/admins")
		grp.GET("", adminAPI.index, LoadJWT(), NeedsAdminAuth(app))
		grp.GET("/:id", adminAPI.detail, LoadJWT(), NeedsAdminAuth(app))
		grp.POST("", adminAPI.save, LoadJWT(), NeedsAdminAuth(app))
		grp.DELETE("/:id", adminAPI.delete, LoadJWT(), NeedsAdminAuth(app))

		grp.POST("/auth-with-password", adminAPI.authWithPassword)
		grp.POST("/reset-password", adminAPI.resetPassword)
		grp.POST("/confirm-reset-password", adminAPI.confirmResetPassword)
	}

	collAPI := CollectionAPI{app: app}
	{
		grp := e.Group("/collections", LoadJWT(), NeedsAdminAuth(app))
		grp.GET("", collAPI.index)
		grp.POST("", collAPI.save)
		grp.GET("/:collection", collAPI.detail)
		grp.DELETE("/:collection", collAPI.delete)
	}

	recordAPI := RecordAPI{app: app}
	{
		grp := e.Group("/collections/:collection/records", LoadCollectionContextFromPath(app))
		grp.GET("", recordAPI.index)
		grp.POST("", recordAPI.save, LoadJWT())
		grp.GET("/:record", recordAPI.detail)
		grp.DELETE("/:record", recordAPI.delete, LoadJWT())

		grp.POST("/auth-with-password", recordAPI.authWithPassword)
		grp.POST("/reset-password", recordAPI.resetPassword)
		grp.POST("/confirm-reset-password", recordAPI.confirmResetPassword)
	}

	realtimeAPI := RealtimeAPI{app: app}
	{
		grp := e.Group("/realtime")
		grp.POST("", realtimeAPI.setSubscription)
		grp.GET("", realtimeAPI.subscribe)
	}

	// TODO: Add the route list printer as a helper function which can be accessed using a flag.
	_, err := json.MarshalIndent(e.Routes(), "", " ")
	if err != nil {
		panic(err)
	}
	// log.Printf("Routes: %v\n", string(data))
}

func SetupServer(app core.App) *echo.Echo {
	e := echo.New()
	e.Validator = utils.NewCustomValidator(validator.New())
	e.HTTPErrorHandler = CustomHTTPErrorHandler
	e.Pre(middleware.RemoveTrailingSlash())
	e.HideBanner = true
	LoadAllAPIRoutes(app, e)
	return e
}
