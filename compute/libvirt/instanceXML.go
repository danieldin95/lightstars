package libvirtdriver

import (
	"encoding/xml"
	"github.com/danieldin95/lightstar/libstar"
)

type DomainXML struct {
	XMLName xml.Name   `xml:"domain" json:"-"`
	Id      string     `xml:"id,attr" json:"id"`
	Type    string     `xml:"type,attr" json:"type"`
	Name    string     `xml:"name" json:"name"`
	Uuid    string     `xml:"uuid" json:"uuid"`
	Devices DevicesXML `xml:"devices" json:"devices"`
	Memory  MemXML     `xml:"memory" json:"memory"`
	CurMem  CurMemXML  `xml:"currentMemory" json:"currentMemory"`
	VCPUXml VCPUXML    `xml:"vcpu" json:"vcpu"`
}

func NewDomainXMLFromDom(dom *Domain, secure bool) *DomainXML {
	var xmlData string
	var err error

	if dom == nil {
		return nil
	}
	xmlData, err = dom.GetXMLDesc(secure)
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

func (domain *DomainXML) Encode() string {
	data, err := xml.Marshal(domain)
	if err != nil {
		libstar.Error("DomainXML.Encode %s", err)
		return ""
	}
	return string(data)
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

type VCPUXML struct {
	XMLName xml.Name `xml:"vcpu" json:"-"`
	Type    string   `xml:"placement,attr" json:"placement"`
	Value   string   `xml:",chardata" json:"Value"`
}

func (cpu *VCPUXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), cpu); err != nil {
		libstar.Error("VCPUXML.Decode %s", err)
		return err
	}
	return nil
}

type MemXML struct {
	XMLName xml.Name `xml:"memory" json:"-"`
	Type    string   `xml:"unit,attr" json:"unit"`
	Value   string   `xml:",chardata" json:"value"`
}

func (mem *MemXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), mem); err != nil {
		libstar.Error("MemXML.Decode %s", err)
		return err
	}
	return nil
}

type CurMemXML struct {
	XMLName xml.Name `xml:"currentMemory" json:"-"`
	Type    string   `xml:"unit,attr" json:"unit"`
	Value   string   `xml:",chardata" json:"value"`
}

func (cmem *CurMemXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), cmem); err != nil {
		libstar.Error("CurMemXML.Decode %s", err)
		return err
	}
	return nil
}

type DevicesXML struct {
	XMLName  xml.Name      `xml:"devices" json:"-"`
	Graphics []GraphicsXML `xml:"graphics" json:"graphics"`
	Disks    []DiskXML     `xml:"disk" json:"disk"`
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

type DiskXML struct {
	XMLName xml.Name      `xml:"disk" json:"-"`
	Type    string        `xml:"type,attr" json:"type"`
	Device  string        `xml:"device,attr" json:"device"`
	Driver  DiskDriverXML `xml:"driver" json:"driver"`
	Source  DiskSourceXML `xml:"source" json:"source"`
	Target  DiskTargetXML `xml:"target" json:"target"`
}

func (disk *DiskXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), disk); err != nil {
		libstar.Error("DiskXML.Decode %s", err)
		return err
	}
	return nil
}

type DiskDriverXML struct {
	XMLName xml.Name `xml:"driver" json:"-"`
	Type    string   `xml:"type,attr" json:"type"`
	Name    string   `xml:"name,attr" json:"name"`
}

func (drv *DiskDriverXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), drv); err != nil {
		libstar.Error("DiskDriverXML.Decode %s", err)
		return err
	}
	return nil
}

type DiskSourceXML struct {
	XMLName xml.Name `xml:"source" json:"-"`
	Name    string   `xml:"file,attr" json:"file"`
	Device  string   `xml:"device,attr,omitempty" json:"device,omitempty"`
}

func (src *DiskSourceXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), src); err != nil {
		libstar.Error("DiskSourceXML.Decode %s", err)
		return err
	}
	return nil
}

type DiskTargetXML struct {
	XMLName xml.Name `xml:"target" json:"-"`
	Bus     string   `xml:"bus,attr" json:"bus"`
	Dev     string   `xml:"dev,attr" json:"dev"`
}

func (tgt *DiskTargetXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), tgt); err != nil {
		libstar.Error("DiskTargetXML.Decode %s", err)
		return err
	}
	return nil
}
