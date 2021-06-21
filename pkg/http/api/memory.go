package api

import (
	"github.com/danieldin95/lightstar/pkg/compute/libvirtc"
	"github.com/danieldin95/lightstar/pkg/schema"
	"github.com/gorilla/mux"
	"net/http"
)

type Memory struct {
}

func (mem Memory) Router(router *mux.Router) {
	router.HandleFunc("/api/instance/{id}/memory", mem.Get).Methods("GET")
	router.HandleFunc("/api/instance/{id}/memory", mem.Put).Methods("PUT")
}

func (mem Memory) Get(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, nil)
}

func (mem Memory) Post(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (mem Memory) Put(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dom.Free()

	conf := &schema.Memory{}
	if err := GetData(r, conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := dom.SetMemory(conf.Size, conf.Unit); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "")
}

func (mem Memory) Delete(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}
