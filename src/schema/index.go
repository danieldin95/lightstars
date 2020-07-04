package schema

import "github.com/danieldin95/lightstar/src/libstar"

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
	Host       string  `json:"host"`
	CpuNum     uint    `json:"cpuNum"`
	CpuVendor  string  `json:"cpuVendor"`
	CpuUtils   uint64  `json:"cpuUtils"`
	MemTotal   uint64  `json:"memTotal"`
	MemFree    uint64  `json:"memFree"`
	MemCached  uint64  `json:"memCached"`
	MemPercent float64 `json:"memPercent"`
}

type Index struct {
	Version    Version     `json:"version"`
	User       User        `json:"user"`
	Hyper      Hyper       `json:"hyper"`
	Instances  []Instance  `json:"instances"`
	DataStores []DataStore `json:"datastores"`
	Networks   []Network   `json:"networks"`
	Default    string      `json:"default"`
}
