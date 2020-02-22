package schema

import (
	"github.com/danieldin95/lightstar/compute/libvir"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/libvirt/libvirt-go"
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
	CpuNum     int     `json:"cpuNum"`
	CpuVendor  string  `json:"cpuVendor"`
	MemTotal   uint64  `json:"memTotal"`
	MemFree    uint64  `json:"memFree"`
	MemPercent float64 `json:"memPercent"`
}

func NewHyper() (hs Hyper) {
	hyper, _ := libvir.GetHyper()

	hs.CpuNum, hs.CpuVendor = hyper.GetCPU()
	hs.MemTotal, hs.MemFree, hs.MemPercent = hyper.GetMem()

	return hs
}

func InstanceState2Str(state libvirt.DomainState) string {
	switch state {
	case libvirt.DOMAIN_NOSTATE:
		return "nostate"
	case libvirt.DOMAIN_RUNNING:
		return "running"
	case libvirt.DOMAIN_BLOCKED:
		return "blocked"
	case libvirt.DOMAIN_PAUSED:
		return "paused"
	case libvirt.DOMAIN_SHUTDOWN:
		return "shutdown"
	case libvirt.DOMAIN_CRASHED:
		return "crashed"
	case libvirt.DOMAIN_PMSUSPENDED:
		return "pmsuspended"
	case libvirt.DOMAIN_SHUTOFF:
		return "shutoff"
	default:
		return "unknown"
	}
}

type Index struct {
	Version   Version    `json:"version"`
	Hyper     Hyper      `json:"hyper"`
	Instances []Instance `json:"instances"`
}

type InstanceConf struct {
	Action     string `json:"action"` // If is "", means not action.
	Name       string `json:"name"`
	Arch       string `json:"arch"`
	Boots      string `json:"boots"`
	DataStore  string `json:"datastore"`
	Cpu        string `json:"cpu"`
	MemorySize string `json:"memorySize"`
	MemoryUnit string `json:"memoryUnit"`
	DiskSize   string `json:"diskSize"`
	DiskUnit   string `json:"diskUnit"`
	IsoFile    string `json:"isoFile"`
	Interface  string `json:"interface"`
}
