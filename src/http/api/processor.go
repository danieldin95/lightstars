package api

import (
	"github.com/danieldin95/lightstar/src/compute/libvirtc"
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/gorilla/mux"
	"net/http"
)

type Processor struct {
}

func (proc Processor) Router(router *mux.Router) {
	router.HandleFunc("/api/instance/{id}/processor", proc.GET).Methods("GET")
	router.HandleFunc("/api/instance/{id}/processor", proc.PUT).Methods("PUT")
}

func (proc Processor) GET(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, nil)
}

func (proc Processor) POST(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (proc Processor) PUT(w http.ResponseWriter, r *http.Request) {
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

func (proc Processor) DELETE(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}
