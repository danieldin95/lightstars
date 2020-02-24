package libvirtc

import (
	"encoding/xml"
	"github.com/danieldin95/lightstar/libstar"
)

type DomainXML struct {
	XMLName xml.Name   `xml:"domain" json:"-"`
	Id      string     `xml:"id,attr" json:"id"`
	Type    string     `xml:"type,attr" json:"type"` // kvm
	Name    string     `xml:"name" json:"name"`
	UUID    string     `xml:"uuid" json:"uuid"`
	OS      OSXML      `xml:"os" json:"os"`
	CPUXml  CPUXML     `xml:"vcpu" json:"cpu"`
	Memory  MemXML     `xml:"memory" json:"memory"`
	CurMem  CurMemXML  `xml:"currentMemory" json:"currentMemory"`
	Devices DevicesXML `xml:"devices" json:"devices"`
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

type CPUXML struct {
	XMLName   xml.Name `xml:"vcpu" json:"-"`
	Placement string   `xml:"placement,attr" json:"placement"` // static
	Value     string   `xml:",chardata" json:"Value"`
}

func (cpu *CPUXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), cpu); err != nil {
		libstar.Error("CPUXML.Decode %s", err)
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

type OSXML struct {
	XMLName xml.Name    `xml:"os" json:"-"`
	Type    OSTypeXML   `xml:"type" json:"type"`
	Boot    []OSBootXML `xml:"boot" json:"boot"`
}

func (osx *OSXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), osx); err != nil {
		libstar.Error("OSXML.Decode %s", err)
		return err
	}
	return nil
}

type OSTypeXML struct {
	XMLName xml.Name `xml:"type" json:"-"`
	Arch    string   `xml:"arch,attr" json:"arch"` // x86_64
	Machine string   `xml:"machine,attr" json:"machine"`
	Value   string   `xml:",chardata" json:"value"` // hvm
}

type OSBootXML struct {
	XMLName xml.Name `xml:"boot" json:"-"`
	Dev     string   `xml:"dev,attr" json:"dev"` // hd, cdrom, network
}

type DevicesXML struct {
	XMLName     xml.Name        `xml:"devices" json:"-"`
	Graphics    []GraphicsXML   `xml:"graphics" json:"graphics"`
	Disks       []DiskXML       `xml:"disk" json:"disk"`
	Interfaces  []InterfaceXML  `xml:"interface" json:"interface"`
	Controllers []ControllerXML `xml:"controller" json:"controller"`
}

func (devices *DevicesXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), devices); err != nil {
		libstar.Error("DevicesXML.Decode %s", err)
		return err
	}
	return nil
}

// Bus 0: root
// Bus 1: disk, ide
// Bus 2: interface
// Bus 3: reverse
type ControllerXML struct {
	XMLName xml.Name   `xml:"controller" json:"-"`
	Type    string     `xml:"type,attr" json:"type"`
	Index   string     `xml:"index,attr" json:"port"`
	Model   string     `xml:"model,attr" json:"model"` // pci-root, pci-bridge.
	Address AddressXML `xml:"address" json:"address"`
}

func (ctl *ControllerXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), ctl); err != nil {
		libstar.Error("ControllerXML.Decode %s", err)
		return err
	}
	return nil
}

type AddressXML struct {
	XMLName  xml.Name `xml:"address" json:"-"`
	Type     string   `xml:"type,attr,omitempty" json:"type,omitempty"`
	Domain   string   `xml:"domain,attr,omitempty" json:"domain,omitempty"`
	Bus      string   `xml:"bus,attr,omitempty" json:"bus,omitempty"`
	Slot     string   `xml:"slot,attr,omitempty" json:"slot,omitempty"`
	Function string   `xml:"function,attr,omitempty" json:"function,omitempty"`
}

type GraphicsXML struct {
	XMLName xml.Name `xml:"graphics" json:"-"`
	Type    string   `xml:"type,attr" json:"type"` // vnc, spice
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

type DiskSourceXML struct {
	XMLName xml.Name `xml:"source" json:"-"`
	File    string   `xml:"file,attr,omitempty" json:"file,omitempty"`
	Device  string   `xml:"device,attr,omitempty" json:"device,omitempty"`
}

type DiskTargetXML struct {
	XMLName xml.Name `xml:"target" json:"-"`
	Bus     string   `xml:"bus,attr,omitempty" json:"bus,omitempty"`
	Dev     string   `xml:"dev,attr,omitempty" json:"dev,omitempty"`
}

type InterfaceXML struct {
	XMLName     xml.Name                `xml:"interface" json:"-"`
	Type        string                  `xml:"type,attr" json:"type"`
	Mac         InterfaceMacXML         `xml:"mac" json:"mac"`
	Source      InterfaceSourceXML      `xml:"source" json:"source"`
	Model       InterfaceModelXML       `xml:"model" json:"model"`
	Target      InterfaceTargetXML      `xml:"target" json:"tatget"`
	VirtualPort InterfaceVirtualPortXML `xml:"virtualport" json:"virtualport"`
}

func (int *InterfaceXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), int); err != nil {
		libstar.Error("InterfaceXML.Decode %s", err)
		return err
	}
	return nil
}

type InterfaceMacXML struct {
	XMLName xml.Name `xml:"mac" json:"-"`
	Address string   `xml:"address,attr,omitempty" json:"address,omitempty"`
}

type InterfaceSourceXML struct {
	XMLName xml.Name `xml:"source" json:"-"`
	Bridge  string   `xml:"bridge,attr,omitempty" json:"bridge,omitempty"`
}

type InterfaceModelXML struct {
	XMLName xml.Name `xml:"model" json:"-"`
	Type    string   `xml:"type,attr" json:"type"` //rtl8139, virtio, e1000
}

type InterfaceTargetXML struct {
	XMLName xml.Name `xml:"target" json:"-"`
	Bus     string   `xml:"bus,attr,omitempty" json:"bus,omitempty"`
	Dev     string   `xml:"dev,attr,omitempty" json:"dev,omitempty"`
}

type InterfaceVirtualPortXML struct {
	XMLName xml.Name `xml:"virtualport" json:"-"`
	Type    string   `xml:"type,attr,omitempty" json:"type,omitempty"` //openvswitch
}
