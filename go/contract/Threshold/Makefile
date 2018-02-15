
all: contracts/KStart.sol contracts/KStart_sol_KStart.abi

migrate: contracts/KStart.sol
	truffle migrate --reset

contracts/KStart.sol: src/KStart.raw.sol ../GenID/addr.m4
	m4 -P src/KStart.raw.sol >contracts/KStart.sol

contracts/KStart_sol_KStart.abi: contracts/KStart.sol
	( cd contracts ; solc --abi KStart.sol )
	( cd contracts ; abigen --abi KStart_sol_KStart.abi --pkg KStart --out KStart.go )
	cp ./contracts/KStart.go ../../interface/m4/KStart

t_dev:
	truffle develop log




# can act as a build/server and pacakager.
# 		http://www.ziggify.com/blog/truffle-3-ethereum-devops-truffle-default-builder-blockchain-stack-4c/
#
#t1: 
#	abigen --solc /Users/corwin/bin/solc2 --sol contracts/ProofOfExistence1.sol --pkg main --out test_t1.go
#
# See: https://github.com/ethereum/py-solc/issues/42
#
# contracts/token_sol_MyToken.abi  contracts/token_sol_tokenRecipient.abi

