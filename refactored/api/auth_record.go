package api

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/orangeseeds/blitzbase/refactored/core"
	model "github.com/orangeseeds/blitzbase/refactored/models"
	"github.com/orangeseeds/blitzbase/utils"
)

func LoadAuthRecordAPI(app core.App, e *echo.Echo) {
	api := AuthRecordAPI{app: app}

	grp := e.Group("/collections/:collection", LoadCollectionContextFromPath(app))

	grp.POST("/auth-with-password", api.authWithPassword)
	grp.POST("/reset-password", api.resetPassword)
	grp.POST("/confirm-reset-password", api.confirmResetPassword)

}

type AuthRecordAPI struct {
	app core.App
}

func (a *AuthRecordAPI) authWithPassword(c echo.Context) error {
	var authReq struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	err := c.Bind(&authReq)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	err = c.Validate(authReq)
	if err != nil {
		return c.JSON(400, err.Error())
	}

	col := c.Get("collection").(*model.Collection)

	record, err := a.app.Store().FindAuthRecordByEmail(a.app.Store().DB(), col.GetName(), authReq.Email)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	log.Println(record)
	valid := record.ValidatePassword(authReq.Password)
	if !valid {
		return c.JSON(400, "Password didnot match!")
	}

	// give claims
	authClaims := utils.JWTAuthClaims{
		Id:         record.Id,
		Type:       "collection",
		Collection: record.TableName(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	// generate token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)
	// encode token using secret
	encoded, err := jwtToken.SignedString([]byte("secret"))
	return c.JSON(200, map[string]any{
		"message": "auth with password success",
		"token":   encoded,
	})
}

func (a *AuthRecordAPI) resetPassword(c echo.Context) error { return nil }

func (a *AuthRecordAPI) confirmResetPassword(c echo.Context) error { return nil }

func (a *AuthRecordAPI) requestVerification(c echo.Context) error { return nil }

func (a *AuthRecordAPI) confirmRequestVeritication(c echo.Context) error { return nil }
