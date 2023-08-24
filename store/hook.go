package store

import (
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

type HookData struct {
	RecordID       int    `json:"record_id"`
	CollectionName string `json:"collection_name"`
}

func initDriverWithUpdateHook(publisher *Publisher) string {
	driverName := "sqlite3_with_create_update_delete_hooks"
	sql.Register(driverName,
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				conn.RegisterUpdateHook(func(op int, db string, table string, rowid int64) {

					data := cleanHookData(table, rowid)
					switch op {
					case sqlite3.SQLITE_INSERT:
						handleInsertHook(publisher, data)
					case sqlite3.SQLITE_UPDATE:
						handleUpdateHook(publisher, data)
					case sqlite3.SQLITE_DELETE:
						handleDeleteHook(publisher, data)
					}
				})
				return nil
			},
		})

	return driverName
}

func cleanHookData(table string, rowID int64) *HookData {
	return &HookData{
		RecordID:       int(rowID),
		CollectionName: table,
	}
}

func handleUpdateHook(pub *Publisher, data *HookData) {
	pub.Dispatch(DBHookEvent{
		Message: Message{
			Action: UPDATE,
			Record: struct {
				ID         string
				Collection string
			}{
				fmt.Sprintf("%d", data.RecordID),
				data.CollectionName,
			},
		},
	})
}
func handleDeleteHook(pub *Publisher, data *HookData) {
	pub.Dispatch(DBHookEvent{
		Message: Message{
			Action: DELETE,
			Record: struct {
				ID         string
				Collection string
			}{
				fmt.Sprintf("%d", data.RecordID),
				data.CollectionName,
			},
		},
	})
}
func handleInsertHook(pub *Publisher, data *HookData) {
	pub.Dispatch(DBHookEvent{
		Message: Message{
			Action: INSERT,
			Record: struct {
				ID         string
				Collection string
			}{
				fmt.Sprintf("%d", data.RecordID),
				data.CollectionName,
			},
		},
	})
}
