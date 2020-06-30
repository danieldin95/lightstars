package libvirts

import (
	"encoding/xml"
	"github.com/danieldin95/lightstar/libstar"
)

type VolumeXML struct {
	XMLName      xml.Name        `xml:"volume" json:"-"`
	Name         string          `xml:"name" json:"name"`
	Key          string          `xml:"key" json:"key"`
	Source       string          `xml:"source" json:"source"`
	Capacity     CapacityXML     `xml:"capacity" json:"capacity"`
	Allocation   AllocationXML   `xml:"allocation" json:"allocation"`
	Physical     PhysicalXML     `xml:"physical" json:"physical"`
	Target       TargetXML       `xml:"target" json:"target"`
	BackingStore BackingStoreXML `xml:"backingStore" json:"backingStore"`
}

func (vol *VolumeXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), vol); err != nil {
		libstar.Error("VolumeXML.Decode %s", err)
		return err
	}
	return nil
}

func (vol *VolumeXML) Encode() string {
	data, err := xml.Marshal(&vol)
	if err != nil {
		libstar.Error("VolumeXML.Encode %s", err)
		return ""
	}
	return string(data)
}

type CapacityXML struct {
	XMLName xml.Name `xml:"capacity" json:"-"`
	Unit    string   `xml:"unit,attr" json:"unit"`
	Value   string   `xml:",chardata" json:"value"`
}

type AllocationXML struct {
	XMLName xml.Name `xml:"allocation" json:"-"`
	Unit    string   `xml:"unit,attr" json:"unit"`
	Value   string   `xml:",chardata" json:"value"`
}

type PhysicalXML struct {
	XMLName xml.Name `xml:"physical" json:"-"`
	Unit    string   `xml:"unit,attr" json:"unit"`
	Value   string   `xml:",chardata" json:"value"`
}

type FormatXML struct {
	XMLName xml.Name `xml:"format" json:"-"`
	Type    string   `xml:"type,attr" json:"type"`
}

type TargetXML struct {
	XMLName xml.Name  `xml:"target" json:"-"`
	Path    string    `xml:"path" json:"path"`
	Format  FormatXML `xml:"format" json:"format"`
}

type BackingStoreXML struct {
	XMLName xml.Name  `xml:"backingStore" json:"-"`
	Path    string    `xml:"path" json:"path"`
	Format  FormatXML `xml:"format" json:"format"`
}
