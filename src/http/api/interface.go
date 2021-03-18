package api

import (
	"github.com/danieldin95/lightstar/src/compute"
	"github.com/danieldin95/lightstar/src/compute/libvirtc"
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

type Interface struct {
}

func (in Interface) Router(router *mux.Router) {
	router.HandleFunc("/api/instance/{id}/interface", in.Get).Methods("GET")
	router.HandleFunc("/api/instance/{id}/interface", in.Post).Methods("POST")
	router.HandleFunc("/api/instance/{id}/interface/{dev}", in.Get).Methods("GET")
	router.HandleFunc("/api/instance/{id}/interface/{dev}", in.Delete).Methods("DELETE")
}

func (in Interface) Get(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	dev, ok := GetArg(r, "dev")
	if !ok {
		dom, err := libvirtc.LookupDomainByUUIDString(uuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		defer dom.Free()
		instance := compute.NewInstance(*dom)
		list := schema.ListInterface{
			Items: make([]schema.Interface, 0, 32),
		}
		for _, inf := range instance.Interfaces {
			list.Items = append(list.Items, inf)
		}
		sort.SliceStable(list.Items, func(i, j int) bool {
			return list.Items[i].Device < list.Items[j].Device
		})
		list.Metadata.Size = len(list.Items)
		list.Metadata.Total = len(list.Items)
		ResponseJson(w, list)
		return
	}
	ResponseMsg(w, 0, dev)
}

func (in Interface) Post(w http.ResponseWriter, r *http.Request) {
	conf := &schema.Interface{}
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

	xmlObj := Interface2XML(conf.Source, conf.Model, conf.Seq, conf.Type, "", "")
	libstar.Debug("Interface.Post: %s", libstar.XML.Encode(xmlObj))

	flags := libvirtc.DomainDeviceModifyPersistent
	if active, _ := dom.IsActive(); !active {
		flags = libvirtc.DomainDeviceModifyConfig
	}
	if err := dom.AttachDeviceFlags(libstar.XML.Encode(xmlObj), flags); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}

func (in Interface) Put(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (in Interface) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dom.Free()

	address, _ := GetArg(r, "dev")
	xml := libvirtc.NewDomainXMLFromDom(dom, true)
	if xml == nil {
		http.Error(w, "Cannot get domain's descXML", http.StatusNotFound)
		return
	}
	if xml.Devices.Interfaces == nil {
		http.Error(w, "Cannot get domain's interface", http.StatusNotFound)
		return
	}
	for _, port := range xml.Devices.Interfaces {
		if port.Mac.Address != address {
			continue
		}
		// found device
		flags := libvirtc.DomainDeviceModifyPersistent
		if active, _ := dom.IsActive(); !active {
			flags = libvirtc.DomainDeviceModifyConfig
		}
		if err := dom.DetachDeviceFlags(libstar.XML.Encode(&port), flags); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	ResponseMsg(w, 0, "success")
}
