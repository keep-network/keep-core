# Change VERSION below before commit to create new release
VERSION := $(shell scripts/version.sh;)
REVISION := $(shell scripts/revision.sh;)
DEP_INSTALLED := $(shell command -v dep;)


.PHONY: bootstrap
bootstrap:
ifndef DEP_INSTALLED
	@go get -u github.com/golang/dep/cmd/dep
endif
	@dep ensure -v -vendor-only


.PHONY: build
build: build_mac # build_linux

build_mac: export GOARCH=amd64
build_mac: export CGO_ENABLED=1
build_mac:
	@GOOS=darwin go build -v --ldflags="-w -X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
		-o bin/darwin/amd64/keep main.go

#build_linux: export GOARCH=amd64
#build_linux: export CGO_ENABLED=1
#build_linux:
#	@GOOS=linux go build -v --ldflags="-w -X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
#		-o bin/linux/amd64/keep main.go


.PHONY: run
run:
	@bin/darwin/amd64/keep --debug smoketest
