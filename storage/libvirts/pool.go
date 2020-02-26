package libvirts

import (
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
	pool, err := hyper.Conn.StoragePoolCreateXML(polXml.Encode(), libvirt.STORAGE_POOL_CREATE_WITH_BUILD)
	if err != nil {
		return err
	}
	defer pool.Free()
	return nil
}

func (pol *Pool) Remove() error {
	hyper, err := GetHyper()
	if err != nil {
		return err
	}
	if pool, err := hyper.Conn.LookupStoragePoolByName(pol.Name); err == nil {
		vols, err := pool.ListAllStorageVolumes(0)
		if err != nil {
			return err
		}
		for _, vol := range vols {
			if err := vol.Delete(0); err != nil {
				return err
			}
			vol.Free()
		}
		//pool.Delete(0)
		pool.Destroy()
		defer pool.Free()
	}
	return nil
}
