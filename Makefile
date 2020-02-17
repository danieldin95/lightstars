.PHONY: lightstar test

PKG = main

LDFLAGS += -X $(PKG).Commit=$$(git rev-list -1 HEAD)
LDFLAGS += -X $(PKG).Date=$$(date +%FT%T%z)
LDFLAGS += -X $(PKG).Version=$$(cat VERSION)


lightstar:
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o resource/lightstar main.go


devel/requirements:
	yum install libvirt-devel

test:
	@echo "TODO"
