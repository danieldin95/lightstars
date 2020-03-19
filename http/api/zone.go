package api

import (
	"github.com/danieldin95/lightstar/schema"
	"github.com/danieldin95/lightstar/service"
	"github.com/gorilla/mux"
	"net/http"
)

type Zone struct {
}

func (z Zone) Router(router *mux.Router) {
	router.HandleFunc("/api/zone", z.GET).Methods("GET")
}

func (z Zone) GET(w http.ResponseWriter, r *http.Request) {
	hosts := make([]schema.Host, 0, 32)
	for h := range service.SERVICE.Zone.List() {
		if h == nil {
			break
		}
		hosts = append(hosts, *h)
	}
	ResponseJson(w, hosts)
}

func (z Zone) POST(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (z Zone) PUT(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (z Zone) DELETE(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}
