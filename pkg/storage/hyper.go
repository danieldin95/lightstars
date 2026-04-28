package storage

import (
	"strings"

	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/virsh"
)

var (
	PoolAll = 0
)

type HyperListener struct {
	Opened func(h *HyperVisor) error
	Closed func(h *HyperVisor) error
}

type HyperVisor struct {
	Name      string
	Connected bool
	Listener  []HyperListener
}

func (h *HyperVisor) AddListener(listen HyperListener) {
	libstar.Debug("HyperVisor.AddListener %v", listen)
	h.Listener = append(h.Listener, listen)
}

func (h *HyperVisor) Open() error {
	if _, err := virsh.Run(h.Name, "uri"); err != nil {
		return err
	}
	if !h.Connected {
		h.Connected = true
		for _, listen := range h.Listener {
			if listen.Opened != nil {
				_ = listen.Opened(h)
			}
		}
	}
	return nil
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
}

func (h *HyperVisor) ListAllPools() ([]Pool, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	out, err := virsh.Run(h.Name, "pool-list", "--all", "--name")
	if err != nil {
		return nil, err
	}
	lines := virsh.Lines(out)
	pools := make([]Pool, 0, len(lines))
	for _, name := range lines {
		if !IsStorePool(name) {
			continue
		}
		uuidOut, _ := virsh.Run(h.Name, "pool-uuid", name)
		pools = append(pools, Pool{Name: name, UUID: strings.TrimSpace(uuidOut), hyper: h})
	}
	return pools, nil
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
