package compute

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/schema"
	"github.com/danieldin95/lightstar/storage"
)

func NewHyper() (hs schema.Hyper) {
	hyper, _ := libvirtc.GetHyper()

	hs.Name = hyper.Url
	hs.Host = hyper.Host
	hs.CpuNum, hs.CpuVendor, hs.CpuUtils = hyper.GetCPU()
	hs.MemTotal, hs.MemFree, hs.MemCached = hyper.GetMem()
	return hs
}

func NewFromInterfaceXML(xml libvirtc.InterfaceXML) (int schema.Interface) {
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

func NewFromDiskXML(xml libvirtc.DiskXML) (disk schema.Disk) {
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

func NewFromControllerXML(xml libvirtc.ControllerXML) (ctl schema.Controller) {
	ctl.Type = xml.Type
	ctl.Model = xml.Model
	ctl.Index = xml.Index
	if xml.Address != nil {
		ctl.Address = NewFromAddressXML(*xml.Address)
	}
	return ctl
}

func NewInstance(dom libvirtc.Domain) schema.Instance {
	obj := schema.Instance{
		Disks:       make([]schema.Disk, 0, 32),
		Interfaces:  make([]schema.Interface, 0, 32),
		Controllers: make([]schema.Controller, 0, 32),
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
			g := schema.Graphics{
				Type:     x.Type,
				Listen:   x.Listen,
				Password: x.Password,
				Port:     x.Port,
			}
			obj.Graphics = append(obj.Graphics, g)
		}
	}
	return obj
}

func NewFromAddressXML(xml libvirtc.AddressXML) (addr schema.Address) {
	addr.Type = xml.Type
	addr.Domain = xml.Domain
	addr.Bus = xml.Bus
	addr.Slot = xml.Slot
	addr.Function = xml.Function
	return addr
}
