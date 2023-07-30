package api

import (
	"log"
	"net/http"

	"github.com/orangeseeds/blitzbase/core"
)

func logger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"%v %v %v\n",
			r.Method, r.URL, r.RemoteAddr,
		)
		fn(w, r)
	}
}

func Serve(app core.App, addr string) {

	mux := http.NewServeMux()
	rtApp := rtServer{
		app: app,
	}
	mux.HandleFunc("/realtime", logger(rtApp.handleRealtime))
	mux.HandleFunc("/register", logger(rtApp.createUser))

	log.Println("serving on port :3300")
	log.Fatal(http.ListenAndServe(addr, mux))
}
