package api

import (
	"fmt"
	"github.com/danieldin95/lightstar/src/compute"
	"github.com/danieldin95/lightstar/src/compute/libvirtc"
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/danieldin95/lightstar/src/network/libvirtn"
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/danieldin95/lightstar/src/storage"
	"github.com/danieldin95/lightstar/src/storage/libvirts"
	"github.com/gorilla/mux"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"
)

type Instance struct {
}

func GetTypeByVolume(file string) (string, string) {
	if file == "" {
		return "", ""
	}
	vol := libvirts.Volume{
		Pool: path.Dir(file),
		Name: path.Base(file),
	}
	desc, err := vol.GetXMLObj()
	if err != nil {
		libstar.Warn("GetTypeByVolume: %s %s", vol.Pool, err)
		vol.Pool = path.Base(vol.Pool)
		desc, err = vol.GetXMLObj()
		if err != nil {
			libstar.Warn("GetTypeByVolume: %s %s", vol.Pool, err)
			return "disk", ""
		}
	}
	libstar.Debug("GetTypeByVolume: %s:%v", file, desc.Target.Format)
	format := desc.Target.Format.Type
	if format == "iso" {
		return "cdrom", "raw"
	}
	return "disk", format
}

type diskSeq struct {
	HdId uint8
	VdId uint8
}

func NewCdXML(file, family string, seq *diskSeq) libvirtc.DiskXML {
	seq.HdId++
	return libvirtc.DiskXML{
		Type:   "block",
		Device: "cdrom",
		Driver: libvirtc.DiskDriverXML{
			Name: "qemu",
			Type: "raw",
		},
		Source: libvirtc.DiskSourceXML{
			Device: file,
		},
		Target: libvirtc.DiskTargetXML{
			Bus: "ide",
			Dev: libvirtc.DISK.Slot2Dev("ide", seq.HdId),
		},
	}
}

func NewIsoXML(file, family string, seq *diskSeq) libvirtc.DiskXML {
	xml := libvirtc.DiskXML{
		Type:   "file",
		Device: "disk",
		Driver: libvirtc.DiskDriverXML{
			Type: "raw",
			Name: "qemu",
		},
		Source: libvirtc.DiskSourceXML{
			File: file,
		},
	}
	name := strings.ToUpper(file)
	device, format := GetTypeByVolume(file)
	xml.Device = device
	xml.Driver.Type = format
	if family == "linux" && !strings.HasSuffix(name, ".ISO") {
		seq.VdId++
		xml.Target = libvirtc.DiskTargetXML{
			Bus: "virtio",
			Dev: libvirtc.DISK.Slot2Dev("virtio", seq.VdId),
		}
	} else {
		seq.HdId++
		xml.Target = libvirtc.DiskTargetXML{
			Bus: "ide",
			Dev: libvirtc.DISK.Slot2Dev("ide", seq.HdId),
		}
	}
	return xml
}

func NewDiskXML(format, file, bus string, seq *diskSeq) libvirtc.DiskXML {
	disk := libvirtc.DiskXML{
		Type:   "file",
		Device: "disk",
		Driver: libvirtc.DiskDriverXML{
			Name: "qemu",
			Type: format,
		},
		Source: libvirtc.DiskSourceXML{
			File: file,
		},
	}
	switch bus {
	case "virtio":
		seq.VdId++
		disk.Target = libvirtc.DiskTargetXML{
			Bus: bus,
			Dev: libvirtc.DISK.Slot2Dev(bus, seq.VdId),
		}
		disk.Address = &libvirtc.AddressXML{
			Type:     "pci",
			Domain:   libvirtc.PciDomain,
			Bus:      libvirtc.PciDiskBus,
			Slot:     fmt.Sprintf("0x%x", seq.VdId),
			Function: libvirtc.PciFunc,
		}
	case "ide", "scsi":
		seq.HdId++
		disk.Target = libvirtc.DiskTargetXML{
			Bus: bus,
			Dev: libvirtc.DISK.Slot2Dev(bus, seq.HdId),
		}
	}
	return disk
}

func NewFileXML(disk *schema.Disk, conf *schema.Instance, seq *diskSeq) (libvirtc.DiskXML, error) {
	obj := libvirtc.DiskXML{}
	file := storage.PATH.Unix(disk.Source)
	name := libvirtc.DISK.Slot2Name(seq.VdId)
	size := libstar.ToBytes(disk.Size, disk.SizeUnit)
	device, format := GetTypeByVolume(file)
	if file == "" {
		vol, err := NewVolumeAndPool(conf.DataStore, conf.Name, name, size)
		if err != nil {
			return obj, err
		}
		file = vol.Target.Path
		format = vol.Target.Format.Type
	} else if device == "disk" && (format == "raw" || format == "qcow2" || format == "qcow") {
		name := path.Base(file)
		vol, err := NewBackingVolumeAndPool(conf.DataStore, conf.Name, name, file, format)
		if err != nil {
			return obj, err
		}
		file = vol.Target.Path
		format = vol.Target.Format.Type
	}
	switch conf.Family {
	case "linux":
		obj = NewDiskXML(format, file, "virtio", seq)
	case "windows": // not scsi.
		obj = NewDiskXML(format, file, "ide", seq)
	default:
		obj = NewDiskXML(format, file, "ide", seq)
	}
	return obj, nil
}

func Instance2XML(conf *schema.Instance) (libvirtc.DomainXML, error) {
	dom := libvirtc.DomainXML{
		Type: "kvm",
		Name: conf.Name,
		Devices: libvirtc.DevicesXML{
			Disks:       make([]libvirtc.DiskXML, 0, 2),
			Graphics:    make([]libvirtc.GraphicsXML, 1), // vnc/spice
			Interfaces:  make([]libvirtc.InterfaceXML, 0, 1),
			Controllers: make([]libvirtc.ControllerXML, 4),
			Inputs:      make([]libvirtc.InputXML, 1), // <input type="tablet" bus="usb"/>
		},
		OS: libvirtc.OSXML{
			Type: libvirtc.OSTypeXML{
				Arch:  conf.Arch,
				Value: "hvm",
			},
			Boot: make([]libvirtc.OSBootXML, 3),
			BootMenu: libvirtc.OSBootMenuXML{
				Enable: "yes",
			},
		},
	}
	if dom.OS.Type.Arch == "" {
		dom.OS.Type.Arch = "x86_64" // i386
	}
	// boot sequence.
	if conf.Boots == "" {
		conf.Boots = "hd,cdrom,network"
	}
	for i, v := range strings.Split(conf.Boots, ",") {
		if i < 3 {
			dom.OS.Boot[i] = libvirtc.OSBootXML{
				Dev: v,
			}
		}
	}
	// features
	dom.Features = libvirtc.FeaturesXML{
		Apic: &libvirtc.ApicXML{},
		Acpi: &libvirtc.AcpiXML{},
		Pae:  &libvirtc.PaeXML{},
	}
	// cpu and memory
	if conf.CpuMode != "" {
		dom.CPUXml = libvirtc.CPUXML{
			Mode:  conf.CpuMode,
			Check: "full",
		}
	}
	dom.VCPUXml = libvirtc.VCPUXML{
		Placement: "static",
		Value:     fmt.Sprintf("%d", conf.MaxCpu),
	}
	dom.Memory = libvirtc.MemXML{
		Value: fmt.Sprintf("%d", conf.MaxMem),
		Type:  "KiB",
	}
	dom.CurMem = libvirtc.CurMemXML{
		Value: fmt.Sprintf("%d", conf.MaxMem),
		Type:  "KiB",
	}
	// vnc
	dom.Devices.Graphics[0] = libvirtc.GraphicsXML{
		Type:     "vnc",
		Listen:   "0.0.0.0",
		Port:     "-1",
		AutoPort: "yes",
		Password: libstar.GenToken(32),
	}
	// controllers
	dom.Devices.Controllers[0] = libvirtc.ControllerXML{
		Type:  "pci",
		Index: "0",
		Model: "pci-root",
	}
	for i := 1; i < len(dom.Devices.Controllers); i++ {
		dom.Devices.Controllers[i] = libvirtc.ControllerXML{
			Type:  "pci",
			Index: strconv.Itoa(i),
			Model: "pci-bridge",
		}
	}
	// disks
	seq := &diskSeq{}
	for _, disk := range conf.Disks {
		file := disk.Source
		obj := libvirtc.DiskXML{}
		if strings.HasPrefix(file, "/dev") {
			obj = NewCdXML(file, conf.Family, seq)
		} else if strings.HasSuffix(file, ".iso") || strings.HasSuffix(file, ".ISO") {
			obj = NewIsoXML(storage.PATH.Unix(file), conf.Family, seq)
		} else {
			var err error
			if obj, err = NewFileXML(&disk, conf, seq); err != nil {
				return dom, err
			}
		}
		dom.Devices.Disks = append(dom.Devices.Disks, obj)
	}

	// interfaces
	for i, inf := range conf.Interfaces {
		seq := fmt.Sprintf("0x%x", i+1)
		source := inf.Source
		br, _ := libvirtn.BRIDGE.Get(source)
		obj := Interface2XML(source, "virtio", seq, br.Type)
		switch conf.Family {
		case "linux":
			obj.Model = libvirtc.InterfaceModelXML{
				Type: "virtio",
			}
		case "windows":
			obj.Model = libvirtc.InterfaceModelXML{
				Type: "rtl8139", //e1000,rtl8139
			}
		default:
			obj.Model = libvirtc.InterfaceModelXML{
				Type: "rtl8139",
			}
		}
		dom.Devices.Interfaces = append(dom.Devices.Interfaces, obj)
	}
	// inputs
	dom.Devices.Inputs[0] = libvirtc.InputXML{
		Type: "tablet",
		Bus:  "usb",
	}
	// sound
	dom.Devices.Sound = libvirtc.SoundXML{
		Model: "ich6",
	}
	return dom, nil
}

func (ins Instance) Router(router *mux.Router) {
	router.HandleFunc("/api/instance", ins.GET).Methods("GET")
	router.HandleFunc("/api/instance", ins.POST).Methods("POST")
	router.HandleFunc("/api/instance/{id}", ins.GET).Methods("GET")
	router.HandleFunc("/api/instance/{id}", ins.PUT).Methods("PUT")
	router.HandleFunc("/api/instance/{id}", ins.DELETE).Methods("DELETE")
}

func (ins Instance) HasPermission(user *schema.User, instance string) bool {
	has := false
	if user.Type == "admin" {
		return true
	}
	if strings.HasPrefix(instance, user.Name+".") {
		has = true
	} else {
		for _, name := range user.Instances {
			if instance == name {
				has = true
				break
			}
		}
	}
	return has
}

func (ins Instance) GetByUser(user *schema.User, list *schema.ListInstance) {
	if domains, err := libvirtc.ListDomains(); err == nil {
		for _, d := range domains {
			inst := compute.NewInstance(d)
			if ins.HasPermission(user, inst.Name) {
				list.Items = append(list.Items, inst)
			}
			_ = d.Free()
		}
	}
}

func (ins Instance) GET(w http.ResponseWriter, r *http.Request) {
	uuid, ok := GetArg(r, "id")
	if !ok {
		user, _ := GetUser(r)
		list := schema.ListInstance{
			Items: make([]schema.Instance, 0, 32),
		}
		// list all instances.
		ins.GetByUser(&user, &list)
		sort.SliceStable(list.Items, func(i, j int) bool {
			return list.Items[i].Name < list.Items[j].Name
		})
		list.Metadata.Size = len(list.Items)
		list.Metadata.Total = len(list.Items)
		ResponseJson(w, list)
		return
	}

	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer dom.Free()
	format := GetQueryOne(r, "format")
	if format == "xml" {
		xmlDesc, err := dom.GetXMLDesc(true)
		if err == nil {
			ResponseXML(w, xmlDesc)
		} else {
			ResponseXML(w, "<error>"+err.Error()+"</error>")
		}
	} else {
		ResponseJson(w, compute.NewInstance(*dom))
	}
}

func (ins Instance) POST(w http.ResponseWriter, r *http.Request) {
	hyper, err := libvirtc.GetHyper()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	conf := &schema.Instance{}
	if err := GetData(r, conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if dom, _ := hyper.LookupDomainByName(conf.Name); dom != nil {
		http.Error(w, conf.Name+" already existed", http.StatusConflict)
		return
	}
	xmlObj, err := Instance2XML(conf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	xmlData := xmlObj.Encode()
	if xmlData == "" {
		// If name already existed, will be clear.
		//RemovePool(conf.Name)
		http.Error(w, "DomainXML.Encode has error.", http.StatusInternalServerError)
		return
	}
	file := storage.PATH.RootXML() + conf.Name + ".xml"
	_ = libstar.XML.MarshalSave(xmlObj, file, true)

	dom, err := hyper.DomainDefineXML(xmlData)
	if err != nil {
		//RemovePool(conf.Name)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dom.Free()
	if conf.Start == "true" {
		if err := dom.Create(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	domXML := libvirtc.NewDomainXMLFromDom(dom, true)
	if domXML != nil {
		ResponseJson(w, domXML)
	} else {
		ResponseJson(w, xmlObj)
	}
}

func (ins Instance) PUT(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")

	hyper, err := libvirtc.GetHyper()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dom, err := hyper.LookupDomainByUUIDName(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer dom.Free()

	conf := &schema.Instance{}
	if err := GetData(r, conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch conf.Action {
	case "start":
		xmlData, err := dom.GetXMLDesc(true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := dom.Undefine(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if domNew, err := hyper.DomainDefineXML(xmlData); err == nil {
			if err := dom.Create(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = domNew.Free()
			if err := dom.SetAutostart(true); err != nil {
				libstar.Warn("Instance.PUT: start %s", err)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "shutdown":
		if err := dom.ShutdownFlags(libvirtc.DomainShutdownAcpi); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := dom.SetAutostart(false); err != nil {
			libstar.Warn("Instance.PUT: shutdown %s", err)
		}
	case "suspend":
		if err := dom.Suspend(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "reset":
		if err := dom.Reset(0); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "destroy":
		if err := dom.DestroyFlags(libvirtc.DomainDestroyGraceful); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "resume":
		if err := dom.Resume(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "undefine":
		if err := dom.Undefine(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	ResponseMsg(w, 0, "success")
}

func (ins Instance) DELETE(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")

	hyper, err := libvirtc.GetHyper()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dom, err := hyper.LookupDomainByUUIDName(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer dom.Free()

	if ok, _ := dom.IsActive(); ok {
		http.Error(w, "not allowed with active instance", http.StatusInternalServerError)
		return
	}
	name, err := dom.GetName()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := dom.Undefine(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := CleanPool(libvirts.ToDomainPool(name)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}
