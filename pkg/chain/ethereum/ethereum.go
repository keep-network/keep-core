package ethereum

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/chaintype"
	"github.com/keep-network/keep-core/pkg/beacon/entry"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

// ThresholdRelay converts from ethereumChain to beacon.ChainInterface.
func (ec *ethereumChain) ThresholdRelay() relaychain.Interface {
	return ec
}

// GetConfig get the GroupSize and Threshold for groups and returns it
// in a becaon.Config struct.
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

// SubmitGroupPublicKey sets up the callback functions for the submission of a
// public key for the group.
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

func (ec *ethereumChain) SubmitRelayEntry(
	entry *relay.Entry,
) *async.RelayEntryPromise {
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

			err := relayEntryPromise.Fulfill(&relay.Entry{
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

	groupSignature := big.NewInt(int64(0)).SetBytes(entry.Value[:])
	_, err = ec.keepRandomBeaconContract.SubmitRelayEntry(
		entry.RequestID,
		entry.GroupID,
		entry.PreviousEntry,
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

// OnRelayEntryRequested registers a callback function for a new
// relay request on chain.
func (ec *ethereumChain) OnRelayEntryRequested(
	handle func(request entry.Request),
) {
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

// AddStaker is a temporary function for Milestone 1 that
// adds a staker to the group contract.
func (ec *ethereumChain) AddStaker(
	groupMemberID string,
) *async.OnStakerAddedPromise {
	onStakerAddedPromise := &async.OnStakerAddedPromise{}

	if len(groupMemberID) != 32 {
		err := onStakerAddedPromise.Fail(
			fmt.Errorf(
				"groupMemberID wrong length, need 32, got %d",
				len(groupMemberID),
			),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Promise Fail failed [%v].\n", err)
		}
		return onStakerAddedPromise
	}

	success := func(
		Index int,
		GroupMemberID []byte,
	) {
		err := onStakerAddedPromise.Fulfill(&chaintype.OnStakerAdded{
			Index:         Index,
			GroupMemberID: string(GroupMemberID),
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Promise Fulfill failed [%v].\n", err)
		}
	}

	fail := func(err error) error {
		err = onStakerAddedPromise.Fail(err)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Promise Fail failed [%v].\n", err)
		}
		return nil
	}

	err := ec.keepGroupContract.WatchOnStakerAdded(
		success,
		fail,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to watch OnStakerAdded [%v].\n", err)
		return onStakerAddedPromise
	}

	index, err := ec.keepGroupContract.GetNStaker()
	if err != nil {
		fmt.Printf("Error: failed on call to GetNStaker: [%v]\n", err)
		return onStakerAddedPromise
	}

	tx, err := ec.keepGroupContract.AddStaker(index, groupMemberID)
	if err != nil {
		fmt.Printf(
			"on staker added failed with: [%v]",
			err,
		)
	}

	ec.tx = tx

	return onStakerAddedPromise
}

// GetStakerList is a temporary function for Milestone 1 that
// gets back the list of stakers.
func (ec *ethereumChain) GetStakerList() ([]string, error) {
	max, err := ec.keepGroupContract.GetNStaker()
	if err != nil {
		err = fmt.Errorf("failed on call to GetNStaker: [%v]", err)
		return []string{}, err
	}

	if max == 0 {
		return []string{}, nil
	}

	listOfStakers := make([]string, 0, max)
	for ii := 0; ii < max; ii++ {
		aStaker, err := ec.keepGroupContract.GetStaker(ii)
		if err != nil {
			return []string{},
				fmt.Errorf("at postion %d out of %d error: [%v]", ii, max, err)
		}
		listOfStakers = append(listOfStakers, string(aStaker))
	}

	return listOfStakers, nil
}
