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

func (ec *ethereumChain) GetConfig() (*relayconfig.Chain, error) {
	groupSize, err := ec.keepGroupContract.GroupSize()
	if err != nil {
		return nil, fmt.Errorf("error calling GroupSize: [%v]", err)
	}

	threshold, err := ec.keepGroupContract.GroupThreshold()
	if err != nil {
		return nil, fmt.Errorf("error calling GroupThreshold: [%v]", err)
	}

	ticketInitialSubmissionTimeout, err :=
		ec.keepGroupContract.TicketInitialSubmissionTimeout()
	if err != nil {
		return nil, fmt.Errorf(
			"error calling TicketInitialSubmissionTimeout: [%v]",
			err,
		)
	}

	ticketReactiveSubmissionTimeout, err :=
		ec.keepGroupContract.TicketReactiveSubmissionTimeout()
	if err != nil {
		return nil, fmt.Errorf(
			"error calling TicketReactiveSubmissionTimeout: [%v]",
			err,
		)
	}

	ticketChallengeTimeout, err :=
		ec.keepGroupContract.TicketChallengeTimeout()
	if err != nil {
		return nil, fmt.Errorf(
			"error calling TicketChallengeTimeout: [%v]",
			err,
		)
	}

	minimumStake, err := ec.keepGroupContract.MinimumStake()
	if err != nil {
		return nil, fmt.Errorf("error calling MinimumStake: [%v]", err)
	}

	tokenSupply, err := ec.keepGroupContract.TokenSupply()
	if err != nil {
		return nil, fmt.Errorf("error calling TokenSupply: [%v]", err)
	}

	naturalThreshold, err := ec.keepGroupContract.NaturalThreshold()
	if err != nil {
		return nil, fmt.Errorf("error calling NaturalThreshold: [%v]", err)
	}

	return &relayconfig.Chain{
		GroupSize:                       groupSize,
		Threshold:                       threshold,
		TicketInitialSubmissionTimeout:  ticketInitialSubmissionTimeout,
		TicketReactiveSubmissionTimeout: ticketReactiveSubmissionTimeout,
		TicketChallengeTimeout:          ticketChallengeTimeout,
		MinimumStake:                    minimumStake,
		TokenSupply:                     tokenSupply,
		NaturalThreshold:                naturalThreshold,
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

func (ec *ethereumChain) IsDKGResultPublished(requestID *big.Int) (bool, error) {
	return ec.keepGroupContract.IsDkgResultSubmitted(requestID)
}

func (ec *ethereumChain) OnDKGResultPublished(
	handler func(dkgResultPublication *event.DKGResultPublication),
) (event.Subscription, error) {
	return ec.keepGroupContract.WatchDKGResultPublishedEvent(
		func(requestID *big.Int) {
			handler(&event.DKGResultPublication{RequestID: requestID})
		},
		func(err error) error {
			return err
		},
	)
}

func (ec *ethereumChain) SubmitDKGResult(
	requestID *big.Int,
	result *relaychain.DKGResult,
) *async.DKGResultPublicationPromise {
	resultPublicationPromise := &async.DKGResultPublicationPromise{}

	failPromise := func(err error) {
		failErr := resultPublicationPromise.Fail(err)
		if failErr != nil {
			fmt.Fprintf(
				os.Stderr,
				"failing promise because of: [%v] failed with: [%v].\n",
				err,
				failErr,
			)
		}
	}

	published := make(chan *event.DKGResultPublication)

	subscription, err := ec.OnDKGResultPublished(
		func(event *event.DKGResultPublication) {
			published <- event
		},
	)
	if err != nil {
		close(published)
		failPromise(err)
		return resultPublicationPromise
	}

	go func() {
		for {
			select {
			case event := <-published:
				if event.RequestID == requestID {
					subscription.Unsubscribe()
					close(published)

					err := resultPublicationPromise.Fulfill(event)
					if err != nil {
						fmt.Fprintf(
							os.Stderr,
							"fulfilling promise failed with: [%v].\n",
							err,
						)
					}

					return
				}
			}
		}
	}()

	_, err = ec.keepGroupContract.SubmitDKGResult(requestID, result)
	if err != nil {
		subscription.Unsubscribe()
		close(published)
		failPromise(err)
	}

	return resultPublicationPromise
}
