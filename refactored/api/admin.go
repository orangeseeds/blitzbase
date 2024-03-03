package api

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/orangeseeds/blitzbase/refactored/core"
	model "github.com/orangeseeds/blitzbase/refactored/models"
	"github.com/orangeseeds/blitzbase/utils"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func LoadAdminAPI(app core.App, e *echo.Echo) {
	api := AdminAPI{app: app}
	grp := e.Group("/admins")

	grp.GET("", api.index, LoadJWT(), NeedsAdminAuth(app))
	grp.GET("/:id", api.detail, LoadJWT(), NeedsAdminAuth(app))

	grp.POST("/auth-with-password", api.authWithPassword)
	grp.POST("/reset-password", api.resetPassword)
	grp.POST("/confirm-reset-password", api.confirmResetPassword)

	grp.POST("", api.save, LoadJWT(), NeedsAdminAuth(app))
	grp.DELETE("/:collection", api.delete, LoadJWT(), NeedsAdminAuth(app))
}

type AdminAPI struct {
	app core.App
}

func (a *AdminAPI) index(c echo.Context) error {
	var admin model.Admin
	var admins []model.Admin
	err := a.app.Store().DB().Select().From(admin.TableName()).All(&admins)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	event := core.AdminEvent{
		Type:    core.IndexEvent,
		Request: &c,
	}
	a.app.OnAdminIndex().Trigger(&event)

	return c.JSON(200, admins)
}

func (a *AdminAPI) detail(c echo.Context) error {
	id := c.Param("id")
	admin, err := a.app.Store().FindAdminById(a.app.Store().DB(), id)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	event := core.AdminEvent{
		Type:    core.DetailEvent,
		Admin:   admin,
		Request: &c,
	}
	a.app.OnAdminDetail().Trigger(&event)
	return c.JSON(200, map[string]any{
		"admin": admin,
	})
}

func (a *AdminAPI) save(c echo.Context) error {
	var admin model.Admin
	err := c.Bind(&admin)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	admin.SetPassword(admin.Password)
	admin.RefreshToken()

	err = a.app.Store().SaveAdmin(a.app.Store().DB(), &admin)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	event := core.AdminEvent{
		Type:    core.CreateEvent,
		Admin:   &admin,
		Request: &c,
	}
	a.app.OnAdminDetail().Trigger(&event)

	return c.JSON(200, map[string]any{
		"message": "saved successfully",
		"admin":   admin,
	})
}

func (a *AdminAPI) delete(c echo.Context) error {
	var admin model.Admin
	id := c.Param("id")
	admin.SetID(id)
	err := a.app.Store().DeleteAdmin(a.app.Store().DB(), &admin)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	event := core.AdminEvent{
		Type:    core.DeleteEvent,
		Admin:   &admin,
		Request: &c,
	}
	a.app.OnAdminDelete().Trigger(&event)

	return c.JSON(200, map[string]any{
		"message": "deleted successfully",
		"admin":   admin,
	})
}

func (a *AdminAPI) authWithPassword(c echo.Context) error {
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

	admin, err := a.app.Store().FindAdminByEmail(a.app.Store().DB(), authReq.Email)
	if err != nil {
		return c.JSON(400, err.Error())
	}

	valid := admin.ValidatePassword(authReq.Password)
	if !valid {
		return c.JSON(400, "Password didnot match!")
	}

	// give claims
	authClaims := utils.JWTAuthClaims{
		Id:         admin.Id,
		Type:       "admin",
		Collection: admin.TableName(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	// generate token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)
	// encode token using secret
	encoded, err := jwtToken.SignedString([]byte("secret"))

	event := core.AdminEvent{
		Type:    core.AuthEvent,
		Admin:   admin,
		Request: &c,
	}
	a.app.OnAdminAuth().Trigger(&event)

	return c.JSON(200, map[string]any{
		"message": "auth with password success",
		"token":   encoded,
	})
}

func (a *AdminAPI) resetPassword(c echo.Context) error {
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

	admin, err := a.app.Store().FindAdminByEmail(a.app.Store().DB(), resetReq.Email)
	if err != nil {
		return c.JSON(400, err.Error())
	}

	// email to admin.Email
	return c.JSON(200, map[string]any{
		"token": admin.Token,
	})
}

func (a *AdminAPI) confirmResetPassword(c echo.Context) error {
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
		return c.JSON(400, fmt.Errorf("password and confirm_password not equal."))
	}

	admin, err := a.app.Store().FindAdminByToken(a.app.Store().DB(), confirmReq.Token)
	if err != nil {
		return c.JSON(400, err.Error())
	}

	admin.SetPassword(confirmReq.ConfirmPassword)
	admin.RefreshToken()

	err = a.app.Store().UpdateAdmin(a.app.Store().DB(), admin)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	event := core.AdminEvent{
		Type:    core.UpdateEvent,
		Admin:   admin,
		Request: &c,
	}
	a.app.OnAdminUpdate().Trigger(&event)

	return c.JSON(200, map[string]any{
		"message": "password updated successfully!",
	})
}
