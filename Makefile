.PHONY: all

all: generate build cmd-help 

generate:
	go generate ./...

build:
	go build -o keep-client -a . 

cmd-help: build
	@echo '$$ keep-client start --help' > docs/development/cmd-help
	./keep-client start --help >> docs/development/cmd-help

# TODO: Consider extracting `download_artifacts` step from go generate command.
