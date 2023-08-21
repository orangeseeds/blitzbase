package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/go-playground/validator/v10"
	"github.com/orangeseeds/blitzbase/core"
)

var (
	validate  *validator.Validate
	tokenAuth *jwtauth.JWTAuth
)

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	validate = validator.New()
}

type authServer struct {
	app core.App
}

func (a *authServer) Router() http.Handler {
	r := chi.NewRouter()
	r.Post("/login", a.login)
	r.Post("/register", a.register)
	return r
}

func (a *authServer) register(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var data struct {
		Username string `validate:"required"`
		Email    string `validate:"required,email"`
		Password string `validate:"required" json:",omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		RespondErr(w, r, 400, err.Error())
		return
	}

	if err := validate.Struct(data); err != nil {
		RespondErr(w, r, 400, err.Error())
		return
	}

	// query := fmt.Sprintf("Insert into users (username, email, password) values ('%s', '%s', '%s')", data.Username, data.Email, data.Password)
	row, err := a.app.Store.DB.Insert("users", dbx.Params{
		"username": data.Username,
		"email":    data.Email,
		"password": data.Password,
	}).Execute()
	if err != nil {
		RespondErr(w, r, 500, err.Error())
		return
	}

	id, _ := row.LastInsertId()

	data.Password = ""
	message := map[string]any{
		"success": "true",
		"message": fmt.Sprintf("successfully created new user id:%d", id),
		"data":    data,
	}

	Respond(w, r, 200, message)
}

func (a *authServer) login(w http.ResponseWriter, r *http.Request) {
	type req struct {
		ID       int `json:",omitempty"`
		Username string
		Email    string `validate:"required,email"`
		Password string `validate:"required" json:",omitempty"`
	}

	var (
		data   req
		dbData req
	)

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		RespondErr(w, r, 500, err.Error())
		return
	}

	q := a.app.Store.DB.Select("id", "username", "email", "password").From("users").Where(dbx.HashExp{"email": data.Email})
	err = q.One(&dbData)
	if err != nil || dbData.Password != data.Password {
		RespondErr(w, r, 403, "email or password didnot match the records.")
		return
	}
	_, tokenString, _ := tokenAuth.Encode(map[string]any{"record_email": dbData.ID})
	dbData.Password = ""
	dbData.ID = 0
	resp := map[string]any{
		"jwt":  tokenString,
		"data": dbData,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		RespondErr(w, r, 500, err.Error())
	}
}
