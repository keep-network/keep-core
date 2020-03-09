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
# *ImplV1.go files will get generated into clean Keep contract bindings, the
# corresponding contract filenames will drop the ImplV1, if it exists, and live
# in the contract/ directory.
clean_contract_stems := $(filter %ImplV1,$(contract_stems)) $(filter %Operator,$(contract_stems)) $(filter TokenStaking, $(contract_stems)) $(filter TokenGrant, $(contract_stems))
contract_files := $(addprefix contract/,$(addsuffix .go,$(subst ImplV1,,$(clean_contract_stems))))

all: gen_contract_go gen_abi_go

clean:
	rm -r abi/*
	rm -r contract/*
	mkdir tmp && mv cmd/cmd*.go tmp
	rm -r cmd/*
	mv tmp/* cmd && rm -r tmp

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

contract/%.go cmd/%.go: abi/%ImplV1.abi abi/%ImplV1.go abi/%.go *.go
	go run github.com/keep-network/keep-common/tools/generators/ethereum $< contract/$*.go cmd/$*.go

contract/%Operator.go cmd/%Operator.go: abi/%Operator.abi abi/%Operator.go *.go
	go run github.com/keep-network/keep-common/tools/generators/ethereum $< contract/$*Operator.go cmd/$*Operator.go

contract/TokenStaking.go cmd/TokenStaking.go: abi/TokenStaking.abi abi/TokenStaking.go *.go
	go run github.com/keep-network/keep-common/tools/generators/ethereum $< contract/TokenStaking.go cmd/TokenStaking.go

contract/TokenGrant.go cmd/TokenGrant.go: abi/TokenGrant.abi abi/TokenGrant.go *.go
	go run github.com/keep-network/keep-common/tools/generators/ethereum $< contract/TokenGrant.go cmd/TokenGrant.go
