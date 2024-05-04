package utils

import (
	"github.com/go-playground/validator/v10"
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
	return cv.validator.Struct(i)
}
