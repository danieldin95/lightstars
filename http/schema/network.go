package schema

import (
	"github.com/danieldin95/lightstar/network/libvirtn"
)

type Network struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	State      string `json:"state"`
	Address    string `json:"address"`
	Netmask    string `json:"netmask,omitempty"`
	Prefix     string `json:"prefix,omitempty"`
	RangeStart string `json:"rangeStart,omitempty"`
	RangeEnd   string `json:"rangeEnd,omitempty"`
	DHCP       string `json:"dhcp,omitempty"`
	Mode       string `json:"mode"` // nat, router.
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
	if xml.IPv4 != nil {
		obj.Address = xml.IPv4.Address
		obj.Netmask = xml.IPv4.Netmask
		obj.Prefix = xml.IPv4.Prefix
	}
	if xml.Bridge.Name != "" {
		obj.Name = xml.Bridge.Name
	}
	return obj
}
