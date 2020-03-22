package api

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/network/libvirtn"
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

	type Tgt struct {
		Name   string `json:"name"`
		Target string `json:"target"`
	}
	tgts := make([]Tgt, 0, 32)

	hyper, err := libvirtc.GetHyper()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Instance{}.GetByUser(&user, &list)
	sort.SliceStable(list.Items, func(i, j int) bool {
		return list.Items[i].(schema.Instance).Name < list.Items[j].(schema.Instance).Name
	})

	leases, err := libvirtn.ListLeases("")
	if err != nil {
		libstar.Warn("ProxyTcp.GET %s", err)
	}
	for _, item := range list.Items {
		inst := item.(schema.Instance)

		for _, graphic := range inst.Graphics {
			if graphic.Port == "" || graphic.Port == "-1" {
				continue
			}
			if libstar.IsDigit(graphic.Port) {

				tgts = append(tgts, Tgt{
					Name:   inst.Name,
					Target: hyper.Address + ":" + graphic.Port,
				})
			}
		}
		if leases != nil {
			for _, inf := range inst.Interfaces {
				libstar.Debug("ProxyTcp.GET %s", inf.Address)
				if le, ok := leases[inf.Address]; ok {
					tgts = append(tgts, Tgt{
						Name:   inst.Name,
						Target: le.IPAddr + ":22",
					}) // ssh
					tgts = append(tgts, Tgt{
						Name:   inst.Name,
						Target: le.IPAddr + ":3389",
					}) // rdp
					break
				}
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
