package api

import (
	"github.com/danieldin95/lightstar/src/compute"
	"github.com/danieldin95/lightstar/src/compute/libvirtc"
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"strings"
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
	format := GetQueryOne(r, "format")
	if format == "vv" {
		instance := compute.NewInstance(*dom)
		var spice schema.Graphics
		for _, gra := range instance.Graphics {
			if gra.Type == "vnc" || gra.Type == "spice" {
				spice = gra
			}
			if spice.Type == "spice" {
				break
			}
		}
		libstar.Debug("Graphics.Get %s, %s", r.URL.Hostname(), r.Host)
		if spice.Type != "vnc" && spice.Type != "spice" {
			return
		}
		context := ""
		filename := ""
		os := GetQueryOne(r, "os")
		if os == "linux" {
			if spice.Type == "spice" {
				context += "[connection]"
				context += "\nhost=" + strings.SplitN(r.Host, ":", 2)[0]
				context += "\nport=" + spice.Port
				context += "\npassword=" + spice.Password
				context += "\nfullscreen=1"
			} else {
				context += "[Connection]"
				context += "\nHost=" + strings.SplitN(r.Host, ":", 2)[0]
				context += "\nPort=" + spice.Port
				context += "\nPassword=" + spice.Password
				context += "\nFullscreen=1"
			}
			filename = instance.Name + "." + spice.Type
		} else {
			context += "[virt-viewer]"
			context += "\ntype=" + spice.Type
			context += "\nhost=" + strings.SplitN(r.Host, ":", 2)[0]
			context += "\nport=" + spice.Port
			context += "\npassword=" + spice.Password
			context += "\nfullscreen=1"
			filename = instance.Name + ".vv"
		}
		w.Header().Set("Content-Type", "application/x-download")
		w.Header().Set("Content-Disposition", "attachment;filename="+filename)
		_, _ = w.Write([]byte(context))

	} else {
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
	if conf.Password == "" {
		conf.Password = libstar.GenToken(32)
	}
	if conf.Listen == "" {
		conf.Listen = "0.0.0.0"
	}
	if conf.AutoPort == "yes" {
		conf.Port = "-1"
	}
	xmlObj := libvirtc.GraphicsXML{
		Type:     conf.Type,
		Listen:   conf.Listen,
		Port:     conf.Port,
		AutoPort: conf.AutoPort,
		Password: conf.Password,
	}

	flags := libvirtc.DomainDeviceModifyConfig
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
