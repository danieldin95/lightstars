package libvirts

import (
	"encoding/xml"
	"github.com/danieldin95/lightstar/libstar"
)

type PoolXML struct {
	XMLName    xml.Name      `xml:"pool" json:"-"`
	Type       string        `xml:"type,attr" json:"type"`
	Name       string        `xml:"name" json:"name"`
	UUID       string        `xml:"uuid" json:"uuid"`
	Source     string        `xml:"source" json:"source"`
	Capacity   CapacityXML   `xml:"capacity", json:"capacity"`
	Allocation AllocationXML `xml:"allocation" json:"allocation"`
	Physical   PhysicalXML   `xml:"physical" json:"physical"`
	Target     TargetXML     `xml:"target" json:"target"`
}

func (pol *PoolXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), pol); err != nil {
		libstar.Error("PoolXML.Decode %s", err)
		return err
	}
	return nil
}

func (pol *PoolXML) Encode() string {
	data, err := xml.Marshal(pol)
	if err != nil {
		libstar.Error("PoolXML.Encode %s", err)
		return ""
	}
	return string(data)
}
