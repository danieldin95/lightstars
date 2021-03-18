package libvirts

import (
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/libvirt/libvirt-go"
	"strconv"
)

type Volume struct {
	Path          string
	Pool          string
	Name          string
	Size          uint64
	Format        string
	BackingFile   string
	BackingFormat string
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

func CreateBackingVolume(pool, name, backingFle, backingFmt string) (*Volume, error) {
	vol := &Volume{
		Pool:          pool,
		Name:          name,
		BackingFile:   backingFle,
		BackingFormat: backingFmt,
		Format:        "qcow2",
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

func (vol *Volume) Create() error {
	hyper, err := GetHyper()
	if err != nil {
		return err
	}
	volXml := &VolumeXML{
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
		BackingStore: BackingStoreXML{
			Path: vol.BackingFile,
			Format: FormatXML{
				Type: vol.BackingFormat,
			},
		},
	}
	pool, err := hyper.Conn.LookupStoragePoolByName(vol.Pool)
	if err != nil {
		pool, err = hyper.Conn.LookupStoragePoolByTargetPath(vol.Pool)
		if err != nil {
			return err
		}
	}
	defer pool.Free()
	volume, err := pool.StorageVolCreateXML(libstar.XML.Encode(volXml), 0)
	if err != nil {
		return err
	}
	defer volume.Free()
	if vol.Path, err = volume.GetPath(); err != nil {
		return err
	}
	return nil
}

func (vol *Volume) GetXMLObj() (*VolumeXML, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	pool, err := hyper.Conn.LookupStoragePoolByName(vol.Pool)
	if err != nil {
		pool, err = hyper.Conn.LookupStoragePoolByTargetPath(vol.Pool)
		if err != nil {
			return nil, err
		}
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
	return xmlObj, libstar.XML.Decode(xmlObj, xmlData)
}

func (vol *Volume) Remove() error {
	hyper, err := GetHyper()
	if err != nil {
		return err
	}
	pool, err := hyper.Conn.LookupStoragePoolByName(vol.Pool)
	if err != nil {
		pool, err = hyper.Conn.LookupStoragePoolByTargetPath(vol.Pool)
		if err != nil {
			return err
		}
	}
	defer pool.Free()

	volume, err := pool.LookupStorageVolByName(vol.Name)
	if err != nil {
		return err
	}
	_ = volume.Delete(0)
	defer volume.Free()

	return nil
}

func VolumeType(t libvirt.StorageVolType) string {
	switch t {
	case libvirt.STORAGE_VOL_FILE:
		return "file"
	case libvirt.STORAGE_VOL_BLOCK:
		return "block"
	case libvirt.STORAGE_VOL_DIR:
		return "dir"
	case libvirt.STORAGE_VOL_NETDIR:
		return "netdir"
	case libvirt.STORAGE_VOL_NETWORK:
		return "network"
	case libvirt.STORAGE_VOL_PLOOP:
		return "ploop"
	default:
		return ""
	}
}

type VolumeInfo struct {
	Pool       string
	Name       string
	Type       string
	Capacity   uint64
	Allocation uint64
}
