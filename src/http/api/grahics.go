package api

import (
	"github.com/danieldin95/lightstar/src/compute"
	"github.com/danieldin95/lightstar/src/compute/libvirtc"
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

type Graphics struct {
}

func (gra Graphics) Router(router *mux.Router) {
	router.HandleFunc("/api/instance/{id}/graphics", gra.GET).Methods("GET")
	router.HandleFunc("/api/instance/{id}/graphics", gra.POST).Methods("POST")
	router.HandleFunc("/api/instance/{id}/graphics/{dev}", gra.DELETE).Methods("DELETE")
}

func (gra Graphics) GET(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")

	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer dom.Free()
	instance := compute.NewInstance(*dom)
	list := schema.ListGraphics{
		Items: make([]schema.Graphics, 0, 32),
	}
	for _, gra := range instance.Graphics {
		list.Items = append(list.Items, gra)
	}
	sort.SliceStable(list.Items, func(i, j int) bool {
		return list.Items[i].Type > list.Items[j].Type
	})
	list.Metadata.Size = len(list.Items)
	list.Metadata.Total = len(list.Items)
	ResponseJson(w, list)
}

func (gra Graphics) POST(w http.ResponseWriter, r *http.Request) {
	conf := &schema.Graphics{}
	if err := GetData(r, conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uuid, _ := GetArg(r, "id")
	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dom.Free()

	xmlObj := libvirtc.GraphicsXML{
		Type:   conf.Type,
		Listen: conf.Listen,
		//Port:     conf.Port,
		AutoPort: "yes",
		Password: conf.Password,
	}
	flags := libvirtc.DomainDeviceModifyPersistent
	if active, _ := dom.IsActive(); !active {
		flags = libvirtc.DomainDeviceModifyConfig
	}
	if err := dom.AttachDeviceFlags(xmlObj.Encode(), flags); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}

func (gra Graphics) PUT(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (gra Graphics) DELETE(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "success")
}
