package libvirtn

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/libvirt/libvirt-go"
)

var (
	NETWORK_ALL = libvirt.CONNECT_LIST_NETWORKS_ACTIVE | libvirt.CONNECT_LIST_NETWORKS_INACTIVE
)

type HyperListener struct {
	Opened func(Conn *libvirt.Connect) error
	Closed func(Conn *libvirt.Connect) error
}

type HyperVisor struct {
	Name     string
	Conn     *libvirt.Connect
	Listener []HyperListener
}

func (h *HyperVisor) AddListener(listen HyperListener) {
	h.Listener = append(h.Listener, listen)
}

func (h *HyperVisor) Open() error {
	if hyper.Conn != nil {
		if _, err := hyper.Conn.GetVersion(); err != nil {
			libstar.Error("HyperVisor.Open %s", err)
			hyper.Conn.Close()
			hyper.Conn = nil
		}
	}
	if hyper.Conn == nil {
		conn, err := libvirt.NewConnect(hyper.Name)
		if err != nil {
			return err
		}
		hyper.Conn = conn
		for _, listen := range h.Listener {
			if listen.Opened != nil {
				listen.Opened(h.Conn)
			}
		}
	}
	if hyper.Conn == nil {
		return libstar.NewErr("Not connect.")
	}
	return nil
}

func (h *HyperVisor) Close() {
	if h.Conn == nil {
		return
	}
	for _, listen := range h.Listener {
		if listen.Closed != nil {
			listen.Closed(h.Conn)
		}
	}
	h.Conn = nil
}

func (h *HyperVisor) ListAllNetworks() ([]Network, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}

	nets, err := h.Conn.ListAllNetworks(NETWORK_ALL)
	if err != nil {
		return nil, err
	}
	newNets := make([]Network, 0, 32)
	for _, n := range nets {
		newNets = append(newNets, *NewNetworkFromVir(&n))
	}
	return newNets, nil
}

func (h *HyperVisor) NetworkDefineXML(xml string) (*Network, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}

	net, err := h.Conn.NetworkDefineXML(xml)
	if err != nil {
		return nil, err
	}
	return &Network{Network: *net}, nil
}

// name: uuid, name
func (h *HyperVisor) LookupNetwork(name string) (*Network, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}

	net, err := h.Conn.LookupNetworkByUUIDString(name)
	if err != nil {
		net, err = h.Conn.LookupNetworkByName(name)
	}
	if err != nil {
		return nil, err
	}
	return &Network{Network: *net}, nil
}

var hyper = HyperVisor{
	Name:     "qemu:///system",
	Listener: make([]HyperListener, 0, 32),
}

func GetHyper() (*HyperVisor, error) {
	return &hyper, hyper.Open()
}

func SetHyper(name string) (*HyperVisor, error) {
	if name == hyper.Name {
		return &hyper, nil
	}
	hyper.Close()
	hyper.Name = name
	return &hyper, hyper.Open()
}

func AddHyperListener(listen HyperListener) {
	hyper.AddListener(listen)
}
