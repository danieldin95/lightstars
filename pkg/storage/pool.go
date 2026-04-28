package storage

import (
	"strings"

	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/virsh"
)

func ToDomainPool(domain string) string {
	return "." + domain
}

func IsDomainPool(name string) bool {
	return strings.HasPrefix(name, ".")
}

func IsStorePool(name string) bool {
	if IsDomainPool(name) {
		return false
	}
	if len(name) <= 0 || len(name) > 2 {
		return false
	}
	if libstar.IsDigit(name) {
		return true
	}
	return false
}

type StoragePoolState int

const (
	STORAGE_POOL_INACTIVE StoragePoolState = iota
	STORAGE_POOL_BUILDING
	STORAGE_POOL_RUNNING
	STORAGE_POOL_DEGRADED
	STORAGE_POOL_INACCESSIBLE
)

type PoolInfo struct {
	State      StoragePoolState
	Capacity   uint64
	Allocation uint64
	Available  uint64
}

type Pool struct {
	Type string
	Name string
	UUID string
	Size uint64
	Path string
	XML  string

	hyper *HyperVisor
}

func (pol *Pool) run(args ...string) (string, error) {
	h, err := GetHyper()
	if err != nil {
		return "", err
	}
	if pol.hyper == nil {
		pol.hyper = h
	}
	return virsh.Run(pol.hyper.Name, args...)
}

func (pol *Pool) Free() error { return nil }

func LookupPoolByUUID(uuid string) (*Pool, error) {
	h, err := GetHyper()
	if err != nil {
		return nil, err
	}
	nameOut, err := virsh.Run(h.Name, "pool-name", uuid)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(nameOut)
	return &Pool{Name: name, UUID: strings.TrimSpace(uuid), hyper: h}, nil
}

func LookupPoolByUUIDOrName(uuid string) (*Pool, error) {
	h, err := GetHyper()
	if err != nil {
		return nil, err
	}
	name := uuid
	if out, err := virsh.Run(h.Name, "pool-name", uuid); err == nil {
		name = strings.TrimSpace(out)
	}
	uuidOut, _ := virsh.Run(h.Name, "pool-uuid", name)
	return &Pool{Name: name, UUID: strings.TrimSpace(uuidOut), hyper: h}, nil
}

func LookupPoolByTargetPath(target string) (*Pool, error) {
	pools, err := ListPools()
	if err != nil {
		return nil, err
	}
	for _, p := range pools {
		xmlData, err := p.GetXMLDesc(0)
		if err != nil {
			continue
		}
		obj := &PoolXML{}
		if err := libstar.XML.Decode(obj, xmlData); err != nil {
			continue
		}
		if obj.Target.Path == target {
			p := p
			return &p, nil
		}
	}
	return nil, libstar.NewErr("pool not found by target path")
}

func NewPoolFromVir(name, uuid string, hyper *HyperVisor) *Pool {
	return &Pool{Name: name, UUID: uuid, hyper: hyper}
}

func CreatePool(name, target string) (*Pool, error) {
	pol := &Pool{Type: "dir", Name: name, Path: target}
	xmlObj := &PoolXML{Type: pol.Type, Name: pol.Name, Target: TargetXML{Path: pol.Path}}
	pol.XML = libstar.XML.Encode(xmlObj)
	return pol, pol.Create()
}

func (pol *Pool) Create() error {
	h, err := GetHyper()
	if err != nil {
		return err
	}
	pol.hyper = h
	if _, err := virsh.Run(h.Name, "pool-info", pol.Name); err == nil {
		return nil
	}
	file, cleanup, err := virsh.TempXML("pool-define", pol.XML)
	if err != nil {
		return err
	}
	defer cleanup()
	if _, err := virsh.Run(h.Name, "pool-define", file); err != nil {
		return err
	}
	_, _ = virsh.Run(h.Name, "pool-build", pol.Name)
	if _, err := virsh.Run(h.Name, "pool-start", pol.Name); err != nil {
		return err
	}
	if _, err := virsh.Run(h.Name, "pool-autostart", pol.Name); err != nil {
		libstar.Warn("Pool.Create SetAutoStart %s", err)
	}
	return nil
}

func (pol *Pool) Clean() error {
	p, err := LookupPoolByUUIDOrName(pol.Name)
	if err != nil {
		return nil
	}
	vols, err := p.listVolumeNames()
	if err == nil {
		for _, vol := range vols {
			if _, err := p.run("vol-delete", vol, "--pool", p.Name); err != nil {
				return err
			}
		}
	}
	if _, err := p.run("pool-destroy", p.Name); err != nil {
		libstar.Warn("Pool.Remove %s", err)
	}
	if _, err := p.run("pool-delete", p.Name); err != nil {
		libstar.Warn("Pool.Delete %s", err)
	}
	if _, err := p.run("pool-undefine", p.Name); err != nil {
		libstar.Warn("Pool.Remove %s", err)
	}
	return nil
}

func (pol *Pool) Remove() error {
	p, err := LookupPoolByUUIDOrName(pol.Name)
	if err != nil {
		return nil
	}
	if _, err := p.run("pool-destroy", p.Name); err != nil {
		libstar.Warn("Pool.Remove %s", err)
	}
	if _, err := p.run("pool-delete", p.Name); err != nil {
		libstar.Warn("Pool.Delete %s", err)
	}
	if _, err := p.run("pool-undefine", p.Name); err != nil {
		libstar.Warn("Pool.Remove %s", err)
	}
	return nil
}

func (pol *Pool) listVolumeNames() ([]string, error) {
	_, _ = pol.run("pool-refresh", pol.Name)
	out, err := pol.run("vol-list", "--pool", pol.Name)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0, 32)
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Name") || strings.HasPrefix(line, "-") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		ret = append(ret, fields[0])
	}
	return ret, nil
}

func (pol *Pool) list() (map[string]VolumeInfo, error) {
	_, _ = pol.run("pool-refresh", pol.Name)
	out, err := pol.run("vol-list", "--pool", pol.Name)
	if err != nil {
		return nil, err
	}
	infos := make(map[string]VolumeInfo, 32)
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Name") || strings.HasPrefix(line, "-") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		name := fields[0]
		path := fields[len(fields)-1]
		infoOut, err := pol.run("vol-info", name, "--pool", pol.Name)
		if err != nil {
			continue
		}
		kv := virsh.ParseKV(infoOut)
		infos[path] = VolumeInfo{
			Pool:       pol.Name,
			Name:       name,
			Type:       strings.ToLower(kv["type"]),
			Allocation: virsh.ParseBytes(kv["allocation"]),
			Capacity:   virsh.ParseBytes(kv["capacity"]),
		}
	}
	return infos, nil
}

func (pol *Pool) List() (map[string]VolumeInfo, error) {
	p, err := LookupPoolByUUIDOrName(pol.Name)
	if err != nil {
		return nil, err
	}
	return p.list()
}

func (pol *Pool) ListByTarget() (map[string]VolumeInfo, error) {
	p, err := LookupPoolByTargetPath(pol.Path)
	if err != nil {
		return nil, err
	}
	return p.list()
}

func ListPools() ([]Pool, error) {
	h, err := GetHyper()
	if err != nil {
		return nil, err
	}
	return h.ListAllPools()
}

func poolStateFromText(text string) StoragePoolState {
	s := strings.ToLower(strings.TrimSpace(text))
	switch {
	case strings.Contains(s, "running"):
		return STORAGE_POOL_RUNNING
	case strings.Contains(s, "building"):
		return STORAGE_POOL_BUILDING
	case strings.Contains(s, "degraded"):
		return STORAGE_POOL_DEGRADED
	case strings.Contains(s, "inaccessible"):
		return STORAGE_POOL_INACCESSIBLE
	default:
		return STORAGE_POOL_INACTIVE
	}
}

func (pol *Pool) GetInfo() (PoolInfo, error) {
	out, err := pol.run("pool-info", pol.Name)
	if err != nil {
		return PoolInfo{}, err
	}
	kv := virsh.ParseKV(out)
	return PoolInfo{
		State:      poolStateFromText(kv["state"]),
		Capacity:   virsh.ParseBytes(kv["capacity"]),
		Allocation: virsh.ParseBytes(kv["allocation"]),
		Available:  virsh.ParseBytes(kv["available"]),
	}, nil
}

func (pol *Pool) GetXMLDesc(flags uint32) (string, error) {
	args := []string{"pool-dumpxml", pol.Name}
	if flags != 0 {
		args = append(args, "--inactive")
	}
	out, err := pol.run(args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (pol *Pool) IsActive() (bool, error) {
	info, err := pol.GetInfo()
	if err != nil {
		return false, err
	}
	return info.State == STORAGE_POOL_RUNNING, nil
}

func (pol *Pool) IsAutostart() (bool, error) {
	out, err := pol.run("pool-info", pol.Name)
	if err != nil {
		return false, err
	}
	val := strings.ToLower(strings.TrimSpace(virsh.ParseKV(out)["autostart"]))
	return val == "yes", nil
}

func (pol *Pool) Start() error {
	_, err := pol.run("pool-start", pol.Name)
	return err
}

func (pol *Pool) Destroy() error {
	_, err := pol.run("pool-destroy", pol.Name)
	return err
}

func (pol *Pool) Refresh() error {
	_, err := pol.run("pool-refresh", pol.Name)
	return err
}

func (pol *Pool) SetAutostart(enable bool) error {
	args := []string{"pool-autostart", pol.Name}
	if !enable {
		args = append(args, "--disable")
	}
	_, err := pol.run(args...)
	return err
}

func PoolState2Str(state StoragePoolState) string {
	switch state {
	case STORAGE_POOL_BUILDING:
		return "building"
	case STORAGE_POOL_INACTIVE:
		return "inactive"
	case STORAGE_POOL_RUNNING:
		return "running"
	case STORAGE_POOL_DEGRADED:
		return "degraded"
	case STORAGE_POOL_INACCESSIBLE:
		return "inaccessible"
	default:
		return "unknown"
	}
}
