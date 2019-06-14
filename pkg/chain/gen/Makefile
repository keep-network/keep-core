# Environment provides the solidity directory as a potentially-relative path,
# which we resolve. Then we resolve the Solidity files in a contracts/ directory
# at that path.
solidity_dir=$(realpath ${SOLIDITY_DIR})
solidity_files := $(wildcard ${solidity_dir}/contracts/*.sol)

# Bare Solidity filenames without .sol or Solidity directory prefix.
contract_stems := $(notdir $(basename $(solidity_files)))
# Go abigen bindings in abi/ subdirectory with .go suffix, alongside solc ABI
# files with .abi suffix.
abi_files := $(addprefix abi/,$(addsuffix .abi,$(contract_stems)))
abigen_files := $(addprefix abi/,$(addsuffix .go,$(contract_stems)))
# *ImplV1.go files will get generated into clean Keep contract bindings, as will
# the Staking*.go contract(s). The corresponding contract filenames will drop
# the ImplV1, if it exists, and live in the contract/ directory.
clean_contract_stems := $(filter %ImplV1,$(contract_stems)) $(filter Staking%,$(contract_stems)) $(filter TokenGrant, $(contract_stems))
contract_files := $(addprefix contract/,$(addsuffix .go,$(subst ImplV1,,$(clean_contract_stems))))

all: gen_contract_go gen_abi_go

clean:
	rm -r abi/*
	rm -r contract/*

gen_abi_go: $(abigen_files)

gen_contract_go: $(contract_files)

abi/%.abi: ${solidity_dir}/contracts/%.sol
	solc solidity-bytes-utils/=${solidity_dir}/node_modules/solidity-bytes-utils/ \
		 openzeppelin-solidity/=${solidity_dir}/node_modules/openzeppelin-solidity/ \
		 --allow-paths ${solidity_dir} \
		 --overwrite \
		 --abi \
		 -o abi $<

abi/%.go: abi/%.abi
	abigen --abi $< --pkg abi --type $* --out $@

contract/%.go cmd/%.go: abi/%ImplV1.abi abi/%ImplV1.go abi/%.go *.tmpl *.go
	go run ./ $< contract/$*.go cmd/$*.go

contract/Staking%.go cmd/Staking%.go: abi/Staking%.abi abi/Staking%.go *.tmpl *.go
	go run ./ $< contract/Staking$*.go cmd/Staking$*.go

contract/TokenGrant.go cmd/TokenGrant.go: abi/TokenGrant.abi abi/TokenGrant.go *.tmpl *.go
	go run ./ $< contract/TokenGrant.go cmd/TokenGrant.go
