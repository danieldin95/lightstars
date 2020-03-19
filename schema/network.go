package schema

import (
	"github.com/danieldin95/lightstar/network/libvirtn"
)

type Range struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
type Network struct {
	UUID    string  `json:"uuid"`
	Name    string  `json:"name"`
	State   string  `json:"state"`
	Address string  `json:"address"`
	Netmask string  `json:"netmask,omitempty"`
	Prefix  string  `json:"prefix,omitempty"`
	Range   []Range `json:"range"`
	DHCP    string  `json:"dhcp,omitempty"`
	Mode    string  `json:"mode"` // nat, router.
}

func NewNetwork(net libvirtn.Network) Network {
	obj := Network{
		Range: make([]Range, 0, 32),
	}
	obj.Name, _ = net.GetName()
	obj.UUID, _ = net.GetUUIDString()
	if ok, _ := net.IsActive(); ok {
		obj.State = "active"
	} else {
		obj.State = "inactive"
	}

	xml := libvirtn.NewNetworkXMLFromNet(&net)
	if xml.Forward != nil {
		obj.Mode = xml.Forward.Mode
	}
	if xml.IPv4 != nil {
		obj.Address = xml.IPv4.Address
		obj.Netmask = xml.IPv4.Netmask
		obj.Prefix = xml.IPv4.Prefix
		if xml.IPv4.DHCP != nil {
			for _, addr := range xml.IPv4.DHCP.Range {
				obj.Range = append(obj.Range, Range{
					Start: addr.Start,
					End:   addr.End,
				})
			}
		}
	}
	if xml.Bridge.Name != "" {
		obj.Name = xml.Bridge.Name
	}
	return obj
}
