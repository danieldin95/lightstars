package libvirts

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/libvirt/libvirt-go"
	"strings"
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

type Pool struct {
	libvirt.StoragePool
	Type string
	Name string
	Size uint64
	Path string
	XML  string
}

func NewPool(name, target string) Pool {
	return Pool{
		Type: "dir",
		Name: name,
		Path: target,
	}
}

func LookupPoolByUUID(uuid string) (*Pool, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	pool, err := hyper.Conn.LookupStoragePoolByUUIDString(uuid)
	if err != nil {
		return nil, err
	}
	return NewPoolFromVir(pool), nil
}

func NewPoolFromVir(pool *libvirt.StoragePool) *Pool {
	return &Pool{StoragePool: *pool}
}

func CreatePool(name, target string) (*Pool, error) {
	pol := &Pool{
		Type: "dir",
		Name: name,
		Path: target,
	}
	polXml := PoolXML{
		Type: pol.Type,
		Name: pol.Name,
		Target: TargetXML{
			Path: pol.Path,
		},
	}
	pol.XML = polXml.Encode()
	return pol, pol.Create()
}

func RemovePool(name string) error {
	pol := &Pool{
		Name: name,
	}
	return pol.Remove()
}

func (pol *Pool) Create() error {
	hyper, err := GetHyper()
	if err != nil {
		return err
	}
	if _, err := hyper.Conn.LookupStoragePoolByName(pol.Name); err == nil {
		return nil
	}
	pool, err := hyper.Conn.StoragePoolDefineXML(pol.XML, 0)
	if err != nil {
		return err
	}
	if err := pool.Create(libvirt.STORAGE_POOL_CREATE_WITH_BUILD); err != nil {
		return err
	}
	if err := pool.SetAutostart(true); err != nil {
		libstar.Warn("Pool.Create SetAutoStart %s", err)
	}
	_ = pool.Free()
	return nil
}

func (pol *Pool) Clean() error {
	hyper, err := GetHyper()
	if err != nil {
		return err
	}
	pool, err := hyper.Conn.LookupStoragePoolByUUIDString(pol.Name)
	if err != nil {
		pool, err = hyper.Conn.LookupStoragePoolByName(pol.Name)
	}
	if err == nil {
		vols, err := pool.ListAllStorageVolumes(0)
		if err == nil {
			for _, vol := range vols {
				if err := vol.Delete(0); err != nil {
					_ = vol.Free()
					return err
				}
				_ = vol.Free()
			}
		}
		defer pool.Free()
	}
	return nil
}

func (pol *Pool) Remove() error {
	hyper, err := GetHyper()
	if err != nil {
		return err
	}
	pool, err := hyper.Conn.LookupStoragePoolByUUIDString(pol.Name)
	if err != nil {
		pool, err = hyper.Conn.LookupStoragePoolByName(pol.Name)
	}
	if err == nil {
		if err := pool.Destroy(); err != nil {
			libstar.Warn("Pool.Remove %s", err)
		}
		if err := pool.Undefine(); err != nil {
			libstar.Warn("Pool.Remove %s", err)
		}
		defer pool.Free()
	}
	return nil
}

func (pol *Pool) list(pool *libvirt.StoragePool) (map[string]VolumeInfo, error) {
	vols, err := pool.ListAllStorageVolumes(0)
	if err != nil {
		return nil, err
	}
	infos := make(map[string]VolumeInfo, 32)
	for _, vol := range vols {
		name, err := vol.GetPath()
		if err != nil {
			continue
		}
		info, err := vol.GetInfo()
		if err != nil {
			vol.Free()
			continue
		}
		path, _ := vol.GetPath()
		infos[path] = VolumeInfo{
			Pool:       pol.Name,
			Name:       name,
			Type:       VolumeType(info.Type),
			Allocation: info.Allocation,
			Capacity:   info.Capacity,
		}
		_ = vol.Free()
	}
	return infos, nil
}

func (pol *Pool) List() (map[string]VolumeInfo, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	pool, err := hyper.Conn.LookupStoragePoolByUUIDString(pol.Name)
	if err != nil {
		pool, err = hyper.Conn.LookupStoragePoolByName(pol.Name)
	}
	if err != nil {
		return nil, err
	}
	defer pool.Free()
	return pol.list(pool)
}

func (pol *Pool) ListByTarget() (map[string]VolumeInfo, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	pool, err := hyper.Conn.LookupStoragePoolByTargetPath(pol.Path)
	if err != nil {
		return nil, err
	}
	defer pool.Free()
	return pol.list(pool)
}

func ListPools() ([]Pool, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	return hyper.ListAllPools()
}

func PoolState2Str(state libvirt.StoragePoolState) string {
	switch state {
	case libvirt.STORAGE_POOL_BUILDING:
		return "building"
	case libvirt.STORAGE_POOL_INACTIVE:
		return "inactive"
	case libvirt.STORAGE_POOL_RUNNING:
		return "running"
	case libvirt.STORAGE_POOL_DEGRADED:
		return "degraded"
	case libvirt.STORAGE_POOL_INACCESSIBLE:
		return "inaccessible"
	default:
		return "unknown"
	}
}
