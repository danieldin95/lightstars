package api

import (
	"github.com/danieldin95/lightstar/src/compute/libvirtc"
	"github.com/danieldin95/lightstar/src/network/libvirtn"
)

func Interface2XML(source, model, seq, typ string) libvirtc.InterfaceXML {
	if br, err := libvirtn.BRIDGE.Get(source); err == nil {
		typ = br.Type
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
	}
	if typ == "openvswitch" {
		xmlObj.VirtualPort = &libvirtc.InterfaceVirPortXML{
			Type: typ,
		}
	}
	return xmlObj
}
