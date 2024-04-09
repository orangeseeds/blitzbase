package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	Code int    `json:"code"`
	Info string `json:"info"`
	Data any    `json:"data"`
}

func (e *ApiError) Error() string {
	return e.Info
}

func NewApiError(code int, msg string, data any) *ApiError {
	return &ApiError{
		Code: code,
		Info: msg,
		Data: parseErr(data),
	}
}

func NewNotFoundError(msg string, data any) *ApiError {
	if msg == "" {
		msg = "Requested resource not found."
	}
	return NewApiError(http.StatusNotFound, msg, data)
}
func NewBadRequestError(msg string, data any) *ApiError {
	if msg == "" {
		msg = "Error processing the request."
	}
	return NewApiError(http.StatusBadRequest, msg, data)
}
func NewForbiddenError(msg string, data any) *ApiError {
	if msg == "" {
		msg = "Not allowed to perform this request."
	}
	return NewApiError(http.StatusForbidden, msg, data)
}
func NewUnauthorizedError(msg string, data any) *ApiError {
	if msg == "" {
		msg = "Missing or invalid token passed."
	}
	return NewApiError(http.StatusUnauthorized, msg, data)
}

func parseErr(data any) map[string]any {
	errorMsg := map[string]any{}

	if data == nil {
		return errorMsg
	}

	if errs, ok := data.(validator.ValidationErrors); ok {
		for _, e := range errs {
			errorMsg[strings.ToLower(e.Field())] = fmt.Sprintf("field '%s' need to be %s %s", e.Field(), e.Tag(), e.Param())
		}
	} else {
		errorMsg["message"] = data.(error).Error()
	}

	return errorMsg
}
