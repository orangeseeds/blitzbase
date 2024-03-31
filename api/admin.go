package api

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/orangeseeds/blitzbase/core"
	model "github.com/orangeseeds/blitzbase/models"
	"github.com/orangeseeds/blitzbase/request"
	"github.com/orangeseeds/blitzbase/store"
	"github.com/orangeseeds/blitzbase/utils"
)

type AdminAPI struct {
	app core.App
}

func (a *AdminAPI) index(c echo.Context) error {
	var admin model.Admin
	var admins []model.Admin
	err := a.app.Store().DB().Select().From(admin.TableName()).All(&admins)
	if err != nil {
		return NewApiError(500, "something went wrong!", nil)
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
		return NewNotFoundError("", err)
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
	req, err := request.JsonValidate[model.Admin, request.AdminSaveRequest](c)
	if err != nil {
		return NewBadRequestError("", err)
	}

	admin := req.Model()
	admin.SetPassword(admin.Password)
	admin.RefreshToken()

	err = a.app.Store().SaveAdmin(a.app.Store().DB(), &admin)
	if err != nil {
		return NewBadRequestError("Error saving admin.", err)
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
		return NewBadRequestError("Error deleting admin.", err)
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
	req, err := request.JsonValidate[model.Admin, request.AdminSaveRequest](c)
	if err != nil {
		return NewBadRequestError("", err)
	}

	exec := store.Wrap(a.app.Store().DB())
	admin, err := a.app.Store().FindAdminByEmail(exec, req.Email)
	if err != nil {
		return NewNotFoundError("admin with given email not found.", err)
	}

	valid := admin.ValidatePassword(req.Password)
	if !valid {
		return NewNotFoundError("admin with given email and password not found.", err)
	}

	// give claims
	authClaims := utils.JWTAuthClaims{
		Id:         admin.Id,
		Type:       utils.JwtTypeAdmin,
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
	req, err := request.JsonValidate[model.Admin, request.AdminSaveRequest](c)
	if err != nil {
		return NewBadRequestError("", err)
	}

	admin, err := a.app.Store().FindAdminByEmail(
		store.Wrap(a.app.Store().DB()),
		req.Email,
	)
	if err != nil {
		return NewNotFoundError("admin with given email not found.", err)
	}

	// email to admin.Email
	return c.JSON(200, map[string]any{
		"token": admin.Token,
	})
}

func (a *AdminAPI) confirmResetPassword(c echo.Context) error {
	confirmReq, err := request.JsonValidate[model.Admin, request.AdminConfirmResetPasswordRequest](c)
	if err != nil {
		return NewBadRequestError("", err)
	}
	exec := store.Wrap(a.app.Store().DB())
	admin, err := a.app.Store().FindAdminByToken(exec, confirmReq.Token)
	if err != nil {
		return NewNotFoundError("admin with given token not found.", err)
	}

	admin.SetPassword(confirmReq.ConfirmPassword)
	admin.RefreshToken()

	err = a.app.Store().UpdateAdmin(exec, admin)
	if err != nil {
		return NewBadRequestError("Error updating admin.", err)
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
