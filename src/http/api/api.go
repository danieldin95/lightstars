package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Api interface {
	Router(router *mux.Router)
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
}
