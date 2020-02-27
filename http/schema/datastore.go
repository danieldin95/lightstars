package schema

import (
	"github.com/danieldin95/lightstar/storage/libvirts"
)

type DataStore struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	State      string `json:"state"`
	Capacity   uint64 `json:"capacity"`  // bytes
	Allocation uint64 `json:"capacity"`  // bytes
	Available  uint64 `json:"available"` // Bytes
}

func NewDataStore(pol libvirts.Pool) DataStore {
	obj := DataStore{}
	obj.Name, _ = pol.GetName()
	if info, err := pol.GetInfo(); err == nil {
		obj.State = libvirts.PoolState2Str(info.State)
		obj.Capacity = info.Capacity
		obj.Available = info.Available
		obj.Allocation = info.Allocation
	}
	return obj
}
