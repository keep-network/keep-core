package ethereum

import (
	"fmt"
	"math/big"

	"github.com/ipfs/go-log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/chain/gen/options"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/subscription"
)

var logger = log.Logger("keep-chain-ethereum")

// ThresholdRelay converts from ethereumChain to beacon.ChainInterface.
func (ec *ethereumChain) ThresholdRelay() relaychain.Interface {
	return ec
}

func (ec *ethereumChain) GetKeys() (*operator.PrivateKey, *operator.PublicKey) {
	return operator.EthereumKeyToOperatorKey(ec.accountKey)
}

func (ec *ethereumChain) GetConfig() (*relayconfig.Chain, error) {
	groupSize, err := ec.keepRandomBeaconOperatorContract.GroupSize()
	if err != nil {
		return nil, fmt.Errorf("error calling GroupSize: [%v]", err)
	}

	threshold, err := ec.keepRandomBeaconOperatorContract.GroupThreshold()
	if err != nil {
		return nil, fmt.Errorf("error calling GroupThreshold: [%v]", err)
	}

	ticketSubmissionTimeout, err :=
		ec.keepRandomBeaconOperatorContract.TicketSubmissionTimeout()
	if err != nil {
		return nil, fmt.Errorf(
			"error calling TicketSubmissionTimeout: [%v]",
			err,
		)
	}

	resultPublicationBlockStep, err := ec.keepRandomBeaconOperatorContract.ResultPublicationBlockStep()
	if err != nil {
		return nil, fmt.Errorf(
			"error calling ResultPublicationBlockStep: [%v]",
			err,
		)
	}

	minimumStake, err := ec.keepRandomBeaconOperatorContract.MinimumStake()
	if err != nil {
		return nil, fmt.Errorf("error calling MinimumStake: [%v]", err)
	}

	relayEntryTimeout, err := ec.keepRandomBeaconOperatorContract.RelayEntryTimeout()
	if err != nil {
		return nil, fmt.Errorf("error calling RelayEntryTimeout: [%v]", err)
	}

	return &relayconfig.Chain{
		GroupSize:                  int(groupSize.Int64()),
		HonestThreshold:            int(threshold.Int64()),
		TicketSubmissionTimeout:    ticketSubmissionTimeout.Uint64(),
		ResultPublicationBlockStep: resultPublicationBlockStep.Uint64(),
		MinimumStake:               minimumStake,
		RelayEntryTimeout:          relayEntryTimeout.Uint64(),
	}, nil
}

// HasMinimumStake returns true if the specified address is staked.  False will
// be returned if not staked.  If err != nil then it was not possible to determine
// if the address is staked or not.
func (ec *ethereumChain) HasMinimumStake(address common.Address) (bool, error) {
	return ec.keepRandomBeaconOperatorContract.HasMinimumStake(address)
}

func (ec *ethereumChain) SubmitTicket(ticket *chain.Ticket) *async.EventGroupTicketSubmissionPromise {
	submittedTicketPromise := &async.EventGroupTicketSubmissionPromise{}

	failPromise := func(err error) {
		failErr := submittedTicketPromise.Fail(err)
		if failErr != nil {
			logger.Errorf(
				"failing promise because of: [%v] failed with: [%v].",
				err,
				failErr,
			)
		}
	}

	ticketBytes := ec.packTicket(ticket)

	_, err := ec.keepRandomBeaconOperatorContract.SubmitTicket(
		ticketBytes,
		options.TransactionOptions{
			GasLimit: 250000,
		},
	)
	if err != nil {
		failPromise(err)
	}

	// TODO: fulfill when submitted

	return submittedTicketPromise
}

func (ec *ethereumChain) packTicket(ticket *relaychain.Ticket) [32]uint8 {
	ticketBytes := []uint8{}
	ticketBytes = append(ticketBytes, common.LeftPadBytes(ticket.Value.Bytes(), 32)[:8]...)
	ticketBytes = append(ticketBytes, ticket.Proof.StakerValue.Bytes()[0:20]...)
	ticketBytes = append(ticketBytes, common.LeftPadBytes(ticket.Proof.VirtualStakerIndex.Bytes(), 4)[0:4]...)

	ticketFixedArray := [32]uint8{}
	copy(ticketFixedArray[:], ticketBytes[:32])

	return ticketFixedArray
}

func (ec *ethereumChain) GetSubmittedTicketsCount() (*big.Int, error) {
	return ec.keepRandomBeaconOperatorContract.SubmittedTicketsCount()
}

func (ec *ethereumChain) GetSelectedParticipants() (
	[]chain.StakerAddress,
	error,
) {
	selectedParticipants, err := ec.keepRandomBeaconOperatorContract.SelectedParticipants()
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
	entry []byte,
) *async.EventEntrySubmittedPromise {
	relayEntryPromise := &async.EventEntrySubmittedPromise{}

	failPromise := func(err error) {
		failErr := relayEntryPromise.Fail(err)
		if failErr != nil {
			logger.Errorf(
				"failed to fail promise for [%v]: [%v]",
				err,
				failErr,
			)
		}
	}

	generatedEntry := make(chan *event.EntrySubmitted)

	subscription, err := ec.OnRelayEntrySubmitted(
		func(onChainEvent *event.EntrySubmitted) {
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

				subscription.Unsubscribe()
				close(generatedEntry)

				err := relayEntryPromise.Fulfill(event)
				if err != nil {
					logger.Errorf(
						"failed to fulfill promise: [%v]",
						err,
					)
				}

				return
			}
		}
	}()

	_, err = ec.keepRandomBeaconOperatorContract.RelayEntry(entry)
	if err != nil {
		subscription.Unsubscribe()
		close(generatedEntry)
		failPromise(err)
	}

	return relayEntryPromise
}

func (ec *ethereumChain) OnRelayEntrySubmitted(
	handle func(entry *event.EntrySubmitted),
) (subscription.EventSubscription, error) {
	return ec.keepRandomBeaconOperatorContract.WatchRelayEntrySubmitted(
		func(blockNumber uint64) {
			handle(&event.EntrySubmitted{
				BlockNumber: blockNumber,
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
	return ec.keepRandomBeaconOperatorContract.WatchRelayEntryRequested(
		func(
			previousEntry []byte,
			groupPublicKey []byte,
			blockNumber uint64,
		) {
			handle(&event.Request{
				PreviousEntry:  previousEntry,
				GroupPublicKey: groupPublicKey,
				BlockNumber:    blockNumber,
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

func (ec *ethereumChain) OnGroupSelectionStarted(
	handle func(groupSelectionStart *event.GroupSelectionStart),
) (subscription.EventSubscription, error) {
	return ec.keepRandomBeaconOperatorContract.WatchGroupSelectionStarted(
		func(
			newEntry *big.Int,
			blockNumber uint64,
		) {
			handle(&event.GroupSelectionStart{
				NewEntry:    newEntry,
				BlockNumber: blockNumber,
			})
		},
		func(err error) error {
			return fmt.Errorf(
				"watch group selection started failed with [%v]",
				err,
			)
		},
	)
}

func (ec *ethereumChain) OnGroupRegistered(
	handle func(groupRegistration *event.GroupRegistration),
) (subscription.EventSubscription, error) {
	return ec.keepRandomBeaconOperatorContract.WatchDkgResultPublishedEvent(
		func(
			groupPublicKey []byte,
			blockNumber uint64,
		) {
			handle(&event.GroupRegistration{
				GroupPublicKey: groupPublicKey,
				BlockNumber:    blockNumber,
			})
		},
		func(err error) error {
			return fmt.Errorf("entry of group key failed with: [%v]", err)
		},
	)
}

func (ec *ethereumChain) IsGroupRegistered(groupPublicKey []byte) (bool, error) {
	return ec.keepRandomBeaconOperatorContract.IsGroupRegistered(groupPublicKey)
}

func (ec *ethereumChain) IsStaleGroup(groupPublicKey []byte) (bool, error) {
	return ec.keepRandomBeaconOperatorContract.IsStaleGroup(groupPublicKey)
}

func (ec *ethereumChain) OnDKGResultSubmitted(
	handler func(dkgResultPublication *event.DKGResultSubmission),
) (subscription.EventSubscription, error) {
	return ec.keepRandomBeaconOperatorContract.WatchDkgResultPublishedEvent(
		func(groupPubKey []byte, blockNumber uint64) {
			handler(&event.DKGResultSubmission{
				GroupPublicKey: groupPubKey,
				BlockNumber:    blockNumber,
			})
		},
		func(err error) error {
			return fmt.Errorf(
				"watch DKG result published failed with: [%v]",
				err,
			)
		},
	)
}

func (ec *ethereumChain) ReportRelayEntryTimeout() error {
	_, err := ec.keepRandomBeaconOperatorContract.ReportRelayEntryTimeout()
	if err != nil {
		return err
	}

	return nil
}

func (ec *ethereumChain) SubmitDKGResult(
	participantIndex group.MemberIndex,
	result *relaychain.DKGResult,
	signatures map[group.MemberIndex][]byte,
) *async.EventDKGResultSubmissionPromise {
	resultPublicationPromise := &async.EventDKGResultSubmissionPromise{}

	failPromise := func(err error) {
		failErr := resultPublicationPromise.Fail(err)
		if failErr != nil {
			logger.Errorf(
				"failed to fail promise for [%v]: [%v]",
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

				subscription.Unsubscribe()
				close(publishedResult)

				err := resultPublicationPromise.Fulfill(event)
				if err != nil {
					logger.Errorf(
						"failed to fulfill promise: [%v]",
						err,
					)
				}

				return
			}
		}
	}()

	membersIndicesOnChainFormat, signaturesOnChainFormat, err :=
		convertSignaturesToChainFormat(signatures)
	if err != nil {
		close(publishedResult)
		failPromise(fmt.Errorf("converting signatures failed [%v]", err))
		return resultPublicationPromise
	}

	if _, err = ec.keepRandomBeaconOperatorContract.SubmitDkgResult(
		participantIndex.Int(),
		result.GroupPublicKey,
		result.Disqualified,
		result.Inactive,
		signaturesOnChainFormat,
		membersIndicesOnChainFormat,
	); err != nil {
		subscription.Unsubscribe()
		close(publishedResult)
		failPromise(err)
	}

	return resultPublicationPromise
}

// convertSignaturesToChainFormat converts signatures map to two slices. First
// slice contains indices of members from the map, second slice is a slice of
// concatenated signatures. Signatures and member indices are returned in the
// matching order. It requires each signature to be exactly 65-byte long.
func convertSignaturesToChainFormat(
	signatures map[group.MemberIndex][]byte,
) ([]*big.Int, []byte, error) {
	var membersIndices []*big.Int
	var signaturesSlice []byte

	for memberIndex, signature := range signatures {
		if len(signatures[memberIndex]) != SignatureSize {
			return nil, nil, fmt.Errorf(
				"invalid signature size for member [%v] got [%d]-bytes but required [%d]-bytes",
				memberIndex,
				len(signatures[memberIndex]),
				SignatureSize,
			)
		}
		membersIndices = append(membersIndices, memberIndex.Int())
		signaturesSlice = append(signaturesSlice, signature...)
	}

	return membersIndices, signaturesSlice, nil
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

	// Encode DKG result to the format matched with Solidity keccak256(abi.encodePacked(...))
	hash := crypto.Keccak256(dkgResult.GroupPublicKey, dkgResult.Disqualified, dkgResult.Inactive)

	return relaychain.DKGResultHashFromBytes(hash)
}
