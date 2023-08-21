package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/orangeseeds/blitzbase/core"
	"github.com/orangeseeds/blitzbase/store"
)

type rtServer struct {
	app core.App
}

func (api *rtServer) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", api.handleRealtime)
	r.Post("/", api.setSubscriptions)
	return r
}

func (api *rtServer) handleRealtime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/event-stream")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE is not supported", http.StatusInternalServerError)
	}

	sub := store.NewSubscriber(5)
	// sub.AddTopics("create")

	api.app.Store.Publisher.Subscribe(sub)

	http.SetCookie(w, &http.Cookie{
		Name:   "subID",
		Value:  sub.ID(),
		Path:   "/",
		MaxAge: 3600,
	})
	flusher.Flush()

	life := time.Minute * 5
	ctx, cancel := context.WithTimeout(context.Background(), life)
	defer cancel()

	for {
		select {
		case data, ok := <-sub.Listen():
			{
				if !ok {
					break
				}
				msg, _ := data.FormatSSE(sub.ID())
				_, err := w.Write([]byte(msg))
				if err != nil {
					log.Print("Write Error:", err)
					return
				}
				flusher.Flush()

				go func(ctx context.Context, cancel context.CancelFunc) {
					ctx, cancel = context.WithTimeout(context.Background(), life)
				}(ctx, cancel)
			}
		case <-ctx.Done():
			{
				if ctx.Err() == context.DeadlineExceeded {
					log.Printf("SSE connection closed due to inactivity, ID: %s", sub.ID())
					return
				} else if ctx.Err() == context.Canceled {
					log.Printf("SSE connection cancelled, ID: %s", sub.ID())
					return
				} else {
					continue
				}
			}
		}

	}

}

func (api *rtServer) setSubscriptions(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	// var reqData struct {
	// 	SubID  string
	// 	Topics []string
	// }

	err := r.ParseForm()

	// err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("subID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//
	topic := r.PostForm.Get("collection")

	sub, err := api.app.Store.Publisher.SubByID(cookie.Value)
	json.NewEncoder(w).Encode(map[string]any{
		"val":  cookie.Value,
		"Subs": sub.ID(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if topic == "none" {
		sub.Deactivate()
	} else {
		sub.Activate()
	}

	sub.AddTopics(topic)
}
