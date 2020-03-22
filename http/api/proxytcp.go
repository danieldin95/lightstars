package api

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/schema"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

type ProxyTcp struct {
}

func (pro ProxyTcp) Router(router *mux.Router) {
	router.HandleFunc("/api/proxy/tcp", pro.GET).Methods("GET")
}

func (pro ProxyTcp) GET(w http.ResponseWriter, r *http.Request) {
	user, _ := GetUser(r)
	list := schema.List{
		Items: make([]interface{}, 0, 32),
	}
	tgts := make([]string, 0, 32)

	hyper, err := libvirtc.GetHyper()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Instance{}.GetByUser(&user, &list)
	sort.SliceStable(list.Items, func(i, j int) bool {
		return list.Items[i].(schema.Instance).Name < list.Items[j].(schema.Instance).Name
	})
	for _, item := range list.Items {
		inst := item.(schema.Instance)

		for _, graphic := range inst.Graphics {
			if graphic.Port == "" || graphic.Port == "-1" {
				continue
			}
			if libstar.IsDigit(graphic.Port) {
				tgts = append(tgts, hyper.Address+":"+graphic.Port)
			}
		}
	}
	ResponseJson(w, tgts)
}

func (pro ProxyTcp) POST(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (pro ProxyTcp) PUT(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (pro ProxyTcp) DELETE(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}
