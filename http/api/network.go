package api

import (
	"github.com/danieldin95/lightstar/network/libvirtn"
	"github.com/gorilla/mux"
	"net/http"
)

type Network struct {
}

func (net Network) Router(router *mux.Router) {
	router.HandleFunc("/api/network", net.GET).Methods("GET")
}

func (net Network) GET(w http.ResponseWriter, r *http.Request) {
	nets, _ := libvirtn.ListNetworks()
	ResponseJson(w, nets)
}
