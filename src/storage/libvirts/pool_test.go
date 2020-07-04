package libvirts

import (
	"fmt"
	"testing"
)

func TestPool(t *testing.T) {
	_, err := CreatePool(".test.xx", "/lightstar/datastore/test.xx")
	fmt.Println(err)
	_, err = CreateVolume(".test.xx", "disk0.qcow2", 1024*1024)
	fmt.Println(err)

	pol := &Pool{Name: ".text.xx"}
	err = pol.Remove()
	fmt.Println(err)
}
