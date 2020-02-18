package http

import "C"
import (
	"github.com/danieldin95/lightstar/compute/libvirt"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/libvirt/libvirt-go"
)

type VersionSchema struct {
	Version string `json:"version"`
	Date    string `json:"date"`
	Commit  string `json:"commit"`
}

func NewVersionSchema() VersionSchema {
	return VersionSchema{
		Version: libstar.Version,
		Date:    libstar.Date,
		Commit:  libstar.Commit,
	}
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

type InstanceSchema struct {
	UUID   string  `json:"uuid"`
	Name   string  `json:"name"`
	State  string  `json:"state"`
	MaxCpu uint    `json:"maxCpu"`
	MaxMem uint64  `json:"maxMem"` // Kbytes
	Memory uint64  `json:"memory"` // KBytes
	CpuTime uint64 `json:"cpuTime"` // MicroSeconds
}

func NewInstanceSchema(dom libvirtdriver.Domain) InstanceSchema {
	object := InstanceSchema{
	}
	object.UUID, _ = dom.GetUUIDString()
	object.Name, _ = dom.GetName()
	if info, err := dom.GetInfo(); err == nil {
		object.State = InstanceState2Str(info.State)
		object.MaxMem = info.MaxMem
		object.Memory = info.Memory
		object.MaxCpu = info.NrVirtCpu
		object.CpuTime = info.CpuTime / 1000000
	}
	return object
}

type IndexSchema struct {
	Version   VersionSchema    `json:"version"`
	Instances []InstanceSchema `json:"instances"`
}
