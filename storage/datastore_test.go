package storage

import (
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
}
