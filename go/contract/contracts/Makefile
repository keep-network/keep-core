
all: KStart_sol_KStart.abi


t1:
	( cd .. ; make t1 )

t4:
	( cd .. ; make t4 )

build_abi:


HelloWorld_sol_HelloWorld.abi: HelloWorld.sol
	solc --abi HelloWorld.sol

build_go_01:
	solc --abi HelloWorld.sol
	abigen --abi HelloWorld_sol_HelloWorld.abi --pkg Bob --out x.go

KStart_sol_KStart.abi: KStart.sol
	solc --abi KStart.sol
	abigen --abi KStart_sol_KStart.abi --pkg KStart --out KStart.go


# Note: https://github.com/ethereum/go-ethereum/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts
# https://ethereum.stackexchange.com/questions/12436/how-to-communicate-with-a-remote-node
