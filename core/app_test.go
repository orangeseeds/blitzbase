package core

import (
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestApp(t *testing.T) {
	app := DBApp{
		onStart:     &Hook[*AppEvent]{},
		onAdminAuth: &Hook[*AdminEvent]{},
	}

	t.Log("before", time.Now())
	app.onAdminAuth.Trigger(&AdminEvent{})
	t.Log("after", time.Now())

	t.Log("Doing other things...")
}
