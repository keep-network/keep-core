
all: GenRequestID_sol_GenRequestID.abi

GenRequestID_sol_GenRequestID.abi: GenRequestID.sol
	solc --abi GenRequestID.sol
	abigen --abi GenRequestID_sol_GenRequestID.abi --pkg GenRequestID --out GenRequestID.go
	mkdir -p ../../../interface/m4/GenRequestID
	cp GenRequestID.go ../../../interface/m4/GenRequestID

m4_changequote([[[,]]])m4_dnl
m4_define([[[m4_comment]]],[[[]]])m4_dnl
m4_comment([[[
	@title This is the address of the RequestID generation contract. 
	@author Philip Schlump
	@notice Set for internal development test environment in geth
]]])m4_dnl
m4_define([[[GEN_REQUEST_ID_ADDR]]],[[["0x9fbda871d559710256a2502a2517b794b482db40"]]])m4_dnl
