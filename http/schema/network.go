package schema

import (
	"github.com/danieldin95/lightstar/network/libvirtn"
)

type Network struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Address string `json:"address"`
	Netmask string `json:"netmask"`
	Prefix  string `json:"prefix,omitempty"`
	DHCP    string `json:"dhcp,omitempty"`
	Mode    string `json:"mode"` // nat, router.
}

func NewNetwork(net libvirtn.Network) Network {
	obj := Network{}
	obj.Name, _ = net.GetName()
	obj.UUID, _ = net.GetUUIDString()
	if ok, _ := net.IsActive(); ok {
		obj.State = "active"
	} else {
		obj.State = "inactive"
	}

	xml := libvirtn.NewNetworkXMLFromNet(&net)
	obj.Mode = xml.Forward.Mode
	obj.Address = xml.IPv4.Address
	obj.Netmask = xml.IPv4.Netmask
	if xml.Bridge.Name != "" {
		obj.Name = xml.Bridge.Name
	}
	return obj
}
