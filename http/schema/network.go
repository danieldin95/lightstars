package schema

import (
	"github.com/danieldin95/lightstar/network/libvirtn"
)

type Network struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Address string `json:"network"`
	Mode    string `json:"mode"` // nat, router.
}

func NewNetwork(net libvirtn.Network) Network {
	obj := Network{}
	obj.Name, _ = net.GetName()
	if ok, _ := net.IsActive(); ok {
		obj.State = "active"
	} else {
		obj.State = "inactive"
	}
	return obj
}
