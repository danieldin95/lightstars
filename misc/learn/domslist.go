package main

import (
	"flag"
	"fmt"
	"github.com/libvirt/libvirt-go"
	"time"
)

func cpuUtils(conn *libvirt.Connect) {
	olderSts, _ := conn.GetCPUStats(-1, 0)
	fmt.Println(olderSts)
	older := olderSts.User + olderSts.Idle + olderSts.Kernel + olderSts.Intr + olderSts.Iowait
	for {
		time.Sleep(time.Second * 2)

		newerSts, _ := conn.GetCPUStats(-1, 0)
		newer := newerSts.User + newerSts.Idle + newerSts.Kernel + newerSts.Intr + newerSts.Iowait

		fmt.Println(1000 * (newerSts.User - olderSts.User) / (newer - older))
		fmt.Println(1000 * (newerSts.Kernel - olderSts.Kernel) / (newer - older))
		fmt.Println(1000 * (newerSts.Idle - olderSts.Idle) / (newer - older))
		fmt.Println(1000 * (newerSts.Idle - olderSts.Idle) / (2261 * 1000 * 1000 * 16))
		fmt.Printf("\n--%d----\n", newer-older)

		older = newer
		olderSts = newerSts
	}
}

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
	fmt.Println(conn.GetNodeInfo())
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
		if name, err := dom.GetName(); err == nil {
			fmt.Printf("%s\n", name)
		}
		//if os, err := dom.GetOSType(); err == nil {
		//	fmt.Printf("  %s\n", os)
		//}
		//if xml, err := dom.GetXMLDesc(libvirt.DOMAIN_XML_INACTIVE); err == nil {
		//	fmt.Printf("  %s", xml)
		//}
		cfg := libvirt.DOMAIN_AFFECT_CONFIG
		if is, err := dom.IsActive(); err == nil && is {
			cfg |= libvirt.DOMAIN_AFFECT_LIVE
		}
		if err := dom.SetMetadata(libvirt.DOMAIN_METADATA_TITLE, "hi", "", "", cfg); err != nil {
			fmt.Printf("Metadata error  %s\n", err)
		}
		desc, err := dom.GetMetadata(libvirt.DOMAIN_METADATA_TITLE, "", cfg)
		fmt.Printf("Metadata desc %s %s\n", desc, err)
		dom.Free()
	}
	cpuUtils(conn)
}
