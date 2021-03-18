package compute

import (
	"fmt"
	"github.com/danieldin95/lightstar/src/compute/libvirtc"
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/danieldin95/lightstar/src/schema"
)

func NewHyper() (hs schema.Hyper) {
	hyper, _ := libvirtc.GetHyper()
	hs.Name = hyper.Url
	hs.Host = hyper.Host
	hs.CpuNum, hs.CpuVendor, hs.CpuUtils = hyper.GetCPU()
	hs.MemTotal, hs.MemFree, hs.MemCached = hyper.GetMem()
	hs.UpTime = hyper.UpTime()
	return hs
}

func NewFromInterfaceXML(xml libvirtc.InterfaceXML, domain schema.Instance) (int schema.Interface) {
	int.Domain = schema.Instance{
		UUID: domain.UUID,
		Name: domain.Name,
	}
	int.Source = xml.Source.Bridge
	int.Network = xml.Source.Network
	addr := xml.Source.Address
	if addr != nil && addr.Type == "pci" {
		int.HostDev = addr.Type
		int.HostDev += fmt.Sprintf(":%02x", libstar.H2D16(addr.Bus))
		int.HostDev += fmt.Sprintf(":%02x", libstar.H2D16(addr.Slot))
		int.HostDev += fmt.Sprintf(".%x", libstar.H2D16(addr.Function))
	}
	int.Address = xml.Mac.Address
	int.Model = xml.Model.Type
	if int.Model == "" && xml.Driver != nil {
		int.Model = xml.Driver.Name
	}
	int.Device = xml.Target.Dev
	addr = xml.Address
	if addr != nil {
		int.AddrType = addr.Type
		if int.AddrType == "pci" {
			int.AddrDomain = fmt.Sprintf("%04x", libstar.H2D16(addr.Domain))
			int.AddrBus = fmt.Sprintf("%02x", libstar.H2D16(addr.Bus))
			int.AddrSlot = fmt.Sprintf("%02x", libstar.H2D16(addr.Slot))
			int.AddrFunc = fmt.Sprintf("%x", libstar.H2D16(addr.Function))
		}
	}
	return int
}

func NewFromDiskXML(xml libvirtc.DiskXML, domain schema.Instance) (disk schema.Disk) {
	disk.Domain = schema.Instance{
		UUID: domain.UUID,
		Name: domain.Name,
	}
	disk.Device = xml.Target.Dev
	disk.Bus = xml.Target.Bus
	if xml.Source.File != "" {
		disk.Source = xml.Source.File
		disk.Name = xml.Source.File
	} else if xml.Source.Device != "" {
		disk.Source = xml.Source.Device
		disk.Name = xml.Source.Device
	}
	disk.Format = xml.Driver.Type
	addr := xml.Address
	if addr != nil {
		disk.AddrType = addr.Type
		switch disk.AddrType {
		case "pci":
			disk.AddrDomain = fmt.Sprintf("%04x", libstar.H2D16(addr.Domain))
			disk.AddrBus = fmt.Sprintf("%02x", libstar.H2D16(addr.Bus))
			disk.AddrSlot = fmt.Sprintf("%02x", libstar.H2D16(addr.Slot))
			disk.AddrFunc = fmt.Sprintf("%x", libstar.H2D16(addr.Function))
		case "drive":
			disk.AddrCtl = fmt.Sprintf("%04x", libstar.H2D16(addr.Controller))
			disk.AddrBus = fmt.Sprintf("%02x", libstar.H2D16(xml.Address.Bus))
			disk.AddrTgt = fmt.Sprintf("%02x", libstar.H2D16(xml.Address.Target))
			disk.AddrUnit = fmt.Sprintf("%x", libstar.H2D16(xml.Address.Unit))
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
	obj := schema.Instance{}
	obj.UUID, _ = dom.GetUUIDString()
	obj.Name, _ = dom.GetName()
	if info, err := dom.GetInfo(); err == nil {
		obj.State = libvirtc.DomainState2Str(info.State)
		obj.MaxMem = info.MaxMem
		obj.Memory = info.Memory
		obj.MaxCpu = info.NrVirtCpu
		obj.CpuTime = info.CpuTime / 1000000
	}
	obj.Title, _ = dom.GetMetadataTitle(true)
	obj.Description, _ = dom.GetMetadataTitle(false)
	xmlObj := libvirtc.NewDomainXMLFromDom(&dom, true)
	if xmlObj == nil {
		return obj
	}
	obj.Arch = xmlObj.OS.Type.Arch
	obj.CpuMode = xmlObj.CPU.Mode
	obj.Type = xmlObj.Type
	for _, x := range xmlObj.Devices.Disks {
		obj.Disks = append(obj.Disks, NewFromDiskXML(x, obj))
	}
	for _, x := range xmlObj.Devices.Interfaces {
		obj.Interfaces = append(obj.Interfaces, NewFromInterfaceXML(x, obj))
	}
	for _, x := range xmlObj.Devices.Controllers {
		obj.Controllers = append(obj.Controllers, NewFromControllerXML(x))
	}
	for _, x := range xmlObj.Devices.Graphics {
		obj.Graphics = append(obj.Graphics, schema.Graphics{
			Type:     x.Type,
			Listen:   x.Listen,
			Password: x.Password,
			Port:     x.Port,
		})
	}
	for _, x := range xmlObj.Devices.Channels {
		obj.Channels = append(obj.Channels, schema.Channel{
			Type:          x.Type,
			TargetName:    x.Target.Name,
			TargetType:    x.Target.Type,
			SourceChannel: x.Source.Channel,
		})
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
