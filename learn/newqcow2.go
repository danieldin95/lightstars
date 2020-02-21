package main

import (
	"github.com/quadrifoglio/go-qemu"
	"log"
)

func main() {
	GiB := uint64(1024 * 1024 * 1024)
	img := qemu.NewImage("vm.qcow2", qemu.ImageFormatQCOW2, 1024*GiB)
	img.SetBackingFile("debian.qcow2")

	err := img.Create()
	if err != nil {
		log.Fatal(err)
	}
}
