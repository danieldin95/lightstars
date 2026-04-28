package compute

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/virsh"
)

const (
	DomainAll                    = 0
	DomainDeviceModifyConfig     = 1 << 0
	DomainDeviceModifyLive       = 1 << 1
	DomainDeviceModifyPersistent = DomainDeviceModifyConfig | DomainDeviceModifyLive
	DomainDestroyGraceful        = 1 << 2
	DomainShutdownAcpi           = 1 << 3
	DomainCpuMaximum             = 1 << 4
	DomainCpuConfig              = 1 << 5
	DomainMemMaximum             = 1 << 6
	DomainMemConfig              = 1 << 7
)

type DomainState uint

const (
	DOMAIN_NOSTATE DomainState = iota
	DOMAIN_RUNNING
	DOMAIN_BLOCKED
	DOMAIN_PAUSED
	DOMAIN_SHUTDOWN
	DOMAIN_SHUTOFF
	DOMAIN_CRASHED
	DOMAIN_PMSUSPENDED
)

type DomainInfo struct {
	State     DomainState
	MaxMem    uint64
	Memory    uint64
	NrVirtCpu uint
	CpuTime   uint64
}

type Domain struct {
	Name  string
	UUID  string
	hyper *HyperVisor
}

func NewDomainFromVir(name, uuid string, hyper *HyperVisor) *Domain {
	return &Domain{Name: name, UUID: uuid, hyper: hyper}
}

func (d *Domain) domainRef() string {
	if d == nil {
		return ""
	}
	if d.Name != "" {
		return d.Name
	}
	return d.UUID
}

func (d *Domain) run(args ...string) (string, error) {
	if d == nil || d.hyper == nil {
		return "", libstar.NewErr("hyper is nil")
	}
	return virsh.Run(d.hyper.Url, args...)
}

func (d *Domain) Free() error { return nil }

func (d *Domain) GetName() (string, error) {
	if d.Name != "" {
		return d.Name, nil
	}
	out, err := d.run("domname", d.domainRef())
	if err != nil {
		return "", err
	}
	d.Name = strings.TrimSpace(out)
	return d.Name, nil
}

func (d *Domain) GetUUIDString() (string, error) {
	if d.UUID != "" {
		return d.UUID, nil
	}
	out, err := d.run("domuuid", d.domainRef())
	if err != nil {
		return "", err
	}
	d.UUID = strings.TrimSpace(out)
	return d.UUID, nil
}

func domStateFromText(state string) DomainState {
	state = strings.ToLower(strings.TrimSpace(state))
	switch state {
	case "running", "in shutdown":
		return DOMAIN_RUNNING
	case "blocked":
		return DOMAIN_BLOCKED
	case "paused", "pmsuspended":
		return DOMAIN_PAUSED
	case "shutdown":
		return DOMAIN_SHUTDOWN
	case "shut off", "shutoff", "off":
		return DOMAIN_SHUTOFF
	case "crashed":
		return DOMAIN_CRASHED
	case "no state":
		return DOMAIN_NOSTATE
	default:
		return DOMAIN_NOSTATE
	}
}

func parseUintField(v string) uint64 {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}
	fields := strings.Fields(v)
	if len(fields) == 0 {
		return 0
	}
	n, _ := strconv.ParseUint(fields[0], 10, 64)
	return n
}

func (d *Domain) GetInfo() (DomainInfo, error) {
	out, err := d.run("dominfo", d.domainRef())
	if err != nil {
		return DomainInfo{}, err
	}
	kv := virsh.ParseKV(out)
	info := DomainInfo{}
	info.State = domStateFromText(kv["state"])
	info.NrVirtCpu = uint(parseUintField(kv["cpu(s)"]))
	info.MaxMem = parseUintField(kv["max memory"])
	info.Memory = parseUintField(kv["used memory"])
	info.CpuTime = virsh.ParseSecondsNS(kv["cpu time"])
	return info, nil
}

func (d *Domain) document() (*etree.Document, error) {
	xmlData, err := d.GetXMLDesc(true)
	if err != nil {
		return nil, err
	}
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlData); err != nil {
		return nil, err
	}
	return doc, nil
}

func (d *Domain) SetCpu(max, mode string) error {
	doc, err := d.document()
	if err != nil {
		return err
	}
	domEle := doc.FindElement("/domain")
	if domEle == nil {
		return libstar.NewErr("domain tag not found")
	}
	vCpuEle := domEle.FindElement("./vcpu")
	if vCpuEle == nil {
		vCpuEle = domEle.CreateElement("vcpu")
	}
	vCpuEle.SetText(max)
	cpuEle := domEle.FindElement("./cpu")
	if cpuEle == nil {
		cpuEle = domEle.CreateElement("cpu")
	}
	cpuEle.RemoveAttr("mode")
	cpuEle.CreateAttr("mode", mode)
	newXML, err := doc.WriteToString()
	if err != nil {
		return err
	}
	return d.reDefine(newXML)
}

func (d *Domain) SetMemory(size, unit string) error {
	doc, err := d.document()
	if err != nil {
		return err
	}
	domEle := doc.FindElement("/domain")
	if domEle == nil {
		return libstar.NewErr("domain tag not found")
	}
	memEle := domEle.FindElement("./memory")
	if memEle == nil {
		memEle = domEle.CreateElement("memory")
	}
	memEle.RemoveAttr("unit")
	memEle.CreateAttr("unit", unit)
	memEle.SetText(size)

	curMemEle := domEle.FindElement("./currentMemory")
	if curMemEle == nil {
		curMemEle = domEle.CreateElement("currentMemory")
	}
	curMemEle.RemoveAttr("unit")
	curMemEle.CreateAttr("unit", unit)
	curMemEle.SetText(size)
	newXML, err := doc.WriteToString()
	if err != nil {
		return err
	}
	return d.reDefine(newXML)
}

func (d *Domain) reDefine(newXML string) error {
	oldXML, err := d.GetXMLDesc(true)
	if err != nil {
		return err
	}
	if err := d.Undefine(); err != nil {
		return err
	}
	path, cleanup, err := virsh.TempXML("domain-define", newXML)
	if err != nil {
		return err
	}
	defer cleanup()
	if _, err := d.run("define", path); err != nil {
		oldPath, oldCleanup, oldErr := virsh.TempXML("domain-rollback", oldXML)
		if oldErr == nil {
			defer oldCleanup()
			_, _ = d.run("define", oldPath)
		}
		return err
	}
	return nil
}

func (d *Domain) GetXMLDesc(secure bool) (string, error) {
	args := []string{"dumpxml", d.domainRef()}
	if secure {
		args = append(args, "--security-info")
	} else {
		args = append(args, "--inactive")
	}
	out, err := d.run(args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (d *Domain) IsActive() (bool, error) {
	out, err := d.run("domstate", d.domainRef())
	if err != nil {
		return false, err
	}
	state := strings.ToLower(strings.TrimSpace(out))
	return strings.Contains(state, "running") || strings.Contains(state, "paused"), nil
}

func (d *Domain) GetMetadataTitle(title bool) (string, error) {
	args := []string{"desc", d.domainRef(), "--current"}
	if title {
		args = append(args, "--title")
	}
	out, err := d.run(args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (d *Domain) SetMetadataTitle(value string, title bool) error {
	args := []string{"desc", d.domainRef(), value}
	if title {
		args = append(args, "--title")
	}
	active, _ := d.IsActive()
	if active {
		if _, err := d.run(append(args, "--live")...); err != nil {
			return err
		}
	}
	_, err := d.run(append(args, "--config")...)
	return err
}

func (d *Domain) Create() error {
	_, err := d.run("start", d.domainRef())
	return err
}

func (d *Domain) ShutdownFlags(_ int) error {
	_, err := d.run("shutdown", d.domainRef(), "--mode", "acpi")
	if err != nil {
		_, err = d.run("shutdown", d.domainRef())
	}
	return err
}

func (d *Domain) Suspend() error {
	_, err := d.run("suspend", d.domainRef())
	return err
}

func (d *Domain) Reset(_ int) error {
	_, err := d.run("reset", d.domainRef())
	return err
}

func (d *Domain) DestroyFlags(_ int) error {
	_, err := d.run("destroy", d.domainRef())
	return err
}

func (d *Domain) Resume() error {
	_, err := d.run("resume", d.domainRef())
	return err
}

func (d *Domain) Undefine() error {
	_, err := d.run("undefine", d.domainRef())
	return err
}

func (d *Domain) SetAutostart(enable bool) error {
	args := []string{"autostart", d.domainRef()}
	if !enable {
		args = append(args, "--disable")
	}
	_, err := d.run(args...)
	return err
}

func optionsFromFlags(flags int) []string {
	opts := make([]string, 0, 3)
	if flags&DomainDeviceModifyConfig != 0 {
		opts = append(opts, "--config")
	}
	if flags&DomainDeviceModifyLive != 0 {
		opts = append(opts, "--live")
	}
	if len(opts) == 0 {
		opts = append(opts, "--current")
	}
	return opts
}

func (d *Domain) AttachDeviceFlags(xmlData string, flags int) error {
	path, cleanup, err := virsh.TempXML("domain-attach", xmlData)
	if err != nil {
		return err
	}
	defer cleanup()
	args := []string{"attach-device", d.domainRef(), path}
	args = append(args, optionsFromFlags(flags)...)
	_, err = d.run(args...)
	return err
}

func (d *Domain) DetachDeviceFlags(xmlData string, flags int) error {
	path, cleanup, err := virsh.TempXML("domain-detach", xmlData)
	if err != nil {
		return err
	}
	defer cleanup()
	args := []string{"detach-device", d.domainRef(), path}
	args = append(args, optionsFromFlags(flags)...)
	_, err = d.run(args...)
	return err
}

type DomainSnapshot struct {
	Domain string
	Name   string
	hyper  *HyperVisor
}

func (s *DomainSnapshot) run(args ...string) (string, error) {
	if s == nil || s.hyper == nil {
		return "", libstar.NewErr("snapshot hyper is nil")
	}
	return virsh.Run(s.hyper.Url, args...)
}

func (s *DomainSnapshot) GetXMLDesc(_ int) (string, error) {
	out, err := s.run("snapshot-dumpxml", s.Domain, s.Name)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (s *DomainSnapshot) IsCurrent(_ int) (bool, error) {
	out, err := s.run("snapshot-current", s.Domain, "--name")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) == s.Name, nil
}

func (s *DomainSnapshot) RevertToSnapshot(_ int) error {
	_, err := s.run("snapshot-revert", s.Domain, s.Name)
	return err
}

func (s *DomainSnapshot) Delete(_ int) error {
	_, err := s.run("snapshot-delete", s.Domain, s.Name)
	return err
}

func (s *DomainSnapshot) Free() error { return nil }

func (d *Domain) CreateSnapshotXML(xmlData string, _ int) (*DomainSnapshot, error) {
	path, cleanup, err := virsh.TempXML("domain-snapshot", xmlData)
	if err != nil {
		return nil, err
	}
	defer cleanup()
	if _, err := d.run("snapshot-create", d.domainRef(), path); err != nil {
		return nil, err
	}
	obj := &SnapshotXML{}
	_ = libstar.XML.Decode(obj, xmlData)
	if obj.Name == "" {
		if out, err := d.run("snapshot-current", d.domainRef(), "--name"); err == nil {
			obj.Name = strings.TrimSpace(out)
		}
	}
	return &DomainSnapshot{Domain: d.domainRef(), Name: obj.Name, hyper: d.hyper}, nil
}

func (d *Domain) SnapshotLookupByName(name string, _ int) (*DomainSnapshot, error) {
	if _, err := d.run("snapshot-info", d.domainRef(), name); err != nil {
		return nil, err
	}
	return &DomainSnapshot{Domain: d.domainRef(), Name: name, hyper: d.hyper}, nil
}

func (d *Domain) ListAllSnapshots(_ int) ([]DomainSnapshot, error) {
	out, err := d.run("snapshot-list", d.domainRef(), "--name")
	if err != nil {
		return nil, err
	}
	lines := virsh.Lines(out)
	ret := make([]DomainSnapshot, 0, len(lines))
	for _, line := range lines {
		ret = append(ret, DomainSnapshot{Domain: d.domainRef(), Name: line, hyper: d.hyper})
	}
	return ret, nil
}

func ListDomains() ([]Domain, error) {
	h, err := GetHyper()
	if err != nil {
		return nil, err
	}
	return h.ListAllDomains()
}

func DomainState2Str(state DomainState) string {
	switch state {
	case DOMAIN_NOSTATE:
		return "no-state"
	case DOMAIN_RUNNING:
		return "running"
	case DOMAIN_BLOCKED:
		return "blocked"
	case DOMAIN_PAUSED:
		return "paused"
	case DOMAIN_SHUTDOWN:
		return "shutdown"
	case DOMAIN_CRASHED:
		return "crashed"
	case DOMAIN_PMSUSPENDED:
		return "pm-suspended"
	case DOMAIN_SHUTOFF:
		return "shutoff"
	default:
		return fmt.Sprintf("unknown(%d)", state)
	}
}
