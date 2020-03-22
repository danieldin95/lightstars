package libvirtn

import (
	"github.com/libvirt/libvirt-go"
)

type DHCPLease struct {
	Type     int    `json:"type"`
	Mac      string `json:"mac"`
	IPAddr   string `json:"ipAddr"`
	Prefix   uint   `json:"prefix"`
	Hostname string `json:"hostname"`
}

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

func ListLeases() (map[string]DHCPLease, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	return hyper.GetLeases(), nil
}
