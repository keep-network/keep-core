package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

func (euc *ethereumUtilityChain) Genesis() error {
	// expressed in gas units
	dkgGasEstimate, err := euc.keepRandomBeaconOperatorContract.DkgGasEstimate()
	if err != nil {
		return err
	}

	// expressed in wei
	gasPrice, err := euc.keepRandomBeaconOperatorContract.GasPriceCeiling()
	if err != nil {
		return err
	}

	dkgFee := new(big.Int).Mul(dkgGasEstimate, gasPrice)

	_, err = euc.keepRandomBeaconOperatorContract.Genesis(dkgFee)
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

	onWatchError := func(err error) error {
		promise.Fail(err)
		return err
	}

	// In the rare case relay entry submission happens before relay request in
	// the same block, we need to make sure we install relay entry generated
	// callback after relay entry request tx has been confirmed to do not
	// react on the previous relay entry.
	_, err = euc.keepRandomBeaconServiceContract.WatchRelayEntryRequested(
		func(requestId *big.Int, blockNumber uint64) {
			logger.Infof(
				"Relay request with id [%v] created at block [%v]",
				requestId,
				blockNumber,
			)
			euc.keepRandomBeaconServiceContract.WatchRelayEntryGenerated(
				func(_, entry *big.Int, blockNumber uint64) {
					promise.Fulfill(&event.EntryGenerated{
						Value:       entry,
						BlockNumber: blockNumber,
					})
				},
				onWatchError,
			)
		},
		onWatchError,
	)

	_, err = euc.keepRandomBeaconServiceContract.RequestRelayEntry(
		common.BytesToAddress([]byte{}),
		"",
		callbackGas,
		payment,
	)
	if err != nil {
		promise.Fail(err)
		return promise
	}

	return promise
}
