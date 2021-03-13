package api

import (
	"github.com/danieldin95/lightstar/src/storage"
	"github.com/gorilla/mux"
	"net/http"
)

type ISO struct {
}

func (iso ISO) Router(router *mux.Router) {
	router.HandleFunc("/api/iso", iso.Get).Methods("GET")
}

func (iso ISO) Get(w http.ResponseWriter, r *http.Request) {
	store := GetQueryOne(r, "datastore")
	if store == "" {
		store = "datastore@01"
	}
	path := storage.PATH.Unix(store)
	ResponseJson(w, storage.ISO.ListFiles(path))
}
