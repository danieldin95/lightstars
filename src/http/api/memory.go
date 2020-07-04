package api

import (
	"github.com/danieldin95/lightstar/src/compute/libvirtc"
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/gorilla/mux"
	"net/http"
)

type Memory struct {
}

func (mem Memory) Router(router *mux.Router) {
	router.HandleFunc("/api/instance/{id}/memory", mem.GET).Methods("GET")
	router.HandleFunc("/api/instance/{id}/memory", mem.PUT).Methods("PUT")
}

func (mem Memory) GET(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, nil)
}

func (mem Memory) POST(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (mem Memory) PUT(w http.ResponseWriter, r *http.Request) {
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
	size := libstar.ToKiB(conf.Size, conf.Unit)
	if err := dom.SetMemoryFlags(size, libvirtc.DOMAIN_MEM_MAXIMUM); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := dom.SetMemoryFlags(size, libvirtc.DOMAIN_MEM_CONFIG); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "")
}

func (mem Memory) DELETE(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}
