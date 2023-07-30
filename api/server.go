package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/orangeseeds/blitzbase/core"
)

func Serve(app core.App, addr string) {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Use(middleware.AllowContentType("application/json"))
	mux.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Content-Type", "application/json")

			h.ServeHTTP(w, r)
		})
	})

	rtApp := rtServer{
		app: app,
	}

	authApp := authServer{
		app: app,
	}

	mux.Mount("/realtime", rtApp.Router())
	mux.Mount("/auth", authApp.Router())

	log.Println("serving on port :3300")
	log.Fatal(http.ListenAndServe(addr, mux))
}
