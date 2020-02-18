package libvirtdriver

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/libvirt/libvirt-go"
	"strings"
)

type HyperVisor struct {
	Name    string
	Schema  string
	Address string
	Path    string
	Conn    *libvirt.Connect
}

func (h *HyperVisor) Init() {
	if h.Name != "" {
		h.Schema = strings.SplitN(h.Name, ":", 2)[0]
		switch h.Schema {
		case "qemu+ssh":
			if strings.Contains(h.Name, "://") {
				address := strings.SplitN(h.Name, "://", 2)[1]
				h.Address = strings.SplitN(address, "/", 2)[0]
				if strings.Contains(address, "/") {
					h.Path = strings.SplitN(address, "/", 2)[1]
				}
				if strings.Contains(h.Address, "@") {
					h.Address = strings.SplitN(h.Address, "@", 2)[1]
				}
			}
		case "qemu+tcp", "qemu+tls":
			if strings.Contains(h.Name, "://") {
				address := strings.SplitN(h.Name, "://", 2)[1]
				h.Address = strings.SplitN(address, "/", 2)[0]
				if strings.Contains(address, "/") {
					h.Path = strings.SplitN(address, "/", 2)[1]
				}
			}
		default:
			h.Address = "localhost"
			h.Path = "system"
		}
		if strings.Contains(h.Address, ":") {
			h.Address = strings.SplitN(h.Address, ":", 2)[0]
		}
	}
}

func (h *HyperVisor) ListAllDomains() ([]Domain, error) {
	if h.Conn == nil {
		return nil, libstar.NewErr("not connected")
	}
	domains, err := h.Conn.ListAllDomains(DOMAIN_ALL)
	if err != nil {
		return nil, err
	}
	newDomains := make([]Domain, 0, 32)
	for _, m := range domains {
		newDomains = append(newDomains, Domain{m})
	}
	return newDomains, nil
}

func (h *HyperVisor) LookupDomainByUUIDString(id string) (*Domain, error) {
	if h.Conn == nil {
		return nil, libstar.NewErr("not connected")
	}

	domain, err := hyper.Conn.LookupDomainByUUIDString(id)
	if err != nil {
		return nil, err
	}
	return &Domain{*domain}, nil
}

var hyper = HyperVisor{
	Name:    "qemu:///system",
	Address: "localhost",
}

func GetHyper(name string) (*HyperVisor, error) {
	if hyper.Conn == nil {
		conn, err := libvirt.NewConnect(hyper.Name)
		if err != nil {
			return &hyper, err
		}
		hyper.Conn = conn
	}
	return &hyper, nil
}

func SetHyper(name string) (*HyperVisor, error) {
	if name == hyper.Name {
		return &hyper, nil
	}
	hyper.Name = name
	hyper.Init()

	conn, err := libvirt.NewConnect(hyper.Name)
	if err != nil {
		return &hyper, err
	}
	hyper.Conn = conn
	return &hyper, nil
}

func CloseHyper(name string) {
	hyper.Conn.Close()
}

func init() {
	hyper.Init()
}
