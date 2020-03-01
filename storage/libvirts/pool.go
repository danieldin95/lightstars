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

type Pool struct {
	libvirt.StoragePool
	Type string
	Name string
	Size uint64
	Path string
}

func NewPool(name, target string) Pool {
	return Pool{
		Type: "dir",
		Name: name,
		Path: target,
	}
}

func LookupPoolByUUID(uuid string) (*Pool, error) {
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
	polXml := PoolXML{
		Type: pol.Type,
		Name: pol.Name,
		Target: TargetXML{
			Path: pol.Path,
		},
	}
	xml := polXml.Encode()
	pool, err := hyper.Conn.StoragePoolDefineXML(xml, 0)
	if err != nil {
		return err
	}
	if err := pool.Create(libvirt.STORAGE_POOL_CREATE_WITH_BUILD); err != nil {
		return err
	}
	if err := pool.SetAutostart(true); err != nil {
		libstar.Warn("Pool.Create SetAutoStart %s", err)
	}
	defer pool.Free()
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
		vols, err := pool.ListAllStorageVolumes(0)
		if err == nil {
			for _, vol := range vols {
				if err := vol.Delete(0); err != nil {
					return err
				}
				vol.Free()
			}
		}
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
