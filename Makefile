#
# github.com/danieldin95/lightstar
#

SHELL := /bin/bash

.ONESHELL:
.PHONY: lightstar lightpix test

LSB = $(shell lsb_release -i -s)$(shell lsb_release -r -s)
VER = $(shell cat VERSION)

## version
MOD = github.com/danieldin95/lightstar/libstar
LDFLAGS += -X $(MOD).Commit=$(shell git rev-list -1 HEAD)
LDFLAGS += -X $(MOD).Date=$(shell date +%FT%T%z)
LDFLAGS += -X $(MOD).Version=$(VER)

## directory
SD = $(shell pwd)
BD = build
LD = lightstar-$(LSB)-$(VER)
WD = lightpix-Windows-$(VER)

# prepare
env:
	@mkdir -p $(BD)

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
rpm/lightutils: env
	./packaging/auto.sh
	rpmbuild -ba packaging/lightutils.spec
	cp -rvf ~/rpmbuild/RPMS/x86_64/lightutils-*.rpm $(BD)

rpm/lightstar:
	./packaging/auto.sh
	rpmbuild -ba packaging/lightstar.spec
	cp -rvf ~/rpmbuild/RPMS/x86_64/lightstar-*.rpm $(BD)

rpm/lightsim:
	@./packaging/auto.sh
	rpmbuild -ba packaging/lightsim.spec
	cp -rvf ~/rpmbuild/RPMS/x86_64/lightsim-*.rpm $(BD)


linux/zip: env lightstar
	@pushd $(BD)
	@rm -rf $(LD) && mkdir -p $(LD)

	@cp $(SD)/packaging/README.md $(LD)
	@mkdir -p $(LD)/etc/lightstar
	@cp -rvf $(SD)/resource/auth.json.example $(LD)/etc/lightstar
	@cp -rvf $(SD)/resource/zone.json.example $(LD)/etc/lightstar
	@cp -rvf $(SD)/resource/permission.json.example $(LD)/etc/lightstar

	@mkdir -p $(LD)/etc/sysconfig
	@echo OPTIONS="-static:dir /var/lightstar/static -crt:dir /var/lightstar/ca -conf /etc/lightstar" > $(LD)/etc/sysconfig/lightstar.cfg

	@mkdir -p $(LD)/var/lightstar
	@cp -R $(SD)/resource/ca $(LD)/var/lightstar
	@cp -R $(SD)/http/static $(LD)/var/lightstar

	@mkdir -p $(LD)/usr/bin
	@cp -rvf $(SD)/lightstar $(LD)/usr/bin

	@mkdir -p $(LD)/usr/lib/systemd/system
	@cp $(SD)/packaging/lightstar.service $(LD)/usr/lib/systemd/system

	zip -r ./$(LD).zip $(LD) > /dev/null
	popd

centos/devel:
	yum install libvirt-devel

ubuntu/devel:
	apt-get install libvirt-dev

## cross build for windows
windows/lightpix:
	GOOS=windows GOARCH=amd64 go build -mod=vendor -o lightpix.windows.x86_64.exe lightpix.go

### packaging light pix for windows
windows/zip: env
	@pushd $(BD)
	@rm -rf $(WD) && mkdir -p $(WD)

	@cp -rvf $(SD)/resource/lightpix.json.example $(WD)/lightpix.json
	@cp -rvf $(SD)/lightpix.windows.x86_64.exe $(WD)

	zip -r $(WD).zip $(WD) > /dev/null
	popd

## unit test
test:
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/libstar
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/storage
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/compute/libvirtc
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/storage/libvirts
