solidity_dir=$(realpath ${SOLIDITY_DIR})
solidity_files := $(wildcard ${solidity_dir}/contracts/*.sol)

abi_files := $(foreach file, $(solidity_files), abi/$(notdir $(basename $(file))).go)

all: gen_abi_go

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
