package api

import (
	"github.com/danieldin95/lightstar/src/compute/libvirtc"
	"github.com/danieldin95/lightstar/src/network/libvirtn"
)

func Interface2XML(source, model, seq, typ string) libvirtc.InterfaceXML {
	br, _ := libvirtn.BRIDGE.Get(source)
	xml := libvirtc.InterfaceXML{
		Type: "bridge",
		Source: libvirtc.InterfaceSourceXML{
			Bridge: source,
		},
		Model: libvirtc.InterfaceModelXML{
			Type: model,
		},
		Address: &libvirtc.AddressXML{
			Type:     "pci",
			Domain:   libvirtc.PCI_DOMAIN,
			Bus:      libvirtc.PCI_INTERFACE_BUS,
			Slot:     seq,
			Function: libvirtc.PCI_FUNC,
		},
	}
	if br.Type == "openvswitch" || typ == "openvswitch" {
		xml.VirtualPort = &libvirtc.InterfaceVirtualPortXML{
			Type: typ,
		}
	}
	return xml
}
