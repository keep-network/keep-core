package ethereum

import (
	"fmt"
	"math/big"
	"os"
	"time"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
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
) *async.GroupRegistrationPromise {
	groupRegistrationPromise := &async.GroupRegistrationPromise{}

	err := ec.keepRandomBeaconContract.WatchSubmitGroupPublicKeyEvent(
		func(
			groupPublicKey []byte,
			requestID *big.Int,
			activationBlockHeight *big.Int,
		) {
			err := groupRegistrationPromise.Fulfill(&event.GroupRegistration{
				GroupPublicKey:        groupPublicKey,
				RequestID:             requestID,
				ActivationBlockHeight: activationBlockHeight,
			})
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"fulfilling promise failed with: [%v].\n",
					err,
				)
			}
		},
		func(err error) error {
			return groupRegistrationPromise.Fail(
				fmt.Errorf(
					"entry of group key failed with: [%v]",
					err,
				),
			)
		},
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"watch group public key event failed with: [%v].\n",
			err,
		)
		return groupRegistrationPromise
	}

	_, err = ec.keepRandomBeaconContract.SubmitGroupPublicKey(key[:], big.NewInt(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to submit GroupPublicKey [%v].\n", err)
		return groupRegistrationPromise
	}

	return groupRegistrationPromise
}

func (ec *ethereumChain) SubmitRelayEntry(
	newEntry *event.Entry,
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

			err := relayEntryPromise.Fulfill(&event.Entry{
				RequestID:     requestID,
				Value:         value,
				GroupID:       requestGroupID,
				PreviousEntry: previousEntry,
				Timestamp:     time.Now().UTC(),
			})
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
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
			fmt.Fprintf(
				os.Stderr,
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

func (ec *ethereumChain) OnRelayEntryGenerated(handle func(entry *event.Entry)) {
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

			handle(&event.Entry{
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

// OnRelayEntryRequested registers a callback function for a new
// relay request on chain.
func (ec *ethereumChain) OnRelayEntryRequested(
	handle func(request *event.Request),
) {
	err := ec.keepRandomBeaconContract.WatchRelayEntryRequested(
		func(
			requestID *big.Int,
			payment *big.Int,
			blockReward *big.Int,
			seed *big.Int,
			blockNumber *big.Int,
		) {
			handle(&event.Request{
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
		fmt.Fprintf(
			os.Stderr,
			"watch relay request failed with: [%v]",
			err,
		)
	}
}

// AddStaker is a temporary function for Milestone 1 that adds a
// staker to the group contract.
func (ec *ethereumChain) AddStaker(
	groupMemberID string,
) *async.StakerRegistrationPromise {
	onStakerAddedPromise := &async.StakerRegistrationPromise{}

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

	err := ec.keepGroupContract.WatchOnStakerAdded(
		func(index int, groupMemberID []byte) {
			err := onStakerAddedPromise.Fulfill(&event.StakerRegistration{
				Index:         index,
				GroupMemberID: string(groupMemberID),
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Promise Fulfill failed [%v].\n", err)
			}
		},
		func(err error) error {
			return onStakerAddedPromise.Fail(
				fmt.Errorf(
					"adding new staker failed with: [%v]",
					err,
				),
			)
		},
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

	_, err = ec.keepGroupContract.AddStaker(index, groupMemberID)
	if err != nil {
		fmt.Printf(
			"on staker added failed with: [%v]",
			err,
		)
	}

	return onStakerAddedPromise
}

// GetStakerList is a temporary function for Milestone 1 that
// gets back the list of stakers.
func (ec *ethereumChain) GetStakerList() ([]string, error) {
	count, err := ec.keepGroupContract.GetNStaker()
	if err != nil {
		err = fmt.Errorf("failed on call to GetNStaker: [%v]", err)
		return nil, err
	}

	listOfStakers := make([]string, 0, count)
	for i := 0; i < count; i++ {
		staker, err := ec.keepGroupContract.GetStaker(i)
		if err != nil {
			return nil,
				fmt.Errorf("at postion %d out of %d error: [%v]", i, count, err)
		}
		listOfStakers = append(listOfStakers, string(staker))
	}

	return listOfStakers, nil
}

func (ec *ethereumChain) OnGroupRegistered(
	handle func(key *event.GroupRegistration),
) {
	err := ec.keepRandomBeaconContract.WatchSubmitGroupPublicKeyEvent(
		func(
			groupPublicKey []byte,
			requestID *big.Int,
			activationBlockHeight *big.Int,
		) {
			handle(&event.GroupRegistration{
				GroupPublicKey:        groupPublicKey,
				RequestID:             requestID,
				ActivationBlockHeight: activationBlockHeight,
			})
		},
		func(err error) error {
			return fmt.Errorf("entry of group key failed with: [%v]", err)
		},
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"watch group public key event failed with: [%v].\n",
			err,
		)
	}
}

// Preliminary - the interface for creating a relay request.

// Will call .../contracts/solidity/contracts/KeepRandomBeaconImplV1.sol;
// function requestRelayEntry(uint256 _blockReward, uint256 _seed) public payable returns (uint256 requestID) {
// through the proxy.
// Note: using 'seed' as a big.Int will loose/drop any leading 0's in the seed.
// Note: Having a return value `returns (uint256 requestID) {` will not return this value through the
//       transaction.  The requestID can be found in the event:
//       RelayEntryRequested(requestID, msg.value, _blockReward, _seed, block.number);
func (ec *ethereumChain) RequestRelayEntry(
	blockReward, seed *big.Int,
) *async.RelayEntryRequestedPromise {
	promise := &async.RelayEntryRequestedPromise{}
	err := ec.keepRandomBeaconContract.WatchRelayEntryRequested(
		func(
			requestID *big.Int,
			payment *big.Int,
			blockReward *big.Int,
			seed *big.Int,
			blockNumber *big.Int,
		) {
			promise.Fulfill(&event.RelayEntryRequested{
				RequestID:   requestID,
				Payment:     payment,
				BlockReward: blockReward,
				Seed:        seed,
				BlockNumber: blockNumber,
			})
		},
		func(err error) error {
			return fmt.Errorf("relay request failed with %v", err)
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "relay request failed with: [%v]", err)
		return promise
	}
	_, err = ec.keepRandomBeaconContract.RequestRelayEntry(blockReward, seed.Bytes())
	if err != nil {
		promise.Fail(fmt.Errorf("failure to request a relay entry: [%v]", err))
		return promise
	}
	return promise
}
