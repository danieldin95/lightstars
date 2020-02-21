package libstar

import (
	"fmt"
	"testing"
)

func TestInstanceXML(t *testing.T) {
	if files, err := DIR.ListFiles("/lightstar/datastore/01", ".iso"); err == nil {
		fmt.Println(files)
	}
	if files, err := DIR.ListFiles("/lightstar/datastore/01/iso", ".iso"); err == nil {
		fmt.Println(files)
	}
}
