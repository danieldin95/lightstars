package schema

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/storage"
)

type Instance struct {
	UUID        string       `json:"uuid"`
	Name        string       `json:"name"`
	State       string       `json:"state"`
	Arch        string       `json:"arch"`
	Type        string       `json:"type"`
	MaxCpu      uint         `json:"maxCpu"`
	MaxMem      uint64       `json:"maxMem"`  // Kbytes
	Memory      uint64       `json:"memory"`  // KBytes
	CpuTime     uint64       `json:"cpuTime"` // MicroSeconds
	Disks       []Disk       `json:"disks,omitempty"`
	Interfaces  []Interface  `json:"interfaces,omitempty"`
	Controllers []Controller `json:"controllers,omitempty"`
}

func NewInstance(dom libvirtc.Domain) Instance {
	obj := Instance{
		Disks:       make([]Disk, 0, 32),
		Interfaces:  make([]Interface, 0, 32),
		Controllers: make([]Controller, 0, 32),
	}
	obj.UUID, _ = dom.GetUUIDString()
	obj.Name, _ = dom.GetName()
	if info, err := dom.GetInfo(); err == nil {
		obj.State = InstanceState2Str(info.State)
		obj.MaxMem = info.MaxMem
		obj.Memory = info.Memory
		obj.MaxCpu = info.NrVirtCpu
		obj.CpuTime = info.CpuTime / 1000000
	}
	xmlObj := libvirtc.NewDomainXMLFromDom(&dom, true)
	if xmlObj != nil {
		obj.Arch = xmlObj.OS.Type.Arch
		obj.Type = xmlObj.Type
		for _, x := range xmlObj.Devices.Disks {
			obj.Disks = append(obj.Disks, NewFromDiskXML(x))
		}
		for _, x := range xmlObj.Devices.Interfaces {
			obj.Interfaces = append(obj.Interfaces, NewFromInterfaceXML(x))
		}
		for _, x := range xmlObj.Devices.Controllers {
			obj.Controllers = append(obj.Controllers, NewFromControllerXML(x))
		}
	}
	return obj
}

type Disk struct {
	Format string `json:"format"`
	Source string `json:"source"`
	Device string `json:"device"`
	Bus    string `json:"bus"`
}

func NewFromDiskXML(xml libvirtc.DiskXML) (disk Disk) {
	disk.Device = xml.Target.Dev
	disk.Bus = xml.Target.Bus
	disk.Source = storage.PATH.Fmt(xml.Source.File)
	disk.Format = xml.Driver.Type
	return disk
}

type Interface struct {
	Address string `json:"address"`
	Source  string `json:"source"`
	Model   string `json:"model"`
	Device  string `json:"device"`
}

func NewFromInterfaceXML(xml libvirtc.InterfaceXML) (int Interface) {
	int.Source = xml.Source.Bridge
	int.Address = xml.Mac.Address
	int.Model = xml.Model.Type
	int.Device = xml.Target.Dev
	return int
}

type Controller struct {
	Type    string  `json:"source"`
	Model   string  `json:"model"`
	Index   string  `json:"device"`
	Address Address `json:"address"`
}

func NewFromControllerXML(xml libvirtc.ControllerXML) (ctl Controller) {
	ctl.Type = xml.Type
	ctl.Model = xml.Model
	ctl.Index = xml.Index
	ctl.Address = NewFromAddressXML(xml.Address)
	return ctl
}

type Address struct {
	Type     string `json:"type"`
	Domain   string `json:"domain"`
	Bus      string `json:"bus"`
	Slot     string `json:"slot"`
	Function string `json:"function"`
}

func NewFromAddressXML(xml libvirtc.AddressXML) (addr Address) {
	addr.Type = xml.Type
	addr.Domain = xml.Domain
	addr.Bus = xml.Bus
	addr.Slot = xml.Slot
	addr.Function = xml.Function
	return addr
}
