package api

import (
	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/network"
	"github.com/danieldin95/lightstar/pkg/schema"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"sort"
)

type Network struct {
}

func IsUniCast(address string) bool {
	addr := net.ParseIP(address)
	if addr == nil {
		return false
	}
	if addr.IsMulticast() || addr.IsLoopback() || addr.IsUnspecified() {
		return false
	}
	return true
}

func Network2XML(conf schema.Network) *network.NetworkXML {
	xmlObj := &network.NetworkXML{
		Name: conf.Name,
		Bridge: network.BridgeXML{
			Name: conf.Bridge,
		},
	}
	if conf.Mode != "" {
		xmlObj.Forward = &network.ForwardXML{
			Mode: conf.Mode,
		}
	}
	if conf.Mode == "nat" {
		xmlObj.Bridge.Stp = "on"
		xmlObj.Bridge.Delay = "0"
	}
	if conf.Address != "" {
		xmlObj.IPv4 = &network.IPv4XML{
			Address: conf.Address,
		}
		if conf.Prefix != "" {
			xmlObj.IPv4.Prefix = conf.Prefix
		}
		if conf.Netmask != "" {
			xmlObj.IPv4.Netmask = conf.Netmask
		}
		// DHCP address range
		xmlObj.IPv4.DHCP = &network.DHCPXML{
			Range: make([]network.DHCPRangeXML, 0, 32),
		}
		for _, addr := range conf.Range {
			xmlObj.IPv4.DHCP.Range = append(xmlObj.IPv4.DHCP.Range,
				network.DHCPRangeXML{
					Start: addr.Start,
					End:   addr.End,
				})
		}
	}
	if conf.Type == "openvswitch" {
		xmlObj.VirtualPort = &network.VirtualPortXML{
			Type: conf.Type,
		}
	}
	return xmlObj
}

func (net Network) Router(router *mux.Router) {
	router.HandleFunc("/api/network", net.Get).Methods("GET")
	router.HandleFunc("/api/network/{id}", net.Get).Methods("GET")
	router.HandleFunc("/api/network", net.Post).Methods("POST")
	router.HandleFunc("/api/network/{id}/start", net.Start).Methods("PUT")
	router.HandleFunc("/api/network/{id}/destroy", net.Destroy).Methods("PUT")
	router.HandleFunc("/api/network/{id}/autostart", net.Autostart).Methods("PUT")
	router.HandleFunc("/api/network/{id}/remove", net.Remove).Methods("PUT")
	router.HandleFunc("/api/network/{id}", net.Delete).Methods("DELETE")
}

func (net Network) Get(w http.ResponseWriter, r *http.Request) {
	uuid, ok := GetArg(r, "id")
	if !ok {
		// list all instances.
		list := schema.ListNetwork{
			Items: make([]schema.Network, 0, 32),
		}
		if obj, err := network.ListNetworks(); err == nil {
			for _, n := range obj {
				sn := network.NewNetwork(n)
				list.Items = append(list.Items, sn)
				_ = n.Free()
			}
			sort.SliceStable(list.Items, func(i, j int) bool {
				return list.Items[i].Name < list.Items[j].Name
			})
			list.Metadata.Size = len(list.Items)
			list.Metadata.Total = len(list.Items)
		}
		ResponseJson(w, list)
	} else {
		if n, err := network.LookupNetwork(uuid); err == nil {

			format := GetQueryOne(r, "format")
			if format == "xml" {
				obj := network.NewNetworkXMLFromNet(n)
				ResponseXML(w, libstar.XML.Encode(obj))
			} else {
				obj := network.NewNetwork(*n)
				ResponseJson(w, obj)
			}
			_ = n.Free()
		} else {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
	}
}

func (net Network) Post(w http.ResponseWriter, r *http.Request) {
	conf := schema.Network{}
	if err := GetData(r, &conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	xmlObj := Network2XML(conf)
	libstar.Debug("Network.Post %s", libstar.XML.Encode(xmlObj))
	hyper, err := network.GetHyper()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	netVir, err := hyper.NetworkDefineXML(libstar.XML.Encode(xmlObj))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer netVir.Free()
	if err := netVir.Create(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := netVir.SetAutostart(true); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, conf.Name+" success")
}

func (net Network) execute(w http.ResponseWriter, uuid, action string) {
	hyper, err := network.GetHyper()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	netvir, err := hyper.LookupNetwork(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer netvir.Free()

	switch action {
	case "start":
		if err := netvir.Create(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "destroy":
		if err := netvir.Destroy(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "remove":
		if err := netvir.Destroy(); err != nil {
			libstar.Warn("Network.Remove destroy %s", err)
		}
		if err := netvir.Undefine(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "unsupported action", http.StatusBadRequest)
		return
	}
	ResponseMsg(w, 0, "success")
}

func (net Network) Start(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	net.execute(w, uuid, "start")
}

func (net Network) Destroy(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	net.execute(w, uuid, "destroy")
}

func (net Network) Autostart(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	enable := GetQueryOne(r, "enable")
	on := !(enable == "false" || enable == "0" || enable == "no" || enable == "off")
	hyper, err := network.GetHyper()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	netvir, err := hyper.LookupNetwork(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer netvir.Free()
	if err := netvir.SetAutostart(on); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}

func (net Network) Remove(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	net.execute(w, uuid, "remove")
}

func (net Network) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	hyper, err := network.GetHyper()
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
		libstar.Warn("Network.Delete %s", err)
	}
	if err := netvir.Undefine(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}
