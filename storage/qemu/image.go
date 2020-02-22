package qemuimgdriver

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/storage"
	"github.com/quadrifoglio/go-qemu"
	"path"
)

const GiB = uint64(1024 * 1024 * 1024)

type Image struct {
	Path     string
	Size     uint64
	BackFile string
	Format   string // qemu.ImageFormatQCOW2
}

func NewImage(path string, size uint64) *Image {
	return &Image{
		Path:   path,
		Size:   size,
		Format: qemu.ImageFormatQCOW2,
	}
}

func (i *Image) Create() error {
	img := qemu.NewImage(i.Path, i.Format, i.Size)
	if i.BackFile != "" {
		img.SetBackingFile(i.BackFile)
	}

	return img.Create()
}

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
