package ethereum

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ipfs/go-log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/subscription"
)

var logger = log.Logger("keep-chain-ethereum")

// ThresholdRelay converts from ethereumChain to beacon.ChainInterface.
func (c *Chain) ThresholdRelay() relayChain.Interface {
	return c
}

func (c *Chain) GetKeys() (*operator.PrivateKey, *operator.PublicKey) {
	return operator.ChainKeyToOperatorKey(c.accountKey)
}

func (c *Chain) Signing() chain.Signing {
	return ethutil.NewSigner(c.accountKey.PrivateKey)
}

func (c *Chain) GetConfig() *relayChain.Config {
	return c.chainConfig
}

func (c *Chain) MinimumStake() (*big.Int, error) {
	return c.stakingContract.MinimumStake()
}

// HasMinimumStake returns true if the specified address is staked.  False will
// be returned if not staked.  If err != nil then it was not possible to determine
// if the address is staked or not.
func (c *Chain) HasMinimumStake(address common.Address) (bool, error) {
	return c.keepRandomBeaconOperatorContract.HasMinimumStake(address)
}

func (c *Chain) SubmitTicket(ticket *relayChain.Ticket) *async.EventGroupTicketSubmissionPromise {
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

	ticketBytes := c.packTicket(ticket)

	_, err := c.keepRandomBeaconOperatorContract.SubmitTicket(
		ticketBytes,
		ethutil.TransactionOptions{
			GasLimit: 250000,
		},
	)
	if err != nil {
		failPromise(err)
	}

	// TODO: fulfill when submitted

	return submittedTicketPromise
}

func (c *Chain) packTicket(ticket *relayChain.Ticket) [32]uint8 {
	ticketBytes := []uint8{}
	ticketBytes = append(ticketBytes, ticket.Value[:]...)
	ticketBytes = append(ticketBytes, common.LeftPadBytes(ticket.Proof.StakerValue.Bytes(), 20)[0:20]...)
	ticketBytes = append(ticketBytes, common.LeftPadBytes(ticket.Proof.VirtualStakerIndex.Bytes(), 4)[0:4]...)

	ticketFixedArray := [32]uint8{}
	copy(ticketFixedArray[:], ticketBytes[:32])

	return ticketFixedArray
}

func (c *Chain) GetSubmittedTickets() ([]uint64, error) {
	return c.keepRandomBeaconOperatorContract.SubmittedTickets()
}

func (c *Chain) GetSelectedParticipants() ([]relayChain.StakerAddress, error) {
	var stakerAddresses []relayChain.StakerAddress
	fetchParticipants := func() error {
		participants, err := c.keepRandomBeaconOperatorContract.SelectedParticipants()
		if err != nil {
			return err
		}

		stakerAddresses = make([]relayChain.StakerAddress, len(participants))
		for i, participant := range participants {
			stakerAddresses[i] = participant.Bytes()
		}

		return nil
	}

	// The reason behind a retry functionality is Infura's load balancer synchronization
	// problem. Whenever a Keep client is connected to Infura, it might experience
	// a slight delay with block updates between ethereum clients. One or more
	// clients might stay behind and report a block number 'n-1', whereas the
	// actual block number is already 'n'. This delay results in error triggering
	// a new group selection. To mitigate Infura's sync issue, a Keep client will
	// retry calling for selected participants up to 4 times.
	// Synchronization issue can occur on any setup where we have more than one
	// Ethereum clients behind a load balancer.
	const numberOfRetries = 10
	const delay = time.Second

	for i := 1; ; i++ {
		err := fetchParticipants()
		if err != nil {
			if i == numberOfRetries {
				return nil, err
			}
			time.Sleep(delay)
			logger.Infof(
				"Retrying getting selected participants; attempt [%v]",
				i,
			)
		} else {
			return stakerAddresses, nil
		}
	}
}

func (c *Chain) SubmitRelayEntry(
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

	subscription := c.OnRelayEntrySubmitted(
		func(onChainEvent *event.EntrySubmitted) {
			generatedEntry <- onChainEvent
		},
	)

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

	gasEstimate, err := c.keepRandomBeaconOperatorContract.RelayEntryGasEstimate(entry)
	if err != nil {
		logger.Errorf("failed to estimate gas [%v]", err)
	}

	gasEstimateWithMargin := float64(gasEstimate) * float64(1.2) // 20% more than original
	_, err = c.keepRandomBeaconOperatorContract.RelayEntry(
		entry,
		ethutil.TransactionOptions{
			GasLimit: uint64(gasEstimateWithMargin),
		},
	)
	if err != nil {
		subscription.Unsubscribe()
		close(generatedEntry)
		failPromise(err)
	}

	return relayEntryPromise
}

func (c *Chain) OnRelayEntrySubmitted(
	handle func(entry *event.EntrySubmitted),
) subscription.EventSubscription {
	onEvent := func(blockNumber uint64) {
		handle(&event.EntrySubmitted{
			BlockNumber: blockNumber,
		})
	}

	subscription := c.keepRandomBeaconOperatorContract.RelayEntrySubmitted(
		nil,
	).OnEvent(onEvent)

	return subscription
}

func (c *Chain) OnRelayEntryRequested(
	handle func(request *event.Request),
) subscription.EventSubscription {
	onEvent := func(
		previousEntry []byte,
		groupPublicKey []byte,
		blockNumber uint64,
	) {
		handle(&event.Request{
			PreviousEntry:  previousEntry,
			GroupPublicKey: groupPublicKey,
			BlockNumber:    blockNumber,
		})
	}

	subscription := c.keepRandomBeaconOperatorContract.RelayEntryRequested(
		nil,
	).OnEvent(onEvent)

	return subscription
}

func (c *Chain) OnGroupSelectionStarted(
	handle func(groupSelectionStart *event.GroupSelectionStart),
) subscription.EventSubscription {
	onEvent := func(
		newEntry *big.Int,
		blockNumber uint64,
	) {
		handle(&event.GroupSelectionStart{
			NewEntry:    newEntry,
			BlockNumber: blockNumber,
		})
	}

	subscription := c.keepRandomBeaconOperatorContract.GroupSelectionStarted(
		nil,
	).OnEvent(onEvent)

	return subscription
}

func (c *Chain) OnGroupRegistered(
	handle func(groupRegistration *event.GroupRegistration),
) subscription.EventSubscription {
	onEvent := func(
		memberIndex *big.Int,
		groupPublicKey []byte,
		misbehaved []byte,
		blockNumber uint64,
	) {
		handle(&event.GroupRegistration{
			GroupPublicKey: groupPublicKey,
			BlockNumber:    blockNumber,
		})
	}

	subscription := c.keepRandomBeaconOperatorContract.DkgResultSubmittedEvent(
		nil,
	).OnEvent(onEvent)

	return subscription
}

func (c *Chain) IsGroupRegistered(groupPublicKey []byte) (bool, error) {
	return c.keepRandomBeaconOperatorContract.IsGroupRegistered(groupPublicKey)
}

func (c *Chain) IsStaleGroup(groupPublicKey []byte) (bool, error) {
	return c.keepRandomBeaconOperatorContract.IsStaleGroup(groupPublicKey)
}

func (c *Chain) GetGroupMembers(groupPublicKey []byte) (
	[]relayChain.StakerAddress,
	error,
) {
	members, err := c.keepRandomBeaconOperatorContract.GetGroupMembers(
		groupPublicKey,
	)
	if err != nil {
		return nil, err
	}

	stakerAddresses := make([]relayChain.StakerAddress, len(members))
	for i, member := range members {
		stakerAddresses[i] = member.Bytes()
	}

	return stakerAddresses, nil
}

func (c *Chain) OnDKGResultSubmitted(
	handler func(dkgResultPublication *event.DKGResultSubmission),
) subscription.EventSubscription {
	onEvent := func(
		memberIndex *big.Int,
		groupPublicKey []byte,
		misbehaved []byte,
		blockNumber uint64,
	) {
		handler(&event.DKGResultSubmission{
			MemberIndex:    uint32(memberIndex.Uint64()),
			GroupPublicKey: groupPublicKey,
			Misbehaved:     misbehaved,
			BlockNumber:    blockNumber,
		})
	}

	subscription := c.keepRandomBeaconOperatorContract.DkgResultSubmittedEvent(
		nil,
	).OnEvent(onEvent)

	return subscription
}

func (c *Chain) ReportRelayEntryTimeout() error {
	_, err := c.keepRandomBeaconOperatorContract.ReportRelayEntryTimeout()
	if err != nil {
		return err
	}

	return nil
}

func (c *Chain) IsEntryInProgress() (bool, error) {
	return c.keepRandomBeaconOperatorContract.IsEntryInProgress()
}

func (c *Chain) CurrentRequestStartBlock() (*big.Int, error) {
	return c.keepRandomBeaconOperatorContract.CurrentRequestStartBlock()
}

func (c *Chain) CurrentRequestPreviousEntry() ([]byte, error) {
	return c.keepRandomBeaconOperatorContract.CurrentRequestPreviousEntry()
}

func (c *Chain) CurrentRequestGroupPublicKey() ([]byte, error) {
	currentRequestGroupIndex, err := c.keepRandomBeaconOperatorContract.CurrentRequestGroupIndex()
	if err != nil {
		return nil, err
	}

	return c.keepRandomBeaconOperatorContract.GetGroupPublicKey(currentRequestGroupIndex)
}

func (c *Chain) SubmitDKGResult(
	participantIndex relayChain.GroupMemberIndex,
	result *relayChain.DKGResult,
	signatures map[relayChain.GroupMemberIndex][]byte,
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

	subscription := c.OnDKGResultSubmitted(
		func(onChainEvent *event.DKGResultSubmission) {
			publishedResult <- onChainEvent
		},
	)

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

	if _, err = c.keepRandomBeaconOperatorContract.SubmitDkgResult(
		big.NewInt(int64(participantIndex)),
		result.GroupPublicKey,
		result.Misbehaved,
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
	signatures map[relayChain.GroupMemberIndex][]byte,
) ([]*big.Int, []byte, error) {
	var membersIndices []*big.Int
	var signaturesSlice []byte

	for memberIndex, signature := range signatures {
		if len(signatures[memberIndex]) != ethutil.SignatureSize {
			return nil, nil, fmt.Errorf(
				"invalid signature size for member [%v] got [%d]-bytes but required [%d]-bytes",
				memberIndex,
				len(signatures[memberIndex]),
				ethutil.SignatureSize,
			)
		}
		membersIndices = append(membersIndices, big.NewInt(int64(memberIndex)))
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
func (c *Chain) CalculateDKGResultHash(
	dkgResult *relayChain.DKGResult,
) (relayChain.DKGResultHash, error) {

	// Encode DKG result to the format matched with Solidity keccak256(abi.encodePacked(...))
	hash := crypto.Keccak256(dkgResult.GroupPublicKey, dkgResult.Misbehaved)

	return relayChain.DKGResultHashFromBytes(hash)
}

func (c *Chain) Address() common.Address {
	return c.accountKey.Address
}
