package core

import (
	"github.com/orangeseeds/blitzbase/refactored/store"
)

type App interface {
	Store() *store.Store
	OnStart() *Hook[*AppStartEvent]
	Addr() string
}

type DBApp struct {
	store   *store.Store
	dbDir   string
	addr    string
	onStart *Hook[*AppStartEvent]
}

func (a *DBApp) Start(addr string) {
	a.addr = addr
	a.onStart.Trigger(&AppStartEvent{App: a})
}

func (a *DBApp) OnStart() *Hook[*AppStartEvent] { return a.onStart }

func (a *DBApp) Store() *store.Store { return a.store }

func (a *DBApp) Addr() string { return a.addr }
