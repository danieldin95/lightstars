package libvirtdriver

import (
	"github.com/libvirt/libvirt-go"
)

var DOMAIN_ALL = libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE

type Domain struct {
	libvirt.Domain
}
