package libvirts

import (
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/libvirt/libvirt-go"
)

var (
	PoolAll = libvirt.CONNECT_LIST_STORAGE_POOLS_ACTIVE | libvirt.CONNECT_LIST_STORAGE_POOLS_INACTIVE
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
	libstar.Debug("HyperVisor.AddListener %v", listen)
	h.Listener = append(h.Listener, listen)
}

func (h *HyperVisor) Open() error {
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

func (h *HyperVisor) ListAllPools() ([]Pool, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}

	pools, err := h.Conn.ListAllStoragePools(PoolAll)
	if err != nil {
		return nil, err
	}
	newPools := make([]Pool, 0, 32)
	for _, p := range pools {
		name, err := p.GetName()
		if err != nil || !IsStorePool(name) {
			_ = p.Free()
			continue
		}
		newPools = append(newPools, *NewPoolFromVir(&p))
	}
	return newPools, nil
}

var hyper = HyperVisor{
	Name:     "qemu:///system",
	Listener: make([]HyperListener, 0, 32),
}

func GetHyper() (*HyperVisor, error) {
	return &hyper, hyper.Open()
}

func AddHyperListener(listen HyperListener) {
	hyper.AddListener(listen)
}

func SetHyper(name string) (*HyperVisor, error) {
	if name == hyper.Name {
		return &hyper, nil
	}
	hyper.Close()
	hyper.Name = name
	return &hyper, hyper.Open()
}
