
# Direcotry where the solidity contracts are.
SOL_DIR=../../../contracts/solidity/contracts

# Relative directory where the Go code needs to go.  This is relative to the
# ${SOL_DIR} directory.
OP_DIR=../../../pkg/chain/gen

generate_solidity_go:
	( cd ${SOL_DIR} ; make -f ${OP_DIR}/Makefile gen_go )

all:
	( cd ${SOL_DIR} ; make -f ${OP_DIR}/Makefile gen_go )
	( cd ethereum ; go build )

SOL=$(shell ls *.sol | sed -e 's/.sol/.go/' )

gen_go: 
	make -f ${OP_DIR}/Makefile ${SOL}

%.go : %.sol
	( cd ${SOL_DIR} ; solc zeppelin-solidity/=../node_modules/zeppelin-solidity/ --allow-paths . --overwrite --abi -o . $<  )
	( cd ${SOL_DIR} ; abigen --abi $*.abi --pkg ethereum --type $* --out $@ )
	mv $@ ${OP_DIR}/ethereum/$@

