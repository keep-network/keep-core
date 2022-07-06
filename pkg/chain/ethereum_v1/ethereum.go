package ethereum_v1

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"

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
func (ec *ethereumChain) ThresholdRelay() relayChain.Interface {
	return ec
}

func (ec *ethereumChain) GetKeys() (
	*operator.PrivateKey,
	*operator.PublicKey,
	error,
) {
	privateKey, publicKey, err := ChainPrivateKeyToOperatorKeyPair(
		ec.accountKey.PrivateKey,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"cannot convert chain private key to operator key pair: [%v]",
			err,
		)
	}

	return privateKey, publicKey, nil
}

func (ec *ethereumChain) Signing() chain.Signing {
	return newSigner(ec.accountKey)
}

func (ec *ethereumChain) GetConfig() *relayChain.Config {
	return ec.chainConfig
}

func (ec *ethereumChain) MinimumStake() (*big.Int, error) {
	return ec.stakingContract.MinimumStake()
}

// HasMinimumStake returns true if the specified address is staked.  False will
// be returned if not staked.  If err != nil then it was not possible to determine
// if the address is staked or not.
func (ec *ethereumChain) HasMinimumStake(address common.Address) (bool, error) {
	return ec.keepRandomBeaconOperatorContract.HasMinimumStake(address)
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

	subscription := ec.OnRelayEntrySubmitted(
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

	gasEstimate, err := ec.keepRandomBeaconOperatorContract.RelayEntryGasEstimate(entry)
	if err != nil {
		logger.Errorf("failed to estimate gas [%v]", err)
	}

	gasEstimateWithMargin := float64(gasEstimate) * float64(1.2) // 20% more than original
	_, err = ec.keepRandomBeaconOperatorContract.RelayEntry(
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

func (ec *ethereumChain) OnRelayEntrySubmitted(
	handle func(entry *event.EntrySubmitted),
) subscription.EventSubscription {
	onEvent := func(blockNumber uint64) {
		handle(&event.EntrySubmitted{
			BlockNumber: blockNumber,
		})
	}

	subscription := ec.keepRandomBeaconOperatorContract.RelayEntrySubmitted(
		nil,
	).OnEvent(onEvent)

	return subscription
}

func (ec *ethereumChain) OnRelayEntryRequested(
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

	subscription := ec.keepRandomBeaconOperatorContract.RelayEntryRequested(
		nil,
	).OnEvent(onEvent)

	return subscription
}

// TODO: Implement a real SelectGroup function once it is possible on the
//       contract side. The current implementation just return a group
//       where all members belong to the chain operator.
func (ec *ethereumChain) SelectGroup(seed *big.Int) ([]relayChain.StakerAddress, error) {
	groupSize := ec.GetConfig().GroupSize
	groupMembers := make([]relayChain.StakerAddress, groupSize)

	for index := range groupMembers {
		groupMembers[index] = ec.accountKey.Address.Bytes()
	}

	return groupMembers, nil
}

func (ec *ethereumChain) OnGroupRegistered(
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

	subscription := ec.keepRandomBeaconOperatorContract.DkgResultSubmittedEvent(
		nil,
	).OnEvent(onEvent)

	return subscription
}

func (ec *ethereumChain) IsGroupRegistered(groupPublicKey []byte) (bool, error) {
	return ec.keepRandomBeaconOperatorContract.IsGroupRegistered(groupPublicKey)
}

func (ec *ethereumChain) IsStaleGroup(groupPublicKey []byte) (bool, error) {
	return ec.keepRandomBeaconOperatorContract.IsStaleGroup(groupPublicKey)
}

func (ec *ethereumChain) GetGroupMembers(groupPublicKey []byte) (
	[]relayChain.StakerAddress,
	error,
) {
	members, err := ec.keepRandomBeaconOperatorContract.GetGroupMembers(
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

// TODO: Implement a real DkgStarted event subscription once it is possible
//       on the contract side. The current implementation generate a fake
//       event every 500th block where the seed is the keccak256 of the
//       block number.
func (ec *ethereumChain) OnDKGStarted(
	handler func(event *event.DKGStarted),
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	blocksChan := ec.blockCounter.WatchBlocks(ctx)

	go func() {
		for {
			select {
			case block := <-blocksChan:
				// Generate an event every 500th block.
				if block%500 == 0 {
					// The seed is keccak256(block).
					blockBytes := make([]byte, 8)
					binary.BigEndian.PutUint64(blockBytes, block)
					seedBytes := crypto.Keccak256(blockBytes)
					seed := new(big.Int).SetBytes(seedBytes)

					handler(&event.DKGStarted{
						Seed:        seed,
						BlockNumber: block,
					})
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return subscription.NewEventSubscription(func() {
		cancelCtx()
	})
}

// TODO: Implement a real OnDKGResultSubmitted event subscription once it is
//       possible on the contract side. The current implementation just pipes
//       the DKG submission event generated within SubmitDKGResult to the
//       handlers registered in the dkgResultSubmissionHandlers map.
func (ec *ethereumChain) OnDKGResultSubmitted(
	handler func(dkgResultPublication *event.DKGResultSubmission),
) subscription.EventSubscription {
	ec.dkgResultSubmissionHandlersMutex.Lock()
	defer ec.dkgResultSubmissionHandlersMutex.Unlock()

	// #nosec G404 (insecure random number source (rand))
	// Temporary test implementation doesn't require secure randomness.
	handlerID := rand.Int()

	ec.dkgResultSubmissionHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		ec.dkgResultSubmissionHandlersMutex.Lock()
		defer ec.dkgResultSubmissionHandlersMutex.Unlock()

		delete(ec.dkgResultSubmissionHandlers, handlerID)
	})
}

func (ec *ethereumChain) ReportRelayEntryTimeout() error {
	_, err := ec.keepRandomBeaconOperatorContract.ReportRelayEntryTimeout()
	if err != nil {
		return err
	}

	return nil
}

func (ec *ethereumChain) IsEntryInProgress() (bool, error) {
	return ec.keepRandomBeaconOperatorContract.IsEntryInProgress()
}

func (ec *ethereumChain) CurrentRequestStartBlock() (*big.Int, error) {
	return ec.keepRandomBeaconOperatorContract.CurrentRequestStartBlock()
}

func (ec *ethereumChain) CurrentRequestPreviousEntry() ([]byte, error) {
	return ec.keepRandomBeaconOperatorContract.CurrentRequestPreviousEntry()
}

func (ec *ethereumChain) CurrentRequestGroupPublicKey() ([]byte, error) {
	currentRequestGroupIndex, err := ec.keepRandomBeaconOperatorContract.CurrentRequestGroupIndex()
	if err != nil {
		return nil, err
	}

	return ec.keepRandomBeaconOperatorContract.GetGroupPublicKey(currentRequestGroupIndex)
}

// TODO: Implement a real SubmitDKGResult action once it is possible on the
//       contract side. The current implementation just creates and pipes
//       the DKG submission event to the handlers registered in the
//       dkgResultSubmissionHandlers map. Consider getting rid of the result
//       promise in favor of the fire-and-forget style.
func (ec *ethereumChain) SubmitDKGResult(
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

	subscription := ec.OnDKGResultSubmitted(
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

	ec.dkgResultSubmissionHandlersMutex.Lock()
	defer ec.dkgResultSubmissionHandlersMutex.Unlock()

	blockNumber, err := ec.blockCounter.CurrentBlock()
	if err != nil {
		close(publishedResult)
		subscription.Unsubscribe()
		failPromise(err)
		return resultPublicationPromise
	}

	for _, handler := range ec.dkgResultSubmissionHandlers {
		go func(handler func(*event.DKGResultSubmission)) {
			handler(&event.DKGResultSubmission{
				MemberIndex:    uint32(participantIndex),
				GroupPublicKey: result.GroupPublicKey,
				Misbehaved:     result.Misbehaved,
				BlockNumber:    blockNumber,
			})
		}(handler)
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
func (ec *ethereumChain) CalculateDKGResultHash(
	dkgResult *relayChain.DKGResult,
) (relayChain.DKGResultHash, error) {

	// Encode DKG result to the format matched with Solidity keccak256(abi.encodePacked(...))
	hash := crypto.Keccak256(dkgResult.GroupPublicKey, dkgResult.Misbehaved)

	return relayChain.DKGResultHashFromBytes(hash)
}

func (ec *ethereumChain) Address() common.Address {
	return ec.accountKey.Address
}
