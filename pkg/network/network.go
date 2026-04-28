package network

import (
	"strconv"
	"strings"

	"github.com/danieldin95/lightstar/pkg/virsh"
)

type DHCPLease struct {
	Type     int    `json:"type"`
	Mac      string `json:"mac"`
	IPAddr   string `json:"ipAddr"`
	Prefix   uint   `json:"prefix"`
	Hostname string `json:"hostname"`
}

type Network struct {
	Name  string
	UUID  string
	hyper *HyperVisor
}

func NewNetworkFromVir(name, uuid string, hyper *HyperVisor) *Network {
	return &Network{Name: name, UUID: uuid, hyper: hyper}
}

func (n *Network) networkRef() string {
	if n.Name != "" {
		return n.Name
	}
	return n.UUID
}

func (n *Network) run(args ...string) (string, error) {
	return virsh.Run(n.hyper.Name, args...)
}

func (n *Network) Free() error { return nil }

func (n *Network) GetName() (string, error) {
	if n.Name != "" {
		return n.Name, nil
	}
	out, err := n.run("net-name", n.networkRef())
	if err != nil {
		return "", err
	}
	n.Name = strings.TrimSpace(out)
	return n.Name, nil
}

func (n *Network) GetUUIDString() (string, error) {
	if n.UUID != "" {
		return n.UUID, nil
	}
	out, err := n.run("net-uuid", n.networkRef())
	if err != nil {
		return "", err
	}
	n.UUID = strings.TrimSpace(out)
	return n.UUID, nil
}

func (n *Network) GetXMLDesc(_ int) (string, error) {
	out, err := n.run("net-dumpxml", n.networkRef())
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (n *Network) GetBridgeName() (string, error) {
	xmlObj := NewNetworkXMLFromNet(n)
	if xmlObj == nil {
		return "", nil
	}
	return xmlObj.Bridge.Name, nil
}

func (n *Network) IsActive() (bool, error) {
	out, err := n.run("net-info", n.networkRef())
	if err != nil {
		return false, err
	}
	state := strings.ToLower(virsh.ParseKV(out)["active"])
	return state == "yes", nil
}

func (n *Network) IsAutostart() (bool, error) {
	out, err := n.run("net-info", n.networkRef())
	if err != nil {
		return false, err
	}
	val := strings.ToLower(virsh.ParseKV(out)["autostart"])
	return val == "yes", nil
}

func (n *Network) Create() error {
	_, err := n.run("net-start", n.networkRef())
	return err
}

func (n *Network) Destroy() error {
	_, err := n.run("net-destroy", n.networkRef())
	return err
}

func (n *Network) Undefine() error {
	_, err := n.run("net-undefine", n.networkRef())
	return err
}

func (n *Network) SetAutostart(enable bool) error {
	args := []string{"net-autostart", n.networkRef()}
	if !enable {
		args = append(args, "--disable")
	}
	_, err := n.run(args...)
	return err
}

func (n *Network) GetDHCPLeases() ([]DHCPLease, error) {
	out, err := n.run("net-dhcp-leases", n.networkRef())
	if err != nil {
		return nil, err
	}
	ret := make([]DHCPLease, 0, 32)
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Expiry") || strings.HasPrefix(line, "-") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		lease := DHCPLease{Mac: fields[1], Hostname: fields[4]}
		if strings.Contains(strings.ToLower(fields[2]), "ipv6") {
			lease.Type = 1
		}
		ip := fields[3]
		if strings.Contains(ip, "/") {
			arr := strings.SplitN(ip, "/", 2)
			lease.IPAddr = arr[0]
			if p, err := strconv.ParseUint(arr[1], 10, 32); err == nil {
				lease.Prefix = uint(p)
			}
		} else {
			lease.IPAddr = ip
		}
		ret = append(ret, lease)
	}
	return ret, nil
}

func ListNetworks() ([]Network, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	return hyper.ListAllNetworks()
}

func LookupNetwork(uuid string) (*Network, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	return hyper.LookupNetwork(uuid)
}

func ListLeases() (map[string]DHCPLease, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	return hyper.GetLeases(), nil
}

func LookupLeases(uuid string) ([]DHCPLease, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	return hyper.LookupLeases(uuid)
}
