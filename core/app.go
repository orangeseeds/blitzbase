package core

type App struct {
	Store *Storage
    Publisher *Publisher
}

func NewApp(store *Storage, pub *Publisher) *App {
	return &App{
		Store: store,
        Publisher: pub,
	}
}
