#
# github.com/danieldin95/lightstar
#

.PHONY: lightstar lightpix test

## version
PKG = github.com/danieldin95/lightstar/libstar
LDFLAGS += -X $(PKG).Commit=$$(git rev-list -1 HEAD)
LDFLAGS += -X $(PKG).Date=$$(date +%FT%T%z)
LDFLAGS += -X $(PKG).Version=$$(cat VERSION)

## all light software
all: lightstar lightpix windows/lightpix

pkg: rpm windows/zip

rpm: rpm/lightutils rpm/lightstar rpm/lightsim

## light star
lightstar:
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o lightstar lightstar.go

## light pix to proxy tcp
lightpix:
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o lightpix lightpix.go

### linux packaging
rpm/lightutils:
	./packaging/auto.sh
	rpmbuild -ba packaging/lightutils.spec
	cp -rvf ~/rpmbuild/RPMS/x86_64/lightutils-*.rpm .

rpm/lightstar:
	./packaging/auto.sh
	rpmbuild -ba packaging/lightstar.spec
	cp -rvf ~/rpmbuild/RPMS/x86_64/lightstar-*.rpm .

rpm/lightsim:
	@./packaging/auto.sh
	rpmbuild -ba packaging/lightsim.spec
	cp -rvf ~/rpmbuild/RPMS/x86_64/lightsim-*.rpm .

LIN_DIR=lightstar-$$(lsb_release -i -s)$$(lsb_release -r -s)-$$(cat VERSION)

linux/zip: lightstar
	@rm -rf $(LIN_DIR) && mkdir -p $(LIN_DIR)
	@cp ./packaging/README.md $(LIN_DIR)

	@mkdir -p $(LIN_DIR)/etc/lightstar
	@cp -rvf resource/auth.json.example $(LIN_DIR)/etc/lightstar
	@cp -rvf resource/zone.json.example $(LIN_DIR)/etc/lightstar
	@cp -rvf resource/permission.json.example $(LIN_DIR)/etc/lightstar

	@mkdir -p $(LIN_DIR)/etc/sysconfig
	@echo OPTIONS="-static:dir /var/lightstar/static -crt:dir /var/lightstar/ca -conf /etc/lightstar" > $(LIN_DIR)/etc/sysconfig/lightstar.cfg

	@mkdir -p $(LIN_DIR)/var/lightstar
	@cp -R ./resource/ca $(LIN_DIR)/var/lightstar
	@cp -R ./http/static $(LIN_DIR)/var/lightstar

	@mkdir -p $(LIN_DIR)/usr/bin
	@cp -rvf lightstar $(LIN_DIR)/usr/bin

	@mkdir -p $(LIN_DIR)/usr/lib/systemd/system
	@cp ./packaging/lightstar.service $(LIN_DIR)/usr/lib/systemd/system

	zip -r $(LIN_DIR).zip $(LIN_DIR)

centos/devel:
	yum install libvirt-devel

ubuntu/devel:
	apt-get install libvirt-dev

## cross build for windows
windows/lightpix:
	GOOS=windows GOARCH=amd64 go build -mod=vendor -o lightpix.windows.x86_64.exe lightpix.go

### packaging light pix for windows
WIN_DIR = lightpix-windows-$$(cat VERSION)

windows/zip:
	@rm -rf $(WIN_DIR) && mkdir -p $(WIN_DIR)
	@cp -rvf resource/lightpix.json.example $(WIN_DIR)/lightpix.json
	@cp -rvf lightpix.windows.x86_64.exe $(WIN_DIR)

	zip -r $(WIN_DIR).zip $(WIN_DIR)

## unit test
test:
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/libstar
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/storage
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/compute/libvirtc
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/storage/libvirts
