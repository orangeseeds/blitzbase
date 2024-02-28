package core

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestApp(t *testing.T) {
	app := DBApp{
		onStart: &Hook[*AppStartEvent]{},
	}

	id := app.OnStart().Add(func(e *AppStartEvent) error {
		t.Log("App started on port:", e.App.Addr())
		return nil
	})
	app.OnStart().Remove(id)

	app.Start(":8080")

}
