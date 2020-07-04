package network

import (
	"github.com/danieldin95/lightstar/src/network/libvirtn"
	"github.com/danieldin95/lightstar/src/schema"
)

func NewNetwork(net libvirtn.Network) schema.Network {
	obj := schema.Network{
		Range: make([]schema.Range, 0, 32),
	}
	obj.Name, _ = net.GetName()
	obj.UUID, _ = net.GetUUIDString()
	obj.Bridge, _ = net.GetBridgeName()
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
				obj.Range = append(obj.Range, schema.Range{
					Start: addr.Start,
					End:   addr.End,
				})
			}
		}
	}
	return obj
}
