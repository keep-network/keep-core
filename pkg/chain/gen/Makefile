
# direcotry where the solidity contracts are
CDIR=../../../contracts/solidity/contracts

all: \
	KeepRandomBeacon/KeepRandomBeacon.go \
	KeepRandomBeaconImplV1/KeepRandomBeaconImplV1.go \
	KeepToken/KeepToken.go \
	StakingProxy/StakingProxy.go \
	TokenGrant/TokenGrant.go \
	TokenStaking/TokenStaking.go

KeepRandomBeacon/KeepRandomBeacon.go: ../../../contracts/solidity/contracts/KeepRandomBeacon.sol
	( cd ${CDIR} ; solc zeppelin-solidity/=../node_modules/zeppelin-solidity/ --allow-paths . --overwrite --abi -o . KeepRandomBeacon.sol  )
	( cd ${CDIR} ; abigen --abi KeepRandomBeacon.abi --pkg KeepRandomBeacon --out KeepRandomBeacon.go )
	mkdir -p ./KeepRandomBeacon
	mv ${CDIR}/KeepRandomBeacon.go ./KeepRandomBeacon
	( cd ./KeepRandomBeacon ; go build )

KeepRandomBeaconImplV1/KeepRandomBeaconImplV1.go: ../../../contracts/solidity/contracts/KeepRandomBeaconImplV1.sol
	( cd ${CDIR} ; solc zeppelin-solidity/=../node_modules/zeppelin-solidity/ --allow-paths . --overwrite --abi -o . KeepRandomBeaconImplV1.sol  )
	( cd ${CDIR} ; abigen --abi KeepRandomBeaconImplV1.abi --pkg KeepRandomBeaconImplV1 --out KeepRandomBeaconImplV1.go )
	mkdir -p ./KeepRandomBeaconImplV1
	mv ${CDIR}/KeepRandomBeaconImplV1.go ./KeepRandomBeaconImplV1
	( cd ./KeepRandomBeaconImplV1 ; go build )

KeepToken/KeepToken.go: ../../../contracts/solidity/contracts/KeepToken.sol
	( cd ${CDIR} ; solc zeppelin-solidity/=../node_modules/zeppelin-solidity/ --allow-paths . --overwrite --abi -o . KeepToken.sol  )
	( cd ${CDIR} ; abigen --abi KeepToken.abi --pkg KeepToken --out KeepToken.go )
	mkdir -p ./KeepToken
	mv ${CDIR}/KeepToken.go ./KeepToken
	( cd ./KeepToken ; go build )

StakingProxy/StakingProxy.go: ../../../contracts/solidity/contracts/StakingProxy.sol
	( cd ${CDIR} ; solc zeppelin-solidity/=../node_modules/zeppelin-solidity/ --allow-paths . --overwrite --abi -o . StakingProxy.sol  )
	( cd ${CDIR} ; abigen --abi StakingProxy.abi --pkg StakingProxy --out StakingProxy.go )
	mkdir -p ./StakingProxy
	mv ${CDIR}/StakingProxy.go ./StakingProxy
	( cd ./StakingProxy ; go build )

TokenGrant/TokenGrant.go: ../../../contracts/solidity/contracts/TokenGrant.sol
	( cd ${CDIR} ; solc zeppelin-solidity/=../node_modules/zeppelin-solidity/ --allow-paths . --overwrite --abi -o . TokenGrant.sol  )
	( cd ${CDIR} ; abigen --abi TokenGrant.abi --pkg TokenGrant --out TokenGrant.go )
	mkdir -p ./TokenGrant
	mv ${CDIR}/TokenGrant.go ./TokenGrant
	( cd ./TokenGrant ; go build )

TokenStaking/TokenStaking.go: ../../../contracts/solidity/contracts/TokenStaking.sol
	( cd ${CDIR} ; solc zeppelin-solidity/=../node_modules/zeppelin-solidity/ --allow-paths . --overwrite --abi -o . TokenStaking.sol  )
	( cd ${CDIR} ; abigen --abi TokenStaking.abi --pkg TokenStaking --out TokenStaking.go )
	mkdir -p ./TokenStaking
	mv ${CDIR}/TokenStaking.go ./TokenStaking
	( cd ./TokenStaking ; go build )







test:
	go run gen.go

generate: all

clean:
	( cd ${CDIR} ; rm -f ./KeepRandomBeacon.go )
	( cd ${CDIR} ; rm -f ./KeepRandomBeaconImplV1.go )
	( cd ${CDIR} ; rm -f ./KeepToken.go )
	( cd ${CDIR} ; rm -f ./StakingProxy.go )
	( cd ${CDIR} ; rm -f ./TokenGrant.go )
	( cd ${CDIR} ; rm -f ./TokenStaking.go )


# Sample - change YYY to name of contract.

#../../../contracts/solidity/contracts/YYY.go: ../../../contracts/solidity/contracts/YYY.sol
#	( cd ${CDIR} ; solc zeppelin-solidity/=../node_modules/zeppelin-solidity/ --allow-paths . --overwrite --abi -o . YYY.sol  )
#	( cd ${CDIR} ; abigen --abi YYY.abi --pkg YYY --out YYY.go )
#	mkdir -p ./YYY
#	mv ${CDIR}/YYY.go ./YYY
#	( cd ./YYY ; go build )


