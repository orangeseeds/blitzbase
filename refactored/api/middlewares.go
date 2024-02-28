package api

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/orangeseeds/blitzbase/refactored/core"
	"github.com/orangeseeds/blitzbase/utils"
)

// loads the jwt claims into "user" key in the request context
func LoadJWT() echo.MiddlewareFunc {
	requireJWTAdminAuth := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(utils.JWTAuthClaims)
		},
		SigningKey: []byte("secret"),
	}
	return echojwt.WithConfig(requireJWTAdminAuth)
}

func NeedsAdminAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return c.JSON(400, "cannot convert to jwt.Token")
			}
			authClaims := token.Claims.(*utils.JWTAuthClaims)

			if authClaims.Type != utils.JWTAdmin {
				return c.JSON(400, "jwt type is not admin")
			}
			return next(c)
		}
	}
}

func LoadAuthContextFromToken(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			authClaims := token.Claims.(*utils.JWTAuthClaims)

			switch authClaims.Type {
			case utils.JWTAdmin:
				admin, err := app.Store().FindAdminById(app.Store().DB(), authClaims.Id)
				if err != nil {
					return c.JSON(400, err.Error())
				}
				c.Set(string(utils.JWTAdmin), admin)
			case utils.JWTCollection:
				record, err := app.Store().FindRecordById(app.Store().DB(), authClaims.Collection, authClaims.Id)
				if err != nil {
					return c.JSON(400, err.Error())
				}
				c.Set(string(utils.JWTAdmin), record)
			}
			return next(c)
		}
	}
}

func LoadCollectionContextFromPath(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			name := c.Param(string(utils.JWTCollection))
			collection, err := app.Store().FindCollectionByNameorId(app.Store().DB(), name)
			if err != nil {
				return c.JSON(400, err.Error())
			}

			c.Set(string(utils.JWTCollection), collection)
			return next(c)
		}
	}
}
