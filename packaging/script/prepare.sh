#!/bin/bash

set -e

systemctl disable firewalld
systemctl stop firewalld

systemctl disable NetworkManager
systemctl stop NetworkManager

systemctl enable network
systemctl start network

cat > /etc/sysconfig/selinux <<EOF

SELINUX=disabled
SELINUXTYPE=targeted

EOF
