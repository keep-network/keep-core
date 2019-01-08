package ethereum

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/groupselection"
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

// HasMinimumStake returns true if the specified address is staked.  False will
// be returned if not staked.  If err != nil then it was not possible to determine
// if the address is staked or not.
func (ec *ethereumChain) HasMinimumStake(address common.Address) (bool, error) {
	return ec.keepRandomBeaconContract.HasMinimumStake(address)
}

func (ec *ethereumChain) SubmitGroupPublicKey(
	requestID *big.Int,
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

	_, err = ec.keepRandomBeaconContract.SubmitGroupPublicKey(key[:], requestID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to submit GroupPublicKey [%v].\n", err)
		return groupRegistrationPromise
	}

	return groupRegistrationPromise
}

// TODO: implement
func (ec *ethereumChain) SubmitTicket(ticket *groupselection.Ticket) *async.GroupTicketPromise {
	return &async.GroupTicketPromise{}
}

// TODO: implement
func (ec *ethereumChain) SubmitChallenge(ticket *groupselection.Challenge) *async.GroupChallengePromise {
	return &async.GroupChallengePromise{}
}

// TODO: implement
func (ec *ethereumChain) GetOrderedTickets() []*groupselection.Ticket {
	return make([]*groupselection.Ticket, 0)
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

			err := relayEntryPromise.Fulfill(&event.Entry{
				RequestID:     requestID,
				Value:         requestResponse,
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

	_, err = ec.keepRandomBeaconContract.SubmitRelayEntry(
		newEntry.RequestID,
		newEntry.GroupID,
		newEntry.PreviousEntry,
		newEntry.Value,
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
			handle(&event.Entry{
				RequestID:     requestID,
				Value:         requestResponse,
				GroupID:       requestGroupID,
				PreviousEntry: previousEntry,
				Timestamp:     time.Now().UTC(),
			})
		},
		func(err error) error {
			fmt.Fprintf(
				os.Stderr,
				"watch relay entry failed with: [%v]",
				err,
			)
			return err
		},
	)
	if err != nil {
		fmt.Printf(
			"watch relay entry failed with: [%v]",
			err,
		)
	}
}

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
		func(index int, stakerID []byte) bool {
			if string(stakerID) == groupMemberID {
				// not what we are waiting for, please continue watching
				return false
			}

			err := onStakerAddedPromise.Fulfill(&event.StakerRegistration{
				Index:         index,
				GroupMemberID: groupMemberID,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Promise Fulfill failed [%v].\n", err)
			}

			return true
		},
		func(err error) error {
			failErr := onStakerAddedPromise.Fail(
				fmt.Errorf(
					"adding new staker failed with: [%v]",
					err,
				),
			)

			if failErr != nil {
				fmt.Fprintf(
					os.Stderr,
					"Staker added promise failing failed: [%v]\n",
					err,
				)
			}

			return failErr
		},
	)
	if err != nil {
		err = onStakerAddedPromise.Fail(
			fmt.Errorf(
				"on staker added failed with: [%v]",
				err,
			),
		)
		if err != nil {
			fmt.Printf("Staker added promise failing failed: [%v]\n", err)
		}

		return onStakerAddedPromise
	}

	_, err = ec.keepGroupContract.AddStaker(groupMemberID)
	if err != nil {
		err = onStakerAddedPromise.Fail(
			fmt.Errorf(
				"on staker added failed with: [%v]",
				err,
			),
		)
		if err != nil {
			fmt.Printf("Staker added promise failing failed: [%v]\n", err)
		}

		return onStakerAddedPromise
	}

	return onStakerAddedPromise
}

func (ec *ethereumChain) OnStakerAdded(
	handle func(staker *event.StakerRegistration),
) {
	err := ec.keepGroupContract.WatchOnStakerAdded(
		func(index int, groupMemberID []byte) bool {
			handle(&event.StakerRegistration{
				Index:         index,
				GroupMemberID: string(groupMemberID),
			})

			// we are not interested in any particular staker,
			// so we stop watching here
			return true
		},
		func(err error) error {
			return fmt.Errorf("staker event failed with %v", err)
		},
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"failed to watch OnStakerAdded [%v].\n",
			err,
		)
	}
}

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
	handle func(registration *event.GroupRegistration),
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

func (ec *ethereumChain) RequestRelayEntry(
	blockReward, seed *big.Int,
) *async.RelayRequestPromise {
	promise := &async.RelayRequestPromise{}
	err := ec.keepRandomBeaconContract.WatchRelayEntryRequested(
		func(
			requestID *big.Int,
			payment *big.Int,
			blockReward *big.Int,
			seed *big.Int,
			blockNumber *big.Int,
		) {
			promise.Fulfill(&event.Request{
				RequestID:   requestID,
				Payment:     payment,
				BlockReward: blockReward,
				Seed:        seed,
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

// IsDKGResultPublished checks if the result is already published to a chain.
func (ec *ethereumChain) IsDKGResultPublished(requestID *big.Int) bool {
	// TODO Implement
	return false
}

// SubmitDKGResult sends DKG result to a chain.
func (ec *ethereumChain) SubmitDKGResult(
	requestID *big.Int, resultToPublish *relaychain.DKGResult,
) *async.DKGResultPublicationPromise {
	// TODO Implement
	return nil
}

func (ec *ethereumChain) OnDKGResultPublished(
	handler func(dkgResultPublication *event.DKGResultPublication),
) event.Subscription {
	// TODO Implement
	return event.NewSubscription(func() {})
}
