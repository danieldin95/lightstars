package main

import (
	"flag"
	"fmt"
	"github.com/libvirt/libvirt-go"
)

func main() {
	hypervisor := "qemu:///system"

	flag.StringVar(&hypervisor, "hyper", hypervisor, "hypervisor connecting to.")
	flag.Parse()

	//fmt.Printf("%d\n", time.Now().Unix())
	conn, err := libvirt.NewConnect(hypervisor)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	//fmt.Println("----:")
	//fmt.Println(conn.GetMaxVcpus(""))
	//fmt.Println("----:")
	//fmt.Println(conn.GetCapabilities())
	fmt.Println("----:")
	fmt.Println(conn.GetFreeMemory())
	fmt.Println("----:")
	fmt.Println(conn.GetMemoryStats(-1, 0))
	//fmt.Println("----:")
	//fmt.Println(conn.GetSysinfo(0))
	//fmt.Println("----:")
	//fmt.Println(conn.GetHostname())

	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%d running domains:\n", len(doms))
	for _, dom := range doms {
		//if name, err := dom.GetName(); err == nil {
		//	fmt.Printf("  %s\n", name)
		//}
		//if os, err := dom.GetOSType(); err == nil {
		//	fmt.Printf("  %s\n", os)
		//}
		if xml, err := dom.GetXMLDesc(libvirt.DOMAIN_XML_INACTIVE); err == nil {
			fmt.Printf("  %s", xml)
		}
		dom.Free()
	}
}
