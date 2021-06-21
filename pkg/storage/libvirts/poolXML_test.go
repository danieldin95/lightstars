package libvirts

import (
	"fmt"
	"github.com/danieldin95/lightstar/pkg/libstar"
	"testing"
)

func TestPoolXML(t *testing.T) {
	xmlData := `
<pool type='dir'>
  <name>.centos.32</name>
  <uuid>b6e10dde-4280-4aaf-9ea3-99c8803f07b0</uuid>
  <capacity unit='bytes'>82650857472</capacity>
  <allocation unit='bytes'>70308085760</allocation>
  <available unit='bytes'>12342771712</available>
  <source>
  </source>
  <target>
    <path>/lightstar/datastore/01/centos.32</path>
    <permissions>
      <mode>0755</mode>
      <owner>0</owner>
      <group>0</group>
    </permissions>
  </target>
</pool>
`
	polXml := &PoolXML{}
	_ = libstar.XML.Decode(polXml, xmlData)

	fmt.Println(polXml)
}
