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
	"github.com/keep-network/keep-core/pkg/subscription"
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
	return ec.keepGroupContract.HasMinimumStake(address)
}

func (ec *ethereumChain) SubmitGroupPublicKey(
	requestID *big.Int,
	key [96]byte,
) *async.GroupRegistrationPromise {
	groupRegistrationPromise := &async.GroupRegistrationPromise{}

	subscription, err := ec.keepRandomBeaconContract.WatchSubmitGroupPublicKeyEvent(
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
		subscription.Unsubscribe()
		fmt.Fprintf(
			os.Stderr,
			"watch group public key event failed with: [%v].\n",
			err,
		)
		return groupRegistrationPromise
	}

	_, err = ec.keepRandomBeaconContract.SubmitGroupPublicKey(key[:], requestID)
	if err != nil {
		subscription.Unsubscribe()
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
func (ec *ethereumChain) SubmitChallenge(
	ticket *groupselection.TicketChallenge,
) *async.GroupTicketChallengePromise {
	return &async.GroupTicketChallengePromise{}
}

// TODO: implement
func (ec *ethereumChain) GetOrderedTickets() []*groupselection.Ticket {
	return make([]*groupselection.Ticket, 0)
}

func (ec *ethereumChain) SubmitRelayEntry(
	newEntry *event.Entry,
) *async.RelayEntryPromise {
	relayEntryPromise := &async.RelayEntryPromise{}

	subscription, err := ec.keepRandomBeaconContract.WatchRelayEntryGenerated(
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
		subscription.Unsubscribe()
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
		subscription.Unsubscribe()
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

func (ec *ethereumChain) OnRelayEntryGenerated(
	handle func(entry *event.Entry),
) subscription.EventSubscription {
	relayEntryGeneratedSubscription, err := ec.keepRandomBeaconContract.WatchRelayEntryGenerated(
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

	return relayEntryGeneratedSubscription
}

func (ec *ethereumChain) OnRelayEntryRequested(
	handle func(request *event.Request),
) subscription.EventSubscription {
	relayEntryRequestedSubscription, err := ec.keepRandomBeaconContract.WatchRelayEntryRequested(
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

	return relayEntryRequestedSubscription
}

func (ec *ethereumChain) OnGroupRegistered(
	handle func(registration *event.GroupRegistration),
) subscription.EventSubscription {
	groupRegisteredSubscription, err :=
		ec.keepRandomBeaconContract.WatchSubmitGroupPublicKeyEvent(
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
	return groupRegisteredSubscription
}

func (ec *ethereumChain) RequestRelayEntry(
	blockReward, seed *big.Int,
) *async.RelayRequestPromise {
	promise := &async.RelayRequestPromise{}
	subscription, err := ec.keepRandomBeaconContract.WatchRelayEntryRequested(
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
		subscription.Unsubscribe()
		fmt.Fprintf(os.Stderr, "relay request failed with: [%v]", err)
		return promise
	}
	_, err = ec.keepRandomBeaconContract.RequestRelayEntry(blockReward, seed.Bytes())
	if err != nil {
		subscription.Unsubscribe()
		promise.Fail(fmt.Errorf("failure to request a relay entry: [%v]", err))
		return promise
	}
	return promise
}

// IsDKGResultPublished checks if the result is already published to a chain.
func (ec *ethereumChain) IsDKGResultPublished(requestID *big.Int) (bool, error) {
	return false, nil
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
) subscription.EventSubscription {
	// TODO Implement
	return subscription.NewEventSubscription(func() {})
}
