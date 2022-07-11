
# Tag of the npm package version for contracts artifacts. If not set by the 
# `CONTRACTS_NPM_PACKAGE_TAG` environment variable it defaults to `development`.
npm_package_tag:=$(CONTRACTS_NPM_PACKAGE_TAG)
ifndef npm_package_tag
	npm_package_tag=development
endif

contracts_dir =_contracts
artifacts_dir = $(realpath ${contracts_dir}/package/artifacts)

# Go bindings generated for the solidity contracts.
contract_files = $(addprefix contract/,$(addsuffix .go,${required_contracts}))

# Files containing addresses extracted from the artifacts.
address_files = $(addprefix _address/,${required_contracts})

all: clean ${address_files} gen_contract_go

clean:
	rm -rf _address/*
	rm -rf abi/*
	rm -rf contract/*
	mkdir tmp && mv cmd/cmd*.go tmp
	rm -rf cmd/*
	mv tmp/* cmd && rm -rf tmp

# It requires npm of at least 7.x version to support `pack-destination` flag.
download_artifacts:
	rm -rf ${contracts_dir} && mkdir ${contracts_dir}
	npm pack --silent \
		--pack-destination=${contracts_dir} \
		$(shell npm view ${npm_package_name}@${npm_package_tag} _id) \
		| xargs -I{} tar -zxf ${contracts_dir}/{} -C ${contracts_dir}

gen_contract_go: ${contract_files}

abi/%.abi: ${artifacts_dir}/%.json
	jq .abi $< > abi/$*.abi

abi/%.go: abi/%.abi
	go run github.com/ethereum/go-ethereum/cmd/abigen --abi $< --pkg abi --type $* --out $@

# Extract address of a contract from an artifact. Artfifacts that are published
# with `development` tag are not accessible on developers' environments, so we
# replace them with zeros.
_address/%: ${artifacts_dir}/%.json
ifeq ($(npm_package_tag), development)
	@echo "Skipping address extraction for development package"
	@echo "0x0000000000000000000000000000000000000000" > _address/$*
else
	jq -jr .address ${artifacts_dir}/$*.json > _address/$*
endif

contract/%.go cmd/%.go: abi/%.abi abi/%.go _address/% ${artifacts_dir}/%.json
	go run github.com/keep-network/keep-common/tools/generators/ethlike $< contract/$*.go cmd/$*.go

# Don't remove intermediate files that got generated.
.PRECIOUS: abi/%.abi abi/%.go _address/%
