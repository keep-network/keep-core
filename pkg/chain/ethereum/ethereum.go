package ethereum

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/chaintype"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/entry"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

// ThresholdRelay converts from ethereumChain to beacon.ChainInterface.
func (ec *ethereumChain) ThresholdRelay() relaychain.Interface {
	return ec
}

func (ec *ethereumChain) GetConfig() (relayconfig.Chain, error) {
	size, err := ec.keepGroupContract.GroupSize()
	if err != nil {
		return relayconfig.Chain{}, fmt.Errorf("error calling GroupSize: [%v]", err)
	}

	threshold, err := ec.keepGroupContract.GroupThreshold()
	if err != nil {
		return relayconfig.Chain{}, fmt.Errorf("error calling GroupThreshold: [%v]", err)
	}

	return relayconfig.Chain{
		GroupSize: size,
		Threshold: threshold,
	}, nil
}

func (ec *ethereumChain) SubmitGroupPublicKey(
	groupID string,
	key [96]byte,
) *async.GroupPublicKeyPromise {
	groupKeyPromise := &async.GroupPublicKeyPromise{}

	success := func(
		GroupPublicKey []byte,
		RequestID *big.Int,
		ActivationBlockHeight *big.Int,
	) {
		err := groupKeyPromise.Fulfill(&chaintype.GroupPublicKey{
			GroupPublicKey:        GroupPublicKey,
			RequestID:             RequestID,
			ActivationBlockHeight: ActivationBlockHeight,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Promise Fulfill failed [%v].\n", err)
		}
	}

	fail := func(err error) error {
		err = groupKeyPromise.Fail(err)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Promise Fail failed [%v].\n", err)
		}
		return nil
	}

	err := ec.keepRandomBeaconContract.WatchSubmitGroupPublicKeyEvent(
		success,
		fail,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to watch GroupPublicKeyEvent [%v].\n", err)
		return groupKeyPromise
	}

	_, err = ec.keepRandomBeaconContract.SubmitGroupPublicKey(key[:], big.NewInt(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to submit GroupPublicKey [%v].\n", err)
		return groupKeyPromise
	}

	return groupKeyPromise
}

func (ec *ethereumChain) SubmitRelayEntry(newEntry *entry.Entry) *async.RelayEntryPromise {
	relayEntryPromise := &async.RelayEntryPromise{}

	err := ec.keepRandomBeaconContract.WatchRelayEntryGenerated(
		func(
			requestID *big.Int,
			requestResponse *big.Int,
			requestGroupID *big.Int,
			previousEntry *big.Int,
			blockNumber *big.Int,
		) {
			var value [32]byte
			copy(value[:], requestResponse.Bytes()[:32])

			err := relayEntryPromise.Fulfill(&entry.Entry{
				RequestID:     requestID,
				Value:         value,
				GroupID:       requestGroupID,
				PreviousEntry: previousEntry,
				Timestamp:     time.Now().UTC(),
			})
			if err != nil {
				fmt.Printf(
					"execution of fulfilling promise failed with: [%v]",
					err,
				)
			}
		},
		func(err error) error {
			return relayEntryPromise.Fail(
				fmt.Errorf(
					"entry of relay submission failed with: [%v]",
					err,
				),
			)
		},
	)
	if err != nil {
		promiseErr := relayEntryPromise.Fail(
			fmt.Errorf(
				"watch relay entry failed with: [%v]",
				err,
			),
		)
		if promiseErr != nil {
			fmt.Printf(
				"execution of failing promise failed with: [%v]",
				promiseErr,
			)
		}
		return relayEntryPromise
	}

	groupSignature := big.NewInt(int64(0)).SetBytes(newEntry.Value[:])
	_, err = ec.keepRandomBeaconContract.SubmitRelayEntry(
		newEntry.RequestID,
		newEntry.GroupID,
		newEntry.PreviousEntry,
		groupSignature,
	)
	if err != nil {
		promiseErr := relayEntryPromise.Fail(
			fmt.Errorf(
				"submitting relay entry to chain failed with: [%v]",
				err,
			),
		)
		if promiseErr != nil {
			fmt.Printf(
				"execution of failing promise failed with: [%v]",
				promiseErr,
			)
		}
		return relayEntryPromise
	}

	return relayEntryPromise
}

func (ec *ethereumChain) OnRelayEntryGenerated(handle func(entry entry.Entry)) {
	err := ec.keepRandomBeaconContract.WatchRelayEntryGenerated(
		func(
			requestID *big.Int,
			requestResponse *big.Int,
			requestGroupID *big.Int,
			previousEntry *big.Int,
			blockNumber *big.Int,
		) {
			var value [32]byte
			copy(value[:], requestResponse.Bytes()[:32])

			handle(entry.Entry{
				RequestID:     requestID,
				Value:         value,
				GroupID:       requestGroupID,
				PreviousEntry: previousEntry,
				Timestamp:     time.Now().UTC(),
			})
		},
		func(err error) error {
			return fmt.Errorf(
				"watch relay entry failed with: [%v]",
				err,
			)
		},
	)
	if err != nil {
		fmt.Printf(
			"watch relay entry failed with: [%v]",
			err,
		)
	}
}

// OnRelayEntryRequested registers a callback function for a new relay request on
// chain.
func (ec *ethereumChain) OnRelayEntryRequested(handle func(request entry.Request)) {
	err := ec.keepRandomBeaconContract.WatchRelayEntryRequested(
		func(
			requestID *big.Int,
			payment *big.Int,
			blockReward *big.Int,
			seed *big.Int,
			blockNumber *big.Int,
		) {
			handle(entry.Request{
				RequestID:   requestID,
				Payment:     payment,
				BlockReward: blockReward,
				Seed:        seed,
			})
		},
		func(err error) error {
			return fmt.Errorf("relay request event failed with %v", err)
		},
	)
	if err != nil {
		fmt.Printf(
			"watch relay request failed with: [%v]",
			err,
		)
	}
}
