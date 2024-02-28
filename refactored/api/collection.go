package api

import (
	"github.com/labstack/echo"
	model "github.com/orangeseeds/blitzbase/refactored/models"
	"github.com/orangeseeds/blitzbase/refactored/query"
	"github.com/orangeseeds/blitzbase/utils"
)

func listCollection(c echo.Context) error {
	col := model.NewCollection(
		utils.RandStr(10),
		"test_collection",
		model.BASE,
	)
	col.Rule = "$req.auth=user"

	fieldSpec := query.RecordFieldSpecifier{
		Collection:    col,
		AllowedFields: []string{},
		Request:       nil,
	}

	_, err := query.FilterRule(col.Rule).BuildExpr(&fieldSpec)
	if err != nil {
		return err
	}
	return nil
}
