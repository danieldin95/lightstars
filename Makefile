#
# github.com/danieldin95/lightstar
#

#
# git clone https://github.com/danieldin95/freecert packaging/resource/cert
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

help: ## show make targets
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);\
		printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## all light software
bin: lightstar ## build all binary

pkg: rpm ## build all packages

rpm: rpm-lightstar ## build rpm packages

# prepare environment
env:
	@mkdir -p $(BD)
	gofmt -w -s ./src

## light star
.PHONY: lightstar
lightstar: env
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o $(BD)/lightstar ./src/cli/lightstar

### linux packaging
rpm-env:
	@./packaging/spec.sh
	@[ -e "$(BD)"/cert ] || ln -s $(SD)/../freecert $(BD)/cert

rpm-lightstar: rpm-env
	rpmbuild -ba ./build/lightstar.spec
	cp -rf ~/rpmbuild/RPMS/x86_64/lightstar-*.rpm $(BD)


linux-zip: env lightstar ## build linux zip packages
	@pushd $(BD)
	@rm -rf $(LD) && mkdir -p $(LD)

	@cp $(SD)/packaging/README.md $(LD)
	@mkdir -p $(LD)/etc/lightstar
	@cp -rvf $(SD)/packaging/resource/auth.json.example $(LD)/etc/lightstar
	@cp -rvf $(SD)/packaging/resource/zone.json.example $(LD)/etc/lightstar
	@cp -rvf $(SD)/packaging/resource/permission.json.example $(LD)/etc/lightstar

	@mkdir -p $(LD)/etc/sysconfig
	@echo OPTIONS="-static:dir /var/lightstar/static -crt:dir /var/lightstar/cert -conf /etc/lightstar" > $(LD)/etc/sysconfig/lightstar.cfg

	@mkdir -p $(LD)/var/lightstar
	@cp -R $(SD)/packaging/resource/ca $(LD)/var/lightstar
	@cp -R $(SD)/src/http/static $(LD)/var/lightstar

	@mkdir -p $(LD)/usr/bin
	@cp -rvf $(BD)/lightstar $(LD)/usr/bin

	@mkdir -p $(LD)/usr/lib/systemd/system
	@cp $(SD)/packaging/lightstar.service $(LD)/usr/lib/systemd/system

	zip -r ./$(LD).zip $(LD) > /dev/null
	@rm -rf $(LD)
	@popd

centos-devel:
	yum install libvirt-devel

ubuntu-devel:
	apt-get install libvirt-dev

## upgrade
upgrade:
	ansible-playbook ./misc/playbook/upgrade.yaml -e "version=$(VER)"

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
