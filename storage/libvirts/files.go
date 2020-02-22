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
	Conn  *libvirt.Connect
	Files []IsoFile `json:"files"`
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

	pol, err := iso.Conn.LookupStoragePoolByTargetPath(dir)
	if err != nil {
		libstar.Warn("IsoMgr.ListFiles %s", err)
		return images
	}
	if vols, err := pol.ListAllStorageVolumes(0); err == nil {
		for _, vol := range vols {
			file, err := vol.GetPath()
			if err != nil {
				continue
			}
			if strings.HasSuffix(file, ".iso") || strings.HasSuffix(file, ".ISO") {
				stdFile := IsoFile{
					Name: path.Base(file),
					Path: storage.PATH.Fmt(file),
				}
				images = append(images, stdFile)
			}
		}
	}
	return images
}

var ISO = IsoMgr{
	Files: make([]IsoFile, 0, 32),
}

type DataStore struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type DataStoreMgr struct {
	Conn    *libvirt.Connect
	Storage []DataStore `json:"storage"`
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

func (store *DataStoreMgr) List() []DataStore {
	stores := make([]DataStore, 0, 32)

	if err := store.Open(); err != nil {
		libstar.Warn("IsoMgr.ListFiles %s", err)
		return stores
	}
	if pools, err := store.Conn.ListAllStoragePools(0); err == nil {
		for _, pol := range pools {
			name, err := pol.GetName()
			if err != nil || strings.HasPrefix(name, ".") {
				continue
			}
			path := storage.DataStore + name
			stores = append(stores, DataStore{Name: path, Path: path})
		}
	}
	return stores
}

var DATASTOR = DataStoreMgr{
	Storage: make([]DataStore, 0, 32),
}
