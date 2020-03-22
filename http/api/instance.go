package api

import (
	"fmt"
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/network/libvirtn"
	"github.com/danieldin95/lightstar/schema"
	"github.com/danieldin95/lightstar/storage"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Instance struct {
}

func GetTypeBySuffix(name string) string {
	name = strings.ToUpper(name)
	if strings.HasSuffix(name, ".ISO") {
		return "cdrom"
	} else if strings.HasSuffix(name, ".RAW") {
		return "raw"
	} else if strings.HasSuffix(name, ".QCOW2") || strings.HasSuffix(name, ".IMG") {
		return "qcow2"
	} else if strings.HasSuffix(name, ".VMDK") {
		return "vmdk"
	}
	return "raw"
}

func NewCDROMXML(file, family string, seq uint8) libvirtc.DiskXML {
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
			Dev: libvirtc.DISK.Slot2Dev("ide", seq),
		},
	}
}

func NewISOXML(file, family string, seq uint8) libvirtc.DiskXML {
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
	xml.Driver.Type = GetTypeBySuffix(name)
	if family == "linux" && !strings.HasSuffix(name, ".ISO") {
		xml.Target = libvirtc.DiskTargetXML{
			Bus: "virtio",
			Dev: libvirtc.DISK.Slot2Dev("virtio", seq),
		}
	} else {
		xml.Target = libvirtc.DiskTargetXML{
			Bus: "ide",
			Dev: libvirtc.DISK.Slot2Dev("ide", seq),
		}
	}
	return xml
}

func NewDiskXML(format, file, bus string, seq uint8) libvirtc.DiskXML {
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
		disk.Target = libvirtc.DiskTargetXML{
			Bus: bus,
			Dev: libvirtc.DISK.Slot2Dev(bus, seq),
		}
		disk.Address = &libvirtc.AddressXML{
			Type:     "pci",
			Domain:   libvirtc.PCI_DOMAIN,
			Bus:      libvirtc.PCI_DISK_BUS,
			Slot:     fmt.Sprintf("0x%x", seq),
			Function: libvirtc.PCI_FUNC,
		}
	case "ide", "scsi":
		disk.Target = libvirtc.DiskTargetXML{
			Bus: bus,
			Dev: libvirtc.DISK.Slot2Dev(bus, seq),
		}
	}
	return disk
}

func Instance2XML(conf *schema.Instance) (libvirtc.DomainXML, error) {
	dom := libvirtc.DomainXML{
		Type: "kvm",
		Name: conf.Name,
		Devices: libvirtc.DevicesXML{
			Disks:       make([]libvirtc.DiskXML, 0, 2),
			Graphics:    make([]libvirtc.GraphicsXML, 2), // vnc and spice
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
	// spice
	dom.Devices.Graphics[1] = libvirtc.GraphicsXML{
		Type:     "spice",
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
	for i, disk := range conf.Disks {
		file := disk.Source
		size := disk.Size
		unit := disk.SizeUnit
		seq := uint8(i + 1)

		obj := libvirtc.DiskXML{}
		if strings.HasPrefix(file, "/dev") {
			obj = NewCDROMXML(file, conf.Family, seq)
		} else if strings.HasSuffix(file, ".iso") || strings.HasSuffix(file, ".ISO") {
			obj = NewISOXML(storage.PATH.Unix(file), conf.Family, seq)
		} else {
			format := GetTypeBySuffix(file)
			if file == "" {
				size := libstar.ToBytes(size, unit)
				vol, err := NewVolumeAndPool(conf.DataStore, conf.Name, Slot2Disk(seq), size)
				if err != nil {
					return dom, err
				}
				file = vol.Target.Path
				format = vol.Target.Format.Type
			} else {
				file = storage.PATH.Unix(file)
			}
			switch conf.Family {
			case "linux":
				obj = NewDiskXML(format, file, "virtio", seq)
			case "windows": // not scsi.
				obj = NewDiskXML(format, file, "ide", seq)
			default:
				obj = NewDiskXML(format, file, "ide", seq)
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

func (ins Instance) GetByUser(user *schema.User, list *schema.List) {
	if domains, err := libvirtc.ListDomains(); err == nil {
		for _, d := range domains {
			inst := schema.NewInstance(d)
			if ins.HasPermission(user, inst.Name) {
				list.Items = append(list.Items, inst)
			}
			d.Free()
		}
	}
}

func (ins Instance) GET(w http.ResponseWriter, r *http.Request) {
	uuid, ok := GetArg(r, "id")
	if !ok {
		user, _ := GetUser(r)
		list := schema.List{
			Items:    make([]interface{}, 0, 32),
			Metadata: schema.MetaData{},
		}
		// list all instances.
		ins.GetByUser(&user, &list)
		sort.SliceStable(list.Items, func(i, j int) bool {
			return list.Items[i].(schema.Instance).Name < list.Items[j].(schema.Instance).Name
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
	} else if format == "schema" {
		ResponseJson(w, schema.NewInstance(*dom))
	} else {
		ResponseJson(w, libvirtc.NewDomainXMLFromDom(dom, true))
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

	xmlObj, err := Instance2XML(conf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// need release created images if fails.
	xmlData := xmlObj.Encode()
	if xmlData == "" {
		//DelVolumeAndPool(conf.Name)
		http.Error(w, "DomainXML.Encode has error.", http.StatusInternalServerError)
		return
	}
	file := storage.PATH.RootXML() + conf.Name + ".xml"
	libstar.XML.MarshalSave(xmlObj, file, true)

	dom, err := hyper.DomainDefineXML(xmlData)
	if err != nil {
		//DelVolumeAndPool(conf.Name)
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
			domNew.Free()
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "shutdown":
		if err := dom.ShutdownFlags(libvirtc.DOMAIN_SHUTDOWN_ACPI); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
		if err := dom.DestroyFlags(libvirtc.DOMAIN_DESTROY_GRACEFUL); err != nil {
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
	if err := DelVolumeAndPool(name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "")
}
