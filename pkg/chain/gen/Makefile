
# direcotry where the solidity contracts are
SOL_DIR=../../../contracts/solidity/contracts

# relative directory where the Go code nees to go.
OP_DIR=../../../pkg/chain/gen

all:
	( cd ${SOL_DIR} ; make -f ${OP_DIR}/Makefile gen_go )
	go build

SOL=$(shell ls *.sol | sed -e 's/.sol/.go/' )

gen_go: 
	make -f ${OP_DIR}/Makefile ${SOL}

%.go : %.sol
	( cd ${SOL_DIR} ; solc zeppelin-solidity/=../node_modules/zeppelin-solidity/ --allow-paths . --overwrite --abi -o . $<  )
	( cd ${SOL_DIR} ; abigen --abi $*.abi --pkg gen --type $* --out $@ )
	mv $@ ${OP_DIR}/$@

