module github.com/danieldin95/lightstar

go 1.12

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/libvirt/libvirt-go v5.10.0+incompatible
	github.com/pkg/errors v0.9.1 // indirect
	github.com/quadrifoglio/go-qemu v0.0.0-20170212183343-c95730abf426
	github.com/shirou/gopsutil v2.20.1+incompatible
	github.com/stretchr/testify v1.5.1
	github.com/zchee/go-qcow2 v0.0.0-20170102190316-9a991fd172f0 // indirect
	golang.org/x/net v0.0.0
	golang.org/x/sys v0.0.0 // indirect
)

replace (
	golang.org/x/net v0.0.0 => github.com/golang/net v0.0.0-20190812203447-cdfb69ac37fc
	golang.org/x/sys v0.0.0 => github.com/golang/sys v0.0.0-20200219091948-cb0a6d8edb6c
)
