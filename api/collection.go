package api

import (
	"encoding/json"
	"net/http"

	"github.com/orangeseeds/blitzbase/utils/forms"
	"github.com/orangeseeds/blitzbase/utils/migrations"
)

func (api *collectionServer) handleCollList(w http.ResponseWriter, r *http.Request) {}

func (api *collectionServer) handleCollCreate(w http.ResponseWriter, r *http.Request) {

	var form forms.CreateCollectionForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := form.IsValid(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = migrations.CreateNewTable(api.app.Store, form.ToCollection())
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = migrations.AddCollectionRecord(api.app.Store, form.ToCollection())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(form.ToCollection())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
