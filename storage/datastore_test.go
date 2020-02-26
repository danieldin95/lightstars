package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFmtPath(t *testing.T) {
	assert.Equal(t, "datastore@01", PATH.Fmt("/lightstar/datastore/01"), "")
	assert.Equal(t, "datastore@01:/xx", PATH.Fmt("/lightstar/datastore/01/xx"), "")
	assert.Equal(t, "datastore@01:/xx/bb", PATH.Fmt("/lightstar/datastore/01/xx/bb"), "")

	assert.Equal(t, "/lightstar/datastore/01", PATH.Unix("datastore@01"), "")
	assert.Equal(t, "/lightstar/datastore/01/xx", PATH.Unix("datastore@01:/xx"), "")
	assert.Equal(t, "/lightstar/datastore/01/xx/bb", PATH.Unix("datastore@01:/xx/bb"), "")

	fmt.Println(PATH.Fmt("/lightstar/datastore/01/cn_windows_7_ultimate_with_sp1_x64_dvd_u_677408.iso"))
	fmt.Println(PATH.Unix("datastore@01:/cn_windows_7_ultimate_with_sp1_x64_dvd_u_677408.iso"))

	assert.Equal(t, "device:/01", PATH.Fmt("/dev/01"), "")
	assert.Equal(t, "device:/01/xx", PATH.Fmt("/dev/01/xx"), "")
	assert.Equal(t, "/dev", PATH.Fmt("/dev"), "")
	assert.Equal(t, "device:/xx", PATH.Fmt("/dev/xx"), "")
	assert.Equal(t, "device:/", PATH.Fmt("/dev/"), "")
}
