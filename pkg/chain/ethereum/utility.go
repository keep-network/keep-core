package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

func (euc *ethereumUtilityChain) Genesis() error {
	// dkgGasEstimate * priceFeedEstimate * fluctuation margin
	// = 2260000 * 20 Gwei * 1.5
	// = 67800000 * 10^9
	genesisPayment := new(big.Int).Mul(
		big.NewInt(67800000),
		new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil),
	)

	_, err := euc.keepRandomBeaconOperatorContract.Genesis(genesisPayment)
	return err
}

func (euc *ethereumUtilityChain) RequestRelayEntry(seed *big.Int) *async.EventRequestPromise {
	promise := &async.EventRequestPromise{}

	callbackGas := big.NewInt(0) // no callback
	payment, err := euc.keepRandomBeaconServiceContract.EntryFeeEstimate(callbackGas)
	if err != nil {
		promise.Fail(err)
		return promise
	}

	_, err = euc.keepRandomBeaconServiceContract.RequestRelayEntry(
		seed,
		common.BytesToAddress([]byte{}),
		"",
		callbackGas,
		payment,
	)
	if err != nil {
		promise.Fail(err)
	}

	return promise
}
