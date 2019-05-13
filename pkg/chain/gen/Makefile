solidity_dir=$(realpath ${SOLIDITY_DIR})
solidity_files := $(wildcard ${solidity_dir}/contracts/*.sol)

abi_files := $(foreach file, $(solidity_files), abi/$(notdir $(basename $(file))).go)
# *ImplV1.go files will get generated into first-class clean contract bindings.
# The corresponding contract filenames will drop the ImplV1 and live in the
# contract/ directory.
contract_candidates := $(filter %ImplV1.go, $(abi_files))
contract_files := $(foreach file, $(contract_candidates), contract/$(notdir $(subst ImplV1,,$(basename $(file)))).go)

all: gen_contract_go

clean:
	rm -rf ${abi_files}
	rm -rf ${contract_files}

gen_abi_go: $(abi_files)

gen_contract_go: $(contract_files)

abi/%.abi: ${solidity_dir}/contracts/%.sol
	solc solidity-bytes-utils/=${solidity_dir}/node_modules/solidity-bytes-utils/ \
		 openzeppelin-solidity/=${solidity_dir}/node_modules/openzeppelin-solidity/ \
		 --allow-paths ${solidity_dir} \
		 --overwrite \
		 --abi \
		 -o abi $<

abi/%.go : abi/%.abi
	abigen --abi $< --pkg abi --type $* --out $@
	goimports -w $@
	go build $@

contract/%.go: abi/%ImplV1.abi abi/%ImplV1.go abi/%.go
	go run ./ $< $@
	goimports -w $@
	go build $@
