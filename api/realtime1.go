package api

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"time"
//
// 	"github.com/go-chi/chi/v5"
// 	"github.com/orangeseeds/blitzbase/core"
// 	"github.com/orangeseeds/blitzbase/store"
// )
//
// type rtServer struct {
// 	app core.App
// }
//
// func (api *rtServer) Router() http.Handler {
// 	r := chi.NewRouter()
// 	r.Get("/", api.handleRealtime)
// 	r.Post("/", api.setSubscriptions)
// 	return r
// }
//
// func (api *rtServer) handleRealtime(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-type", "text/event-stream")
// 	flusher, ok := w.(http.Flusher)
// 	if !ok {
// 		http.Error(w, "SSE is not supported", http.StatusInternalServerError)
// 	}
//
// 	sub := store.NewSubscriber(1)
// 	api.app.Store.Publisher.Subscribe(sub)
//
// 	http.SetCookie(w, &http.Cookie{
// 		Name:   "subscriber_id",
// 		Value:  sub.ID(),
// 		Path:   "/",
// 		MaxAge: 3600,
// 	})
//
// 	flusher.Flush()
//
// 	life := time.Second * 10
// 	start := time.Now()
// 	timer := time.NewTimer(life)
//
// 	handlerEvent := func(w http.ResponseWriter, flusher http.Flusher, e store.DBHookEvent) error {
// 		msg, _ := e.Message.FormatSSE()
// 		// log.Println(msg, "formatted")
// 		_, err := w.Write([]byte(msg))
// 		if err != nil {
// 			log.Print("Write Error:", err)
// 			return err
// 		}
// 		flusher.Flush()
//
// 		return nil
// 	}
//
// 	for {
// 		select {
// 		case e := <-sub.Listen():
// 			{
// 				timer.Reset(life)
// 				err := handlerEvent(w, flusher, e)
// 				if err != nil {
// 					http.Error(w, err.Error(), 500)
// 					return
// 				}
// 			}
// 		case t := <-timer.C:
// 			{
// 				log.Printf("SSE connection closed due to inactivity after %d, ID: %s", int(t.Sub(start).Seconds()), sub.ID())
// 				api.app.Store.Publisher.Unsubscribe(sub.ID())
// 				// sub.Close()
// 				return
// 			}
// 		case <-r.Context().Done():
// 			{
// 				if r.Context().Err() == context.Canceled {
// 					log.Printf("SSE connection closed by client, ID: %s", sub.ID())
// 				} else {
// 					log.Printf("SSE connection closed, ID: %s, err: %s", sub.ID(), r.Context().Err().Error())
// 				}
// 				api.app.Store.Publisher.Unsubscribe(sub.ID())
// 				// sub.Close()
// 				return
// 			}
// 		}
// 	}
// }
//
// func (api *rtServer) setSubscriptions(w http.ResponseWriter, r *http.Request) {
//
// 	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
//
// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	cookie, err := r.Cookie("subscriber_id")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	topic := r.Form.Get("collection")
// 	log.Println(topic)
//
// 	sub, err := api.app.Store.Publisher.SubByID(cookie.Value)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	if topic == "none" {
// 		sub.AddTopics(topic)
// 		api.app.Store.Publisher.Unsubscribe(sub.ID())
// 	} else {
// 		sub.AddTopics(topic)
// 		api.app.Store.Publisher.Subscribe(sub)
// 	}
// }
