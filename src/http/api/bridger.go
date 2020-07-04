package api

import (
	"github.com/danieldin95/lightstar/src/network/libvirtn"
	"github.com/gorilla/mux"
	"net/http"
)

type Bridger struct {
}

func (br Bridger) Router(router *mux.Router) {
	router.HandleFunc("/api/bridge", br.GET).Methods("GET")
}

func (br Bridger) GET(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, libvirtn.BRIDGE.List())
}

func (br Bridger) POST(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, nil)
}

func (br Bridger) PUT(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, nil)
}

func (br Bridger) DELETE(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, nil)
}
