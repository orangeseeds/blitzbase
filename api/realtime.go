package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/orangeseeds/blitzbase/core"
)

type rtServer struct {
	app core.App
}

func (api *rtServer) handleRealtime(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			flusher, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "SSE is not supported", http.StatusInternalServerError)
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-type", "text/event-stream")

			sub := core.NewSubscriber(5)
			api.app.Publisher.Subscribe(sub, "sqlite_insert")

			for data := range sub.Listen() {
				msg, _ := data.FormatSSE()
				_, err := w.Write([]byte(msg))
				if err != nil {
					log.Print("Write Error:", err)
					return
				}
				flusher.Flush()
			}

			// for data := range api.app.Publisher.Notifier {
			// 	message, _ := formatServerSentEvent("newUserCreated", data)
			// 	_, err := w.Write([]byte(message))
			// 	if err != nil {
			// 		log.Println("Write Error: ", err)
			// 		break
			// 	}
			// 	flusher.Flush()
			// }
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (api *rtServer) createUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")

		var data struct {
			Name string
		}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			log.Println(err)
			return
		}

		db := api.app.Store.DB
		query := fmt.Sprintf("Insert into users (name) values ('%s')", data.Name)
		_, err = db.Exec(query)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		message := map[string]any{
			"status":  "success",
			"message": "successfully created new user " + data.Name,
		}
		json.NewEncoder(w).Encode(message)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}
