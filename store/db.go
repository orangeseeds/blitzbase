package store

import dbx "github.com/go-ozzo/ozzo-dbx"

type DBExector interface {
	Select(s ...string) *dbx.SelectQuery
	Insert(table string, cols dbx.Params) *dbx.Query
	Update(table string, cols dbx.Params, where dbx.Expression) *dbx.Query
	Delete(table string, where dbx.Expression) *dbx.Query
	DropTable(table string) *dbx.Query
	CreateTable(table string, cols map[string]string, options ...string) *dbx.Query
	Model(model any) *dbx.ModelQuery
}

func Wrap[T *dbx.DB | *dbx.Tx](db T) DBExector {
	switch db := any(db).(type) {
	case *dbx.DB:
		return &DBWrapper{
			db: db,
		}
	case *dbx.Tx:
		return &TxWrapper{
			db: db,
		}
	default:
		panic("Unsupported Type")
	}
}

type DBWrapper struct {
	db *dbx.DB
}

func (dbw *DBWrapper) Select(s ...string) *dbx.SelectQuery {
	return dbw.db.Select(s...)
}
func (dbw *DBWrapper) Insert(table string, cols dbx.Params) *dbx.Query {
	return dbw.db.Insert(table, cols)
}
func (dbw *DBWrapper) Update(table string, cols dbx.Params, where dbx.Expression) *dbx.Query {
	return dbw.db.Update(table, cols, where)
}
func (dbw *DBWrapper) Delete(table string, where dbx.Expression) *dbx.Query {
	return dbw.db.Delete(table, where)
}

func (dbw *DBWrapper) DropTable(table string) *dbx.Query {
	return dbw.db.DropTable(table)
}

func (dbw *DBWrapper) CreateTable(table string, cols map[string]string, options ...string) *dbx.Query {
	return dbw.db.CreateTable(table, cols, options...)
}

func (dbw *DBWrapper) Model(model any) *dbx.ModelQuery {
	return dbw.db.Model(model)
}

type TxWrapper struct {
	db *dbx.Tx
}

func (dbw *TxWrapper) Select(s ...string) *dbx.SelectQuery {
	return dbw.db.Select(s...)
}
func (dbw *TxWrapper) Insert(table string, cols dbx.Params) *dbx.Query {
	return dbw.db.Insert(table, cols)
}
func (dbw *TxWrapper) Update(table string, cols dbx.Params, where dbx.Expression) *dbx.Query {
	return dbw.db.Update(table, cols, where)
}
func (dbw *TxWrapper) Delete(table string, where dbx.Expression) *dbx.Query {
	return dbw.db.Delete(table, where)
}

func (dbw *TxWrapper) DropTable(table string) *dbx.Query {
	return dbw.db.DropTable(table)
}

func (dbw *TxWrapper) CreateTable(table string, cols map[string]string, options ...string) *dbx.Query {
	return dbw.db.CreateTable(table, cols, options...)
}

func (dbw *TxWrapper) Model(model any) *dbx.ModelQuery {
	return dbw.db.Model(model)
}
