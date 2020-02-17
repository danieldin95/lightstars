package http

import "C"
import (
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
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	State  string `json:"state"`
	MaxCpu uint   `json:"maxCpu"`
	MaxMem uint64 `json:"maxMem"`
}

func NewInstanceSchema(dom libvirt.Domain) InstanceSchema {
	uuid, _ := dom.GetUUIDString()
	name, _ := dom.GetName()
	state, _, _ := dom.GetState()
	cpu, _ := dom.GetMaxVcpus()
	mem, _ := dom.GetMaxMemory()

	return InstanceSchema{
		UUID:   uuid,
		Name:   name,
		State:  InstanceState2Str(state),
		MaxCpu: cpu,
		MaxMem: mem,
	}
}

type IndexSchema struct {
	Version   VersionSchema    `json:"version"`
	Instances []InstanceSchema `json:"instances"`
}
