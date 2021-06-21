package api

import (
	"github.com/danieldin95/lightstar/pkg/compute/libvirtc"
	"github.com/danieldin95/lightstar/pkg/schema"
	"github.com/gorilla/mux"
	"net/http"
)

type Processor struct {
}

func (proc Processor) Router(router *mux.Router) {
	router.HandleFunc("/api/instance/{id}/processor", proc.Get).Methods("GET")
	router.HandleFunc("/api/instance/{id}/processor", proc.Put).Methods("PUT")
}

func (proc Processor) Get(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, nil)
}

func (proc Processor) Post(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (proc Processor) Put(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dom.Free()
	conf := &schema.Processor{}
	if err := GetData(r, conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := dom.SetCpu(conf.Cpu, conf.Mode); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "")
}

func (proc Processor) Delete(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}
