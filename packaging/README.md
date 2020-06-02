# Install Lightstar

## Copy files to your rootfs
    cp -rvf etc var usr /

## Check libvirt daemon already installed
    virsh list

### Install Dependence on CentOS
    yum install libvirt-daemon libvirt qemu-kvm qemu-img

### Install Dependence on Ubuntu
    apt-get install libvirt-bin qemu-kvm qemu-img

## Enable lightstar service
    systemctl enable lightstar
