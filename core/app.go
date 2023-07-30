package core

import "github.com/orangeseeds/blitzbase/store"

type App struct {
	Store *store.Storage
}

func NewApp(store *store.Storage) *App {
	return &App{
		Store: store,
	}
}
