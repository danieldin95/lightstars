package compute

import (
	"encoding/xml"
	"github.com/danieldin95/lightstar/pkg/libstar"
)

type SnapshotXML struct {
	XMLName   xml.Name `xml:"domainsnapshot" json:"-"`
	Name      string   `xml:"name,omitempty" json:"name"`
	State     string   `xml:"state,omitempty" json:"state"`
	CreateAt  int64    `xml:"creationTime,omitempty" json:"creationTime"`
	IsCurrent bool     `xml:"-" json:"-"`
}

func NewSnapshotXMLFromDom(ds *DomainSnapshot) *SnapshotXML {
	if ds == nil {
		return nil
	}
	xmlData, err := ds.GetXMLDesc(0)
	if err != nil {
		return nil
	}
	obj := &SnapshotXML{}
	if err := libstar.XML.Decode(obj, xmlData); err != nil {
		return nil
	}
	obj.IsCurrent, _ = ds.IsCurrent(0)
	return obj
}
