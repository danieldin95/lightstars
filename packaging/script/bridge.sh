#!/bin/bash

PHY=em1
BR="br-${PHY}"
ADDR=192.168.40.100
NETMASK=255.255.255.0
GATEWAY=192.168.40.1
DNS1=192.168.40.1

## Bridge Interface

cat > /etc/sysconfig/network-scripts/ifcfg-${BR} << EOF
DEVICE=${BR}
BOOTPROTO=static
IPADDR=${ADDR}
NETMASK=${NETMASK}
GATEWAY=${GATEWAY}
DNS1=${DNS1}
ONBOOT=yes
TYPE=Bridge
NM_CONTROLLED=no
EOF

## Physical Interface

cat > /etc/sysconfig/network-scripts/ifcfg-${PHY} << EOF
DEVICE=${PHY}
TYPE=Ethernet
BOOTPROTO=none
ONBOOT=yes
NM_CONTROLLED=no
BRIDGE=${BR}
EOF

ifdown ${BR}; ifdown ${PHY}; ifup ${BR}; ifup ${PHY}; 
