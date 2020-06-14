package storage

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/storage/libvirts"
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

	hyper, err := libvirts.GetHyper()
	if err != nil {
		libstar.Warn("IsoMgr.ListFiles %s", err)
		return images
	}

	pool, err := hyper.Conn.LookupStoragePoolByTargetPath(dir)
	if err != nil {
		name := path.Base(dir)
		libstar.Warn("IsoMgr.ListFiles %s, and try %s", err, name)
		pool, err = hyper.Conn.LookupStoragePoolByName(name)
		if err != nil {
			return images
		}
	}

	defer pool.Free()
	_ = pool.Refresh(0)
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
					Path: PATH.Fmt(file),
				})
			}
			_ = vol.Free()
		}
	}
	return images
}

var ISO = IsoMgr{
	Files: make([]IsoFile, 0, 32),
}

type Store struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	State      int    `json:"state"`
	Capacity   uint64 `json:"capacity"`
	Allocation uint64 `json:"allocation"`
	Available  uint64 `json:"available"`
}

type StoreMgr struct {
	Store []Store `json:"storage"`
}

func (store *StoreMgr) Init() {
	libvirts.AddHyperListener(libvirts.HyperListener{
		Opened: func(Conn *libvirt.Connect) error {
			_, err := libvirts.CreatePool("01", PATH.Unix("datastore@01"))
			if err != nil {
				libstar.Error("StoreMgr.Init CreatePool %s", err)
			}
			return nil
		},
		Closed: nil,
	})
}

func (store *StoreMgr) List() []Store {
	stores := make([]Store, 0, 32)

	hyper, err := libvirts.GetHyper()
	if err != nil {
		libstar.Warn("StoreMgr.List %s", err)
		return stores
	}
	if pools, err := hyper.Conn.ListAllStoragePools(0); err == nil {
		for _, pool := range pools {
			name, err := pool.GetName()
			if err != nil {
				continue
			}
			if libvirts.IsDomainPool(name) {
				_ = pool.Free()
				continue
			}
			info, err := pool.GetInfo()
			if err == nil {
				path := DataStore + name
				stores = append(stores, Store{
					Name:       path,
					Path:       path,
					State:      int(info.State),
					Capacity:   info.Capacity,
					Allocation: info.Allocation,
					Available:  info.Available,
				})
			}
			_ = pool.Free()
		}
	}
	return stores
}

var DATASTOR = StoreMgr{
	Store: make([]Store, 0, 32),
}

type File struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type libvirt.StorageVolType `json:"type"`
}

type FileMgr struct {
	Files []File `json:"files"`
}

func (f *FileMgr) List(dir string) []File {
	files := make([]File, 0, 32)

	hyper, err := libvirts.GetHyper()
	if err != nil {
		libstar.Warn("FileMgr.ListFiles %s", err)
		return files
	}

	pool, err := hyper.Conn.LookupStoragePoolByTargetPath(dir)
	if err != nil {
		name := path.Base(dir)
		libstar.Warn("FileMgr.ListFiles %s, and try %s", err, name)
		pool, err = hyper.Conn.LookupStoragePoolByName(name)
		if err != nil {
			return files
		}
	}

	defer pool.Free()
	_ = pool.Refresh(0)
	if vols, err := pool.ListAllStorageVolumes(0); err == nil {
		for _, vol := range vols {
			i, _ := vol.GetInfo()
			libstar.Info("%s", &i.Type)
			file, err := vol.GetPath()
			if err != nil {
				continue
			}

			files = append(files, File{
				Name: path.Base(file),
				Path: PATH.Fmt(file),
				Type: i.Type,
			})

			_ = vol.Free()
		}
	}
	return files
}

var FILE = FileMgr{
	Files: make([]File, 0, 32),
}