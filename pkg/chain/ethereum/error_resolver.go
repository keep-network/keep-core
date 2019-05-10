package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ABI for errors bubbled out from revert calls. Not used directly as errors are
// neither encoded strictly as method calls nor strictly as return values, nor
// strictly as events, but some various bits of it are used for unpacking the
// errors. See ResolveError below.
const errorABIString = "[{\"constant\":true,\"outputs\":[{\"type\":\"string\"}],\"inputs\":[{\"name\":\"message\", \"type\":\"string\"}],\"name\":\"Error\"}]"

var errorABI abi.ABI

func init() {
	var err error
	errorABI, err = abi.JSON(strings.NewReader(errorABIString))
	if err != nil {
		panic(fmt.Sprintf("Failed to parse error ABI string: [%v] (ABI was: [%v])", err, errorABIString))
	}
}

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
func (er *errorResolver) ResolveError(originalErr error, value *big.Int, methodName string, parameters ...interface{}) error {
	packed, err := er.abi.Pack(methodName, parameters...)
	msg := ethereum.CallMsg{To: er.address, Data: packed, Value: value}

	response, err := er.client.CallContract(context.TODO(), msg, nil)
	if err != nil {
		return fmt.Errorf("got error [%v] while resolving original error [%v]", err, originalErr)
	}

	// An error is returned as a 4-byte error id (same encoding as a method id)
	// followed by a set of ABI-encoded values as if the error were a method
	// that returned those values.
	//
	// Current spec-ish @ https://github.com/ethereum/EIPs/issues/838#issuecomment-458919375
	// Bless Ethereum's heart.
	errorID, encodedReturns := response[0:3], response[4:]

	errorMethod, err := errorABI.MethodById(errorID)
	if err != nil {
		return fmt.Errorf("got [%v] while resolving original error [%v] on return [%v]", err, originalErr, response)
	}

	errorValues, err := errorMethod.Outputs.UnpackValues(encodedReturns)
	if err != nil {
		return fmt.Errorf("got [%v] while resolving original error [%v] on return [%v]", err, originalErr, response)
	}

	return fmt.Errorf(
		"contract failed with: [%v] (original error [%v])",
		errorValues,
		originalErr,
	)
}
