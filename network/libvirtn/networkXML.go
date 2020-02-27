package libvirtn

import (
	"encoding/xml"
	"github.com/danieldin95/lightstar/libstar"
)

type NetworkXML struct {
	XMLName xml.Name   `xml:"network" json:"-"`
	Name    string     `xml:"name" json:"name"`
	UUID    string     `xml:"uuid" json:"uuid"`
	Forward ForwardXML `xml:"forward" json:"forward"`
	IPv4    IPv4XML    `xml:"ip" json:"ipv4"`
	Bridge  BridgeXML  `xml:"bridge" json:"bridge"`
}

func NewNetworkXMLFromNet(net *Network) *NetworkXML {
	obj := &NetworkXML{}
	if net == nil {
		return nil
	}
	if desc, err := net.GetXMLDesc(0); err == nil {
		if err := obj.Decode(desc); err != nil {
			return obj
		}
	}
	return obj
}

func (net *NetworkXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), net); err != nil {
		libstar.Error("NetworkXML.Decode %s", err)
		return err
	}
	return nil
}

func (net *NetworkXML) Encode() string {
	data, err := xml.Marshal(net)
	if err != nil {
		libstar.Error("NetworkXML.Encode %s", err)
		return ""
	}
	return string(data)
}

type ForwardXML struct {
	XMLName xml.Name `xml:"forward" json:"-"`
	Mode    string   `xml:"mode,attr" json:"mode"`
}

type IPv4XML struct {
	XMLName xml.Name `xml:"ip" json:"-"`
	Address string   `xml:"address,attr" json:"address"`
	Netmask string   `xml:"netmask,attr" json:"netmask"`
}

type BridgeXML struct {
	XMLName xml.Name `xml:"bridge" json:"-"`
	Name    string   `xml:"name,attr" json:"name"`
	Stp     string   `xml:"stp,attr,omitempty" json:"stp,omitempty"`     // on,off
	Delay   string   `xml:"delay,attr,omitempty" json:"delay,omitempty"` // 0-32
}
