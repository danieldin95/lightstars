package libvirtn

import (
	"github.com/danieldin95/lightstar/src/libstar"
)

type Bridge struct {
	Network string `json:"network"`
	Name    string `json:"name"`
	Type    string `json:"type"` // bridge, ovs etc.
}

type BridgeMgr struct {
	Bridges []Bridge `json:"bridge"`
}

func (br *BridgeMgr) List() []Bridge {
	brs := make([]Bridge, 0, 32)

	hyper, err := GetHyper()
	if err != nil {
		libstar.Warn("BridgeMgr.List %s", err)
		return brs
	}
	if nets, err := hyper.Conn.ListAllNetworks(0); err == nil {
		for _, net := range nets {
			if is, _ := net.IsActive(); !is {
				_ = net.Free()
				continue
			}
			br := NewNetworkXMLFromNet(NewNetworkFromVir(&net))
			if br != nil {
				if br.VirtualPort != nil {
					brs = append(brs, Bridge{
						Network: br.Name,
						Name:    br.Bridge.Name,
						Type:    br.VirtualPort.Type,
					})
				} else {
					brs = append(brs, Bridge{
						Network: br.Name,
						Name:    br.Bridge.Name,
						Type:    "bridge",
					})
				}
			}
			_ = net.Free()
		}
	}
	return brs
}

func (br *BridgeMgr) Get(name string) (Bridge, error) {
	b := Bridge{}
	hyper, err := GetHyper()
	if err != nil {
		libstar.Warn("BridgeMgr.Get %s", err)
		return b, err
	}
	if net, err := hyper.Conn.LookupNetworkByName(name); err == nil {
		br := NewNetworkXMLFromNet(NewNetworkFromVir(net))
		if br != nil {
			if br.VirtualPort != nil {
				b = Bridge{Name: br.Name, Type: br.VirtualPort.Type}
			} else {
				b = Bridge{Name: br.Name, Type: "bridge"}
			}
		}
		_ = net.Free()
	}
	return b, nil
}

var BRIDGE = BridgeMgr{
	Bridges: make([]Bridge, 0, 32),
}
