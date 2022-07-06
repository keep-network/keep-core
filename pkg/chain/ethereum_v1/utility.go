package ethereum_v1

import (
	"math/big"

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
		failErr := promise.Fail(err)
		if failErr != nil {
			logger.Errorf("could not fail the promise: [%v]", failErr)
		}

		return promise
	}

	// In the rare case relay entry submission happens before relay request in
	// the same block, we need to make sure we install relay entry generated
	// callback after relay entry request tx has been confirmed to do not
	// react on the previous relay entry.
	_ = euc.keepRandomBeaconServiceContract.RelayEntryRequested(nil).OnEvent(
		func(requestId *big.Int, blockNumber uint64) {
			logger.Infof(
				"Relay request with id [%v] created at block [%v]",
				requestId,
				blockNumber,
			)
			_ = euc.keepRandomBeaconServiceContract.RelayEntryGenerated(nil).OnEvent(
				func(_, entry *big.Int, blockNumber uint64) {
					fulfillErr := promise.Fulfill(&event.EntryGenerated{
						Value:       entry,
						BlockNumber: blockNumber,
					})
					if fulfillErr != nil {
						logger.Errorf("could not fulfill the promise: [%v]", fulfillErr)
					}
				},
			)
		},
	)

	_, err = euc.keepRandomBeaconServiceContract.RequestRelayEntry(payment)
	if err != nil {
		failErr := promise.Fail(err)
		if failErr != nil {
			logger.Errorf("could not fail the promise: [%v]", failErr)
		}

		return promise
	}

	return promise
}
