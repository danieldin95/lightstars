package libvirtc

import (
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/libvirt/libvirt-go"
	"strings"
	"sync"
	"time"
)

type HyperListener struct {
	Opened func(Conn *libvirt.Connect) error
	Closed func(Conn *libvirt.Connect) error
}

type HyperVisor struct {
	Url        string
	Schema     string
	Address    string
	Path       string
	Host       string
	Conn       *libvirt.Connect
	Listener   []HyperListener
	ConnedTime int64
	Lock       sync.RWMutex
	Ticker     *time.Ticker
	CpuSts     *libvirt.NodeCPUStats
	Done       chan bool
	IdleUtil   uint64
	DomUtil    map[string]uint64
}

func parseQemuTCP(name string) (address, path string) {
	if strings.Contains(name, "://") {
		addrs := strings.SplitN(name, "://", 2)[1]
		address = strings.SplitN(addrs, "/", 2)[0]
		if strings.Contains(addrs, "/") {
			path = strings.SplitN(addrs, "/", 2)[1]
		}
	}
	return address, path
}

func parseQemuSSH(name string) (address, path string) {
	if strings.Contains(name, "://") {
		addrs := strings.SplitN(name, "://", 2)[1]
		address = strings.SplitN(addrs, "/", 2)[0]
		if strings.Contains(addrs, "/") {
			path = strings.SplitN(addrs, "/", 2)[1]
		}
		if strings.Contains(address, "@") {
			address = strings.SplitN(address, "@", 2)[1]
		}
	}
	return address, path
}

func (h *HyperVisor) open() error {
	if hyper.Conn != nil {
		if _, err := hyper.Conn.GetVersion(); err != nil {
			libstar.Error("HyperVisor.open %s", err)
			_, _ = hyper.Conn.Close()
			hyper.Conn = nil
		}
	}
	if hyper.Conn == nil {
		conn, err := libvirt.NewConnect(hyper.Url)
		if err != nil {
			return err
		}
		hyper.ConnedTime = time.Now().Unix()
		hyper.Conn = conn
		if name, err := conn.GetHostname(); err == nil {
			hyper.Host = name
		}
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

func (h *HyperVisor) UpTime() int64 {
	return time.Now().Unix() - h.ConnedTime
}

func (h *HyperVisor) Open() error {
	h.Lock.Lock()
	defer h.Lock.Unlock()
	return h.open()
}

func (h *HyperVisor) AddListener(listen HyperListener) {
	h.Listener = append(h.Listener, listen)
}

func (h *HyperVisor) SetUrl(name string) {
	hyper.Url = name

	h.Schema = strings.SplitN(h.Url, ":", 2)[0]
	switch h.Schema {
	case "qemu+ssh":
		h.Address, h.Path = parseQemuSSH(h.Url)
	case "qemu+tcp", "qemu+tls":
		h.Address, h.Path = parseQemuTCP(h.Url)
	default:
		h.Address = "localhost"
		h.Path = "system"
	}
	if strings.Contains(h.Address, ":") {
		h.Address = strings.SplitN(h.Address, ":", 2)[0]
	}
	h.Host = h.Address
}

func (h *HyperVisor) FigureCPU() (err error) {
	h.Lock.Lock()
	defer h.Lock.Unlock()

	libstar.Debug("HyperVisor.FigureCpu")
	if err := h.open(); err != nil {
		return err
	}
	if h.CpuSts == nil {
		h.CpuSts, err = h.Conn.GetCPUStats(-1, 0)
		if err != nil {
			libstar.Warn("HyperVisor.FigureCpu %s", err)
			return err
		}
	}
	newerSts, err := h.Conn.GetCPUStats(-1, 0)
	if err != nil {
		libstar.Warn("HyperVisor.FigureCpu %s", err)
		return err
	}
	older := h.CpuSts.User
	older += h.CpuSts.Idle
	older += h.CpuSts.Kernel
	older += h.CpuSts.Intr
	older += h.CpuSts.Iowait
	newer := newerSts.User
	newer += newerSts.Idle
	newer += newerSts.Kernel
	newer += newerSts.Intr
	newer += newerSts.Iowait
	dt := newer - older
	if dt > 0 {
		h.IdleUtil = 1000 * (newerSts.Idle - h.CpuSts.Idle) / dt
		// record last statics
		h.CpuSts = newerSts
	}
	return nil
}

func (h *HyperVisor) LoopForever() {
	for {
		select {
		case <-h.Done:
			return
		case <-h.Ticker.C:
			if err := h.FigureCPU(); err != nil {
				libstar.Warn("HyperVisor.LoopForever %s", err)
			}
		}
	}
}

func (h *HyperVisor) GetCPU() (uint, string, uint64) {
	if err := h.Open(); err != nil {
		return 0, "", 1000
	}

	h.Lock.RLock()
	defer h.Lock.RUnlock()
	if info, err := h.Conn.GetNodeInfo(); err == nil {
		return info.Cpus, info.Model, h.IdleUtil
	}
	return 0, "", 1000
}

func (h *HyperVisor) GetMem() (t uint64, f uint64, c uint64) {
	if err := h.Open(); err != nil {
		return 0, 0, 0
	}
	if stats, err := h.Conn.GetMemoryStats(-1, 0); err == nil {
		if stats.TotalSet {
			t = stats.Total * 1024
		}
		if stats.FreeSet {
			f = stats.Free * 1024
		}
		if stats.CachedSet {
			c = stats.Cached * 1024
		}
	}
	return t, f, c
}

func (h *HyperVisor) GetRootfs() string {
	if err := h.Open(); err != nil {
		return ""
	}
	return ""
}

func (h *HyperVisor) ListAllDomains() ([]Domain, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}

	domains, err := h.Conn.ListAllDomains(DomainAll)
	if err != nil {
		return nil, err
	}
	newDomains := make([]Domain, 0, 32)
	for _, m := range domains {
		newDomains = append(newDomains, *NewDomainFromVir(&m))
	}
	return newDomains, nil
}

func (h *HyperVisor) LookupDomainByUUIDString(id string) (*Domain, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	domain, err := hyper.Conn.LookupDomainByUUIDString(id)
	if err != nil {
		return nil, err
	}
	return NewDomainFromVir(domain), nil
}

func (h *HyperVisor) LookupDomainByUUIDName(id string) (*Domain, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	domain, err := hyper.Conn.LookupDomainByUUIDString(id)
	if err != nil {
		domain, err := hyper.Conn.LookupDomainByName(id)
		if err != nil {
			return nil, err
		}
		return NewDomainFromVir(domain), nil
	}
	return NewDomainFromVir(domain), nil
}

func (h *HyperVisor) LookupDomainByName(id string) (*Domain, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	domain, err := hyper.Conn.LookupDomainByName(id)
	if err != nil {
		return nil, err
	}
	return &Domain{*domain}, nil
}

func (h *HyperVisor) DomainDefineXML(xmlConfig string) (*Domain, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	domain, err := h.Conn.DomainDefineXML(xmlConfig)
	if err != nil {
		return nil, err
	}
	return &Domain{*domain}, nil
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

var hyper = HyperVisor{
	Listener: make([]HyperListener, 0, 32),
	Ticker:   time.NewTicker(2 * time.Second),
	Done:     make(chan bool),
	IdleUtil: 1000,
	DomUtil:  make(map[string]uint64, 32),
}

func GetHyper() (*HyperVisor, error) {
	return &hyper, hyper.Open()
}

func SetHyper(name string) (*HyperVisor, error) {
	if name == hyper.Url {
		return &hyper, nil
	}
	hyper.Close()
	hyper.SetUrl(name)
	return &hyper, hyper.Open()
}

func LookupDomainByUUIDString(uuid string) (*Domain, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	dom, err := hyper.LookupDomainByUUIDString(uuid)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func LookupDomainByUUIDName(uuid string) (*Domain, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	dom, err := hyper.LookupDomainByUUIDName(uuid)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func AddHyperListener(listen HyperListener) {
	hyper.AddListener(listen)
}

func init() {
	hyper.SetUrl("qemu:///system")
	go hyper.LoopForever()
}
