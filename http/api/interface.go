package api

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/http/schema"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

type Interface struct {
}

func Interface2XML(conf *schema.Interface) (*libvirtc.InterfaceXML, error) {
	xml := libvirtc.InterfaceXML{
		Type: "bridge",
		Source: libvirtc.InterfaceSourceXML{
			Bridge: conf.Source,
		},
		Model: libvirtc.InterfaceModelXML{
			Type: conf.Model,
		},
		Address: &libvirtc.AddressXML{
			Type:     "pci",
			Domain:   libvirtc.PCI_DOMAIN,
			Bus:      libvirtc.PCI_INTERFACE_BUS,
			Slot:     conf.Seq,
			Function: libvirtc.PCI_FUNC,
		},
	}
	if conf.Type == "openvswitch" {
		xml.VirtualPort = &libvirtc.InterfaceVirtualPortXML{
			Type: conf.Type,
		}
	}
	return &xml, nil
}

func (in Interface) Router(router *mux.Router) {
	router.HandleFunc("/api/instance/{id}/interface", in.GET).Methods("GET")
	router.HandleFunc("/api/instance/{id}/interface", in.POST).Methods("POST")
	router.HandleFunc("/api/instance/{id}/interface/{dev}", in.GET).Methods("GET")
	router.HandleFunc("/api/instance/{id}/interface/{dev}", in.DELETE).Methods("DELETE")
}

func (in Interface) GET(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	dev, ok := GetArg(r, "dev")
	format := GetQueryOne(r, "format")
	if !ok {
		dom, err := libvirtc.LookupDomainByUUIDString(uuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		defer dom.Free()
		instance := schema.NewInstance(*dom)
		if format == "schema" {
			list := schema.List{
				Items:    make([]interface{}, 0, 32),
				Metadata: schema.MetaData{},
			}
			for _, intf := range instance.Interfaces {
				list.Items = append(list.Items, intf)
			}
			sort.SliceStable(list.Items, func(i, j int) bool {
				return list.Items[i].(schema.Interface).Device < list.Items[j].(schema.Interface).Device
			})
			list.Metadata.Size = len(list.Items)
			list.Metadata.Total = len(list.Items)
			ResponseJson(w, list)
		} else {
			if instance.XMLObj == nil {
				http.Error(w, "Get DescXML failed.", http.StatusInternalServerError)
				return
			}
			ResponseJson(w, instance.XMLObj.Devices.Interfaces)
		}
		return
	}
	ResponseMsg(w, 0, dev)
}

func (in Interface) POST(w http.ResponseWriter, r *http.Request) {
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

	xmlObj, err := Interface2XML(conf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	libstar.Debug("Interface.POST: %s", xmlObj.Encode())
	flags := libvirtc.DOMAIN_DEVICE_MODIFY_PERSISTENT
	if active, _ := dom.IsActive(); !active {
		flags = libvirtc.DOMAIN_DEVICE_MODIFY_CONFIG
	}
	if err := dom.AttachDeviceFlags(xmlObj.Encode(), flags); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}

func (in Interface) PUT(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (in Interface) DELETE(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Cannot get domain's descXML", http.StatusInternalServerError)
		return
	}

	if xml.Devices.Interfaces != nil {
		for _, intf := range xml.Devices.Interfaces {
			if intf.Mac.Address != address {
				continue
			}
			// found deivice
			flags := libvirtc.DOMAIN_DEVICE_MODIFY_PERSISTENT
			if active, _ := dom.IsActive(); !active {
				flags = libvirtc.DOMAIN_DEVICE_MODIFY_CONFIG
			}
			if err := dom.DetachDeviceFlags(intf.Encode(), flags); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	ResponseMsg(w, 0, "success")
}
