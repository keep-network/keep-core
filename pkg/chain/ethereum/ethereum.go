package ethereum

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/subscription"
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
	groupPublicKey []byte,
) *async.GroupRegistrationPromise {
	groupRegistrationPromise := &async.GroupRegistrationPromise{}

	failPromise := func(err error) {
		failErr := groupRegistrationPromise.Fail(err)
		if failErr != nil {
			fmt.Fprintf(
				os.Stderr,
				"failing promise because of [%v] failed with [%v]\n",
				err,
				failErr,
			)
		}
	}

	groupRegistered := make(chan *event.GroupRegistration)

	subscription, err := ec.OnGroupRegistered(
		func(groupRegistrationEvent *event.GroupRegistration) {
			groupRegistered <- groupRegistrationEvent
		})
	if err != nil {
		close(groupRegistered)
		failPromise(err)
		return groupRegistrationPromise
	}

	go func() {
		for {
			select {
			case event, success := <-groupRegistered:
				// Channel is closed when SubmitGroupPublicKey failed.
				// When this happens, event is nil.
				if !success {
					return
				}

				if event.RequestID.Cmp(requestID) == 0 {
					subscription.Unsubscribe()
					close(groupRegistered)

					err := groupRegistrationPromise.Fulfill(event)
					if err != nil {
						fmt.Fprintf(
							os.Stderr,
							"fulfilling promise failed with [%v]\n",
							err,
						)
					}

					return
				}
			}
		}
	}()

	_, err = ec.keepGroupContract.SubmitGroupPublicKey(groupPublicKey, requestID)
	if err != nil {
		subscription.Unsubscribe()
		close(groupRegistered)
		failPromise(err)
	}

	return groupRegistrationPromise
}

func (ec *ethereumChain) SubmitTicket(ticket *chain.Ticket) *async.GroupTicketPromise {
	submittedTicketPromise := &async.GroupTicketPromise{}

	failPromise := func(err error) {
		failErr := submittedTicketPromise.Fail(err)
		if failErr != nil {
			fmt.Fprintf(
				os.Stderr,
				"failing promise because of: [%v] failed with: [%v].\n",
				err,
				failErr,
			)
		}
	}

	_, err := ec.keepGroupContract.SubmitTicket(ticket)
	if err != nil {
		failPromise(err)
	}

	// TODO: fulfill when submitted

	return submittedTicketPromise
}

func (ec *ethereumChain) GetSelectedParticipants() (
	[]chain.StakerAddress,
	error,
) {
	selectedParticipants, err := ec.keepGroupContract.SelectedParticipants()
	if err != nil {
		return nil, err
	}

	stakerAddresses := make([]chain.StakerAddress, len(selectedParticipants))
	for i, selectedParticipant := range selectedParticipants {
		stakerAddresses[i] = selectedParticipant.Bytes()
	}

	return stakerAddresses, nil
}

func (ec *ethereumChain) SubmitRelayEntry(
	newEntry *event.Entry,
) *async.RelayEntryPromise {
	relayEntryPromise := &async.RelayEntryPromise{}

	failPromise := func(err error) {
		failErr := relayEntryPromise.Fail(err)
		if failErr != nil {
			fmt.Fprintf(
				os.Stderr,
				"failing promise because of [%v] failed with [%v]\n",
				err,
				failErr,
			)
		}
	}

	generatedEntry := make(chan *event.Entry)

	subscription, err := ec.OnRelayEntryGenerated(
		func(onChainEvent *event.Entry) {
			generatedEntry <- onChainEvent
		},
	)
	if err != nil {
		close(generatedEntry)
		failPromise(err)
		return relayEntryPromise
	}

	go func() {
		for {
			select {
			case event, success := <-generatedEntry:
				// Channel is closed when SubmitRelayEntry failed.
				// When this happens, event is nil.
				if !success {
					return
				}

				if event.RequestID.Cmp(newEntry.RequestID) == 0 {
					subscription.Unsubscribe()
					close(generatedEntry)

					err := relayEntryPromise.Fulfill(event)
					if err != nil {
						fmt.Fprintf(
							os.Stderr,
							"fulfilling promise failed with [%v]\n",
							err,
						)
					}

					return
				}
			}
		}
	}()

	_, err = ec.keepRandomBeaconContract.SubmitRelayEntry(
		newEntry.RequestID,
		newEntry.GroupPubKey,
		newEntry.PreviousEntry,
		newEntry.Value,
		newEntry.Seed,
	)
	if err != nil {
		subscription.Unsubscribe()
		close(generatedEntry)
		failPromise(err)
	}

	return relayEntryPromise
}

func (ec *ethereumChain) OnRelayEntryGenerated(
	handle func(entry *event.Entry),
) (subscription.EventSubscription, error) {
	return ec.keepRandomBeaconContract.WatchRelayEntryGenerated(
		func(
			requestID *big.Int,
			requestResponse *big.Int,
			requestGroupPubKey []byte,
			previousEntry *big.Int,
			seed *big.Int,
			blockNumber uint64,
		) {
			handle(&event.Entry{
				RequestID:     requestID,
				Value:         requestResponse,
				GroupPubKey:   requestGroupPubKey,
				PreviousEntry: previousEntry,
				Timestamp:     time.Now().UTC(),
				Seed:          seed,
				BlockNumber:   blockNumber,
			})
		},
		func(err error) error {
			return fmt.Errorf(
				"watch relay entry generated failed with [%v]",
				err,
			)
		},
	)
}

func (ec *ethereumChain) OnRelayEntryRequested(
	handle func(request *event.Request),
) (subscription.EventSubscription, error) {
	return ec.keepRandomBeaconContract.WatchRelayEntryRequested(
		func(
			requestID *big.Int,
			payment *big.Int,
			previousEntry *big.Int,
			seed *big.Int,
			blockNumber uint64,
		) {
			handle(&event.Request{
				RequestID:     requestID,
				Payment:       payment,
				PreviousEntry: previousEntry,
				Seed:          seed,
				BlockNumber:   blockNumber,
			})
		},
		func(err error) error {
			return fmt.Errorf(
				"watch relay entry requested failed with [%v]",
				err,
			)
		},
	)
}

func (ec *ethereumChain) OnGroupRegistered(
	handle func(groupRegistration *event.GroupRegistration),
) (subscription.EventSubscription, error) {
	return ec.keepGroupContract.WatchSubmitGroupPublicKeyEvent(
		func(
			groupPublicKey []byte,
			requestID *big.Int,
			blockNumber uint64,
		) {
			handle(&event.GroupRegistration{
				GroupPublicKey: groupPublicKey,
				RequestID:      requestID,
				BlockNumber:    blockNumber,
			})
		},
		func(err error) error {
			return fmt.Errorf("entry of group key failed with: [%v]", err)
		},
	)
}

func (ec *ethereumChain) RequestRelayEntry(
	seed *big.Int,
) *async.RelayRequestPromise {
	relayRequestPromise := &async.RelayRequestPromise{}

	failPromise := func(err error) {
		failErr := relayRequestPromise.Fail(err)
		if failErr != nil {
			fmt.Fprintf(
				os.Stderr,
				"failing promise because of [%v] failed with [%v]\n",
				err,
				failErr,
			)
		}
	}

	requestedEntry := make(chan *event.Request)

	subscription, err := ec.OnRelayEntryRequested(
		func(onChainEvent *event.Request) {
			requestedEntry <- onChainEvent
		},
	)
	if err != nil {
		close(requestedEntry)
		failPromise(err)
		return relayRequestPromise
	}

	go func() {
		for {
			select {
			case event, success := <-requestedEntry:
				// Channel is closed when RequestRelayEntry failed.
				// When this happens, event is nil.
				if !success {
					return
				}

				subscription.Unsubscribe()
				close(requestedEntry)

				err := relayRequestPromise.Fulfill(event)
				if err != nil {
					fmt.Fprintf(
						os.Stderr,
						"fulfilling promise failed with [%v]\n",
						err,
					)
				}

				return
			}
		}
	}()

	_, err = ec.keepRandomBeaconContract.RequestRelayEntry(seed.Bytes())
	if err != nil {
		subscription.Unsubscribe()
		close(requestedEntry)
		failPromise(err)
	}

	return relayRequestPromise
}

func (ec *ethereumChain) IsDKGResultSubmitted(requestID *big.Int) (bool, error) {
	return ec.keepGroupContract.IsDkgResultSubmitted(requestID)
}

func (ec *ethereumChain) OnDKGResultSubmitted(
	handler func(dkgResultPublication *event.DKGResultSubmission),
) (subscription.EventSubscription, error) {
	return ec.keepGroupContract.WatchDKGResultPublishedEvent(
		func(requestID *big.Int, groupPubKey []byte, blockNumber uint64) {
			handler(&event.DKGResultSubmission{
				RequestID:      requestID,
				GroupPublicKey: groupPubKey,
				BlockNumber:    blockNumber,
			})
		},
		func(err error) error {
			return fmt.Errorf(
				"watch DKG result published failed with [%v]",
				err,
			)
		},
	)
}

func (ec *ethereumChain) SubmitDKGResult(
	requestID *big.Int,
	participantIndex uint32,
	result *relaychain.DKGResult,
	signatures map[gjkr.MemberID]operator.Signature,
) *async.DKGResultSubmissionPromise {
	resultPublicationPromise := &async.DKGResultSubmissionPromise{}

	failPromise := func(err error) {
		failErr := resultPublicationPromise.Fail(err)
		if failErr != nil {
			fmt.Fprintf(
				os.Stderr,
				"failing promise because of [%v] failed with [%v]\n",
				err,
				failErr,
			)
		}
	}

	publishedResult := make(chan *event.DKGResultSubmission)

	subscription, err := ec.OnDKGResultSubmitted(
		func(onChainEvent *event.DKGResultSubmission) {
			publishedResult <- onChainEvent
		},
	)
	if err != nil {
		close(publishedResult)
		failPromise(err)
		return resultPublicationPromise
	}

	go func() {
		for {
			select {
			case event, success := <-publishedResult:
				// Channel is closed when SubmitDKGResult failed.
				// When this happens, event is nil.
				if !success {
					return
				}

				if event.RequestID.Cmp(requestID) == 0 {
					subscription.Unsubscribe()
					close(publishedResult)

					err := resultPublicationPromise.Fulfill(event)
					if err != nil {
						fmt.Fprintf(
							os.Stderr,
							"fulfilling promise failed with [%v]\n",
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
		close(publishedResult)
		failPromise(err)
	}

	return resultPublicationPromise
}

// CalculateDKGResultHash calculates Keccak-256 hash of the DKG result. Operation
// is performed off-chain.
//
// It first encodes the result using solidity ABI and then calculates Keccak-256
// hash over it. This corresponds to the DKG result hash calculation on-chain.
// Hashes calculated off-chain and on-chain must always match.
func (ec *ethereumChain) CalculateDKGResultHash(
	dkgResult *relaychain.DKGResult,
) (relaychain.DKGResultHash, error) {
	dkgResultHash := relaychain.DKGResultHash{}

	// Encode DKG result to the format described by Solidity Contract Application
	// Binary Interface (ABI).
	boolType, err := abi.NewType("bool")
	if err != nil {
		return dkgResultHash, fmt.Errorf("bool type creation failed: [%v]", err)
	}
	bytesType, err := abi.NewType("bytes")
	if err != nil {
		return dkgResultHash, fmt.Errorf("bytes type creation failed: [%v]", err)
	}

	arguments := abi.Arguments{
		{Type: boolType},
		{Type: bytesType},
		{Type: bytesType},
		{Type: bytesType},
	}

	encodedDKGResult, err := arguments.Pack(
		dkgResult.Success,
		dkgResult.GroupPublicKey,
		dkgResult.Disqualified,
		dkgResult.Inactive,
	)
	if err != nil {
		return dkgResultHash, fmt.Errorf("encoding failed: [%v]", err)
	}

	// Calculate Keccak-256 hash.
	hash := sha3.NewKeccak256()
	hash.Write(encodedDKGResult)
	hash.Sum(dkgResultHash[:0])

	return dkgResultHash, nil
}
