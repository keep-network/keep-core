package ethereum

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// errorResolver bundles up the bits needed to turn errors like "failed to
// estimate gas needed" that are triggered by contract reverts but don't include
// revert causes into proper revert error messages from a contract by calling
// the contract method without trying to commit it.
//
// It has one method, ResolveError, that does the heavy lifting.
type errorResolver struct {
	client  *ethclient.Client
	abi     *abi.ABI
	address *common.Address
}

// Resolves the given transaction error to a standard error that, if available,
// contains the error message the transaction produced when reverting.
func (er *errorResolver) ResolveError(originalErr error, value *big.Int, method string, parameters ...interface{}) error {
	packed, err := er.abi.Pack(method, parameters...)
	msg := ethereum.CallMsg{To: er.address, Data: packed, Value: value}

	response, err := er.client.CallContract(context.TODO(), msg, nil)
	if err != nil {
		return fmt.Errorf("got error [%v] while resolving original error [%v]", err, originalErr)
	}

	responseString := string(response)
	if len(responseString) < 68 {
		return fmt.Errorf("couldn't interpret contract return while resolving original error [%v]", originalErr)
	}

	contractErrorMessage := responseString[68:]
	return fmt.Errorf(
		"contract failed with: [%v] (original error [%v])",
		contractErrorMessage,
		originalErr,
	)
}
