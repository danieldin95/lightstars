module github.com/danieldin95/lightstar

go 1.12

require (
	github.com/beevik/etree v1.1.0
	github.com/coreos/go-systemd/v22 v22.0.0
	github.com/gorilla/mux v1.7.4
	github.com/libvirt/libvirt-go v5.10.0+incompatible
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0

)

replace golang.org/x/net v0.0.0 => github.com/golang/net v0.0.0-20190812203447-cdfb69ac37fc
replace github.com/libvirt/libvirt-go v5.10.0+incompatible => github.com/danieldin95/libvirt-go v7.4.1+incompatible