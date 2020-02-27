package api

import (
	"github.com/danieldin95/lightstar/http/schema"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/network/libvirtn"
	"github.com/gorilla/mux"
	"net/http"
)

type Network struct{}

func NetworkConf2XML(conf schema.Network) libvirtn.NetworkXML {
	xmlObj := libvirtn.NetworkXML{
		Name: conf.Name,
		Bridge: libvirtn.BridgeXML{
			Name: conf.Name,
		},
		Forward: libvirtn.ForwardXML{
			Mode: conf.Mode,
		},
	}
	if conf.Mode == "nat" {
		xmlObj.Bridge.Stp = "on"
		xmlObj.Bridge.Delay = "0"
	}
	if conf.Address != "" {
		addr, mask := libstar.ParseIP4Netmask(conf.Address, conf.Prefix)
		if addr != nil && mask != nil {
			xmlObj.IPv4 = libvirtn.IPv4XML{
				Address: addr.String(),
				Netmask: mask.String(),
			}
		}
	}
	return xmlObj
}

func (net Network) Router(router *mux.Router) {
	router.HandleFunc("/api/network", net.GET).Methods("GET")
	router.HandleFunc("/api/network", net.POST).Methods("POST")
	router.HandleFunc("/api/network/{id}", net.DELETE).Methods("DELETE")
}

func (net Network) GET(w http.ResponseWriter, r *http.Request) {
	nets, _ := libvirtn.ListNetworks()
	ResponseJson(w, nets)
}

func (net Network) POST(w http.ResponseWriter, r *http.Request) {
	conf := schema.Network{}
	if err := GetData(r, &conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	xmlObj := NetworkConf2XML(conf)
	libstar.Debug("Network.POST %s", xmlObj.Encode())
	hyper, err := libvirtn.GetHyper()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	netvir, err := hyper.NetworkDefineXML(xmlObj.Encode())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer netvir.Free()
	if err := netvir.Create(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := netvir.SetAutostart(true); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "")
}

func (net Network) PUT(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (net Network) DELETE(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	hyper, err := libvirtn.GetHyper()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	netvir, err := hyper.LookupNetwork(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer netvir.Free()
	if err := netvir.Destroy(); err != nil {
		libstar.Warn("Network.DELETE %s", err)
	}
	if err := netvir.Undefine(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "")
}
