package libvirtdriver

import (
	"github.com/libvirt/libvirt-go"
)

var DOMAIN_ALL = libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE

type Domain struct {
	libvirt.Domain
}

func NewDomainFromVir(dom *libvirt.Domain) *Domain {
	return &Domain{*dom}
}

func (d *Domain) GetXMLDesc(secure bool) (string, error) {
	if secure {
		return d.Domain.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	} else {
		return d.Domain.GetXMLDesc(0)
	}
}
