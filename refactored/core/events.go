package core

import (
	"github.com/labstack/echo/v4"
	model "github.com/orangeseeds/blitzbase/refactored/models"
)

const (
	StartEvent     = "StartEvent"
	TerminateEvent = "TerminateEvent"

	IndexEvent  = "IndexEvent"
	DetailEvent = "DetailEvent"
	CreateEvent = "CreateEvent"
	UpdateEvent = "UpdateEvent"
	DeleteEvent = "DeleteEvent"
	AuthEvent   = "AuthEvent"
)

type AppEvent struct {
	Type string // start, terminate
	App  App
}

type AdminEvent struct {
	Type    string
	Admin   *model.Admin
	Request *echo.Context
}

type CollectionEvent struct {
	Type       string
	Collection *model.Collection
	Request    *echo.Context
}

type RecordEvent struct {
	Type    string
	Record  *model.Record
	Request *echo.Context
}
