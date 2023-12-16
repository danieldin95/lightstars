#!/bin/bash

set -ex

BR=$1

if [ "$BR"x == ""x ]; then
	exit 0
fi

cat > /etc/sysconfig/network-scripts/ifcfg-${BR} << EOF
## Generage by LightStar project
NAME="${BR}"
DEVICE="${BR}"
ONBOOT="yes"
TYPE="Bridge"
NM_CONTROLLED="no"
EOF

ifup $BR
