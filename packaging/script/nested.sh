#!/usr/bin/env bash

# Generate configuration
cat > /etc/modprobe.d/kvm.conf <<EOF
# Intel
options kvm_intel nested=1
# AMD
options kvm_amd nested=1
EOF

# Reintall kvm_intel

modprobe -r kvm_intel
modprobe -a kvm_intel

# modprobe -r kvm_amd
# modprobe -a kvm_amd
