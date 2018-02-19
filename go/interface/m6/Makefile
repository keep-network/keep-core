
all:
	go build

all_old: gen_abi  gen_go

gen_abi:
	solc --abi greeter.sol

gen_go:
	abigen --abi greeter_sol_Greeter.abi --type Greeter --pkg main --out greeter.go
