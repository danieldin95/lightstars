package api

import (
	"github.com/danieldin95/lightstar/pkg/network/libvirtn"
	"github.com/gorilla/mux"
	"net/http"
)

type Bridger struct {
}

func (br Bridger) Router(router *mux.Router) {
	router.HandleFunc("/api/bridge", br.Get).Methods("GET")
}

func (br Bridger) Get(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, libvirtn.BRIDGE.List())
}

func (br Bridger) Post(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, nil)
}

func (br Bridger) Put(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, nil)
}

func (br Bridger) Delete(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, nil)
}
