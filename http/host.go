package http

import (
	"github.com/danieldin95/lightstar/http/api"
	"github.com/danieldin95/lightstar/http/service"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/gorilla/mux"
	"net/http"
)

type Host struct {
}

func (h Host) Router(router *mux.Router) {
	router.PathPrefix("/host/{name}/api/").HandlerFunc(h.Handle).Methods("GET")
	router.PathPrefix("/host/{name}/api/instance").HandlerFunc(h.Handle).Methods("POST")
	router.PathPrefix("/host/{name}/api/").HandlerFunc(h.Handle).Methods("PUT")
	router.PathPrefix("/host/{name}/api/").HandlerFunc(h.Handle).Methods("DELETE")
}

func (h Host) Handle(w http.ResponseWriter, r *http.Request) {
	name, _ := api.GetArg(r, "name")
	node := service.SERVICE.Zone.Get(name)
	if node == nil {
		http.Error(w, "not found host", http.StatusNotFound)
		return
	}
	libstar.Debug("Host.Handle %s", node)
	r.Header.Del("cookie")
	pri := &libstar.ProxyUrl{
		Proxy: libstar.Proxy{
			Prefix: "/host/" + name,
			Server: node.Url,
			Auth: libstar.Auth{
				Type:     "basic",
				Username: node.Username,
				Password: node.Password,
			},
		},
	}
	pri.Initialize()
	pri.Handler(w, r)
}
