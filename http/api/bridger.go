package api

import (
	"github.com/gorilla/mux"
	"net/http"
)


type Bridger struct {

}

func (br Bridger) Router(router *mux.Router) {
	router.HandleFunc("/api/bridge", br.GET).Methods("GET")
}

func (br Bridger) GET(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, nil)
}
