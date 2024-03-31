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
		grp.GET("/d/:collection", collAPI.detail)
		grp.DELETE("/d/:collection", collAPI.delete)
	}

	recordAPI := RecordAPI{app: app}
	{
		grp := e.Group("/collections/:collection/records", LoadCollectionContextFromPath(app))
		grp.GET("", recordAPI.index)
		grp.POST("", recordAPI.save)
		grp.GET("/:record", recordAPI.detail)
		grp.DELETE("/:record", recordAPI.delete, LoadJWT())
	}

	authRecAPI := AuthRecordAPI{app: app}
	{
		grp := e.Group("/collections/:collection", LoadCollectionContextFromPath(app))
		grp.POST("/auth-with-password", authRecAPI.authWithPassword)
		grp.POST("/reset-password", authRecAPI.resetPassword)
		grp.POST("/confirm-reset-password", authRecAPI.confirmResetPassword)
	}

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
	LoadAllAPIRoutes(app, e)
	return e
}
