package libvirtc

import (
	"fmt"
	"github.com/danieldin95/lightstar/pkg/libstar"
)

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

func (d *Disk) Slot2Name(slot uint8) string {
	rand := libstar.GenToken(8)
	return fmt.Sprintf("disk-%s-%d.img", rand, slot)
}

var DISK = &Disk{}
