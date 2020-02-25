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
	ROOT_BUS          = "0x00"
	DISK_BUS          = "0x01"
	INTERFACE_BUS     = "0x02"
	ROOT_BUS_DRV      = "0"
	DISK_BUS_DRV      = "1"
	INTERFACE_BUS_DRV = "2"
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
		return prefix + string('a'+slot)
	}
	return ""
}

func (d *Disk) Slot2DiskName(slot uint8) string {
	return fmt.Sprintf("disk%d.img", slot)
}

var DISK = &Disk{}
