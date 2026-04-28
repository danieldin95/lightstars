package network

import (
	"github.com/danieldin95/lightstar/pkg/libstar"
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
	nets, err := hyper.ListAllNetworks()
	if err != nil {
		return brs
	}
	for _, net := range nets {
		if is, _ := net.IsActive(); !is {
			continue
		}
		xmlObj := NewNetworkXMLFromNet(&net)
		if xmlObj == nil {
			continue
		}
		if xmlObj.VirtualPort != nil {
			brs = append(brs, Bridge{Network: xmlObj.Name, Name: xmlObj.Bridge.Name, Type: xmlObj.VirtualPort.Type})
		} else {
			brs = append(brs, Bridge{Network: xmlObj.Name, Name: xmlObj.Bridge.Name, Type: "bridge"})
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
	net, err := hyper.LookupNetwork(name)
	if err != nil {
		return b, nil
	}
	xmlObj := NewNetworkXMLFromNet(net)
	if xmlObj != nil {
		if xmlObj.VirtualPort != nil {
			b = Bridge{Name: xmlObj.Name, Type: xmlObj.VirtualPort.Type}
		} else {
			b = Bridge{Name: xmlObj.Name, Type: "bridge"}
		}
	}
	return b, nil
}

var BRIDGE = BridgeMgr{Bridges: make([]Bridge, 0, 32)}
