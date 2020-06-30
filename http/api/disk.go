package api

import (
	"github.com/danieldin95/lightstar/compute"
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/schema"
	"github.com/danieldin95/lightstar/storage"
	"github.com/danieldin95/lightstar/storage/libvirts"
	"github.com/gorilla/mux"
	"net/http"
	"path"
	"sort"
	"strings"
)

type Disk struct {
}

func (disk Disk) Router(router *mux.Router) {
	router.HandleFunc("/api/instance/{id}/disk", disk.GET).Methods("GET")
	router.HandleFunc("/api/instance/{id}/disk", disk.POST).Methods("POST")
	router.HandleFunc("/api/instance/{id}/disk/{dev}", disk.GET).Methods("GET")
	router.HandleFunc("/api/instance/{id}/disk/{dev}", disk.DELETE).Methods("DELETE")
}

func (disk Disk) Travel(instance schema.Instance) map[string]libvirts.VolumeInfo {
	name := instance.Name
	disks := instance.Disks

	vols := make(map[string]libvirts.VolumeInfo, 32)
	sources := make(map[string]int, 4)
	for _, disk := range disks {
		if _, ok := vols[disk.Name]; ok {
			continue
		}
		dir := path.Dir(disk.Name)
		volsDir, err := (&libvirts.Pool{Path: dir}).ListByTarget()
		if err == nil {
			for file, vol := range volsDir {
				vols[file] = vol
			}
			continue
		}
		if _, ok := sources[dir]; ok {
			sources[dir] += 1
		} else {
			sources[dir] = 1
		}
	}
	libstar.Debug("Disk.Check %v", sources)

	curDir := ""
	curUsed := 0
	for dir, c := range sources {
		if curUsed < c {
			curUsed = c
			curDir = dir
		} else if curUsed == c {
			// select has more deep directory.
			if len(strings.Split(curDir, "/")) < len(strings.Split(dir, "/")) {
				curDir = dir
			}
		}
	}
	if curDir != "" {
		if _, err := libvirts.CreatePool(libvirts.ToDomainPool(name), curDir); err != nil {
			libstar.Warn("Disk.Travel %s", err)
		}
	}
	volsDir, err := (&libvirts.Pool{Name: libvirts.ToDomainPool(name)}).List()
	if err == nil {
		for file, vol := range volsDir {
			vols[file] = vol
		}
	}
	return vols
}

func (disk Disk) GET(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	dev, ok := GetArg(r, "dev")
	if !ok {
		dom, err := libvirtc.LookupDomainByUUIDString(uuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		defer dom.Free()

		instance := compute.NewInstance(*dom)
		vols := disk.Travel(instance)

		list := schema.ListDisk{
			Items: make([]schema.Disk, 0, 32),
		}
		for _, disk := range instance.Disks {
			if vol, ok := vols[disk.Name]; ok {
				disk.Volume = schema.Volume{
					Type:       vol.Type,
					Capacity:   vol.Capacity,
					Allocation: vol.Allocation,
				}
			}
			list.Items = append(list.Items, disk)
		}
		sort.SliceStable(list.Items, func(i, j int) bool {
			return list.Items[i].Device < list.Items[j].Device
		})
		list.Metadata.Size = len(list.Items)
		list.Metadata.Total = len(list.Items)
		ResponseJson(w, list)
		return
	} else {
		//TODO
	}
	ResponseMsg(w, 0, dev)
}

func IsVolume(file string) bool {
	if !strings.HasPrefix(file, storage.Location) {
		return false
	}
	if strings.HasSuffix(file, ".img") || strings.HasSuffix(file, ".qcow2") {
		return true
	}
	return false
}

func Disk2XML(conf *schema.Disk) (*libvirtc.DiskXML, error) {
	// create new disk firstly.
	size := libstar.ToBytes(conf.Size, conf.SizeUnit)
	slot := libstar.H2D8(conf.Seq)
	name := libvirtc.DISK.Slot2Name(slot)
	vol, err := NewVolume(conf.Name, name, size)
	if err != nil {
		return nil, err
	}
	xml := libvirtc.DiskXML{
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
			Bus: conf.Bus,
			Dev: libvirtc.DISK.Slot2Dev(conf.Bus, slot),
		},
	}

	switch conf.Bus {
	case "virtio":
		xml.Address = &libvirtc.AddressXML{
			Type:     "pci",
			Domain:   libvirtc.PCI_DOMAIN,
			Bus:      libvirtc.PCI_DISK_BUS,
			Slot:     conf.Seq,
			Function: libvirtc.PCI_FUNC,
		}
	}
	return &xml, nil
}

func (disk Disk) POST(w http.ResponseWriter, r *http.Request) {
	conf := &schema.Disk{}
	if err := GetData(r, conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uuid, _ := GetArg(r, "id")
	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dom.Free()

	if conf.Name == "" {
		conf.Name, _ = dom.GetName()
	}
	xmlObj, err := Disk2XML(conf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	libstar.Debug("Disk.POST: %s", xmlObj.Encode())
	flags := libvirtc.DOMAIN_DEVICE_MODIFY_PERSISTENT
	if active, _ := dom.IsActive(); !active {
		flags = libvirtc.DOMAIN_DEVICE_MODIFY_CONFIG
	}
	if err := dom.AttachDeviceFlags(xmlObj.Encode(), flags); err != nil {
		file := xmlObj.Source.File
		if IsVolume(file) {
			volume := path.Base(file)
			_ = libvirts.RemoveVolume(libvirts.ToDomainPool(conf.Name), volume)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}

func (disk Disk) PUT(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (disk Disk) DELETE(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dom.Free()

	dev, _ := GetArg(r, "dev")
	xml := libvirtc.NewDomainXMLFromDom(dom, true)
	if xml == nil {
		http.Error(w, "Cannot get domain's descXML", http.StatusInternalServerError)
		return
	}

	if xml.Devices.Disks != nil {
		for _, disk := range xml.Devices.Disks {
			if disk.Target.Dev != dev {
				continue
			}
			// found device
			flags := libvirtc.DOMAIN_DEVICE_MODIFY_PERSISTENT
			if active, _ := dom.IsActive(); !active {
				flags = libvirtc.DOMAIN_DEVICE_MODIFY_CONFIG
			}
			if err := dom.DetachDeviceFlags(disk.Encode(), flags); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			file := disk.Source.File
			if IsVolume(file) {
				dir := path.Dir(file)
				volume := path.Base(file)
				pool := path.Base(dir)
				libvirts.RemoveVolume(libvirts.ToDomainPool(pool), volume)
			}
		}
	}
	ResponseMsg(w, 0, "success")
}
