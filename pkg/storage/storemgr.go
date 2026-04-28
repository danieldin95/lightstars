package storage

import (
	"path"
	"strings"

	"github.com/danieldin95/lightstar/pkg/libstar"
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

	pool, err := LookupPoolByTargetPath(dir)
	if err != nil {
		name := path.Base(dir)
		libstar.Warn("IsoMgr.ListFiles %s, and try %s", err, name)
		pool, err = LookupPoolByUUIDOrName(name)
		if err != nil {
			return images
		}
	}

	vols, err := pool.List()
	if err != nil {
		return images
	}
	for file := range vols {
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
	}
	return images
}

var ISO = IsoMgr{Files: make([]IsoFile, 0, 32)}

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
	AddHyperListener(HyperListener{
		Opened: func(_ *HyperVisor) error {
			_, err := CreatePool("01", PATH.Unix("datastore@01"))
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
	pools, err := ListPools()
	if err != nil {
		libstar.Warn("StoreMgr.List %s", err)
		return stores
	}
	for _, pool := range pools {
		name := pool.Name
		if IsDomainPool(name) {
			continue
		}
		info, err := pool.GetInfo()
		if err != nil {
			continue
		}
		p := DataStore + name
		stores = append(stores, Store{
			Name:       p,
			Path:       p,
			State:      int(info.State),
			Capacity:   info.Capacity,
			Allocation: info.Allocation,
			Available:  info.Available,
		})
	}
	return stores
}

var DATASTOR = StoreMgr{Store: make([]Store, 0, 32)}
