package libvirtc

import (
	"github.com/libvirt/libvirt-go"
)

var (
	DOMAIN_ALL                      = libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE
	DOMAIN_DEVICE_MODIFY_CONFIG     = libvirt.DOMAIN_DEVICE_MODIFY_CONFIG
	DOMAIN_DEVICE_MODIFY_PERSISTENT = libvirt.DOMAIN_DEVICE_MODIFY_LIVE | libvirt.DOMAIN_DEVICE_MODIFY_CONFIG
	DOMAIN_DESTROY_GRACEFUL         = libvirt.DOMAIN_DESTROY_GRACEFUL
	DOMAIN_SHUTDOWN_ACPI            = libvirt.DOMAIN_SHUTDOWN_ACPI_POWER_BTN
)

var (
	PCI_DOMAIN        = "0x00"
	PCI_ROOT_BUS      = "0x00"
	PCI_DISK_BUS      = "0x01"
	PCI_INTERFACE_BUS = "0x02"
	PCI_FUNC          = "0x00"
	DRV_CTL           = "0"
	DRV_ROOT_BUS      = "0"
	DRV_DISK_BUS      = "1"
	DRV_INTERFACE_BUS = "2"
)

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
		return d.Domain.GetXMLDesc(libvirt.DOMAIN_XML_INACTIVE)
	}
}

func ListDomains() ([]Domain, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	return hyper.ListAllDomains()
}

func DomainState2Str(state libvirt.DomainState) string {
	switch state {
	case libvirt.DOMAIN_NOSTATE:
		return "no-state"
	case libvirt.DOMAIN_RUNNING:
		return "running"
	case libvirt.DOMAIN_BLOCKED:
		return "blocked"
	case libvirt.DOMAIN_PAUSED:
		return "paused"
	case libvirt.DOMAIN_SHUTDOWN:
		return "shutdown"
	case libvirt.DOMAIN_CRASHED:
		return "crashed"
	case libvirt.DOMAIN_PMSUSPENDED:
		return "pm-suspended"
	case libvirt.DOMAIN_SHUTOFF:
		return "shutoff"
	default:
		return "unknown"
	}
}
