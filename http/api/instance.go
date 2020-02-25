package api

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/http/schema"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/storage"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type Instance struct {

}

func InstanceConf2XML(conf *schema.InstanceConf) (libvirtc.DomainXML, error) {
	dom := libvirtc.DomainXML{
		Type: "kvm",
		Name: conf.Name,
		Devices: libvirtc.DevicesXML{
			Disks:      make([]libvirtc.DiskXML, 2),
			Graphics:   make([]libvirtc.GraphicsXML, 1),
			Interfaces: make([]libvirtc.InterfaceXML, 1),
		},
		OS: libvirtc.OSXML{
			Type: libvirtc.OSTypeXML{
				Arch:  conf.Arch,
				Value: "hvm",
			},
			Boot: make([]libvirtc.OSBootXML, 3),
		},
	}
	if dom.OS.Type.Arch == "" {
		dom.OS.Type.Arch = "x86_64"
	}
	// create new disk firstly.
	size := libstar.ToBytes(conf.DiskSize, conf.DiskUnit)
	vol, err := NewVolumeAndPool(conf.DataStore, conf.Name, Slot2Disk(0), size)
	if err != nil {
		return dom, err
	}
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
	dom.CPUXml = libvirtc.CPUXML{
		Placement: "static",
		Value:     conf.Cpu,
	}
	dom.Memory = libvirtc.MemXML{
		Value: conf.MemorySize,
		Type:  conf.MemoryUnit,
	}
	dom.CurMem = libvirtc.CurMemXML{
		Value: conf.MemorySize,
		Type:  conf.MemoryUnit,
	}
	dom.Devices.Graphics[0] = libvirtc.GraphicsXML{
		Type:   "vnc",
		Listen: "0.0.0.0",
		Port:   "-1",
	}
	if strings.HasPrefix(conf.IsoFile, "/dev") {
		dom.Devices.Disks[0] = libvirtc.DiskXML{
			Type:   "block",
			Device: "cdrom",
			Driver: libvirtc.DiskDriverXML{
				Name: "qemu",
				Type: "raw",
			},
			Source: libvirtc.DiskSourceXML{
				Device: conf.IsoFile,
			},
			Target: libvirtc.DiskTargetXML{
				Bus: "ide",
				Dev: "hda",
			},
		}
	} else {
		dom.Devices.Disks[0] = libvirtc.DiskXML{
			Type:   "file",
			Device: "cdrom",
			Driver: libvirtc.DiskDriverXML{
				Name: "qemu",
				Type: "raw",
			},
			Source: libvirtc.DiskSourceXML{
				File: storage.PATH.Unix(conf.IsoFile),
			},
			Target: libvirtc.DiskTargetXML{
				Bus: "ide",
				Dev: "hda",
			},
		}
	}
	dom.Devices.Disks[1] = libvirtc.DiskXML{
		Type:   "file",
		Device: "disk",
		Driver: libvirtc.DiskDriverXML{
			Name: "qemu",
			Type: vol.Target.Format.Type,
		},
		Source: libvirtc.DiskSourceXML{
			File: vol.Target.Path,
		},
		Target: libvirtc.DiskTargetXML{
			Bus: "virtio",
			Dev: "vda",
		},
	}
	dom.Devices.Interfaces[0] = libvirtc.InterfaceXML{
		Type: "bridge",
		Source: libvirtc.InterfaceSourceXML{
			Bridge: conf.Interface,
		},
		Model: libvirtc.InterfaceModelXML{
			Type: "virtio",
		},
	}
	return dom, nil
}

func (ins Instance) Router(router *mux.Router) {
	router.HandleFunc("/api/instance", ins.POST).Methods("POST")
	router.HandleFunc("/api/instance/{id}", ins.GET).Methods("GET")
	router.HandleFunc("/api/instance/{id}", ins.PUT).Methods("PUT")
	router.HandleFunc("/api/instance/{id}", ins.DELETE).Methods("DELETE")
}

func (ins Instance) GET(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")

	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer dom.Free()
	format := GetQueryOne(r, "format")
	if format == "xml" {
		xmlDesc, err := dom.GetXMLDesc(false)
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
	conf := &schema.InstanceConf{}
	if err := GetData(r, conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	xmlObj, err := InstanceConf2XML(conf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	xmlData := xmlObj.Encode()
	if xmlData == "" {
		http.Error(w, "DomainXML.Encode has error.", http.StatusInternalServerError)
		return
	}
	file := storage.PATH.RootXML() + conf.Name + ".xml"
	libstar.XML.MarshalSave(xmlObj, file, true)

	dom, err := hyper.DomainDefineXML(xmlData)
	if err != nil {
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

	conf := &schema.InstanceConf{}
	if err := GetData(r, conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch conf.Action {
	case "start":
		xmlData, err := dom.GetXMLDesc(false)
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
			}
			domNew.Free()
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "shutdown":
		xmlData, err := dom.GetXMLDesc(false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := dom.Shutdown(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		domNew, err := hyper.DomainDefineXML(xmlData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		domNew.Free()
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
		if err := dom.Destroy(); err != nil {
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
