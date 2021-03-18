package api

import (
	"github.com/danieldin95/lightstar/src/compute/libvirtc"
	"github.com/danieldin95/lightstar/src/network/libvirtn"
)

func Interface2XML(source, model, seq, typ, drv, que string) libvirtc.InterfaceXML {
	if br, err := libvirtn.BRIDGE.Get(source); err == nil {
		typ = br.Type
	}
	if drv == "" {
		drv = "vhost"
		if que == "" {
			que = "2"
		}
	}
	xmlObj := libvirtc.InterfaceXML{
		Type: "bridge",
		Source: libvirtc.InterfaceSourceXML{
			Bridge: source,
		},
		Model: libvirtc.InterfaceModelXML{
			Type: model,
		},
		Address: &libvirtc.AddressXML{
			Type:     "pci",
			Domain:   libvirtc.PciDomain,
			Bus:      libvirtc.PciInterfaceBus,
			Slot:     seq,
			Function: libvirtc.PciFunc,
		},
		Driver: &libvirtc.InterfaceDriverXML{
			Name:   drv,
			Queues: que,
		},
	}
	if typ == "openvswitch" {
		xmlObj.VirtualPort = &libvirtc.InterfaceVirPortXML{
			Type: typ,
		}
	}
	return xmlObj
}
