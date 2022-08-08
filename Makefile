# environment is used as a tag of the npm packages for contracts artifacts. If
# not overwritten it defaults to `development`.
environment = development

development:
	make all environment=development

goerli:
	make all environment=goerli

# TODO: Mainnet packages have not been published yet.
# mainnet:
# 	make all environment=mainnet

all: download_artifacts generate build cmd-help

# List of NPM packages containing contracts needed by the client contracts bindings
# generation.
npm_packages := @keep-network/random-beacon \
	@keep-network/ecdsa \
	@threshold-network/solidity-contracts \
	@keep-network/tbtc-v2

# Working directory where contracts artifacts should be stored.
contracts_dir := tmp/contracts

# It requires npm of at least 7.x version to support `pack-destination` flag.
define get_npm_package
$(eval npm_package_tag := $(1))
$(eval npm_package_name := $(2))
$(eval destination_dir := ${contracts_dir}/${npm_package_tag}/${npm_package_name})
@rm -rf ${destination_dir}
@mkdir -p ${destination_dir}
@npm pack --silent \
	--pack-destination=${destination_dir} \
	$(shell npm view ${npm_package_name}@${npm_package_tag} _id) \
	| xargs -I{} tar -zxf ${destination_dir}/{} -C ${destination_dir} --strip-components 1 package/artifacts
$(info Downloaded NPM package ${npm_package_name}@${npm_package_tag} to ${contracts_dir})
endef

download_artifacts:
	$(foreach package,$(npm_packages),$(call get_npm_package,$(environment),$(package)))

generate:
	$(info Running Go code generator)
	go generate ./...

build:
	$(info Building Go code)
	go build -o keep-client -a . 

cmd-help: build
	@echo '$$ keep-client start --help' > docs/development/cmd-help
	./keep-client start --help >> docs/development/cmd-help

.PHONY: all
