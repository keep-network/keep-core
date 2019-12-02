package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
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

func (euc *ethereumUtilityChain) RequestRelayEntry() *async.EventEntryGeneratedPromise {
	promise := &async.EventEntryGeneratedPromise{}

	callbackGas := big.NewInt(0) // no callback
	payment, err := euc.keepRandomBeaconServiceContract.EntryFeeEstimate(callbackGas)
	if err != nil {
		promise.Fail(err)
		return promise
	}

	_, err = euc.keepRandomBeaconServiceContract.RequestRelayEntry(
		common.BytesToAddress([]byte{}),
		"",
		callbackGas,
		payment,
	)
	if err != nil {
		promise.Fail(err)
	}

	euc.keepRandomBeaconServiceContract.WatchRelayEntryGenerated(
		func(RequestId *big.Int, Entry *big.Int, blockNumber uint64) {
			promise.Fulfill(&event.EntryGenerated{
				Value:       Entry,
				BlockNumber: blockNumber,
			})
		},
		func(err error) error {
			promise.Fail(err)
			return err
		},
	)

	return promise
}
