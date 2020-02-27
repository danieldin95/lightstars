package api

import (
	"github.com/danieldin95/lightstar/http/schema"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/network/libvirtn"
	"github.com/gorilla/mux"
	"net"
	"net/http"
)

type Network struct {
}

func IsUnicast(address string) bool {
	addr := net.ParseIP(address)
	if addr == nil {
		return false
	}
	if addr.IsMulticast() || addr.IsLoopback() || addr.IsUnspecified() {
		return false
	}
	return true
}
func Network2XML(conf schema.Network) libvirtn.NetworkXML {
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

	if IsUnicast(conf.Address) && conf.Mode == "nat" {
		addr, mask := libstar.ParseIP4Netmask(conf.Address, conf.Prefix)
		if addr != nil && mask != nil {
			xmlObj.IPv4 = &libvirtn.IPv4XML{
				Address: conf.Address,
				Netmask: mask.String(),
			}
		}
		if conf.DHCP != "no" {
			start, end := libstar.IP4Network2Range(addr, mask)
			xmlObj.IPv4.DHCP = &libvirtn.DHCPXML{
				Range: libvirtn.DHCPRangeXML{
					Start: start.String(),
					End:   end.String(),
				},
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
	xmlObj := Network2XML(conf)
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
