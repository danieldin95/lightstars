module github.com/danieldin95/lightstar

go 1.12

require (
	github.com/coreos/go-systemd/v22 v22.0.0
	github.com/gorilla/mux v1.7.4
	github.com/libvirt/libvirt-go v5.10.0+incompatible
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0
)

replace golang.org/x/net v0.0.0 => github.com/golang/net v0.0.0-20190812203447-cdfb69ac37fc
