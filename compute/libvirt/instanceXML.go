package libvirtdriver

import (
	"encoding/xml"
	"github.com/danieldin95/lightstar/libstar"
)

type DomainXML struct {
	XMLName xml.Name   `xml:"domain"`
	Id      string     `xml:"id,attr"`
	Type    string     `xml:"type,attr"`
	Name    string      `xml:"name"`
	Uuid    string      `xml:"uuid"`
	Devices DevicesXML  `xml:"devices"`
}

func (domain *DomainXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), domain); err != nil {
		libstar.Error("DomainXML.Decode %s", err)
		return err
	}
	return nil
}

func (domain *DomainXML) VNCDisplay() (string, string) {
	if len(domain.Devices.Graphics) == 0 {
		return "", ""
	}
	for _, g := range domain.Devices.Graphics {
		if g.Type == "vnc" {
			return g.Listen, g.Port
		}
	}
	return "",""
}

type DevicesXML struct {
	XMLName  xml.Name      `xml:"devices"`
	Graphics []GraphicsXML `xml:"graphics"`
}

func (devices *DevicesXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), devices); err != nil {
		libstar.Error("DevicesXML.Decode %s", err)
		return err
	}
	return nil
}

type GraphicsXML struct {
	XMLName  xml.Name    `xml:"graphics"`
	Type     string      `xml:"type,attr"`
	Port     string      `xml:"port,attr"`
	Listen   string      `xml:"listen,attr"`
}

func (graphics *GraphicsXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), graphics); err != nil {
		libstar.Error("GraphicsXML.Decode %s", err)
		return err
	}
	return nil
}