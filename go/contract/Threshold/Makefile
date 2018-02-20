
all: contracts/KeepRelayBeacon.sol contracts/KeepRelayBeacon_sol_KeepRelayBeacon.abi

# for local testing with truffle - note this will not test events correctly with Go
migrate: contracts/KeepRelayBeacon.sol
	truffle migrate --reset

# put into the testnet system
migrate_testnet: contracts/KeepRelayBeacon.sol
	truffle migrate --reset --network testnet >,addr 2>&1
	cat ,addr

contracts/KeepRelayBeacon.sol: src/KeepRelayBeacon.src.sol ../GenID/addr.m4
	m4 -P src/KeepRelayBeacon.src.sol >contracts/KeepRelayBeacon.sol

contracts/KeepRelayBeacon_sol_KeepRelayBeacon.abi: contracts/KeepRelayBeacon.sol
	( cd contracts ; solc --abi KeepRelayBeacon.sol )
	( cd contracts ; abigen --abi KeepRelayBeacon_sol_KeepRelayBeacon.abi --pkg KeepRelayBeacon --out KeepRelayBeacon.go )
	mkdir -p ../../interface/lib/KeepRelayBeacon
	cp ./contracts/KeepRelayBeacon.go ../../interface/lib/KeepRelayBeacon

t_dev:
	truffle develop log

