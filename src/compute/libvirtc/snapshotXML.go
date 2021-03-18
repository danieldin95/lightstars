package libvirtc

import (
	"encoding/xml"
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/libvirt/libvirt-go"
)

type SnapshotXML struct {
	libstar.XMLBase
	XMLName   xml.Name `xml:"domainsnapshot" json:"-"`
	Name      string   `xml:"name,omitempty" json:"name"`
	State     string   `xml:"state,omitempty" json:"state"`
	CreateAt  int64    `xml:"creationTime,omitempty" json:"creationTime"`
	IsCurrent bool     `xml:"-" json:"-"`
}

func NewSnapshotXMLFromDom(ds *libvirt.DomainSnapshot) *SnapshotXML {
	var err error
	var xmlData string

	if ds == nil {
		return nil
	}
	xmlData, err = ds.GetXMLDesc(0)
	if err != nil {
		return nil
	}
	obj := &SnapshotXML{}
	if err := obj.Decode(xmlData); err != nil {
		return nil
	}
	obj.IsCurrent, _ = ds.IsCurrent(0)
	return obj
}
