package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/orangeseeds/blitzbase/store"
)

type collName string

func WithCollection(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "collection")
		collection := getCollection(name)
		if collection == nil {
			http.Error(w, "Collection not found", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), collName(name), collection)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func getCollection(collection string) store.Collection {
	return &store.BaseColletion{}
}
