package libvirtc

import (
	"fmt"
	"github.com/libvirt/libvirt-go"
)

var (
	DOMAIN_ALL                      = libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE
	DOMAIN_DEVICE_MODIFY_CONFIG     = libvirt.DOMAIN_DEVICE_MODIFY_CONFIG
	DOMAIN_DEVICE_MODIFY_PERSISTENT = libvirt.DOMAIN_DEVICE_MODIFY_LIVE | libvirt.DOMAIN_DEVICE_MODIFY_CONFIG
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

type Disk struct {
	//
}

func (d *Disk) Slot2Dev(bus string, slot uint8) string {
	prefix := "vd"
	if bus == "ide" || bus == "scsi" {
		prefix = "hd"
	}
	if slot <= 26 {
		return prefix + string('a'+slot-1)
	}
	return ""
}

func (d *Disk) Slot2DiskName(slot uint8) string {
	return fmt.Sprintf("disk%d.img", slot)
}

var DISK = &Disk{}
