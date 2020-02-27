package api

import (
	"github.com/danieldin95/lightstar/storage"
	"github.com/danieldin95/lightstar/storage/libvirts"
	"github.com/gorilla/mux"
	"net/http"
)

type ISO struct {
}

func (iso ISO) Router(router *mux.Router) {
	router.HandleFunc("/api/iso", iso.GET).Methods("GET")
}

func (iso ISO) GET(w http.ResponseWriter, r *http.Request) {
	store := GetQueryOne(r, "datastore")
	if store == "" {
		store = "datastore@01"
	}
	path := storage.PATH.Unix(store)
	ResponseJson(w, libvirts.ISO.ListFiles(path))
}

func (iso ISO) POST(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (iso ISO) PUT(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (iso ISO) DELETE(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}
