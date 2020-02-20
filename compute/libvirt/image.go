package libvirtdriver

import (
	"github.com/quadrifoglio/go-qemu"
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
