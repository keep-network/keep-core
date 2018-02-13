
all: GenRequestID_sol_GenRequestID.abi

GenRequestID_sol_GenRequestID.abi: GenRequestID.sol
	solc --abi GenRequestID.sol
	abigen --abi GenRequestID_sol_GenRequestID.abi --pkg GenRequestID --out GenRequestID.go
	mkdir -p ../../../interface/m4/GenRequestID
	cp GenRequestID.go ../../../interface/m4/GenRequestID

