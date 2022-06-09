# Solidity files prefixed with the solidity directory path.
solidity_files := $(addprefix ${solidity_dir}/,${solidity_files_sources}) 

# Go bindings generated for the solidity contracts.
contract_files = $(addprefix contract/,$(addsuffix .go,${required_contracts}))

# Function to run solidity sources comilation.
# Expects two arguments:
# 1 - dependencies mapping for contracts import
# 2 - solidity source file path
solc_generate = solc $(1) \
		--allow-paths ${solidity_dir} \
		--overwrite \
		--abi \
		-o abi $(2)

all: gen_contract_go

clean:
	rm -r abi/*
	rm -r contract/*
	mkdir tmp && mv cmd/cmd*.go tmp
	rm -r cmd/*
	mv tmp/* cmd && rm -r tmp

gen_contract_go: ${contract_files}

abi/%.abi: ${solidity_files}
	$(call solc_generate, ${solidity_dependencies_mapping}, $(filter %$*.sol, ${solidity_files}))

abi/%.go: abi/%.abi
	go run github.com/ethereum/go-ethereum/cmd/abigen --abi $< --pkg abi --type $* --out $@

contract/%.go cmd/%.go: abi/%.abi abi/%.go
	go run github.com/keep-network/keep-common/tools/generators/ethlike $< contract/$*.go cmd/$*.go

# Don't remove intermediate files that got generated.
.PRECIOUS: abi/%.abi abi/%.go
