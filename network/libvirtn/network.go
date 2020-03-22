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

func ListLeases(name string) (map[string]DHCPLease, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	nets := make([]Network, 0, 32)
	if name == "" {
		nets, err = hyper.ListAllNetworks()
		if err != nil {
			return nil, err
		}
	} else {
		if net, err := hyper.LookupNetwork(name); err == nil {
			nets = append(nets, *net)
		}
	}
	leases := make(map[string]DHCPLease, 128)
	for _, net := range nets {
		les, err := net.GetDHCPLeases()
		net.Free()
		if err != nil {
			continue
		}
		for _, le := range les {
			d := DHCPLease{
				Type:     int(le.Type),
				IPAddr:   le.IPaddr,
				Prefix:   le.Prefix,
				Hostname: le.Hostname,
				Mac:      le.Mac,
			}
			leases[d.Mac] = d
		}
	}
	return leases, nil
}
