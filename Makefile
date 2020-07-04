#
# github.com/danieldin95/lightstar
#

SHELL := /bin/bash
.ONESHELL:


LSB = $(shell lsb_release -i -s)$(shell lsb_release -r -s)
VER = $(shell cat VERSION)

## version
MOD = github.com/danieldin95/lightstar/src/libstar
LDFLAGS += -X $(MOD).Commit=$(shell git rev-list -1 HEAD)
LDFLAGS += -X $(MOD).Date=$(shell date +%FT%T%z)
LDFLAGS += -X $(MOD).Version=$(VER)

## directory
SD = $(shell pwd)
BD = $(SD)/build
LD = lightstar-$(LSB)-$(VER)
WD = lightpix-Windows-$(VER)

help: ## show make targets
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);\
		printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## all light software
all: lightstar lightpix windows-lightpix ## build all binary

pkg: rpm windows-zip ## build all packages

rpm: rpm-lightstar rpm-lightsim ## build rpm packages

# prepare environment
env:
	@mkdir -p $(BD)

## light star
.PHONY: lightstar
lightstar: env
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o $(BD)/lightstar ./src/cli/lightstar

## light pix to proxy tcp
.PHONY: lightpix
lightpix: env
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o $(BD)/lightpix ./src/cli/lightpix

### linux packaging
rpm-env:
	@./packaging/spec.sh

rpm-lightstar: rpm-env
	rpmbuild -ba packaging/lightstar.spec
	cp -rf ~/rpmbuild/RPMS/x86_64/lightstar-*.rpm $(BD)

rpm-lightsim: rpm-env
	rpmbuild -ba packaging/lightsim.spec
	cp -rf ~/rpmbuild/RPMS/x86_64/lightsim-*.rpm $(BD)


linux-zip: env lightstar lightpix ## build linux zip packages
	@pushd $(BD)
	@rm -rf $(LD) && mkdir -p $(LD)

	@cp $(SD)/packaging/README.md $(LD)
	@mkdir -p $(LD)/etc/lightstar
	@cp -rvf $(SD)/packaging/resource/auth.json.example $(LD)/etc/lightstar
	@cp -rvf $(SD)/packaging/resource/zone.json.example $(LD)/etc/lightstar
	@cp -rvf $(SD)/packaging/resource/permission.json.example $(LD)/etc/lightstar

	@mkdir -p $(LD)/etc/sysconfig
	@echo OPTIONS="-static:dir /var/lightstar/static -crt:dir /var/lightstar/ca -conf /etc/lightstar" > $(LD)/etc/sysconfig/lightstar.cfg

	@mkdir -p $(LD)/var/lightstar
	@cp -R $(SD)/packaging/resource/ca $(LD)/var/lightstar
	@cp -R $(SD)/src/http/static $(LD)/var/lightstar

	@mkdir -p $(LD)/usr/bin
	@cp -rvf $(BD)/lightstar $(LD)/usr/bin
	@cp -rvf $(BD)/lightpix $(LD)/usr/bin

	@mkdir -p $(LD)/usr/lib/systemd/system
	@cp $(SD)/packaging/lightstar.service $(LD)/usr/lib/systemd/system

	zip -r ./$(LD).zip $(LD) > /dev/null
	@rm -rf $(LD)
	@popd

centos-devel:
	yum install libvirt-devel

ubuntu-devel:
	apt-get install libvirt-dev

## cross build for windows
windows-lightpix: env
	GOOS=windows GOARCH=amd64 go build -mod=vendor -o $(BD)/lightpix.windows.x86_64.exe ./src/cli/lightpix

### packaging light pix for windows
windows-zip: env ## build windows packages
	@pushd $(BD)
	@rm -rf $(WD) && mkdir -p $(WD)

	@cp -rvf $(SD)/packaging/resource/lightpix.json.example $(WD)/lightpix.json
	@cp -rvf $(BD)/lightpix.windows.x86_64.exe $(WD)

	zip -r $(WD).zip $(WD) > /dev/null
	@rm -rf $(WD)
	@popd

## unit test
.PHONY: test
test: ## execute unit test
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/src/libstar
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/src/storage
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/src/compute/libvirtc
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/src/storage/libvirts
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/src/http/client

clean: ## clean cache
	rm -rvf ./build
