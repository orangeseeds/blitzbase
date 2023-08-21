package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/orangeseeds/blitzbase/core"
)

type collectionServer struct {
	app core.App
}

func (api collectionServer) Router() http.Handler {
	r := chi.NewRouter()
	r.Post("/", api.createUser)
	return r
}

func (api *collectionServer) createUser(w http.ResponseWriter, r *http.Request) {

	var data struct {
		Username string
		Email    string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println(err)
		return
	}

	// query := fmt.Sprintf("Insert into users (username, email, password) values ('%s', '%s', '%s')", data.Username, data.Email, data.Password)
	res, err := api.app.Store.DB.Insert("users", dbx.Params{
		"username":  data.Username,
		"email":     data.Email,
		"passsword": data.Password,
	}).Execute()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	message := map[string]any{
		"status":  "success",
		"message": fmt.Sprintf("successfully created new user %d", id),
	}
	json.NewEncoder(w).Encode(message)
}
