package libvirts

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/storage"
	"github.com/libvirt/libvirt-go"
	"path"
	"strings"
)

type IsoFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type IsoMgr struct {
	Files []IsoFile `json:"files"`
}

func (iso *IsoMgr) ListFiles(dir string) []IsoFile {
	images := make([]IsoFile, 0, 32)

	hyper, err := GetHyper()
	if err != nil {
		libstar.Warn("IsoMgr.ListFiles %s", err)
		return images
	}

	pool, err := hyper.Conn.LookupStoragePoolByTargetPath(dir)
	if err != nil {
		libstar.Warn("IsoMgr.ListFiles %s", err)
		return images
	}
	defer pool.Free()
	pool.Refresh(0)
	if vols, err := pool.ListAllStorageVolumes(0); err == nil {
		for _, vol := range vols {
			file, err := vol.GetPath()
			if err != nil {
				continue
			}
			name := strings.ToUpper(file)
			if strings.HasSuffix(name, ".ISO") ||
				strings.HasSuffix(name, ".IMG") ||
				strings.HasSuffix(name, ".QCOW2") ||
				strings.HasSuffix(name, ".RAW") ||
				strings.HasSuffix(name, ".VMDK") {
				images = append(images, IsoFile{
					Name: path.Base(file),
					Path: storage.PATH.Fmt(file),
				})
			}
			vol.Free()
		}
	}
	return images
}

var ISO = IsoMgr{
	Files: make([]IsoFile, 0, 32),
}

type DataStore struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	State      int    `json:"state"`
	Capacity   uint64 `json:"capacity"`
	Allocation uint64 `json:"allocation"`
	Available  uint64 `json:"available"`
}

type DataStoreMgr struct {
	Storage []DataStore `json:"storage"`
}

func (store *DataStoreMgr) Init() {
	AddHyperListener(HyperListener{
		Opened: func(Conn *libvirt.Connect) error {
			_, err := CreatePool("01", storage.PATH.Unix("datastore@01"))
			if err != nil {
				libstar.Error("DataStoreMgr.Init CreatePool %s", err)
			}
			return nil
		},
		Closed: nil,
	})
}

func (store *DataStoreMgr) List() []DataStore {
	stores := make([]DataStore, 0, 32)

	hyper, err := GetHyper()
	if err != nil {
		libstar.Warn("IsoMgr.ListFiles %s", err)
		return stores
	}
	if pools, err := hyper.Conn.ListAllStoragePools(0); err == nil {
		for _, pool := range pools {
			name, err := pool.GetName()
			if err != nil {
				continue
			}
			if IsDomainPool(name) {
				pool.Free()
				continue
			}

			info, err := pool.GetInfo()
			if err == nil {
				path := storage.DataStore + name
				stores = append(stores, DataStore{
					Name:       path,
					Path:       path,
					State:      int(info.State),
					Capacity:   info.Capacity,
					Allocation: info.Allocation,
					Available:  info.Available,
				})
			}
			pool.Free()
		}
	}
	return stores
}

var DATASTOR = DataStoreMgr{
	Storage: make([]DataStore, 0, 32),
}
