package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Api interface {
	Router(router *mux.Router)
	GET(w http.ResponseWriter, r *http.Request)
	POST(w http.ResponseWriter, r *http.Request)
	DELETE(w http.ResponseWriter, r *http.Request)
	PUT(w http.ResponseWriter, r *http.Request)
}
