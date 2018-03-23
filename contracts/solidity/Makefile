
all: contracts/KeepRelayBeacon_sol_KeepRelayBeacon.abi

contracts/KeepRelayBeacon_sol_KeepRelayBeacon.abi: contracts/KeepRelayBeacon.sol
	( cd contracts ; solcjs --abi \
		./zeppelin-solidity/contracts/token/ERC20/StandardToken.sol \
		./zeppelin-solidity/contracts/token/ERC20/BasicToken.sol \
		./zeppelin-solidity/contracts/token/ERC20/ERC20.sol \
		./zeppelin-solidity/contracts/token/ERC20/ERC20Basic.sol \
		./zeppelin-solidity/contracts/token/ERC20/SafeERC20.sol \
		./zeppelin-solidity/contracts/math/SafeMath.sol \
		./TokenStaking.sol \
	 	./KeepRelayBeacon.sol \
		)
	( cd contracts ; abigen --abi __KeepRelayBeacon_sol_KeepRelayBeacon.abi --pkg KeepRelayBeacon --out KeepRelayBeacon.go )
	mkdir -p ./lib/KeepRelayBeacon
	cp ./contracts/KeepRelayBeacon.go ./lib/KeepRelayBeacon
	( cd ./lib/KeepRelayBeacon ; go build )

local_migrate:
	truffle migrate --reset	--network local

