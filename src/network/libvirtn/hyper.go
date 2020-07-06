package libvirtn

import (
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/libvirt/libvirt-go"
	"sync"
	"time"
)

var (
	NetworkAll = libvirt.CONNECT_LIST_NETWORKS_ACTIVE | libvirt.CONNECT_LIST_NETWORKS_INACTIVE
)

type HyperListener struct {
	Opened func(Conn *libvirt.Connect) error
	Closed func(Conn *libvirt.Connect) error
}

type HyperVisor struct {
	Name     string
	Conn     *libvirt.Connect
	Listener []HyperListener
	Lock     sync.RWMutex
	Ticker   *time.Ticker
	Done     chan bool
	Leases   map[string]DHCPLease
}

func (h *HyperVisor) AddListener(listen HyperListener) {
	h.Listener = append(h.Listener, listen)
}

func (h *HyperVisor) OpenNotSafe() error {
	if hyper.Conn != nil {
		if _, err := hyper.Conn.GetVersion(); err != nil {
			libstar.Error("HyperVisor.open %s", err)
			_, _ = hyper.Conn.Close()
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
				_ = listen.Opened(h.Conn)
			}
		}
	}
	if hyper.Conn == nil {
		return libstar.NewErr("Not connect.")
	}
	return nil
}

func (h *HyperVisor) Open() error {
	h.Lock.Lock()
	defer h.Lock.Unlock()

	return h.OpenNotSafe()
}

func (h *HyperVisor) Close() {
	if h.Conn == nil {
		return
	}
	for _, listen := range h.Listener {
		if listen.Closed != nil {
			_ = listen.Closed(h.Conn)
		}
	}
	h.Conn = nil
}

func (h *HyperVisor) ListAllNetworks() ([]Network, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}

	nets, err := h.Conn.ListAllNetworks(NetworkAll)
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

func (h *HyperVisor) SyncLeases() error {
	nets, err := hyper.ListAllNetworks()
	if err != nil {
		return err
	}
	h.Lock.Lock()
	defer h.Lock.Unlock()
	h.Leases = make(map[string]DHCPLease, 128)
	for _, net := range nets {
		les, err := net.GetDHCPLeases()
		_ = net.Free()
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
			h.Leases[d.Mac] = d
		}
	}
	return nil
}

func (h *HyperVisor) GetLeases() map[string]DHCPLease {
	h.Lock.RLock()
	defer h.Lock.RUnlock()
	leases := make(map[string]DHCPLease, 128)
	for name, value := range h.Leases {
		leases[name] = value
	}
	return leases
}

func (h *HyperVisor) LookupLeases(uuid string) ([]DHCPLease, error) {
	n, err := hyper.LookupNetwork(uuid)
	if err != nil {
		return nil, err
	}
	defer n.Free()

	data := make([]DHCPLease, 0, 128)
	leases, err := n.GetDHCPLeases()
	if err != nil {
		return nil, err
	}
	for _, le := range leases {
		data = append(data, DHCPLease{
			Type:     int(le.Type),
			IPAddr:   le.IPaddr,
			Prefix:   le.Prefix,
			Hostname: le.Hostname,
			Mac:      le.Mac,
		})
	}
	return data, nil
}

func (h *HyperVisor) LoopForever() {
	for {
		select {
		case <-h.Done:
			return
		case <-h.Ticker.C:
			_ = h.SyncLeases()
		}
	}
}

var hyper = HyperVisor{
	Name:     "qemu:///system",
	Listener: make([]HyperListener, 0, 32),
	Ticker:   time.NewTicker(2 * time.Second),
	Done:     make(chan bool, 2),
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

func init() {
	go hyper.LoopForever()
}
