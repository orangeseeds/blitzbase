package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
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
	var data struct {
		Username string `validate:"required"`
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
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

	query := fmt.Sprintf("Insert into users (username, email, password) values ('%s', '%s', '%s')", data.Username, data.Email, data.Password)
	_, err = a.app.Store.DB.Exec(query)
	if err != nil {
		RespondErr(w, r, 500, err.Error())
		return
	}

	message := map[string]any{
		"success": "true",
		"message": "successfully created new user " + data.Username,
	}

	Respond(w, r, 200, message)
}

func (a *authServer) login(w http.ResponseWriter, r *http.Request) {
	var (
		email    string
		password string
		data     struct {
			Email    string `validate:"required,email"`
			Password string `validate:"required"`
		}
	)

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		RespondErr(w, r, 500, err.Error())
		return
	}

	query := fmt.Sprintf("Select email, password from users where email='%s'", data.Email)
	row := a.app.Store.DB.QueryRow(query)
	err = row.Scan(&email, &password)
	if err != nil || password != data.Password {
		RespondErr(w, r, 403, "email or password didnot match the records.")
		return
	}
	_, tokenString, _ := tokenAuth.Encode(map[string]any{"record_email": data.Email})
	resp := map[string]any{
		"jwt":  tokenString,
		"data": data,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		RespondErr(w, r, 500, err.Error())
	}
}
