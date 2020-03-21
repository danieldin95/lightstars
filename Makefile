.PHONY: lightstar test

PKG = github.com/danieldin95/lightstar/libstar

LDFLAGS += -X $(PKG).Commit=$$(git rev-list -1 HEAD)
LDFLAGS += -X $(PKG).Date=$$(date +%FT%T%z)
LDFLAGS += -X $(PKG).Version=$$(cat VERSION)


lightstar:
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o lightstar lightstar.go
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o lightprix lightprix.go


rpm:
	./packaging/auto.sh
	rpmbuild -ba packaging/lightsim.spec
	rpmbuild -ba packaging/lightstar.spec
	cp -rvf ~/rpmbuild/RPMS/x86_64/lightsim-*.rpm .
	cp -rvf ~/rpmbuild/RPMS/x86_64/lightstar-*.rpm .


devel/requirements:
	yum install libvirt-devel


windows:
	go build -mod=vendor -o lightprix.windows.x86_64.exe lightprix.go


test:
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/libstar
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/storage
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/compute/libvirtc
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/storage/libvirts
