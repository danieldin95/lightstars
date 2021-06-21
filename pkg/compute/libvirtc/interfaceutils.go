package libvirtc

import "fmt"

type Interface struct {
	//
}

func (int *Interface) Slot2Dev(slot uint8) string {
	prefix := "vnet"
	if slot <= 32 {
		return fmt.Sprintf("%s%d", prefix, slot)
	}
	return ""
}

var INTERFACE = &Interface{}
