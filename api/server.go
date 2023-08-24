package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	// Setting up the fileserver
	{
		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, "views/static/"))
		fileServer(mux, "/static", filesDir)
	}

	// Setting up the API routes
	apiGroup := chi.NewRouter()
	{
		// apiGroup.Use(middleware.AllowContentType("application/json"))
		// apiGroup.Use(middleware.AllowContentType("application/x-www-form-urlencoded"))
		apiGroup.Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				// w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

				h.ServeHTTP(w, r)
			})
		})

		rtApp := rtServer{app: app}
		authApp := authServer{app: app}
		collectionApp := collectionServer{app: app}

		apiGroup.Mount("/realtime", rtApp.Router())
		apiGroup.Mount("/auth", authApp.Router())
		apiGroup.Mount("/collection", collectionApp.Router())
	}


	// Setting up the view routes
	viewGroup := viewServer{app: app}
	mux.Mount("/dashboard", viewGroup.Router())
	mux.Mount("/api", apiGroup)

	log.Println("serving on port :3300")
	log.Fatal(http.ListenAndServe(addr, mux))
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
