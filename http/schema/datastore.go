package schema

import (
	"github.com/danieldin95/lightstar/storage/libvirts"
	"unicode"
)

type DataStore struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	State      string `json:"state"`
	Capacity   uint64 `json:"capacity"`   // bytes
	Allocation uint64 `json:"allocation"` // bytes
	Available  uint64 `json:"available"`  // Bytes
	Source     string `json:"source"`
}

func IsDigit(s string) bool {
	for _, v := range s {
		if unicode.IsDigit(v) {
			continue
		}
		return false
	}
	return true
}

func NewDataStore(pol libvirts.Pool) DataStore {
	obj := DataStore{}
	xml, _ := pol.GetXMLDesc(0)
	xmlObj := &libvirts.PoolXML{}
	xmlObj.Decode(xml)

	obj.Name = xmlObj.Name
	obj.UUID = xmlObj.UUID
	if len(obj.Name) == 2 && IsDigit(obj.Name) {
		obj.Name = "datastore@" + obj.Name
	}
	obj.Type = xmlObj.Type
	switch obj.Type {
	case "netfs":
		if xmlObj.Source.Format.Type == "nfs" {
			obj.Source = "nfs://" + xmlObj.Source.Host.Name + xmlObj.Source.Dir.Path
		}
	default:
		obj.Source = obj.Type + "://" + obj.Name
	}
	if info, err := pol.GetInfo(); err == nil {
		obj.State = libvirts.PoolState2Str(info.State)
		obj.Capacity = info.Capacity
		obj.Available = info.Available
		obj.Allocation = info.Allocation
	}
	return obj
}
