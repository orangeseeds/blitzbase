package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

type customValidator struct {
	validator *validator.Validate
}

func NewCustomValidator(validator *validator.Validate) *customValidator {
	return &customValidator{
		validator: validator,
	}
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
