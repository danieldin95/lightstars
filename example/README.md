# KVM

Linux KVM/Qemu

## CentOS7

The following that start a new CentOS7 instance.

### Datastore

    mkdir -p /lightstar/datastore && /lightstar/datastore
    mkdir -p 09b191af-b82a-4736-b492-d43224bb5379
    ln -s 09b191af-b82a-4736-b492-d43224bb5379 0 && cd 0
   
### Create Image

    mkdir -p centos7 && cd centos7
    qemu-img create -f qcow2 disk0.qcow2 10G
    
### Add New Network

    yum install bridge-utils
    brctl addbr virbr0
    brctl addbr virbr1
    brctl addbr virbr2

### Start Instance 

    virsh create centos7.xml
    virsh start centos7
    
    virsh list --all
    virsh vncdisplay panabit-01
    
## KVM UEFI

    virsh create kvm-uefi.xml

## KVM BIOS

    virsh create kvm-biso.xml
   
## Panabit
 
    qemu-img create -f qcow2 /var/lib/libvirt/images/panabit.disk-0.img 10G
    virsh create panabit.xml
    virsh list --all
    virsh vncdisplay panabit-01

