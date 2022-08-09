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

local:
	make all environment=local

all: get_artifacts generate build cmd-help

modules := beacon \
	ecdsa \
	threshold \
	tbtc

# Required by get_npm_package function.
npm_beacon_package := @keep-network/random-beacon
npm_ecdsa_package := @keep-network/ecdsa
npm_threshold_package := @threshold-network/solidity-contracts
npm_tbtc_package := @keep-network/tbtc-v2

# Required by get_local_package function. The paths can be overwritten when calling
# the make command, e.g.:
#   make local local_threshold_path=/other/path/threshold
local_beacon_path := ./solidity/random-beacon
local_ecdsa_path := ./solidity/ecdsa
local_threshold_path := ../../threshold-network/solidity-contracts
local_tbtc_path := ../tbtc-v2/solidity

# Working directory where contracts artifacts should be stored.
contracts_dir := tmp/contracts

# It requires npm of at least 7.x version to support `pack-destination` flag.
define get_npm_package
$(eval npm_package_name := $(npm_$(1)_package))
$(eval npm_package_tag := $(2))
$(eval destination_dir := ${contracts_dir}/${npm_package_tag}/${npm_package_name})
@rm -rf ${destination_dir}
@mkdir -p ${destination_dir}
@npm pack --silent \
	--pack-destination=${destination_dir} \
	$(shell npm view ${npm_package_name}@${npm_package_tag} _id) \
	| xargs -I{} tar -zxf ${destination_dir}/{} -C ${destination_dir} --strip-components 1 package/artifacts
$(info Downloaded NPM package ${npm_package_name}@${npm_package_tag} to ${contracts_dir})
endef

define get_local_package
$(eval module := $(1))
$(eval local_solidity_path := $(local_$(module)_path))
$(eval npm_package_name := $(npm_$(module)_package))
$(eval destination_dir := ${contracts_dir}/local/${npm_package_name})
@[ -d "$(local_solidity_path)" ] || { echo "$(module) path [$(local_solidity_path)] does not exist!"; exit 1; }
@rm -rf ${destination_dir}
@mkdir -p ${destination_dir}
$(info Fetching local package ${module} from path ${local_solidity_path})
rsync -a $(local_solidity_path)/deployments/development/ ${destination_dir}/artifacts
endef

get_artifacts:
ifeq ($(environment), local)
	$(foreach module,$(modules),$(call get_local_package,$(module)))
else
	$(foreach module,$(modules),$(call get_npm_package,$(module),$(environment)))
endif

generate:
	$(info Running Go code generator)
	go generate ./...

build:
	$(info Building Go code)
	go build -o keep-client -a . 

cmd-help: build
	@echo '$$ keep-client start --help' > docs/development/cmd-help
	./keep-client start --help >> docs/development/cmd-help

.PHONY: all development goerli download_artifacts generate build cmd-help
