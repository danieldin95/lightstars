package libstar

import (
	"fmt"
	"path"
	"testing"
)

func TestInstanceXML(t *testing.T) {
	if files, err := DIR.ListFiles("/lightstar/datastore/01", ".iso"); err == nil {
		fmt.Println(files)
	}
	if files, err := DIR.ListFiles("/lightstar/datastore/01/iso", ".iso"); err == nil {
		fmt.Println(files)
	}

	file := "/lightstar/datastore/09b191af-b82a-4736-b492-d43224bb5379/centos.cc/disk0.qcow2"
	volume := path.Base(file)
	dir := path.Dir(file)
	pool := path.Base(dir)

	fmt.Println(volume)
	fmt.Println(pool)
}
