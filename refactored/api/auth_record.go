package api

import (
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

func (a *AuthRecordAPI) resetPassword(c echo.Context) error {
	var resetReq struct {
		Email string `json:"email" validate:"required"`
	}
	err := c.Bind(&resetReq)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	err = c.Validate(resetReq)
	if err != nil {
		return c.JSON(400, err.Error())
	}

	collection := c.Get(string(utils.JWTCollection)).(*model.Collection)

	record, err := a.app.Store().FindAuthRecordByEmail(a.app.Store().DB(), collection.Name, resetReq.Email)
	if err != nil {
		return c.JSON(400, err.Error())
	}

	// email to admin.Email
	return c.JSON(200, map[string]any{
		"token": record.GetString(model.FieldToken),
	})
}

func (a *AuthRecordAPI) confirmResetPassword(c echo.Context) error {
	var confirmReq struct {
		Token           string `json:"token" validate:"required"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirm_password" validate:"required"`
	}

	err := c.Bind(&confirmReq)
	if err != nil {
		return c.JSON(400, err.Error())
	}

	err = c.Validate(confirmReq)
	if err != nil {
		return c.JSON(400, err.Error())
	}

	if confirmReq.Password != confirmReq.ConfirmPassword {
		return c.JSON(400, "password and confirm_password not equal.")
	}

	coll := c.Get(string(utils.JWTCollection)).(*model.Collection)

	record, err := a.app.Store().FindAuthRecordByToken(a.app.Store().DB(), coll.Name, confirmReq.Token)
	if err != nil {
		return c.JSON(400, err.Error())
	}

	record.SetPassword(confirmReq.ConfirmPassword)
	record.RefreshToken()

	err = a.app.Store().UpdateRecord(a.app.Store().DB(), coll.Name, record)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	return c.JSON(200, map[string]any{
		"message": "password updated successfully!",
	})
}

func (a *AuthRecordAPI) requestVerification(c echo.Context) error { return nil }

func (a *AuthRecordAPI) confirmRequestVeritication(c echo.Context) error { return nil }
