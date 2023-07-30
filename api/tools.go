package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/orangeseeds/blitzbase/store"
)

func printSQliteVersion(db store.Storage) {
	var result string
	row := db.DB.QueryRow("select sqlite_version()")
	_ = row.Scan(&result)

	log.Println(result)
}

func decodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func EncodeBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		EncodeBody(w, r, data)
	}
}

func RespondErr(w http.ResponseWriter, r *http.Request, status int, args ...interface{}) {
	Respond(w, r, status, map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"message": fmt.Sprint(args...),
		},
	})
}

func RespondHTTPErr(w http.ResponseWriter, r *http.Request, status int) {
	RespondErr(w, r, status, http.StatusText(status))
}

