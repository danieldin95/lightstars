package api

import (
	"github.com/danieldin95/lightstar/schema"
	"github.com/gorilla/mux"
	"net/http"
)

type Hyper struct {
}

func (h Hyper) Router(router *mux.Router) {
	router.HandleFunc("/api/hyper", h.GET).Methods("GET")
}

func (h Hyper) GET(w http.ResponseWriter, r *http.Request) {
	index := schema.Index{}
	index.Version = schema.NewVersion()
	index.Hyper = schema.NewHyper()
	index.User, _ = GetUser(r)
	ResponseJson(w, index)
}
