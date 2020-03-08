package api

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/http/schema"
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
	format := GetQueryOne(r, "format")

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
		for _, gra := range instance.Graphics {
			list.Items = append(list.Items, gra)
		}
		sort.SliceStable(list.Items, func(i, j int) bool {
			return list.Items[i].(schema.Graphics).Type < list.Items[j].(schema.Graphics).Type
		})
		list.Metadata.Size = len(list.Items)
		list.Metadata.Total = len(list.Items)
		ResponseJson(w, list)
	} else {
		if instance.XMLObj == nil {
			http.Error(w, "Get DescXML failed.", http.StatusInternalServerError)
			return
		}
		ResponseJson(w, instance.XMLObj.Devices.Graphics)
	}
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

func (gra Graphics) PUT(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (gra Graphics) DELETE(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dom.Free()

	dev, _ := GetArg(r, "dev")
	xml := libvirtc.NewDomainXMLFromDom(dom, true)
	if xml == nil {
		http.Error(w, "Cannot get domain's descXML", http.StatusInternalServerError)
		return
	}
	if xml.Devices.Graphics != nil {
		for _, gra := range xml.Devices.Graphics {
			if gra.Type != dev {
				continue
			}
			// found device
			flags := libvirtc.DOMAIN_DEVICE_MODIFY_PERSISTENT
			if active, _ := dom.IsActive(); !active {
				flags = libvirtc.DOMAIN_DEVICE_MODIFY_CONFIG
			}
			if err := dom.DetachDeviceFlags(gra.Encode(), flags); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	ResponseMsg(w, 0, "success")
}
