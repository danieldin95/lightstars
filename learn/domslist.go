package main

import (
	"flag"
	"fmt"
	"github.com/libvirt/libvirt-go"
	"time"
)

func main() {
	hypervisor := "qemu:///system"

	flag.StringVar(&hypervisor, "hyper", hypervisor, "hypervisor connecting to.")
	flag.Parse()

	fmt.Printf("%d\n", time.Now().Unix())
	conn, err := libvirt.NewConnect(hypervisor)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%d running domains:\n", len(doms))
	for _, dom := range doms {
		if name, err := dom.GetName(); err == nil {
			fmt.Printf("  %s\n", name)
		}
		if os, err := dom.GetOSType(); err == nil {
			fmt.Printf("  %s\n", os)
		}
		if xml, err := dom.GetInfo(); err == nil {
			fmt.Printf("  %s", xml)
		}
		dom.Free()
	}
}
