package tbtc

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/subscription"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
	"golang.org/x/crypto/sha3"
)

var errNilDKGResult = fmt.Errorf("nil DKG result")

type localChain struct {
	dkgResultSubmissionHandlersMutex sync.Mutex
	dkgResultSubmissionHandlers      map[int]func(submission *DKGResultSubmittedEvent)

	resultSubmissionMutex sync.Mutex
	activeWallet          []byte
	resultSubmitterIndex  group.MemberIndex

	blockCounter       chain.BlockCounter
	chainConfig        *ChainConfig
	operatorPrivateKey *operator.PrivateKey
}

// GetConfig returns the chain configuration.
func (lc *localChain) GetConfig() *ChainConfig {
	return lc.chainConfig
}

// BlockCounter returns the block counter associated with the chain.
func (lc *localChain) BlockCounter() (chain.BlockCounter, error) {
	return lc.blockCounter, nil
}

// Signing returns the signing associated with the chain.
func (lc *localChain) Signing() chain.Signing {
	return local_v1.NewSigner(lc.operatorPrivateKey)
}

func (lc *localChain) OperatorKeyPair() (
	*operator.PrivateKey,
	*operator.PublicKey,
	error,
) {
	panic("unsupported")
}

func (lc *localChain) OperatorToStakingProvider() (chain.Address, bool, error) {
	panic("unsupported")
}

func (lc *localChain) EligibleStake(stakingProvider chain.Address) (*big.Int, error) {
	panic("unsupported")
}

func (lc *localChain) IsPoolLocked() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) IsOperatorInPool() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) IsOperatorUpToDate() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) JoinSortitionPool() error {
	panic("unsupported")
}

func (lc *localChain) UpdateOperatorStatus() error {
	panic("unsupported")
}

func (lc *localChain) IsEligibleForRewards() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) CanRestoreRewardEligibility() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) RestoreRewardEligibility() error {
	panic("unsupported")
}

func (lc *localChain) IsChaosnetActive() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) IsBetaOperator() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) SelectGroup() (*GroupSelectionResult, error) {
	panic("not implemented")
}

func (lc *localChain) OnDKGStarted(
	handler func(event *DKGStartedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

// OnDKGResultSubmitted registers a callback that is invoked when an on-chain
// notification of the DKG result submission is seen.
func (lc *localChain) OnDKGResultSubmitted(
	handler func(event *DKGResultSubmittedEvent),
) subscription.EventSubscription {
	lc.dkgResultSubmissionHandlersMutex.Lock()
	defer lc.dkgResultSubmissionHandlersMutex.Unlock()

	handlerID := local_v1.GenerateHandlerID()
	lc.dkgResultSubmissionHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		lc.dkgResultSubmissionHandlersMutex.Lock()
		defer lc.dkgResultSubmissionHandlersMutex.Unlock()

		delete(lc.dkgResultSubmissionHandlers, handlerID)
	})
}

// SubmitDKGResult submits the DKG result to the chain, along with signatures
// over result hash from group participants supporting the result.
func (lc *localChain) SubmitDKGResult(
	memberIndex group.MemberIndex,
	dkgResult *dkg.Result,
	signatures map[group.MemberIndex][]byte,
	groupSelectionResult *GroupSelectionResult,
) error {
	lc.dkgResultSubmissionHandlersMutex.Lock()
	defer lc.dkgResultSubmissionHandlersMutex.Unlock()

	lc.resultSubmissionMutex.Lock()
	defer lc.resultSubmissionMutex.Unlock()

	blockNumber, err := lc.blockCounter.CurrentBlock()
	if err != nil {
		return fmt.Errorf("failed to get the current block")
	}

	groupPublicKeyBytes, err := dkgResult.GroupPublicKeyBytes()
	if err != nil {
		return fmt.Errorf(
			"failed to extract group public key bytes from the result [%v]",
			err,
		)
	}

	for _, handler := range lc.dkgResultSubmissionHandlers {
		handler(&DKGResultSubmittedEvent{
			MemberIndex:         uint32(memberIndex),
			GroupPublicKeyBytes: groupPublicKeyBytes,
			Misbehaved:          dkgResult.MisbehavedMembersIndexes(),
			BlockNumber:         blockNumber,
		})
	}

	lc.activeWallet = groupPublicKeyBytes
	lc.resultSubmitterIndex = memberIndex

	return nil
}

// GetDKGState returns the current state of the DKG procedure.
func (lc *localChain) GetDKGState() (DKGState, error) {
	return AwaitingResult, nil
}

// CalculateDKGResultHash calculates 256-bit hash of DKG result using SHA3-256
// hashing algorithm.
func (lc *localChain) CalculateDKGResultHash(
	result *dkg.Result,
) (dkg.ResultHash, error) {
	if result == nil {
		return dkg.ResultHash{}, errNilDKGResult
	}

	encodedDKGResult := fmt.Sprint(result)
	dkgResultHash := dkg.ResultHash(
		sha3.Sum256([]byte(encodedDKGResult)),
	)
	return dkgResultHash, nil
}

func (lc *localChain) OnSignatureRequested(
	handler func(event *SignatureRequestedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

// Connect sets up the local chain.
func Connect(
	groupSize int,
	groupQuorum int,
	honestThreshold int,
) *localChain {
	operatorPrivateKey, _, err := operator.GenerateKeyPair(local_v1.DefaultCurve)
	if err != nil {
		panic(err)
	}

	return ConnectWithKey(
		groupSize,
		groupQuorum,
		honestThreshold,
		operatorPrivateKey,
	)
}

// ConnectWithKey sets up the local chain using the provided operator private
// key.
func ConnectWithKey(
	groupSize int,
	groupQuorum int,
	honestThreshold int,
	operatorPrivateKey *operator.PrivateKey,
) *localChain {
	blockCounter, _ := local_v1.BlockCounter()

	chainConfig := &ChainConfig{
		GroupSize:       groupSize,
		GroupQuorum:     groupQuorum,
		HonestThreshold: honestThreshold,
	}

	return &localChain{
		dkgResultSubmissionHandlers: make(
			map[int]func(submission *DKGResultSubmittedEvent),
		),
		blockCounter:       blockCounter,
		chainConfig:        chainConfig,
		operatorPrivateKey: operatorPrivateKey,
	}
}
