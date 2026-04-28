package compute

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/virsh"
)

type HyperListener struct {
	Opened func(h *HyperVisor) error
	Closed func(h *HyperVisor) error
}

type HyperVisor struct {
	Url        string
	Schema     string
	Address    string
	Path       string
	Host       string
	Listener   []HyperListener
	ConnedTime int64
	Lock       sync.RWMutex
	Ticker     *time.Ticker
	Done       chan bool
	IdleUtil   uint64
	DomUtil    map[string]uint64
}

func parseUrl(url string) (address, path string) {
	if strings.Contains(url, "://") {
		shortURL := strings.SplitN(url, "://", 2)[1]
		address = strings.SplitN(shortURL, "/", 2)[0]
		if strings.Contains(shortURL, "/") {
			path = strings.SplitN(shortURL, "/", 2)[1]
		}
		if strings.Contains(address, "@") {
			address = strings.SplitN(address, "@", 2)[1]
		}
	}
	return address, path
}

func (h *HyperVisor) open() error {
	if _, err := virsh.Run(h.Url, "uri"); err != nil {
		return err
	}
	if h.ConnedTime <= 0 {
		h.ConnedTime = time.Now().Unix()
		for _, listen := range h.Listener {
			if listen.Opened != nil {
				_ = listen.Opened(h)
			}
		}
	}
	return nil
}

func (h *HyperVisor) UpTime() int64 {
	if h.ConnedTime <= 0 {
		return 0
	}
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
	h.Url = name
	h.Schema = strings.SplitN(h.Url, "://", 2)[0]
	switch h.Schema {
	case "qemu+ssh", "qemu+tcp", "qemu+tls":
		h.Address, h.Path = parseUrl(h.Url)
	default:
		h.Address = "localhost"
		h.Path = "system"
	}
	if strings.Contains(h.Address, ":") {
		h.Address = strings.SplitN(h.Address, ":", 2)[0]
	}
	h.Host = h.Address
}

func (h *HyperVisor) FigureCPU() error {
	if err := h.Open(); err != nil {
		return err
	}
	out, err := virsh.Run(h.Url, "nodecpustats")
	if err != nil {
		return err
	}
	kv := virsh.ParseKV(out)
	idle := parseUint(kv["idle"])
	user := parseUint(kv["user"])
	kernel := parseUint(kv["kernel"])
	iowait := parseUint(kv["iowait"])
	total := idle + user + kernel + iowait
	if total == 0 {
		return nil
	}
	h.Lock.Lock()
	defer h.Lock.Unlock()
	h.IdleUtil = 1000 * idle / total
	return nil
}

func parseUint(v string) uint64 {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}
	n, _ := strconv.ParseUint(strings.Fields(v)[0], 10, 64)
	return n
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
	out, err := virsh.Run(h.Url, "nodeinfo")
	if err != nil {
		return 0, "", 1000
	}
	kv := virsh.ParseKV(out)
	cpus := uint(parseUint(kv["cpu(s)"]))
	model := kv["cpu model"]
	h.Lock.RLock()
	util := h.IdleUtil
	h.Lock.RUnlock()
	return cpus, model, util
}

func readMemInfo() (total, free, cached uint64) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, 0
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "MemTotal:") {
			total = parseUint(line) * 1024
		} else if strings.HasPrefix(line, "MemAvailable:") {
			free = parseUint(line) * 1024
		} else if strings.HasPrefix(line, "Cached:") {
			cached = parseUint(line) * 1024
		}
	}
	if free == 0 {
		free = cached
	}
	return total, free, cached
}

func (h *HyperVisor) GetMem() (t uint64, f uint64, c uint64) {
	if err := h.Open(); err != nil {
		return 0, 0, 0
	}
	return readMemInfo()
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
	out, err := virsh.Run(h.Url, "list", "--all", "--name")
	if err != nil {
		return nil, err
	}
	lines := virsh.Lines(out)
	domains := make([]Domain, 0, len(lines))
	for _, name := range lines {
		uuidOut, err := virsh.Run(h.Url, "domuuid", name)
		uuid := ""
		if err == nil {
			uuid = strings.TrimSpace(uuidOut)
		}
		domains = append(domains, Domain{Name: name, UUID: uuid, hyper: h})
	}
	return domains, nil
}

func (h *HyperVisor) LookupDomainByUUIDString(id string) (*Domain, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	nameOut, err := virsh.Run(h.Url, "domname", id)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(nameOut)
	uuidOut, _ := virsh.Run(h.Url, "domuuid", name)
	return &Domain{Name: name, UUID: strings.TrimSpace(uuidOut), hyper: h}, nil
}

func (h *HyperVisor) LookupDomainByUUIDName(id string) (*Domain, error) {
	if dom, err := h.LookupDomainByUUIDString(id); err == nil {
		return dom, nil
	}
	return h.LookupDomainByName(id)
}

func (h *HyperVisor) LookupDomainByName(id string) (*Domain, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	uuidOut, err := virsh.Run(h.Url, "domuuid", id)
	if err != nil {
		return nil, err
	}
	return &Domain{Name: id, UUID: strings.TrimSpace(uuidOut), hyper: h}, nil
}

func (h *HyperVisor) DomainDefineXML(xmlConfig string) (*Domain, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	path, cleanup, err := virsh.TempXML("domain-define", xmlConfig)
	if err != nil {
		return nil, err
	}
	defer cleanup()
	if _, err := virsh.Run(h.Url, "define", path); err != nil {
		return nil, err
	}
	xmlObj := &DomainXML{}
	_ = libstar.XML.Decode(xmlObj, xmlConfig)
	if xmlObj.Name == "" {
		return nil, libstar.NewErr("domain name missing in xml")
	}
	return h.LookupDomainByName(xmlObj.Name)
}

func (h *HyperVisor) Close() {
	if h.ConnedTime <= 0 {
		return
	}
	for _, listen := range h.Listener {
		if listen.Closed != nil {
			_ = listen.Closed(h)
		}
	}
	h.ConnedTime = 0
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
	h, err := GetHyper()
	if err != nil {
		return nil, err
	}
	return h.LookupDomainByUUIDString(uuid)
}

func LookupDomainByUUIDName(uuid string) (*Domain, error) {
	h, err := GetHyper()
	if err != nil {
		return nil, err
	}
	return h.LookupDomainByUUIDName(uuid)
}

func AddHyperListener(listen HyperListener) {
	hyper.AddListener(listen)
}

func init() {
	hyper.SetUrl("qemu:///system")
	go hyper.LoopForever()
}
