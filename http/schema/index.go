package schema

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/libstar"
)

type Version struct {
	Version string `json:"version"`
	Date    string `json:"date"`
	Commit  string `json:"commit"`
}

func NewVersion() Version {
	return Version{
		Version: libstar.Version,
		Date:    libstar.Date,
		Commit:  libstar.Commit,
	}
}

type Hyper struct {
	Name       string  `json:"name"`
	CpuNum     uint    `json:"cpuNum"`
	CpuVendor  string  `json:"cpuVendor"`
	CpuUtils   uint64  `json:"cpuUtils"`
	MemTotal   uint64  `json:"memTotal"`
	MemFree    uint64  `json:"memFree"`
	MemCached  uint64  `json:"memCached"`
	MemPercent float64 `json:"memPercent"`
}

func NewHyper() (hs Hyper) {
	hyper, _ := libvirtc.GetHyper()

	hs.Name = hyper.Name
	hs.CpuNum, hs.CpuVendor, hs.CpuUtils = hyper.GetCPU()
	hs.MemTotal, hs.MemFree, hs.MemCached = hyper.GetMem()

	return hs
}

type Index struct {
	Version    Version     `json:"version"`
	Hyper      Hyper       `json:"hyper"`
	Instances  []Instance  `json:"instances"`
	DataStores []DataStore `json:"datastores"`
	Networks   []Network   `json:"networks"`
}
