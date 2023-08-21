package api

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/orangeseeds/blitzbase/core"
)

type viewServer struct {
	app core.App
}

func (api *viewServer) Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/", api.handleIndex)
	return r
}

func (api *viewServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	temps := []string{
		"./views/index.html",
	}

	html, err := template.ParseFiles(temps...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = html.Execute(w, map[string]any{
		"Static": "/static",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
