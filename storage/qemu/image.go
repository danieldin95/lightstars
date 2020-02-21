package qemuimgdriver

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/quadrifoglio/go-qemu"
	"path"
	"strings"
)

const GiB = uint64(1024 * 1024 * 1024)
const Location = "/lightstar/"

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

func (iso IsoMgr) ListFiles(store string) []IsoFile {
	images := make([]IsoFile, 0, 32)

	if files, err := libstar.DIR.ListFiles(Location+store, ".iso"); err == nil {
		for _, file := range files {
			images = append(images, IsoFile{Name: path.Base(file), Path: file})
		}
	}
	if files, err := libstar.DIR.ListFiles(Location+store+"/iso", ".iso"); err == nil {
		for _, file := range files {
			images = append(images, IsoFile{Name: path.Base(file), Path: file})
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

	if dirs, err := libstar.DIR.ListDirs(Location + "datastore"); err == nil {
		for _, dir := range dirs {
			path := strings.Replace(dir, Location, "", 1)
			name := strings.Replace(path, "/", ".", 4)
			stores = append(stores, DataStore{Name: name, Path: path})
		}
	}
	return stores
}

var DATASTOR = DataStoreMgr{
	Storage: make([]DataStore, 0, 32),
}
