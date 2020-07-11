# LightStar
[![Build Status](https://travis-ci.org/danieldin95/lightstar.svg?branch=master)](https://travis-ci.org/danieldin95/lightstar)
[![Go Report Card](https://goreportcard.com/badge/github.com/danieldin95/lightstar)](https://goreportcard.com/report/danieldin95/lightstar)
[![GPL 3.0 License](https://img.shields.io/badge/License-GPL%203.0-blue.svg)](LICENSE)

This software makes it easier for you to control virtual compute, network and storage resource.

# CentOS7

## Check Intel VT-x or AMD-V

    lscpu | egrep '(vmx|svm)'

## Disable SElinux firstly

    cat > /etc/sysconfig/selinux <<EOF
    SELINUX=disabled
    SELINUXTYPE=targeted
    EOF
    
    reboot

## Install epel and qemu-ev release

    yum install -y epel centos-release-qemu-ev

## Install by RPM packaging

    yum install -y wget
    wget https://github.com/danieldin95/lightstar/releases/download/v0.2.23/lightstar-0.2.23-1.el7.x86_64.rpm
    
    yum install -y ./lightstar-0.2.23-1.el7.x86_64.rpm


## Start Libvirtd service.

    systemctl enable libvirtd
    systemctl start libvirtd
    

## Start LightStar service.


    systemctl enable lightstar
    systemctl start lightstar


## Upload a Linux ISO file

    cd /lightstar/datastore/01
    wget http://mirrors.aliyun.com/centos/7.7.1908/isos/x86_64/CentOS-7-x86_64-Minimal-1908.iso


# Open UI on browser

    https://your-machine-address:10080


## Get admin's password

    cat /etc/lightstar.auth


# Welcome to attention

WeChat: DanielDin

Email: danieldin95@163.com
