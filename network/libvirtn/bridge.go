package libvirtn

import (
	"github.com/danieldin95/lightstar/libstar"
)

type Bridge struct {
	Name string `json:"name"`
	Type string `json:"type"` // bridge, ovs etc.
}

type BridgeMgr struct {
	Bridges []Bridge `json:"bridge"`
}

func (br *BridgeMgr) List() []Bridge {
	brs := make([]Bridge, 0, 32)

	hyper, err := GetHyper()
	if err != nil {
		libstar.Warn("IsoMgr.ListFiles %s", err)
		return brs
	}
	if nets, err := hyper.Conn.ListAllNetworks(0); err == nil {
		for _, net := range nets {
			if is, _ := net.IsActive(); !is {
				net.Free()
				continue
			}
			if name, err := net.GetBridgeName(); err == nil {
				brs = append(brs, Bridge{Name: name, Type: "bridge"})
			}
			net.Free()
		}
	}
	return brs
}

var BRIDGE = BridgeMgr{
	Bridges: make([]Bridge, 0, 32),
}
