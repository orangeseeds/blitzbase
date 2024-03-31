package request

import (
	"github.com/labstack/echo/v4"
	model "github.com/orangeseeds/blitzbase/models"
)

type Request[T model.Model] interface {
	Model() T
}

// passing a reference of body
func JsonValidate[Y model.Model, T Request[Y]](c echo.Context) (*T, error) {
	var body T
	err := c.Bind(&body)
	if err != nil {
		return nil, err
	}

	err = c.Validate(body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}
