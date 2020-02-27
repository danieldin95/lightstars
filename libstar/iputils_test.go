package libstar

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestIpUtils24(t *testing.T) {
	addr := net.ParseIP("192.168.1.0")
	mask := net.ParseIP("255.255.255.0")

	start, end := IP4Network2Range(addr, mask)
	assert.Equal(t, "192.168.1.1", start.String(), "")
	assert.Equal(t, "192.168.1.254", end.String(), "")
}

func TestIpUtils25(t *testing.T) {

	addr := net.ParseIP("192.168.1.0")
	mask := net.ParseIP("255.255.255.128")

	start, end := IP4Network2Range(addr, mask)
	assert.Equal(t, "192.168.1.1", start.String(), "")
	assert.Equal(t, "192.168.1.126", end.String(), "")
}