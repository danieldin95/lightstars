#!/bin/bash

PHY="$1"

if [ "$PHY"x == ""x  ]; then
  exit 0
fi

yum install bridge-utils -y

BR="br-${PHY}"

## Bridge Interface

cat > /etc/sysconfig/network-scripts/ifcfg-${BR} << EOF
## Generage by LightStar project
NAME="${BR}"
DEVICE="${BR}"
BOOTPROTO="static"
#IPADDR="${ADDR}"
#NETMASK="${NETMASK}"
#GATEWAY="${GATEWAY}"
#DNS1="${DNS1}"
ONBOOT="yes"
TYPE="Bridge"
NM_CONTROLLED="no"
EOF

## Physical Interface

if ip link show $PHY > /dev/null; then
  cat > /etc/sysconfig/network-scripts/ifcfg-${PHY} << EOF
## Generage by LightStar project
NAME="${PHY}"
DEVICE="${PHY}"
TYPE="Ethernet"
BOOTPROTO="none"
ONBOOT="yes"
NM_CONTROLLED="no"
BRIDGE="${BR}"
EOF
fi

## Apply Network
ifdown ${BR}; ifup ${BR}
if ip link show $PHY > /dev/null; then
  ifdown ${PHY}; ifup ${PHY}
fi

