package api

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/orangeseeds/blitzbase/core"
	model "github.com/orangeseeds/blitzbase/models"
	"github.com/orangeseeds/blitzbase/request"
	"github.com/orangeseeds/blitzbase/store"
	"github.com/orangeseeds/blitzbase/utils"
	"github.com/orangeseeds/blitzbase/utils/resolver"
)

type RecordAPI struct {
	app core.App
}

func (a *RecordAPI) index(c echo.Context) error {
	log.Println("here")
	collection, ok := c.Get(string(utils.JwtTypeCollection)).(*model.Collection)
	if !ok {
		return NewApiError(500, "some error occured", nil)
	}
	exec := store.Wrap(a.app.Store().DB())
	records, err := a.app.Store().FindRecordsAll(exec, collection.Name)
	if err != nil {
		return NewApiError(500, "some error occured", err)
	}

	if collection.IndexRule != "" {
		allowed := resolver.CheckPermission(collection.IndexRule, resolver.RequestInfo{
			Data:       map[string]any{},
			Method:     "GET",
			Collection: collection,
			Context:    c,
		})
		log.Println("Permission: ", allowed)
		if !allowed {
			return NewForbiddenError("Rule doesn't allow access to the collection", nil)
		}
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
	col := c.Param(CtxCollectionKey)
	exec := store.Wrap(a.app.Store().DB())
	record, err := a.app.Store().FindRecordById(exec, id, col)
	if err != nil {
		return NewNotFoundError("", err)
	}

	collection, ok := c.Get(string(utils.JwtTypeCollection)).(*model.Collection)
	if !ok {
		return NewApiError(500, "some error occured", nil)
	}

	if collection.DetailRule != "" {
		allowed := resolver.CheckPermission(collection.DetailRule, resolver.RequestInfo{
			Data:       map[string]any{},
			Method:     "GET",
			Collection: collection,
			Context:    c,
		})
		log.Println("Permission: ", allowed)
		if !allowed {
			return NewForbiddenError("Rule doesn't allow access to the collection", nil)
		}
	}

	a.app.OnRecordDetail().Trigger(&core.RecordEvent{
		Type:    core.DetailEvent,
		Record:  record,
		Request: &c,
	})

	return c.JSON(200, map[string]any{
		"record": record,
	})
}

func (a *RecordAPI) save(c echo.Context) error {
	col, _ := c.Get(string(utils.JwtTypeCollection)).(*model.Collection)

	record := model.NewRecord(col)

	err := c.Bind(record)
	if err != nil {
		return NewBadRequestError("", err)
	}

	if col.IsAuth() {
		err := record.SetPassword(record.GetString(model.FieldPassword))
		if err != nil {
			return NewApiError(500, "some error occured", err)
		}
	}

	exec := store.Wrap(a.app.Store().DB())
	err = a.app.Store().SaveRecord(exec, record)
	if err != nil {
		return NewBadRequestError("error occured when saving record.", err)
	}

	if col.UpdateRule != "" {
		allowed := resolver.CheckPermission(col.UpdateRule, resolver.RequestInfo{
			Data:       map[string]any{},
			Method:     "POST",
			Collection: col,
			Context:    c,
		})
		log.Println("Permission: ", allowed)
		if !allowed {
			return NewForbiddenError("Rule doesn't allow access to the collection", nil)
		}
	}

	a.app.OnRecordCreate().Trigger(&core.RecordEvent{
		Type:    core.IndexEvent,
		Record:  record,
		Request: &c,
	})

	return c.JSON(200, map[string]any{
		"message": "saved successfully",
		"record":  record,
	})
}

func (a *RecordAPI) delete(c echo.Context) error {
	col, _ := c.Get(string(utils.JwtTypeCollection)).(*model.Collection)
	record := model.NewRecord(col)

	exec := store.Wrap(a.app.Store().DB())
	err := a.app.Store().DeleteRecord(exec, record)
	if err != nil {
		return NewBadRequestError("error occured when deleting record.", err)
	}

	if col.DeleteRule != "" {
		allowed := resolver.CheckPermission(col.DeleteRule, resolver.RequestInfo{
			Data:       map[string]any{},
			Method:     "DELETE",
			Collection: col,
			Context:    c,
		})
		log.Println("Permission: ", allowed)
		if !allowed {
			return NewForbiddenError("Rule doesn't allow access to the collection", nil)
		}
	}

	a.app.OnRecordDelete().Trigger(&core.RecordEvent{
		Type:    core.DeleteEvent,
		Record:  record,
		Request: &c,
	})

	return c.JSON(200, map[string]any{
		"message": "deleted successfully",
	})
}

func (a *RecordAPI) authWithPassword(c echo.Context) error {
	req, err := request.JsonValidate[model.Record, request.RecordAuthWithPasswordRequest](c)
	if err != nil {
		return NewBadRequestError("", err)
	}

	col, _ := c.Get(CtxCollectionKey).(*model.Collection)

	exec := store.Wrap(a.app.Store().DB())
	record, err := a.app.Store().FindAuthRecordByEmail(exec, col.GetName(), req.Email)
	if err != nil {
		return NewNotFoundError("record with given email not found.", err)
	}
	valid := record.ValidatePassword(req.Password)
	if !valid {
		return NewNotFoundError("record with given email and password not found.", err)
	}

	// give claims
	authClaims := utils.JWTAuthClaims{
		Id:         record.Id,
		Type:       utils.JwtTypeCollection,
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

func (a *RecordAPI) resetPassword(c echo.Context) error {
	req, err := request.JsonValidate[model.Record, request.RecordResetPasswordRequest](c)
	if err != nil {
		return NewBadRequestError("", err)
	}

	coll, _ := c.Get(string(utils.JwtTypeCollection)).(*model.Collection)

	exec := store.Wrap(a.app.Store().DB())
	record, err := a.app.Store().FindAuthRecordByEmail(exec, coll.GetName(), req.Email)
	if err != nil {
		return NewNotFoundError("record with given email not found.", err)
	}

	// email to admin.Email
	return c.JSON(200, map[string]any{
		"token": record.GetString(model.FieldToken),
	})
}

func (a *RecordAPI) confirmResetPassword(c echo.Context) error {
	req, err := request.JsonValidate[model.Record, request.RecordConfirmResetPasswordRequest](c)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	coll, _ := c.Get(string(utils.JwtTypeCollection)).(*model.Collection)

	exec := store.Wrap(a.app.Store().DB())
	record, err := a.app.Store().FindAuthRecordByToken(exec, coll.Name, req.Token)
	if err != nil {
		return NewNotFoundError("record with given token not found.", err)
	}

	record.SetPassword(req.ConfirmPassword)
	record.RefreshToken()

	err = a.app.Store().UpdateRecord(exec, coll.Name, record)
	if err != nil {
		return NewBadRequestError("Error updating record.", err)
	}

	a.app.OnRecordUpdate().Trigger(&core.RecordEvent{
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
