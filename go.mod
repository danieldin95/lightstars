module github.com/danieldin95/lightstar

go 1.12

require (
	github.com/gorilla/mux v1.7.4
	github.com/libvirt/libvirt-go v5.10.0+incompatible
	github.com/pkg/errors v0.9.1 // indirect
	github.com/quadrifoglio/go-qemu v0.0.0-20170212183343-c95730abf426
	github.com/zchee/go-qcow2 v0.0.0-20170102190316-9a991fd172f0 // indirect
	golang.org/x/net v0.0.0
)

replace golang.org/x/net v0.0.0 => github.com/golang/net v0.0.0-20190812203447-cdfb69ac37fc
