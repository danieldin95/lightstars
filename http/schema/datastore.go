package schema

import (
	"github.com/danieldin95/lightstar/storage/libvirts"
	"unicode"
)

type DataStore struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	State      string `json:"state"`
	Capacity   uint64 `json:"capacity"`  // bytes
	Allocation uint64 `json:"allocation"`  // bytes
	Available  uint64 `json:"available"` // Bytes
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
	obj.Name, _ = pol.GetName()
	obj.UUID, _ = pol.GetUUIDString()

	if len(obj.Name) == 2 && IsDigit(obj.Name) {
		obj.Name = "datastore@" + obj.Name
	}

	if info, err := pol.GetInfo(); err == nil {
		obj.State = libvirts.PoolState2Str(info.State)
		obj.Capacity = info.Capacity
		obj.Available = info.Available
		obj.Allocation = info.Allocation
	}
	return obj
}
