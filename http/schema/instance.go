package schema

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/libstar"
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
	Password    string       `json:"password"`
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
		for _, x := range xmlObj.Devices.Graphics {
			if x.Type == "vnc" {
				obj.Password = x.Password
			}
		}
	}
	return obj
}

type Disk struct {
	Format     string `json:"format"`
	Source     string `json:"source"`
	Device     string `json:"device"`
	Bus        string `json:"bus"`
	AddrType   string `json:"addrType"` // pci, and drive
	AddrSlot   uint16 `json:"addrSlot"`
	AddrDomain uint16 `json:"addrDomain"`
	AddrBus    uint16 `json:"addrBus"`
	AddrFunc   uint16 `json:"addrFunc"`
	AddrCtl    uint16 `json"addrCtl"`
	AddrTgt    uint16 `json:"addrTgt"`
	AddrUnit   uint16 `json:"addrUnit"`
}

func NewFromDiskXML(xml libvirtc.DiskXML) (disk Disk) {
	disk.Device = xml.Target.Dev
	disk.Bus = xml.Target.Bus
	disk.Source = storage.PATH.Fmt(xml.Source.File)
	disk.Format = xml.Driver.Type
	if xml.Address != nil {
		disk.AddrType = xml.Address.Type
		if disk.AddrType == "pci" {
			disk.AddrDomain = libstar.H2D16(xml.Address.Domain)
			disk.AddrBus = libstar.H2D16(xml.Address.Bus)
			disk.AddrSlot = libstar.H2D16(xml.Address.Slot)
			disk.AddrFunc = libstar.H2D16(xml.Address.Function)
		} else if xml.Address.Type == "drive" {
			disk.AddrCtl = libstar.H2D16(xml.Address.Controller)
			disk.AddrBus = libstar.H2D16(xml.Address.Bus)
			disk.AddrTgt = libstar.H2D16(xml.Address.Target)
			disk.AddrUnit = libstar.H2D16(xml.Address.Unit)
		}
	}
	return disk
}

type Interface struct {
	Address    string `json:"address"`
	Source     string `json:"source"`
	Model      string `json:"model"`
	Device     string `json:"device"`
	AddrType   string `json:"addrType"` // now only pci.
	AddrSlot   uint16 `json:"addrSlot"`
	AddrDomain uint16 `json:"addrDomain"`
	AddrBus    uint16 `json:"addrBus"`
	AddrFunc   uint16 `json:"addrFunc"`
}

func NewFromInterfaceXML(xml libvirtc.InterfaceXML) (int Interface) {
	int.Source = xml.Source.Bridge
	int.Address = xml.Mac.Address
	int.Model = xml.Model.Type
	int.Device = xml.Target.Dev
	if xml.Address != nil {
		int.AddrType = xml.Address.Type
		if int.AddrType == "pci" {
			int.AddrDomain = libstar.H2D16(xml.Address.Domain)
			int.AddrBus = libstar.H2D16(xml.Address.Bus)
			int.AddrSlot = libstar.H2D16(xml.Address.Slot)
			int.AddrFunc = libstar.H2D16(xml.Address.Function)
		}
	}
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
	if xml.Address != nil {
		ctl.Address = NewFromAddressXML(*xml.Address)
	}
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

type InstanceConf struct {
	Action     string `json:"action"` // If is "", means not action.
	Name       string `json:"name"`
	Family     string `json:"family"'`
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
	Start      string `json:"start"`
}

type DiskConf struct {
	Action string `json:"action"` // If is "", means not action.
	Name   string `json:"name"`
	UUID   string `json:"uuid"`
	Store  string `json:"datastore"`
	Size   string `json:"size"`
	Unit   string `json:"unit"`
	Bus    string `json:"bus"`
	Slot   string `json:"slot"`
}

type InterfaceConf struct {
	Action    string `json:"action"` // If is "", means not action.
	Name      string `json:"name"`
	UUID      string `json:"uuid"`
	Interface string `json:"interface"`
	Type      string `json:"type"`
	Bus       string `json:"bus"`
	Model     string `json:"model"`
	Slot      string `json:"slot"`
}
