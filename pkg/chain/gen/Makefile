
# In .../github.com/keep-network/keep-core/pkg/chain/gen

# Directory where the solidity contracts are.
SOL_DIR=../../../contracts/solidity/contracts

# Relative directory where the Go code needs to go.  This is relative to the
# ${SOL_DIR} directory.
OP_DIR=../../../pkg/chain/gen
GEN_DIR=../../../pkg/chain/gen/abi

generate_solidity_go:
	( cd ${SOL_DIR} ; make -f ${OP_DIR}/Makefile gen_go )

all:
	( cd ${SOL_DIR} ; make -f ${OP_DIR}/Makefile gen_go )
	( cd ${GEN_DIR} ; goimports -w *.go ; go build )

SOL=$(shell ls *.sol | sed -e 's/.sol/.go/' )

gen_go:
	make -f ${OP_DIR}/Makefile ${SOL}

%.go : %.sol
	( cd ${SOL_DIR} ; solc solidity-bytes-utils/=../node_modules/solidity-bytes-utils/ openzeppelin-solidity/=../node_modules/openzeppelin-solidity/ --allow-paths . --overwrite --abi -o . $<  )
	( cd ${SOL_DIR} ; abigen --abi $*.abi --pkg abi --type $* --out ${GEN_DIR}/$@ )

