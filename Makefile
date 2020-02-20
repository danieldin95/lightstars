.PHONY: lightstar test

PKG = github.com/danieldin95/lightstar/libstar

LDFLAGS += -X $(PKG).Commit=$$(git rev-list -1 HEAD)
LDFLAGS += -X $(PKG).Date=$$(date +%FT%T%z)
LDFLAGS += -X $(PKG).Version=$$(cat VERSION)


lightstar:
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o lightstar main.go


rpm:
	./packaging/auto.sh
	rpmbuild -ba packaging/lightstar.spec
	cp -rvf ~/rpmbuild/RPMS/x86_64/lightstar-*.rpm ./packaging


devel/requirements:
	yum install libvirt-devel


test:
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/compute/libvirt
