package libvirtdriver

import (
	"encoding/xml"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/libvirt/libvirt-go"
)

type DomainXML struct {
	XMLName xml.Name   `xml:"domain" json:"-"`
	Id      string     `xml:"id,attr" json:"id"`
	Type    string     `xml:"type,attr" json:"type"`
	Name    string     `xml:"name" json:"name"`
	Uuid    string     `xml:"uuid" json:"uuid"`
	Devices DevicesXML `xml:"devices" json:"devices"`
}

func NewDomainXMLFromDom(dom *Domain, secure bool) *DomainXML {
	var xmlData string
	var err error

	if dom == nil {
		return nil
	}
	if secure {
		xmlData, err = dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	} else {
		xmlData, err = dom.GetXMLDesc(0)
	}
	if err != nil {
		return nil
	}
	instXml := &DomainXML{}
	if err := instXml.Decode(xmlData); err != nil {
		return nil
	}
	return instXml
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
	return "", ""
}

type DevicesXML struct {
	XMLName  xml.Name      `xml:"devices" json:"-"`
	Graphics []GraphicsXML `xml:"graphics" json:"graphics"`
}

func (devices *DevicesXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), devices); err != nil {
		libstar.Error("DevicesXML.Decode %s", err)
		return err
	}
	return nil
}

type GraphicsXML struct {
	XMLName xml.Name `xml:"graphics" json:"-"`
	Type    string   `xml:"type,attr" json:"type"`
	Port    string   `xml:"port,attr" json:"port"`
	Listen  string   `xml:"listen,attr" json:"listen"`
}

func (graphics *GraphicsXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), graphics); err != nil {
		libstar.Error("GraphicsXML.Decode %s", err)
		return err
	}
	return nil
}
