
# Tag of the npm package version for contracts artifacts. If not set by the 
# `CONTRACTS_NPM_PACKAGE_TAG` environment variable it defaults to `development`.
npm_package_tag:=$(CONTRACTS_NPM_PACKAGE_TAG)
ifndef npm_package_tag
	npm_package_tag=development
endif

contracts_dir =_contracts
artifacts_dir = $(realpath ${contracts_dir}/package/export/artifacts)

# Go bindings generated for the solidity contracts.
contract_files = $(addprefix contract/,$(addsuffix .go,${required_contracts}))

all: gen_contract_go

clean:
	rm -rf ${contracts_dir}/*
	rm -r abi/*
	rm -r contract/*
	mkdir tmp && mv cmd/cmd*.go tmp
	rm -r cmd/*
	mv tmp/* cmd && rm -r tmp

# It requires npm of at least 7.x version to support `pack-destination` flag.
download_artifacts:
	rm -rf ${contracts_dir} && mkdir ${contracts_dir}
	npm pack --silent --pack-destination=${contracts_dir} ${npm_package_name}@${npm_package_tag}\
		| xargs -I{} tar -zxf ${contracts_dir}/{} -C ${contracts_dir}

gen_contract_go: ${contract_files}

abi/%.abi: ${artifacts_dir}/%.json
	jq .abi $< > abi/$*.abi

abi/%.go: abi/%.abi
	go run github.com/ethereum/go-ethereum/cmd/abigen --abi $< --pkg abi --type $* --out $@

contract/%.go cmd/%.go: abi/%.abi abi/%.go
	go run github.com/keep-network/keep-common/tools/generators/ethlike $< contract/$*.go cmd/$*.go

# Don't remove intermediate files that got generated.
.PRECIOUS: abi/%.abi abi/%.go
