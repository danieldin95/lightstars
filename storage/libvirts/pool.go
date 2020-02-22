package libvirts

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
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
	Conn *libvirt.Connect
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

func (pol *Pool) Open() error {
	if pol.Conn == nil {
		hyper, err := libvirtc.GetHyper()
		if err != nil {
			return err
		}
		pol.Conn = hyper.Conn
	}
	if pol.Conn == nil {
		return libstar.NewErr("Not found libvirt.Connect")
	}
	return nil
}

func (pol *Pool) Create() error {
	if err := pol.Open(); err != nil {
		return err
	}
	if _, err := pol.Conn.LookupStoragePoolByName(pol.Name); err == nil {
		return nil
	}
	polXml := PoolXML{
		Type: pol.Type,
		Name: pol.Name,
		Target: TargetXML{
			Path: pol.Path,
		},
	}
	if _, err := pol.Conn.StoragePoolCreateXML(polXml.Encode(), 0); err != nil {
		return err
	}
	return nil
}
