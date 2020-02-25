package api

import (
	"github.com/danieldin95/lightstar/storage/libvirts"
	"github.com/gorilla/mux"
	"net/http"
)

type DataStore struct {
}

func (store DataStore) Router(router *mux.Router) {
	router.HandleFunc("/api/datastore", store.GET).Methods("GET")
}

func (store DataStore) GET(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, libvirts.DATASTOR.List())
}
