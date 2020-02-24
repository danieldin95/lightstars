package libvirts

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/libvirt/libvirt-go"
	"strconv"
)

type Volume struct {
	Conn     *libvirt.Connect
	Pool     string
	Name     string
	Size     uint64
	Format   string
	BackFile string
}

func NewVolume(pool, name string, size uint64) Volume {
	return Volume{
		Pool:   pool,
		Name:   name,
		Size:   size,
		Format: "qcow2",
	}
}

func CreateVolume(pool, name string, size uint64) (*Volume, error) {
	vol := &Volume{
		Pool:   pool,
		Name:   name,
		Size:   size,
		Format: "qcow2",
	}
	return vol, vol.Create()
}

func RemoveVolume(pool string, name string) error {
	vol := &Volume{
		Pool: pool,
		Name: name,
	}
	return vol.Remove()
}

func (vol *Volume) Open() error {
	if vol.Conn == nil {
		hyper, err := libvirtc.GetHyper()
		if err != nil {
			return err
		}
		vol.Conn = hyper.Conn
	}
	if vol.Conn == nil {
		return libstar.NewErr("Not found libvirt.Connect")
	}
	return nil
}

func (vol *Volume) Create() error {
	if err := vol.Open(); err != nil {
		return err
	}
	volXml := VolumeXML{
		Name: vol.Name,
		Capacity: CapacityXML{
			Unit:  "bytes",
			Value: strconv.FormatUint(vol.Size, 10),
		},
		Target: TargetXML{
			Format: FormatXML{
				Type: vol.Format,
			},
		},
	}
	pool, err := vol.Conn.LookupStoragePoolByName(vol.Pool)
	if err != nil {
		return err
	}
	defer pool.Free()
	volume, err := pool.StorageVolCreateXML(volXml.Encode(), 0)
	if err != nil {
		return err
	}
	defer volume.Free()

	return nil
}

func (vol *Volume) GetXMLObj() (*VolumeXML, error) {
	if err := vol.Open(); err != nil {
		return nil, err
	}
	pool, err := vol.Conn.LookupStoragePoolByName(vol.Pool)
	if err != nil {
		return nil, err
	}
	defer pool.Free()
	volume, err := pool.LookupStorageVolByName(vol.Name)
	if err != nil {
		return nil, err
	}
	defer volume.Free()
	xmlData, err := volume.GetXMLDesc(0)
	if err != nil {
		return nil, err
	}
	xmlObj := &VolumeXML{}
	return xmlObj, xmlObj.Decode(xmlData)
}

func (vol *Volume) Remove() error {
	if err := vol.Open(); err != nil {
		return err
	}
	pool, err := vol.Conn.LookupStoragePoolByName(vol.Pool)
	if err != nil {
		return err
	}
	defer pool.Free()

	volume, err := pool.LookupStorageVolByName(vol.Name)
	if err != nil {
		return err
	}
	volume.Delete(0)
	defer volume.Free()

	return nil
}
