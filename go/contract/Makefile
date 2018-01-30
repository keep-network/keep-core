
all:

migrate:
	truffle migrate --reset

t_dev:
	truffle develop log

# can act as a build/server and pacakager.
# 		http://www.ziggify.com/blog/truffle-3-ethereum-devops-truffle-default-builder-blockchain-stack-4c/



t1: 
	abigen --solc /Users/corwin/bin/solc2 --sol contracts/ProofOfExistence1.sol --pkg main --out test_t1.go

# See: https://github.com/ethereum/py-solc/issues/42

# contracts/token_sol_MyToken.abi  contracts/token_sol_tokenRecipient.abi

t2: 
	( cd contracts ; solc --abi token.sol )
	abigen --abi contracts/token_sol_Token.abi --type TokenType --pkg main --out test_t2_a.go
	abigen --abi contracts/token_sol_tokenRecipient.abi --type TokenRecipientType --pkg main --out test_t2_b.go
	cp test_t2_?.go /Users/corwin/go/src/www.2c-why.com/job/m4


