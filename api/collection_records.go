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
	// r.Post("/", api.createUser)

	r.Get("/", api.handleCollList)
	r.Post("/", api.handleCollCreate)

	r.Get("/{collection}/records/", api.handleList)
	r.Get("/{collection}/records/{id}", api.handleView)
	r.Post("/{collection}/records", api.handleCreate)
	r.Patch("/{collection}/records", api.handleUpdate)
	r.Delete("/{collection}/records/{id}", api.handleDelete)
	return r
}

func (api *collectionServer) handleList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// var users []struct {
	// 	Username string `json:"username"`
	// 	Email    string `json:"email"`
	// }

	users := []dbx.NullStringMap{}
	user := dbx.NullStringMap{}

	q := api.app.Store.DB.Select("*").From(chi.URLParam(r, "collection"))

	rows, err := q.Rows()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for rows.Next() {
		err := rows.ScanMap(user)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}
func (api *collectionServer) handleView(w http.ResponseWriter, r *http.Request)   {}
func (api *collectionServer) handleCreate(w http.ResponseWriter, r *http.Request) {}
func (api *collectionServer) handleUpdate(w http.ResponseWriter, r *http.Request) {}
func (api *collectionServer) handleDelete(w http.ResponseWriter, r *http.Request) {}

func (api *collectionServer) createUser(w http.ResponseWriter, r *http.Request) {

	// var data struct {
	// 	Username string
	// 	Email    string
	// 	Password string
	// }
	//
	// err := json.NewDecoder(r.Body).Decode(&data)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println(err)
		return
	}

	// query := fmt.Sprintf("Insert into users (username, email, password) values ('%s', '%s', '%s')", data.Username, data.Email, data.Password)
	res, err := api.app.Store.DB.Insert("users", dbx.Params{
		"username": r.Form.Get("username"),
		"email":    r.Form.Get("email"),
		"password": r.Form.Get("password"),
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
