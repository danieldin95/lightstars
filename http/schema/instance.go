package schema

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/storage"
)

type Instance struct {
	Action           string `json:"action,omitempty"` // If is "", means not action.
	UUID             string `json:"uuid"`
	Name             string `json:"name"`
	Family           string `json:"family"'` // configure
	State            string `json:"state"`
	Arch             string `json:"arch"` // configure: x86_64 or i386
	Type             string `json:"type"`
	Boots            string `json:"boots,omitempty"`     // configure: hd,cdrom,network.
	DataStore        string `json:"datastore,omitempty"` // configure
	Cpu              string `json:"cpu,omitempty"`       // configure
	MemSize          string `json:"memSize,omitempty"`   // configure
	MemUnit          string `json:"memUnit,omitempty"`   // configure
	CpuMode          string `json:"cpuMode,omitempty"`   // configure
	Disk0File        string `json:"disk0File,omitempty"` // configure
	Disk1Size        string `json:"disk1Size,omitempty"` // configure
	Disk1Unit        string `json:"disk1Unit,omitempty"` // configure
	Disk1Bus         string `json:"disk1Bus,omitempty"`  // configure
	Interface0Source string `json:"interface0Source,omitempty"`
	Interface0Bus    string `json:"interface0Bus,omitempty"`
	Start            string `json:"start,omitempty"` //configure

	MaxCpu      uint         `json:"maxCpu"`
	MaxMem      uint64       `json:"maxMem"`  // Kbytes
	Memory      uint64       `json:"memory"`  // KBytes
	CpuTime     uint64       `json:"cpuTime"` // MicroSeconds
	Disks       []Disk       `json:"disks,omitempty"`
	Interfaces  []Interface  `json:"interfaces,omitempty"`
	Controllers []Controller `json:"controllers,omitempty"`
	Password    string       `json:"password"`

	XMLObj *libvirtc.DomainXML `json:"-"`
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
		obj.State = libvirtc.DomainState2Str(info.State)
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
	obj.XMLObj = xmlObj
	return obj
}

type Disk struct {
	Action     string `json:"action,omitempty"`
	Seq        string `json:"seq,omitempty"`       // configure
	Name       string `json:"name,omitempty"`      // disk name
	UUID       string `json:"uuid,omitempty"`      // disk UUID
	Store      string `json:"datastore,omitempty"` // disk saved to datastore
	Size       string `json:"size"`                // configure
	SizeUnit   string `json:"sizeUnit,omitempty"`  //configure
	Format     string `json:"format"`
	Source     string `json:"source"`
	Device     string `json:"device"`
	Bus        string `json:"bus"`      //configre
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
	if xml.Source.File != "" {
		disk.Source = storage.PATH.Fmt(xml.Source.File)
	} else if xml.Source.Device != "" {
		disk.Source = storage.PATH.Fmt(xml.Source.Device)
	}
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
	Action     string `json:"action,omitempty"` // If is "", means not action.
	Seq        string `json:"seq,omitempty"`    //configure
	Name       string `json:"name,omitempty"`
	UUID       string `json:"uuid,omitempty"`
	Type       string `json:"type,omitempty"` //bridge or openvswitch
	Address    string `json:"address"`
	Source     string `json:"source"`
	Model      string `json:"model"` // configure
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
