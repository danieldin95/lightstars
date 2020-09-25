package libvirtc

import (
	"github.com/libvirt/libvirt-go"
)

var (
	DomainAll                    = libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE
	DomainDeviceModifyConfig     = libvirt.DOMAIN_DEVICE_MODIFY_CONFIG
	DomainDeviceModifyPersistent = libvirt.DOMAIN_DEVICE_MODIFY_LIVE | libvirt.DOMAIN_DEVICE_MODIFY_CONFIG
	DomainDestroyGraceful        = libvirt.DOMAIN_DESTROY_GRACEFUL
	DomainShutdownAcpi           = libvirt.DOMAIN_SHUTDOWN_ACPI_POWER_BTN
	DomainCpuMaximum             = libvirt.DOMAIN_VCPU_MAXIMUM | libvirt.DOMAIN_VCPU_CONFIG
	DomainCpuConfig              = libvirt.DOMAIN_VCPU_CONFIG
	DomainMemMaximum             = libvirt.DOMAIN_MEM_MAXIMUM | libvirt.DOMAIN_MEM_CONFIG
	DomainMemConfig              = libvirt.DOMAIN_MEM_CONFIG
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
