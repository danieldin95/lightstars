
# VM on VM

## Enable nested on host 

    cat /sys/module/kvm_amd/parameters/nested
    
    cat  > /etc/modprobe.d/kvm.conf <<EOF   
    # Intel
    options kvm_intel nested=1
    #
    # AMD
    options kvm_amd nested=1
    
    EOF
    
    modprobe -r kvm_intel
    modprobe -a kvm_intel
    
    ## 
    cat /sys/module/kvm_intel/parameters/nested
    

## Configure host-mode on guest

    virt edit xxx
    
    <cpu mode="host-passthrough" check="full"/>

