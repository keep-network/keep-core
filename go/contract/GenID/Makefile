
all:

migrate:
	truffle migrate --reset	>most-recent.out
	grep GenRequestID: most-recent.out >/dev/null
	./get-addr.sh most-recent.out >most-recent.addr

t_dev:
	truffle develop log

# can act as a build/server and pacakager.
# 		http://www.ziggify.com/blog/truffle-3-ethereum-devops-truffle-default-builder-blockchain-stack-4c/



t1: 
	abigen --solc /Users/corwin/bin/solc2 --sol contracts/ProofOfExistence1.sol --pkg main --out test_t1.go

# See: https://github.com/ethereum/py-solc/issues/42

# contracts/token_sol_MyToken.abi  contracts/token_sol_tokenRecipient.abi

