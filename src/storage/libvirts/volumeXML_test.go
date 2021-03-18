package libvirts

import (
	"fmt"
	"github.com/danieldin95/lightstar/src/libstar"
	"testing"
)

func TestVolumeXML(t *testing.T) {
	xmlData := `
<volume type='file'>
  <name>disk0.qcow2</name>
  <key>/lightstar/datastore/01/centos.30/disk0.qcow2</key>
  <source>
  </source>
  <capacity unit='bytes'>10737418240</capacity>
  <allocation unit='bytes'>200704</allocation>
  <physical unit='bytes'>197120</physical>
  <target>
    <path>/lightstar/datastore/01/centos.30/disk0.qcow2</path>
    <format type='qcow2'/>
    <permissions>
      <mode>0600</mode>
      <owner>0</owner>
      <group>0</group>
    </permissions>
    <timestamps>
      <atime>1582379803.208613720</atime>
      <mtime>1582379776.103327524</mtime>
      <ctime>1582379776.104327534</ctime>
    </timestamps>
  </target>
</volume>
`
	volXml := &VolumeXML{}
	_ = libstar.XML.Decode(volXml, xmlData)

	fmt.Println(volXml)
}

func TestVolume_Create(t *testing.T) {
	//pol, err := CreatePool(ToDomainPool("centos.test.xx"), "/lightstar/datastore/01/centos.test.xx")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//size := libstar.ToBytes("12", "GiB")
	//vol, err := CreateVolume(pol.Name, "disk0.qcow2", size)
	//fmt.Println(vol)
}
