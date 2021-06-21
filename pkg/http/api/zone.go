package api

import (
	"github.com/danieldin95/lightstar/pkg/schema"
	"github.com/danieldin95/lightstar/pkg/service"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

type Zone struct {
}

func (z Zone) Router(router *mux.Router) {
	router.HandleFunc("/api/zone", z.Get).Methods("GET")
}

func (z Zone) Get(w http.ResponseWriter, r *http.Request) {
	hosts := make([]schema.Host, 0, 32)
	for h := range service.SERVICE.Zone.List() {
		if h == nil {
			break
		}
		hosts = append(hosts, *h)
	}
	sort.SliceStable(hosts, func(i, j int) bool {
		return hosts[i].Name < hosts[j].Name
	})
	ResponseJson(w, hosts)
}

func (z Zone) Post(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (z Zone) Put(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (z Zone) Delete(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}
