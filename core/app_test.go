package core

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestApp(t *testing.T) {
	app := DBApp{
		onStart: &Hook[*AppEvent]{},
	}

	id := app.OnStart().Add(func(e *AppEvent) error {
		t.Log("App started on port:", e.App.Addr())
		return nil
	})
	app.OnStart().Remove(id)

	app.Start(":8080")

}
