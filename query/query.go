package query

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Passin a model
// specify fields
// check persmission and generate query accordingly from the Rules string

type FilterRule string

func (f FilterRule) BuildExpr(r *RecordFieldSpecifier) (dbx.Expression, error) {
	// return dbx.HashExp{
	// 	r.Collection.Name + ".field_one": "test",
	// }, nil
    return nil,nil

}
