package storage

import (
	"strconv"
	"strings"

	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/virsh"
)

type Volume struct {
	Path          string
	Pool          string
	Name          string
	Size          uint64
	Format        string
	BackingFile   string
	BackingFormat string
}

type StorageVolInfo struct {
	Type       string
	Capacity   uint64
	Allocation uint64
}

type StorageVolume struct {
	Pool string
	Name string
	Path string
}

func (v *StorageVolume) Free() error { return nil }

func (v *StorageVolume) run(args ...string) (string, error) {
	h, err := GetHyper()
	if err != nil {
		return "", err
	}
	return virsh.Run(h.Name, args...)
}

func (v *StorageVolume) GetPath() (string, error) {
	if v.Path != "" {
		return v.Path, nil
	}
	out, err := v.run("vol-path", v.Name, "--pool", v.Pool)
	if err != nil {
		return "", err
	}
	v.Path = strings.TrimSpace(out)
	return v.Path, nil
}

func (v *StorageVolume) GetXMLDesc(_ uint32) (string, error) {
	out, err := v.run("vol-dumpxml", v.Name, "--pool", v.Pool)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (v *StorageVolume) Delete(_ uint32) error {
	_, err := v.run("vol-delete", v.Name, "--pool", v.Pool)
	return err
}

func (v *StorageVolume) GetInfo() (StorageVolInfo, error) {
	out, err := v.run("vol-info", v.Name, "--pool", v.Pool)
	if err != nil {
		return StorageVolInfo{}, err
	}
	kv := virsh.ParseKV(out)
	return StorageVolInfo{
		Type:       strings.ToLower(kv["type"]),
		Capacity:   virsh.ParseBytes(kv["capacity"]),
		Allocation: virsh.ParseBytes(kv["allocation"]),
	}, nil
}

func resolvePoolName(pool string) (string, error) {
	if p, err := LookupPoolByUUIDOrName(pool); err == nil {
		return p.Name, nil
	}
	if p, err := LookupPoolByTargetPath(pool); err == nil {
		return p.Name, nil
	}
	return "", libstar.NewErr("pool not found")
}

func CreateVolume(pool, name string, size uint64) (*Volume, error) {
	vol := &Volume{Pool: pool, Name: name, Size: size, Format: "qcow2"}
	return vol, vol.Create()
}

func CreateBackingVolume(pool, name, backingFle, backingFmt string) (*Volume, error) {
	vol := &Volume{Pool: pool, Name: name, BackingFile: backingFle, BackingFormat: backingFmt, Format: "qcow2"}
	return vol, vol.Create()
}

func RemoveVolume(pool string, name string) error {
	vol := &Volume{Pool: pool, Name: name}
	return vol.Remove()
}

func (vol *Volume) Create() error {
	h, err := GetHyper()
	if err != nil {
		return err
	}
	poolName, err := resolvePoolName(vol.Pool)
	if err != nil {
		return err
	}
	volXml := &VolumeXML{
		Name: vol.Name,
		Capacity: CapacityXML{
			Unit:  "bytes",
			Value: strconv.FormatUint(vol.Size, 10),
		},
		Target: TargetXML{Format: FormatXML{Type: vol.Format}},
		BackingStore: BackingStoreXML{
			Path: vol.BackingFile,
			Format: FormatXML{
				Type: vol.BackingFormat,
			},
		},
	}
	xmlData := libstar.XML.Encode(volXml)
	file, cleanup, err := virsh.TempXML("volume-create", xmlData)
	if err != nil {
		return err
	}
	defer cleanup()
	if _, err := virsh.Run(h.Name, "vol-create", "--pool", poolName, file); err != nil {
		return err
	}
	out, err := virsh.Run(h.Name, "vol-path", vol.Name, "--pool", poolName)
	if err != nil {
		return err
	}
	vol.Pool = poolName
	vol.Path = strings.TrimSpace(out)
	return nil
}

func (vol *Volume) GetXMLObj() (*VolumeXML, error) {
	h, err := GetHyper()
	if err != nil {
		return nil, err
	}
	poolName, err := resolvePoolName(vol.Pool)
	if err != nil {
		return nil, err
	}
	xmlData, err := virsh.Run(h.Name, "vol-dumpxml", vol.Name, "--pool", poolName)
	if err != nil {
		return nil, err
	}
	xmlObj := &VolumeXML{}
	return xmlObj, libstar.XML.Decode(xmlObj, xmlData)
}

func (vol *Volume) Remove() error {
	h, err := GetHyper()
	if err != nil {
		return err
	}
	poolName, err := resolvePoolName(vol.Pool)
	if err != nil {
		return err
	}
	_, err = virsh.Run(h.Name, "vol-delete", vol.Name, "--pool", poolName)
	return err
}

func (pol *Pool) LookupStorageVolByName(name string) (*StorageVolume, error) {
	if _, err := pol.run("vol-info", name, "--pool", pol.Name); err != nil {
		return nil, err
	}
	pathOut, _ := pol.run("vol-path", name, "--pool", pol.Name)
	return &StorageVolume{Pool: pol.Name, Name: name, Path: strings.TrimSpace(pathOut)}, nil
}

func VolumeType(t string) string {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "file":
		return "file"
	case "block":
		return "block"
	case "dir":
		return "dir"
	case "netdir":
		return "netdir"
	case "network":
		return "network"
	case "ploop":
		return "ploop"
	default:
		return ""
	}
}

type VolumeInfo struct {
	Pool       string
	Name       string
	Type       string
	Capacity   uint64
	Allocation uint64
}
