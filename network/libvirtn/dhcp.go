package libvirtn

import "github.com/libvirt/libvirt-go"

type DHCP struct {
	libvirt.NetworkDHCPLease
}
