package ethutil_test

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/keep-network/keep-core/pkg/chain/ethereum/ethutil"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var testABIMethods = map[string]abi.Method{
	"Test": abi.Method{
		Const:   false,
		Name:    "Test",
		Inputs:  []abi.Argument{},
		Outputs: []abi.Argument{},
	},
}
var testABI = abi.ABI{
	Constructor: testABIMethods["Test"],
	Events:      map[string]abi.Event{},
	Methods:     testABIMethods,
}
var testAddress = common.Address([20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

var errOriginal = fmt.Errorf("OG")
var errCall = fmt.Errorf("call error")

func assertErrorContains(t *testing.T, err error, substrings ...string) {
	t.Helper()

	if err == nil {
		t.Errorf("\nexpected: {error}\nactual:   nil")
		return
	}

	for _, substring := range substrings {
		if !strings.Contains(err.Error(), substring) {
			t.Errorf(
				"\nexpected: {error containing [%#v]}\nactual:   %v",
				substring,
				err,
			)
		}
	}
}

// Helper caller type that always returns errCall.
type erroringCaller struct{}

func (*erroringCaller) CallContract(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	return nil, errCall
}

// Helper caller type that always returns a nil error and the provided
// returnedBytes.
type fixedReturnCaller struct {
	returnedBytes []byte
}

func (frc *fixedReturnCaller) CallContract(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	return frc.returnedBytes, nil
}

func TestErrorResolverHandlesErrorCall(t *testing.T) {
	caller := &erroringCaller{}
	resolver := ethutil.NewErrorResolver(caller, &testABI, &testAddress)

	err := resolver.ResolveError(errOriginal, common.Address{}, nil, "Test")
	if err == nil {
		t.Errorf("\nexpected: {error}\nactual:   nil")
		return
	}

	assertErrorContains(
		t,
		err,
		errOriginal.Error(),
		errCall.Error(),
	)
}

func TestErrorResolverHandlesShortResponses(t *testing.T) {
	caller := &fixedReturnCaller{}
	resolver := ethutil.NewErrorResolver(caller, &testABI, &testAddress)

	for returnLength := 0; returnLength < 4; returnLength++ {
		caller.returnedBytes = make([]byte, returnLength)
		err := resolver.ResolveError(errOriginal, common.Address{}, nil, "Test")
		assertErrorContains(
			t,
			err,
			errOriginal.Error(),
			fmt.Sprintf("%v", caller.returnedBytes),
			"was not long enough",
		)
	}
}

func TestErrorResolverHandlesUnknownMethodResponses(t *testing.T) {
	caller := &fixedReturnCaller{[]byte{0, 0, 0, 1}}
	resolver := ethutil.NewErrorResolver(caller, &testABI, &testAddress)

	err := resolver.ResolveError(errOriginal, common.Address{}, nil, "Test")
	assertErrorContains(
		t,
		err,
		errOriginal.Error(),
		fmt.Sprintf("%v", caller.returnedBytes),
		"no method with id",
	)
}

func TestErrorResolverHandlesBadParameterResponses(t *testing.T) {
	caller := &fixedReturnCaller{[]byte{8, 195, 121, 160}}
	resolver := ethutil.NewErrorResolver(caller, &testABI, &testAddress)

	// bad response length
	err := resolver.ResolveError(errOriginal, common.Address{}, nil, "Test")
	assertErrorContains(
		t,
		err,
		errOriginal.Error(),
		fmt.Sprintf("%v", caller.returnedBytes),
		"length insufficient",
	)

	// good response length, bad data offset in response
	buildingBlock := [32]byte{}
	caller.returnedBytes = append(caller.returnedBytes, buildingBlock[:]...)
	caller.returnedBytes[len(caller.returnedBytes)-1] = 1
	err = resolver.ResolveError(errOriginal, common.Address{}, nil, "Test")
	assertErrorContains(
		t,
		err,
		errOriginal.Error(),
		fmt.Sprintf("%v", caller.returnedBytes),
		"would go over slice boundary",
	)

	// good response length, good data offset, bad string length in response
	caller.returnedBytes = []byte{8, 195, 121, 160}
	caller.returnedBytes = append(caller.returnedBytes, buildingBlock[:]...)
	caller.returnedBytes[len(caller.returnedBytes)-1] = 32
	caller.returnedBytes = append(caller.returnedBytes, buildingBlock[:]...)
	caller.returnedBytes[len(caller.returnedBytes)-1] = 1
	err = resolver.ResolveError(errOriginal, common.Address{}, nil, "Test")
	assertErrorContains(
		t,
		err,
		errOriginal.Error(),
		fmt.Sprintf("%v", caller.returnedBytes),
		"length insufficient",
	)
}

func TestErrorResolverHandlesGoodErrorResponse(t *testing.T) {
	caller := &fixedReturnCaller{[]byte{8, 195, 121, 160}}

	// Build a blank error message.
	buildingBlock := [32]byte{}
	caller.returnedBytes = append(caller.returnedBytes, buildingBlock[:]...)
	caller.returnedBytes[len(caller.returnedBytes)-1] = 32 // data offset, fixed
	caller.returnedBytes = append(caller.returnedBytes, buildingBlock[:]...)
	caller.returnedBytes[len(caller.returnedBytes)-1] = 0

	resolver := ethutil.NewErrorResolver(caller, &testABI, &testAddress)
	err := resolver.ResolveError(errOriginal, common.Address{}, nil, "Test")
	assertErrorContains(
		t,
		err,
		errOriginal.Error(),
		"[]",
	)

	// Build an error message.
	errorMessage := "Something's gone awry."
	caller.returnedBytes = []byte{8, 195, 121, 160}
	caller.returnedBytes = append(caller.returnedBytes, buildingBlock[:]...)
	caller.returnedBytes[len(caller.returnedBytes)-1] = 32 // data offset, fixed
	caller.returnedBytes = append(caller.returnedBytes, buildingBlock[:]...)
	caller.returnedBytes[len(caller.returnedBytes)-1] = byte(len(errorMessage))
	caller.returnedBytes = append(caller.returnedBytes, errorMessage[:]...)

	err = resolver.ResolveError(errOriginal, common.Address{}, nil, "Test")
	assertErrorContains(
		t,
		err,
		errOriginal.Error(),
		errorMessage,
	)

}
