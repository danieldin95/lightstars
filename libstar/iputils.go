package libstar

import (
	"encoding/binary"
	"net"
)


func ParseIP4Netmask(addr, prefix string) (net.IP, net.IP){
	_, inet, err := net.ParseCIDR(addr + "/"+prefix)
	if len(inet.Mask) != 4 || err != nil {
		return nil, nil
	}
	m := inet.Mask

	return inet.IP, net.IPv4(m[0], m[1], m[2], m[3])
}

func IP4Network2RangeN(addr, mask net.IP, n uint32) (net.IP, net.IP){
	intAddr := binary.BigEndian.Uint32(addr.To4())
	intMask := binary.BigEndian.Uint32(mask.To4())

	intHost := ^intMask

	start := make([]byte, 4)
	binary.BigEndian.PutUint32(start, intAddr+n)
	end := make([]byte, 4)
	binary.BigEndian.PutUint32(end, intAddr+intHost-1)

	return net.IPv4(start[0], start[1], start[2], start[3]),
		net.IPv4(end[0], end[1], end[2], end[3])
}

func IP4Network2Range(addr, mask net.IP) (net.IP, net.IP){
	intAddr := binary.BigEndian.Uint32(addr.To4())
	intMask := binary.BigEndian.Uint32(mask.To4())

	intHost := ^intMask

	start := make([]byte, 4)
	binary.BigEndian.PutUint32(start, intAddr+1)
	end := make([]byte, 4)
	binary.BigEndian.PutUint32(end, intAddr+intHost-1)

	return net.IPv4(start[0], start[1], start[2], start[3]),
		net.IPv4(end[0], end[1], end[2], end[3])
}