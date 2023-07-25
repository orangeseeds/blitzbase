package core

import (
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func initDriverWithUpdateHook(publisher *Publisher) string {
	driverName := "sqlite3_with_hook_example"
	sql.Register(driverName,
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				conn.RegisterUpdateHook(func(op int, db string, table string, rowid int64) {
					switch op {
					case sqlite3.SQLITE_INSERT:
						msg := fmt.Sprintf(" Data inserted on db %s, table %s and row %d", db, table, rowid)

						publisher.Broadcast(
							msg,
							"sqlite_insert",
						)
					}
				})
				return nil
			},
		})

	return driverName
}
