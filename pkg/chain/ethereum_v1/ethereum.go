package ethereum_v1

import (
	"fmt"
	"github.com/ipfs/go-log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
	"github.com/keep-network/keep-core/pkg/beacon/event"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/subscription"
)

var logger = log.Logger("keep-chain-ethereum-v1")

// ThresholdRelay converts from ethereumChain to beacon.ChainInterface.
func (ec *ethereumChain) ThresholdRelay() beaconchain.Interface {
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

func (ec *ethereumChain) GetConfig() *beaconchain.Config {
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

func (ec *ethereumChain) SelectGroup(seed *big.Int) ([]chain.Address, error) {
	panic("unsupported")
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
	panic("unsupported")
}

func (ec *ethereumChain) IsStaleGroup(groupPublicKey []byte) (bool, error) {
	return ec.keepRandomBeaconOperatorContract.IsStaleGroup(groupPublicKey)
}

func (ec *ethereumChain) GetGroupMembers(groupPublicKey []byte) (
	[]chain.Address,
	error,
) {
	members, err := ec.keepRandomBeaconOperatorContract.GetGroupMembers(
		groupPublicKey,
	)
	if err != nil {
		return nil, err
	}

	stakerAddresses := make([]chain.Address, len(members))
	for i, member := range members {
		stakerAddresses[i] = chain.Address(member.String())
	}

	return stakerAddresses, nil
}

func (ec *ethereumChain) OnDKGStarted(
	handler func(event *event.DKGStarted),
) subscription.EventSubscription {
	panic("unsupported")
}

func (ec *ethereumChain) OnDKGResultSubmitted(
	handler func(dkgResultPublication *event.DKGResultSubmission),
) subscription.EventSubscription {
	panic("unsupported")
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

func (ec *ethereumChain) SubmitDKGResult(
	participantIndex beaconchain.GroupMemberIndex,
	result *beaconchain.DKGResult,
	signatures map[beaconchain.GroupMemberIndex][]byte,
) error {
	panic("unsupported")
}

// convertSignaturesToChainFormat converts signatures map to two slices. First
// slice contains indices of members from the map, second slice is a slice of
// concatenated signatures. Signatures and member indices are returned in the
// matching order. It requires each signature to be exactly 65-byte long.
func convertSignaturesToChainFormat(
	signatures map[beaconchain.GroupMemberIndex][]byte,
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
	dkgResult *beaconchain.DKGResult,
) (beaconchain.DKGResultHash, error) {

	// Encode DKG result to the format matched with Solidity keccak256(abi.encodePacked(...))
	hash := crypto.Keccak256(dkgResult.GroupPublicKey, dkgResult.Misbehaved)

	return beaconchain.DKGResultHashFromBytes(hash)
}

func (ec *ethereumChain) Address() common.Address {
	return ec.accountKey.Address
}

func (ec *ethereumChain) OperatorToStakingProvider() (chain.Address, bool, error) {
	panic("unsupported")
}

func (ec *ethereumChain) EligibleStake(stakingProvider chain.Address) (*big.Int, error) {
	panic("unsupported")
}

func (ec *ethereumChain) IsPoolLocked() (bool, error) {
	panic("unsupported")
}

func (ec *ethereumChain) IsOperatorInPool() (bool, error) {
	panic("unsupported")
}

func (ec *ethereumChain) IsOperatorUpToDate() (bool, error) {
	panic("unsupported")
}

func (ec *ethereumChain) JoinSortitionPool() error {
	panic("unsupported")
}

func (ec *ethereumChain) UpdateOperatorStatus() error {
	panic("unsupported")
}
