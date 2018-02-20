# Notes

## To use abigen

`abigen` seems to refuse to work with a --sol flag.   This can be worked around

```
	cd contracts
	solc --abi ContractName.sol
	# this creates ContractName_sol_ContractName.abi
	abigen --abi ContractName_sol_ContractName.abi --pkg ContractName --out ContractName.go 
```


