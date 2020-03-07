# NFS on linux

## Open port on firewall

    iptables -I INPUT -p tcp --dport 110 -j ACCEPT
    iptables -I INPUT -p udp --dport 110 -j ACCEPT
    iptables -I INPUT -p tcp --dport 2049 -j ACCEPT
    iptables -I INPUT -p tcp --dport 662 -j ACCEPT
    iptables -I INPUT -p udp --dport 662 -j ACCEPT
    iptables -I INPUT -p tcp --dport 875 -j ACCEPT
    iptables -I INPUT -p udp --dport 875 -j ACCEPT
    iptables -I INPUT -p tcp --dport 892 -j ACCEPT
    iptables -I INPUT -p udp --dport 892 -j ACCEPT
    iptables -I INPUT -p tcp --dport 32803 -j ACCEPT
    iptables -I INPUT -p udp --dport 32769 -j ACCEPT 

## Install packages

    yum -y install nfs-utils rpcbind

## Configure exportfs

    cat /etc/exportfs
    /lightstar/nfs/01 192.168.4.0/24(rw,no_root_squash,no_all_squash)
    
## Restart NFS.

    systemctl restart nfs