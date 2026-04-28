package api

import (
	"github.com/danieldin95/lightstar/pkg/compute"
	"github.com/danieldin95/lightstar/pkg/network"
)

func Interface2XML(source, model, seq, typ, drv, que string) *compute.InterfaceXML {
	if br, err := network.BRIDGE.Get(source); err == nil {
		typ = br.Type
	}
	if drv == "" {
		drv = "vhost"
		if que == "" {
			que = "2"
		}
	}
	xmlObj := &compute.InterfaceXML{
		Type: "bridge",
		Source: compute.InterfaceSourceXML{
			Bridge: source,
		},
		Model: compute.InterfaceModelXML{
			Type: model,
		},
		Address: &compute.AddressXML{
			Type:     "pci",
			Domain:   compute.PciDomain,
			Bus:      compute.PciInterfaceBus,
			Slot:     seq,
			Function: compute.PciFunc,
		},
		Driver: &compute.InterfaceDriverXML{
			Name:   drv,
			Queues: que,
		},
	}
	if typ == "openvswitch" {
		xmlObj.VirtualPort = &compute.InterfaceVirPortXML{
			Type: typ,
		}
	}
	return xmlObj
}
