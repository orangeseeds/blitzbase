package query

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
	model "github.com/orangeseeds/blitzbase/refactored/models"
)

type RecordFieldSpecifier struct {
	Collection    *model.Collection
	AllowedFields []string
	Request       echo.Context
}

func parse(code string) {
}

func (r *RecordFieldSpecifier) Parse() *dbx.Expression {
	// $request.auth
	// $request.method
	// $request.data
	parse(r.Collection.Rule)
	return nil

}
