package api

import (
	"github.com/danieldin95/lightstar/pkg/compute"
	"github.com/danieldin95/lightstar/pkg/network/libvirtn"
	"github.com/danieldin95/lightstar/pkg/schema"
	"github.com/danieldin95/lightstar/pkg/storage/libvirts"
	"github.com/gorilla/mux"
	"net/http"
)

type Hyper struct {
}

func (h Hyper) Router(router *mux.Router) {
	router.HandleFunc("/api/hyper", h.Get).Methods("GET")
	router.HandleFunc("/api/hyper/statics", h.Statics).Methods("GET")
}

func (h Hyper) Get(w http.ResponseWriter, r *http.Request) {
	index := schema.Index{}
	index.Version = schema.NewVersion()
	index.Hyper = compute.NewHyper()
	index.User, _ = GetUser(r)
	index.User.Password = ""
	ResponseJson(w, index)
}

func (h Hyper) Statics(w http.ResponseWriter, r *http.Request) {
	sts := schema.Statics{}
	list := schema.ListInstance{
		Items: make([]schema.Instance, 0, 32),
	}
	Instance{}.GetByUser(nil, &list)
	for _, obj := range list.Items {
		switch obj.State {
		case "running":
			sts.Instance.Active++
		case "shutdown":
			sts.Instance.Inactive++
		default:
			sts.Instance.Unknown++
		}
		for _, port := range obj.Interfaces {
			if port.Device == "" {
				sts.Ports.Inactive++
			} else {
				sts.Ports.Active++
			}
		}
		sts.Ports.Total += len(obj.Interfaces)
	}
	sts.Instance.Total = len(list.Items)
	if objs, err := libvirtn.ListNetworks(); err == nil {
		for _, obj := range objs {
			if ok, err := obj.IsActive(); err == nil {
				if ok {
					sts.Network.Active++
				}
			} else {
				sts.Network.Unknown++
			}
		}
		sts.Network.Total = len(objs)
	}
	if objs, err := libvirts.ListPools(); err == nil {
		for _, obj := range objs {
			if ok, err := obj.IsActive(); err == nil {
				if ok {
					sts.DataStore.Active++
				}
			} else {
				sts.DataStore.Unknown++
			}
		}
		sts.DataStore.Total = len(objs)
	}
	ResponseJson(w, sts)
}
