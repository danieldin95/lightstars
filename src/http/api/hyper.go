package api

import (
	"github.com/danieldin95/lightstar/src/compute"
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/gorilla/mux"
	"net/http"
)

type Hyper struct {
}

func (h Hyper) Router(router *mux.Router) {
	router.HandleFunc("/api/hyper", h.Get).Methods("GET")
}

func (h Hyper) Get(w http.ResponseWriter, r *http.Request) {
	index := schema.Index{}
	index.Version = schema.NewVersion()
	index.Hyper = compute.NewHyper()
	index.User, _ = GetUser(r)
	index.User.Password = ""
	ResponseJson(w, index)
}
