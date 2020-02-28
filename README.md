# LightStar
[![Build Status](https://travis-ci.org/danieldin95/lightstar.svg?branch=master)](https://travis-ci.org/danieldin95/lightstar)
[![Go Report Card](https://goreportcard.com/badge/github.com/danieldin95/lightstar)](https://goreportcard.com/report/danieldin95/lightstar)
[![GPL 3.0 License](https://img.shields.io/badge/License-GPL%203.0-blue.svg)](LICENSE)

This software makes it easier for you to control compute resource.

# CentOS7

## Install by RPM packaging

    yum install -y wget bridge-utils
    yum install -y lightstar-x.x.x.rpm

## Start Libvirtd and LightStar service.

    systemctl enable libvirtd
    systemctl start libvirtd
    
    systemctl enable lightstar
    systemctl start lightstar
    
## Configure default datastore

    mkdir -p /lighstar/datastore/01
    virsh pool-create-as --name 01 --type dir --target /lightstar/datastore/01
    
## Configure default network

    cat > virbr0.xml <<EOF
        <network>
          <name>virbr0</name>
          <forward mode='nat'>
            <nat>
              <port start='1024' end='65535'/>
            </nat>
          </forward>
          <bridge name='virbr0' stp='on' delay='0'/>
          <ip address='172.16.10.1' netmask='255.255.255.0'>
            <dhcp>
              <range start='172.16.10.10' end='172.16.10.100'/>
            </dhcp>
          </ip>
        </network>
    EOF
    virsh net-create virbr0.xml

## Upload a linux ISO file

    cd /lightstar/datastore/01
    wget http://mirrors.aliyun.com/centos/7.7.1908/isos/x86_64/CentOS-7-x86_64-Minimal-1908.iso
    
# Open UI on browser

    https://your-machine-address:10080

