package libvirtn

import "github.com/libvirt/libvirt-go"

type Network struct {
	libvirt.Network
}

func NewNetworkFromVir(net *libvirt.Network) *Network {
	return &Network{Network: *net}
}

func ListNetworks() ([]Network, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	return hyper.ListAllNetworks()
}
