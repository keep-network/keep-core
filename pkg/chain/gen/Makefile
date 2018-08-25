
# In .../github.com/keep-network/keep-core/pkg/chain/gen

# Directory where the solidity contracts are.
SOL_DIR=../../../contracts/solidity/contracts

# Relative directory where the Go code needs to go.  This is relative to the
# ${SOL_DIR} directory.
OP_DIR=../../../pkg/chain/gen
GEN_DIR=../../../pkg/chain/gen/abi/

generate_solidity_go:
	( cd ${SOL_DIR} ; make -f ${OP_DIR}/Makefile gen_go )

all:
	( cd ${SOL_DIR} ; make -f ${OP_DIR}/Makefile gen_go )
	goimports -w *.go
	go build

GO=$(shell ls *.sol | sed -e 's/.sol/.go/' )

gen_go: 
	make -f ${OP_DIR}/Makefile ${GO}

%.go : %.sol
	( cd ${SOL_DIR} ; mkdir ${GEN_DIR}/$* ; solc openzeppelin-solidity/=../node_modules/openzeppelin-solidity/ \
		--allow-paths . \
		--combined-json abi,bin,devdoc,userdoc,metadata \
		$< | \
		abigen --abi - --pkg $* --out ${GEN_DIR}/$*/$@ \
		)


