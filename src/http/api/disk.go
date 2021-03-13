package api

import (
	"github.com/danieldin95/lightstar/src/compute"
	"github.com/danieldin95/lightstar/src/compute/libvirtc"
	"github.com/danieldin95/lightstar/src/libstar"
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

type Disk struct {
}

func (disk Disk) Router(router *mux.Router) {
	router.HandleFunc("/api/instance/{id}/disk", disk.Get).Methods("GET")
	router.HandleFunc("/api/instance/{id}/disk", disk.Post).Methods("POST")
	router.HandleFunc("/api/instance/{id}/disk/{dev}", disk.Get).Methods("GET")
	router.HandleFunc("/api/instance/{id}/disk/{dev}", disk.Delete).Methods("DELETE")
}

func (disk Disk) Travel(instance schema.Instance) map[string]libvirts.VolumeInfo {
	name := instance.Name
	disks := instance.Disks

	vols := make(map[string]libvirts.VolumeInfo, 32)
	sources := make(map[string]int, 4)
	for _, disk := range disks { // to traver all disks and record it's path.
		if _, ok := vols[disk.Name]; ok {
			continue
		}
		dir := path.Dir(disk.Name)
		if volsDir, err := (&libvirts.Pool{Path: dir}).ListByTarget(); err == nil {
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
	// figure out default pool
	for dir, c := range sources {
		if curUsed < c { // dir has more used wined.
			curUsed = c
			curDir = dir
		} else if curUsed == c {
			// select has more deep directory.
			if len(strings.Split(curDir, "/")) < len(strings.Split(dir, "/")) {
				curDir = dir
			}
		}
	}
	if curDir != "" { // try to create it.
		if _, err := libvirts.CreatePool(libvirts.ToDomainPool(name), curDir); err != nil {
			libstar.Warn("Disk.Travel %s", err)
		}
	}
	pol := &libvirts.Pool{Name: libvirts.ToDomainPool(name)}
	if volsDir, err := pol.List(); err == nil {
		for file, vol := range volsDir {
			vols[file] = vol
		}
	}
	return vols
}

func (disk Disk) Get(w http.ResponseWriter, r *http.Request) {
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
	if strings.HasSuffix(file, ".img") || strings.HasSuffix(file, ".qcow2") ||
		strings.HasSuffix(file, ".raw") || strings.HasSuffix(file, ".qcow") {
		return true
	}
	return false
}

func Disk2XML(conf *schema.Disk) (*libvirtc.DiskXML, error) {
	xml := libvirtc.DiskXML{}
	if conf.Source == "" { // create new disk firstly.
		size := libstar.ToBytes(conf.Size, conf.SizeUnit)
		slot := libstar.H2D8(conf.Seq)
		name := libvirtc.DISK.Slot2Name(slot)
		vol, err := NewVolume(conf.Name, name, size)
		if err != nil {
			return nil, err
		}
		xml = libvirtc.DiskXML{
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
	} else if strings.HasSuffix(conf.Source, ".iso") ||
		strings.HasSuffix(conf.Source, ".ISO") {
		// attach cdrom.
		file := storage.PATH.Unix(conf.Source)
		seq, _ := strconv.Atoi(conf.Seq)
		xml = libvirtc.DiskXML{
			Type:   "file",
			Device: "cdrom",
			Driver: libvirtc.DiskDriverXML{
				Type: "raw",
				Name: "qemu",
			},
			Source: libvirtc.DiskSourceXML{
				File: file,
			},
			Target: libvirtc.DiskTargetXML{
				Bus: "ide",
				Dev: libvirtc.DISK.Slot2Dev("ide", uint8(seq)),
			},
		}
	}
	switch conf.Bus {
	case "virtio":
		xml.Address = &libvirtc.AddressXML{
			Type:     "pci",
			Domain:   libvirtc.PciDomain,
			Bus:      libvirtc.PciDiskBus,
			Slot:     conf.Seq,
			Function: libvirtc.PciFunc,
		}
	}
	return &xml, nil
}

func (disk Disk) Post(w http.ResponseWriter, r *http.Request) {
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
	libstar.Debug("Disk.Post: %s", xmlObj.Encode())
	flags := libvirtc.DomainDeviceModifyPersistent
	if active, _ := dom.IsActive(); !active {
		flags = libvirtc.DomainDeviceModifyConfig
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
	ResponseMsg(w, 0, xmlObj.Target.Dev)
}

func (disk Disk) FindByDev(devices *libvirtc.DevicesXML, dev string) *libvirtc.DiskXML {
	if devices == nil || devices.Disks == nil {
		return nil
	}
	for _, disk := range devices.Disks {
		if disk.Target.Dev != dev {
			continue
		}
		// found device
		return &disk
	}
	return nil
}

func (disk Disk) Delete(w http.ResponseWriter, r *http.Request) {
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
	if d := disk.FindByDev(&xml.Devices, dev); d != nil {
		// found device
		flags := libvirtc.DomainDeviceModifyPersistent
		if active, _ := dom.IsActive(); !active {
			flags = libvirtc.DomainDeviceModifyConfig
		}
		if err := dom.DetachDeviceFlags(d.Encode(), flags); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file := d.Source.File
		if IsVolume(file) {
			dir := path.Dir(file)
			volume := path.Base(file)
			pool := path.Base(dir)
			_ = libvirts.RemoveVolume(libvirts.ToDomainPool(pool), volume)
		}
	}
	ResponseMsg(w, 0, "")
}
