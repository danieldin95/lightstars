package libvirts

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/storage"
	"path"
)

type IsoFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type IsoMgr struct {
	Files []IsoFile `json:"files"`
}

func (iso IsoMgr) ListFiles(dir string) []IsoFile {
	images := make([]IsoFile, 0, 32)

	if files, err := libstar.DIR.ListFiles(dir, ".iso"); err == nil {
		for _, file := range files {
			stdFile := IsoFile{
				Name: path.Base(file),
				Path: storage.PATH.Fmt(file),
			}
			images = append(images, stdFile)
		}
	}
	if files, err := libstar.DIR.ListFiles(dir+"/iso", ".iso"); err == nil {
		for _, file := range files {
			stdFile := IsoFile{
				Name: path.Base(file),
				Path: storage.PATH.Fmt(file),
			}
			images = append(images, stdFile)
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
	Storage []DataStore `json:"storage"`
}

func (store DataStoreMgr) List() []DataStore {
	stores := make([]DataStore, 0, 32)

	if dirs, err := libstar.DIR.ListDirs(storage.Location + "datastore"); err == nil {
		for _, dir := range dirs {
			path := storage.PATH.Fmt(dir)
			stores = append(stores, DataStore{Name: path, Path: path})
		}
	}
	return stores
}

var DATASTOR = DataStoreMgr{
	Storage: make([]DataStore, 0, 32),
}
