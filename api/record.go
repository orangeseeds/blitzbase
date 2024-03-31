package api

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/orangeseeds/blitzbase/core"
	model "github.com/orangeseeds/blitzbase/models"
	"github.com/orangeseeds/blitzbase/utils"
)

type RecordAPI struct {
	app core.App
}

func (a *RecordAPI) index(c echo.Context) error {
	collection, ok := c.Get(string(utils.JwtTypeCollection)).(*model.Collection)
	if !ok {
		return c.JSON(500, fmt.Errorf("couldnt conver to collection").Error())
	}
	records, err := a.app.Store().FindRecordsAll(a.app.Store().DB(), collection.Name)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	a.app.OnRecordIndex().Trigger(&core.RecordEvent{
		Type:    core.IndexEvent,
		Request: &c,
	})

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

	a.app.OnRecordIndex().Trigger(&core.RecordEvent{
		Type:    core.DetailEvent,
		Record:  record,
		Request: &c,
	})

	return c.JSON(200, map[string]any{
		"record": record,
	})
}

func (a *RecordAPI) save(c echo.Context) error {
	col := c.Get(string(utils.JwtTypeCollection)).(*model.Collection)

	record := model.NewRecord(col)

	err := c.Bind(record)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	if col.IsAuth() {
		err := record.SetPassword(record.GetString(model.FieldPassword))
		if err != nil {
			return c.JSON(500, err.Error())
		}
	}

	err = a.app.Store().SaveRecord(a.app.Store().DB(), record)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	a.app.OnRecordIndex().Trigger(&core.RecordEvent{
		Type:    core.CreateEvent,
		Record:  record,
		Request: &c,
	})

	return c.JSON(200, map[string]any{
		"message": "saved successfully",
		"record":  record,
	})
}

func (a *RecordAPI) delete(c echo.Context) error {
	col := c.Get(string(utils.JwtTypeCollection)).(*model.Collection)
	record := model.NewRecord(col)
	err := a.app.Store().DeleteRecord(a.app.Store().DB(), record)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	a.app.OnRecordIndex().Trigger(&core.RecordEvent{
		Type:    core.DetailEvent,
		Record:  record,
		Request: &c,
	})

	return c.JSON(200, map[string]any{
		"message": "deleted successfully",
	})
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

	a.app.OnRecordAuth().Trigger(&core.RecordEvent{
		Type:    core.AuthEvent,
		Record:  record,
		Request: &c,
	})

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

	collection := c.Get(string(utils.JwtTypeCollection)).(*model.Collection)

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

	coll := c.Get(string(utils.JwtTypeCollection)).(*model.Collection)

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

	a.app.OnRecordIndex().Trigger(&core.RecordEvent{
		Type:    core.UpdateEvent,
		Record:  record,
		Request: &c,
	})

	return c.JSON(200, map[string]any{
		"message": "password updated successfully!",
	})
}

// func (a *AuthRecordAPI) requestVerification(c echo.Context) error { return nil }
//
// func (a *AuthRecordAPI) confirmRequestVeritication(c echo.Context) error { return nil }
