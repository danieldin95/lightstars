package api

import (
	"github.com/danieldin95/lightstar/http/schema"
	"github.com/danieldin95/lightstar/storage"
	"github.com/danieldin95/lightstar/storage/libvirts"
	"github.com/gorilla/mux"
	"net/http"
)

type DataStore struct {
}

func (store DataStore) Router(router *mux.Router) {
	router.HandleFunc("/api/datastore", store.GET).Methods("GET")
	router.HandleFunc("/api/datastore", store.POST).Methods("POST")
	router.HandleFunc("/api/datastore/{id}", store.DELETE).Methods("DELETE")
}

func (store DataStore) GET(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, libvirts.DATASTOR.List())
}

func (store DataStore) POST(w http.ResponseWriter, r *http.Request) {
	data := &schema.DataStore{}
	if err := GetData(r, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if storage.PATH.IsDataStore(data.Name) {
		name := storage.PATH.GetStoreID(data.Name)
		path := storage.PATH.Unix(data.Name)
		if _, err := libvirts.CreatePool(name, path); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	ResponseMsg(w, 0, "success")
}

func (store DataStore) PUT(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, nil)
}

func (store DataStore) DELETE(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")

	if err := libvirts.RemovePool(uuid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}
