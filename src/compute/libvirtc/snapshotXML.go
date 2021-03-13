package libvirtc

import (
	"encoding/xml"
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/libvirt/libvirt-go"
)

type SnapshotXML struct {
	XMLName  xml.Name `xml:"domainsnapshot" json:"-"`
	Name     string   `xml:"name,omitempty" json:"name"`
	State    string   `xml:"state,omitempty" json:"state"`
	CreateAt int64    `xml:"creationTime,omitempty" json:"creationTime"`
}

func NewSnapshotXMLFromDom(dom *libvirt.DomainSnapshot) *SnapshotXML {
	var err error
	var xmlData string

	if dom == nil {
		return nil
	}
	xmlData, err = dom.GetXMLDesc(0)
	if err != nil {
		return nil
	}
	obj := &SnapshotXML{}
	if err := obj.Decode(xmlData); err != nil {
		return nil
	}
	return obj
}

func (dom *SnapshotXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), dom); err != nil {
		libstar.Error("SnapshotXML.Decode %s", err)
		return err
	}
	return nil
}

func (dom *SnapshotXML) Encode() string {
	data, err := xml.Marshal(dom)
	if err != nil {
		libstar.Error("SnapshotXML.Encode %s", err)
		return ""
	}
	return string(data)
}
