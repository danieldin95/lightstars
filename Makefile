#
# github.com/danieldin95/lightstar
#

SHELL := /bin/bash
.ONESHELL:


LSB = $(shell lsb_release -i -s | tr '[:upper:]' '[:lower:]')
VER = $(shell cat VERSION)

## version
MOD = github.com/danieldin95/lightstar/pkg/libstar
LDFLAGS += -X $(MOD).Commit=$(shell git rev-list -1 HEAD)
LDFLAGS += -X $(MOD).Date=$(shell date +%FT%T%z)
LDFLAGS += -X $(MOD).Version=$(VER)

## directory
SRC_DIR = $(shell pwd)
BUR_DIR = $(SRC_DIR)/build
LIN_DIR = lightstar-$(LSB)-$(VER)

## all light star software
bin: tar ## build all binary
	@cat $(SRC_DIR)/dist/install.sh > $(BUR_DIR)/$(LIN_DIR).bin && \
	echo "__ARCHIVE_BELOW__:" >> $(BUR_DIR)/$(LIN_DIR).bin && \
	cat $(BUR_DIR)/$(LIN_DIR).tar.gz >> $(BUR_DIR)/$(LIN_DIR).bin && \
	chmod +x $(BUR_DIR)/$(LIN_DIR).bin && \
	echo "Save to $(LIN_DIR).bin"

help: ## show make targets
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);\
		printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# prepare environment
env:
	@mkdir -p $(BUR_DIR)
	gofmt -w -s ./pkg ./cmd

## light star
.PHONY: lightstar

lightstar: env
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o $(BUR_DIR)/lightstar ./cmd/lightstar


tar: env lightstar ## build tar packages
	@pushd $(BUR_DIR)
	@rm -rf $(LIN_DIR) && mkdir -p $(LIN_DIR)

	@cp $(SRC_DIR)/dist/README.md $(LIN_DIR)
	@mkdir -p $(LIN_DIR)/etc/lightstar
	@cp -rvf $(SRC_DIR)/dist/resource/auth.json.example $(LIN_DIR)/etc/lightstar
	@cp -rvf $(SRC_DIR)/dist/resource/zone.json.example $(LIN_DIR)/etc/lightstar
	@cp -rvf $(SRC_DIR)/dist/resource/permission.json.example $(LIN_DIR)/etc/lightstar

	@mkdir -p $(LIN_DIR)/etc/sysconfig
	@echo OPTIONS="-static:dir /var/lightstar/static -crt:dir /var/lightstar/cert -conf /etc/lightstar" > $(LIN_DIR)/etc/sysconfig/lightstar.cfg

	@mkdir -p $(LIN_DIR)/var/lightstar
	@cp -R $(SRC_DIR)/dist/resource/cert/lightstar/ca $(LIN_DIR)/var/lightstar
	@cp -R $(SRC_DIR)/pkg/http/static $(LIN_DIR)/var/lightstar

	@mkdir -p $(LIN_DIR)/usr/bin
	@cp -rvf $(BUR_DIR)/lightstar $(LIN_DIR)/usr/bin

	@mkdir -p $(LIN_DIR)/usr/lib/systemd/system
	@cp $(SRC_DIR)/dist/resource/lightstar.service $(LIN_DIR)/usr/lib/systemd/system

	tar -czf ./$(LIN_DIR).tar.gz $(LIN_DIR) > /dev/null
	@rm -rf $(LIN_DIR)
	@popd


## upgrade
upgrade:
	ansible-playbook ./misc/playbook/upgrade.yaml -e "version=$(VER)"

## unit test
.PHONY: test
test: ## execute unit test
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/pkg/libstar
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/pkg/storage
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/pkg/compute
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/pkg/storage
	go test -v -mod=vendor -bench=. github.com/danieldin95/lightstar/pkg/http/client

clean: ## clean cache
	rm -rvf ./build
