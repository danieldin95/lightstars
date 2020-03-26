.PHONY: lightstar test

PKG = github.com/danieldin95/lightstar/libstar

LDFLAGS += -X $(PKG).Commit=$$(git rev-list -1 HEAD)
LDFLAGS += -X $(PKG).Date=$$(date +%FT%T%z)
LDFLAGS += -X $(PKG).Version=$$(cat VERSION)


lightstar:
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o lightstar lightstar.go
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o lightpix lightpix.go


linux/rpm: linux/rpm/utils linux/rpm/star linux/rpm/sim


linux/rpm/utils:
	./packaging/auto.sh
	rpmbuild -ba packaging/lightutils.spec
	cp -rvf ~/rpmbuild/RPMS/x86_64/lightutils-*.rpm .


linux/rpm/star:
	./packaging/auto.sh
	rpmbuild -ba packaging/lightstar.spec
	cp -rvf ~/rpmbuild/RPMS/x86_64/lightstar-*.rpm .


linux/rpm/sim:
	./packaging/auto.sh
	rpmbuild -ba packaging/lightsim.spec
	cp -rvf ~/rpmbuild/RPMS/x86_64/lightsim-*.rpm .


devel/requirements:
	yum install libvirt-devel


windows:
	go build -mod=vendor -o lightpix.windows.x86_64.exe lightpix.go


WIN_DIR = "lightpix-windows-"$$(cat VERSION)


windows/zip:
	rm -rf $(WIN_DIR) && mkdir -p $(WIN_DIR)
	cp -rvf resource/lightpix.json.example $(WIN_DIR)/lightpix.json
	cp -rvf lightpix.windows.x86_64.exe $(WIN_DIR)
	rm -rf $(WIN_DIR).zip
	zip -r $(WIN_DIR).zip $(WIN_DIR)


test:
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/libstar
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/storage
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/compute/libvirtc
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/storage/libvirts
