package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/orangeseeds/blitzbase/core"
	"github.com/orangeseeds/blitzbase/store"
	"github.com/orangeseeds/blitzbase/utils"
)

const (
	CtxUserKey       string = "user"
	CtxAdminKey      string = "admin"
	CtxAuthRecordKey string = "authRecord"
	CtxCollectionKey string = "collection"
)

func LoadJWT() echo.MiddlewareFunc {
	requireJWTAdminAuth := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(utils.JWTAuthClaims)
		},
		SigningKey: []byte("secret"),
	}
	return echojwt.WithConfig(requireJWTAdminAuth)
}

// check if jwt is of type admin and loads admin into ctx
func NeedsAdminAuth(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get(CtxUserKey).(*jwt.Token)
			if !ok {
				return c.JSON(500, "cannot convert contex.%s to jwt.Token")
			}

			authClaims := token.Claims.(*utils.JWTAuthClaims)
			if authClaims.Type != utils.JwtTypeAdmin {
				return c.JSON(400, "jwt type is not admin")
			}

			admin, err := app.Store().FindAdminById(app.Store().DB(), authClaims.Id)
			if err != nil {
				return c.JSON(500, map[string]any{
					"message": "invalid jwt token",
				})
			}
			c.Set("admin", admin)
			return next(c)
		}
	}
}

func LoadAuthContextFromToken(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get(CtxUserKey).(*jwt.Token)
			claims := token.Claims.(*utils.JWTAuthClaims)

			exec := store.Wrap(app.Store().DB())
			switch claims.Type {
			case utils.JwtTypeAdmin:
				admin, err := app.Store().FindAdminById(exec, claims.Id)
				if err != nil {
					return c.JSON(400, err.Error())
				}
				c.Set(CtxAdminKey, admin)
			case utils.JwtTypeCollection:
				record, err := app.Store().FindRecordById(exec, claims.Collection, claims.Id)
				if err != nil {
					return c.JSON(400, err.Error())
				}
				c.Set(CtxAuthRecordKey, record)
			}
			return next(c)
		}
	}
}

func LoadCollectionContextFromPath(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			name := c.Param(CtxCollectionKey)
			exec := store.Wrap(app.Store().DB())
			collection, err := app.Store().FindCollectionByNameorId(exec, name)
			if err != nil {
				return c.JSON(400, err.Error())
			}

			c.Set(string(CtxCollectionKey), collection)
			return next(c)
		}
	}
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	var apiErr = new(ApiError)
	if errors.As(err, &apiErr) {
		if apiErr.Data != nil {
			log.Println(apiErr.Data)
		}
	} else if v := new(echo.HTTPError); errors.As(err, &v) {
		if v.Internal != nil {
			log.Println(v.Internal)
		}
		msg := fmt.Sprintf("%v", v.Message)
		apiErr = NewApiError(v.Code, msg, v)
	} else {
		if true { //debug
			log.Println("as", err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			apiErr = NewNotFoundError("", err)
		} else {
			apiErr = NewBadRequestError("", err)
		}
	}
	c.Logger().Error(apiErr)
	c.JSON(apiErr.Code, map[string]any{
		"error": apiErr,
	})
}
