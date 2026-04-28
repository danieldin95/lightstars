package network

import (
	"strings"
	"sync"
	"time"

	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/virsh"
)

var (
	NetworkAll = 0
)

type HyperListener struct {
	Opened func(h *HyperVisor) error
	Closed func(h *HyperVisor) error
}

type HyperVisor struct {
	Name       string
	Connected  bool
	Listener   []HyperListener
	Lock       sync.RWMutex
	Ticker     *time.Ticker
	Done       chan bool
	Leases     map[string]DHCPLease
	ConnedTime int64
}

func (h *HyperVisor) AddListener(listen HyperListener) {
	h.Listener = append(h.Listener, listen)
}

func (h *HyperVisor) OpenNotSafe() error {
	if _, err := virsh.Run(h.Name, "uri"); err != nil {
		return err
	}
	if !h.Connected {
		h.Connected = true
		h.ConnedTime = time.Now().Unix()
		for _, listen := range h.Listener {
			if listen.Opened != nil {
				_ = listen.Opened(h)
			}
		}
	}
	return nil
}

func (h *HyperVisor) Open() error {
	h.Lock.Lock()
	defer h.Lock.Unlock()
	return h.OpenNotSafe()
}

func (h *HyperVisor) Close() {
	if !h.Connected {
		return
	}
	for _, listen := range h.Listener {
		if listen.Closed != nil {
			_ = listen.Closed(h)
		}
	}
	h.Connected = false
	h.ConnedTime = 0
}

func (h *HyperVisor) ListAllNetworks() ([]Network, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	out, err := virsh.Run(h.Name, "net-list", "--all", "--name")
	if err != nil {
		return nil, err
	}
	lines := virsh.Lines(out)
	nets := make([]Network, 0, len(lines))
	for _, name := range lines {
		uuidOut, _ := virsh.Run(h.Name, "net-uuid", name)
		nets = append(nets, Network{Name: name, UUID: strings.TrimSpace(uuidOut), hyper: h})
	}
	return nets, nil
}

func (h *HyperVisor) NetworkDefineXML(xmlData string) (*Network, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	file, cleanup, err := virsh.TempXML("network-define", xmlData)
	if err != nil {
		return nil, err
	}
	defer cleanup()
	if _, err := virsh.Run(h.Name, "net-define", file); err != nil {
		return nil, err
	}
	obj := &NetworkXML{}
	_ = libstar.XML.Decode(obj, xmlData)
	if obj.Name == "" {
		return nil, libstar.NewErr("network name missing")
	}
	return h.LookupNetwork(obj.Name)
}

func (h *HyperVisor) LookupNetwork(name string) (*Network, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	netName := name
	if out, err := virsh.Run(h.Name, "net-name", name); err == nil {
		netName = strings.TrimSpace(out)
	}
	uuidOut, err := virsh.Run(h.Name, "net-uuid", netName)
	if err != nil {
		return nil, err
	}
	return &Network{Name: netName, UUID: strings.TrimSpace(uuidOut), hyper: h}, nil
}

func (h *HyperVisor) SyncLeases() error {
	nets, err := h.ListAllNetworks()
	if err != nil {
		return err
	}
	leases := make(map[string]DHCPLease, 128)
	for _, net := range nets {
		les, err := net.GetDHCPLeases()
		if err != nil {
			continue
		}
		for _, le := range les {
			leases[le.Mac] = le
		}
	}
	h.Lock.Lock()
	h.Leases = leases
	h.Lock.Unlock()
	return nil
}

func (h *HyperVisor) GetLeases() map[string]DHCPLease {
	h.Lock.RLock()
	defer h.Lock.RUnlock()
	leases := make(map[string]DHCPLease, len(h.Leases))
	for name, value := range h.Leases {
		leases[name] = value
	}
	return leases
}

func (h *HyperVisor) LookupLeases(uuid string) ([]DHCPLease, error) {
	n, err := h.LookupNetwork(uuid)
	if err != nil {
		return nil, err
	}
	return n.GetDHCPLeases()
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
	Leases:   make(map[string]DHCPLease, 128),
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
