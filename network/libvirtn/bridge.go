package libvirtn

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/libvirt/libvirt-go"
)

type Bridge struct {
	Name       string `json:"name"`
	Type      string    `json:"type"` // bridge, ovs etc.
}

type BridgeMgr struct {
	Conn    *libvirt.Connect `json:"-"`
	Bridges []Bridge      `json:"bridge"`
}

func (br *BridgeMgr) Open() error {
	if br.Conn == nil {
		hyper, err := libvirtc.GetHyper()
		if err != nil {
			return err
		}
		br.Conn = hyper.Conn
	}
	if br.Conn == nil {
		return libstar.NewErr("Not found libvirt.Connect")
	}
	return nil
}

func (br *BridgeMgr) List() []Bridge {
	brs := make([]Bridge, 0, 32)

	if err := br.Open(); err != nil {
		libstar.Warn("IsoMgr.ListFiles %s", err)
		return brs
	}
	if nets, err := br.Conn.ListAllNetworks(0); err == nil {
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

