package api

import (
	"net/http"
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
		Data: data,
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

// func parseErr(data any) map[string]any {
// 	errorMsg := map[string]any{}
//
// 	log.Println("hey", errors.Unwrap(data.(*echo.HTTPError)))
//
// 	if errors.Is(data.(error), validator.ValidationErrors{}) {
// 		for _, e := range data.(validator.ValidationErrors) {
// 			errorMsg[e.Field()] = e.Error()
// 		}
// 	}
// 	// default:
// 	// 	errorMsg["other"] = data
// 	return errorMsg
// }
