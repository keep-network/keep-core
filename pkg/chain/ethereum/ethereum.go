package ethereum

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ipfs/go-log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/internal/byteutils"
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

	ticketInitialSubmissionTimeout, err :=
		ec.keepRandomBeaconOperatorContract.TicketInitialSubmissionTimeout()
	if err != nil {
		return nil, fmt.Errorf(
			"error calling TicketInitialSubmissionTimeout: [%v]",
			err,
		)
	}

	ticketReactiveSubmissionTimeout, err :=
		ec.keepRandomBeaconOperatorContract.TicketReactiveSubmissionTimeout()
	if err != nil {
		return nil, fmt.Errorf(
			"error calling TicketReactiveSubmissionTimeout: [%v]",
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

	tokenSupply, err := ec.keepRandomBeaconOperatorContract.TokenSupply()
	if err != nil {
		return nil, fmt.Errorf("error calling TokenSupply: [%v]", err)
	}

	naturalThreshold, err := ec.keepRandomBeaconOperatorContract.NaturalThreshold()
	if err != nil {
		return nil, fmt.Errorf("error calling NaturalThreshold: [%v]", err)
	}

	relayEntryTimeout, err := ec.keepRandomBeaconOperatorContract.RelayEntryTimeout()
	if err != nil {
		return nil, fmt.Errorf("error calling RelayEntryTimeout: [%v]", err)
	}

	return &relayconfig.Chain{
		GroupSize:                       int(groupSize.Int64()),
		HonestThreshold:                 int(threshold.Int64()),
		TicketInitialSubmissionTimeout:  ticketInitialSubmissionTimeout.Uint64(),
		TicketReactiveSubmissionTimeout: ticketReactiveSubmissionTimeout.Uint64(),
		ResultPublicationBlockStep:      resultPublicationBlockStep.Uint64(),
		MinimumStake:                    minimumStake,
		TokenSupply:                     tokenSupply,
		NaturalThreshold:                naturalThreshold,
		RelayEntryTimeout:               relayEntryTimeout.Uint64(),
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

	_, err := ec.keepRandomBeaconOperatorContract.SubmitTicket(
		ticket.Value,
		ticket.Proof.StakerValue,
		ticket.Proof.VirtualStakerIndex,
	)
	if err != nil {
		failPromise(err)
	}

	// TODO: fulfill when submitted

	return submittedTicketPromise
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
	entryValue *big.Int,
) *async.EventEntryPromise {
	relayEntryPromise := &async.EventEntryPromise{}

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

	generatedEntry := make(chan *event.Entry)

	subscription, err := ec.OnSignatureSubmitted(
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

	_, err = ec.keepRandomBeaconOperatorContract.RelayEntry(entryValue)
	if err != nil {
		subscription.Unsubscribe()
		close(generatedEntry)
		failPromise(err)
	}

	return relayEntryPromise
}

func (ec *ethereumChain) OnSignatureSubmitted(
	handle func(entry *event.Entry),
) (subscription.EventSubscription, error) {
	return ec.keepRandomBeaconOperatorContract.WatchSignatureSubmitted(
		func(
			requestResponse *big.Int,
			requestGroupPubKey []byte,
			previousEntry *big.Int,
			seed *big.Int,
			blockNumber uint64,
		) {
			handle(&event.Entry{
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

func (ec *ethereumChain) OnSignatureRequested(
	handle func(request *event.Request),
) (subscription.EventSubscription, error) {
	return ec.keepRandomBeaconOperatorContract.WatchSignatureRequested(
		func(
			previousEntry *big.Int,
			seed *big.Int,
			groupPublicKey []byte,
			blockNumber uint64,
		) {
			handle(&event.Request{
				PreviousEntry:  previousEntry,
				Seed:           seed,
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

// CombineToSign takes the previous relay entry value and the current
// requests's seed and:
//  - pads them with zeros if their byte length is less than 32 bytes. These
//   values are used later on-chain as `uint256` values and combined with
//   `abi.encodePacked` during signature verification. `uint256` is always
//   packed to 256-bits with leading zeros if needed,
// - combines them into a single slice of bytes.
//
// Function returns an error if previous entry or seed takes more than 32 bytes.
func (ec *ethereumChain) CombineToSign(
	previousEntry *big.Int,
	seed *big.Int,
) ([]byte, error) {
	previousEntryBytes := previousEntry.Bytes()
	seedBytes := seed.Bytes()

	if len(previousEntryBytes) > 32 {
		return nil, fmt.Errorf("entry can not be longer than 32 bytes")
	}
	if len(seedBytes) > 32 {
		return nil, fmt.Errorf("seed can not be longer than 32 bytes")
	}

	previousEntryPadded, err := byteutils.LeftPadTo32Bytes(previousEntryBytes)
	if err != nil {
		return nil, err
	}
	seedPadded, err := byteutils.LeftPadTo32Bytes(seedBytes)
	if err != nil {
		return nil, err
	}

	combinedEntryToSign := make([]byte, 0)
	combinedEntryToSign = append(combinedEntryToSign, previousEntryPadded...)
	combinedEntryToSign = append(combinedEntryToSign, seedPadded...)

	return combinedEntryToSign, nil
}
