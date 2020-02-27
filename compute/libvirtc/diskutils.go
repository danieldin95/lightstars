package libvirtc

import "fmt"

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
