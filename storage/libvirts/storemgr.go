package libvirts

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
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
	Conn  *libvirt.Connect `json:"-"`
	Files []IsoFile        `json:"files"`
}

func (iso *IsoMgr) Open() error {
	if iso.Conn == nil {
		hyper, err := libvirtc.GetHyper()
		if err != nil {
			return err
		}
		iso.Conn = hyper.Conn
	}
	if iso.Conn == nil {
		return libstar.NewErr("Not found libvirt.Connect")
	}
	return nil
}

func (iso *IsoMgr) ListFiles(dir string) []IsoFile {
	images := make([]IsoFile, 0, 32)

	if err := iso.Open(); err != nil {
		libstar.Warn("IsoMgr.ListFiles %s", err)
		return images
	}

	pool, err := iso.Conn.LookupStoragePoolByTargetPath(dir)
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
			if strings.HasSuffix(file, ".iso") || strings.HasSuffix(file, ".ISO") {
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
	Conn    *libvirt.Connect `json:"-"`
	Storage []DataStore      `json:"storage"`
}

func (store *DataStoreMgr) Open() error {
	if store.Conn == nil {
		hyper, err := libvirtc.GetHyper()
		if err != nil {
			return err
		}
		store.Conn = hyper.Conn
	}
	if store.Conn == nil {
		return libstar.NewErr("Not found libvirt.Connect")
	}
	return nil
}

func (store *DataStoreMgr) Init() {
	_, err := CreatePool("01", storage.PATH.Unix("datastore@01"))
	if err != nil {
		libstar.Error("DataStoreMgr.Init CreatePool %s", err)
	}
}

func (store *DataStoreMgr) List() []DataStore {
	stores := make([]DataStore, 0, 32)

	if err := store.Open(); err != nil {
		libstar.Warn("IsoMgr.ListFiles %s", err)
		return stores
	}
	if pools, err := store.Conn.ListAllStoragePools(0); err == nil {
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
