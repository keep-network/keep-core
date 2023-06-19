package coordinator_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/subscription"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
)

type localTbtcChain struct {
	mutex sync.Mutex

	depositRequests                 map[[32]byte]*tbtc.DepositChainRequest
	pastDepositRevealedEvents       map[[32]byte][]*tbtc.DepositRevealedEvent
	pastNewWalletRegisteredEvents   map[[32]byte][]*tbtc.NewWalletRegisteredEvent
	depositParameters               depositParameters
	depositSweepProposalValidations map[[32]byte]bool
	depositSweepProposals           []*tbtc.DepositSweepProposal
}

type depositParameters = struct {
	dustThreshold      uint64
	treasuryFeeDivisor uint64
	txMaxFee           uint64
	revealAheadPeriod  uint32
}

func newLocalTbtcChain() *localTbtcChain {
	return &localTbtcChain{
		depositRequests:                 make(map[[32]byte]*tbtc.DepositChainRequest),
		pastDepositRevealedEvents:       make(map[[32]byte][]*tbtc.DepositRevealedEvent),
		pastNewWalletRegisteredEvents:   make(map[[32]byte][]*tbtc.NewWalletRegisteredEvent),
		depositSweepProposalValidations: make(map[[32]byte]bool),
	}
}

func (lc *localTbtcChain) BlockCounter() (chain.BlockCounter, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) Signing() chain.Signing {
	panic("unsupported")
}

func (lc *localTbtcChain) OperatorKeyPair() (*operator.PrivateKey, *operator.PublicKey, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) GetBlockNumberByTimestamp(timestamp uint64) (uint64, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) OperatorToStakingProvider() (chain.Address, bool, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) EligibleStake(stakingProvider chain.Address) (*big.Int, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) IsPoolLocked() (bool, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) IsOperatorInPool() (bool, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) IsOperatorUpToDate() (bool, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) JoinSortitionPool() error {
	panic("unsupported")
}

func (lc *localTbtcChain) UpdateOperatorStatus() error {
	panic("unsupported")
}

func (lc *localTbtcChain) IsEligibleForRewards() (bool, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) CanRestoreRewardEligibility() (bool, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) RestoreRewardEligibility() error {
	panic("unsupported")
}

func (lc *localTbtcChain) IsChaosnetActive() (bool, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) IsBetaOperator() (bool, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) SelectGroup() (*tbtc.GroupSelectionResult, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) OnDKGStarted(
	func(event *tbtc.DKGStartedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localTbtcChain) PastDKGStartedEvents(
	filter *tbtc.DKGStartedEventFilter,
) ([]*tbtc.DKGStartedEvent, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) OnDKGResultSubmitted(
	func(event *tbtc.DKGResultSubmittedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localTbtcChain) OnDKGResultChallenged(
	func(event *tbtc.DKGResultChallengedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localTbtcChain) OnDKGResultApproved(
	func(event *tbtc.DKGResultApprovedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localTbtcChain) AssembleDKGResult(
	submitterMemberIndex group.MemberIndex,
	groupPublicKey *ecdsa.PublicKey,
	operatingMembersIndexes []group.MemberIndex,
	misbehavedMembersIndexes []group.MemberIndex,
	signatures map[group.MemberIndex][]byte,
	groupSelectionResult *tbtc.GroupSelectionResult,
) (*tbtc.DKGChainResult, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) SubmitDKGResult(dkgResult *tbtc.DKGChainResult) error {
	panic("unsupported")
}

func (lc *localTbtcChain) GetDKGState() (tbtc.DKGState, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) CalculateDKGResultSignatureHash(
	groupPublicKey *ecdsa.PublicKey,
	misbehavedMembersIndexes []group.MemberIndex,
	startBlock uint64,
) (dkg.ResultSignatureHash, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) IsDKGResultValid(dkgResult *tbtc.DKGChainResult) (bool, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) ChallengeDKGResult(dkgResult *tbtc.DKGChainResult) error {
	panic("unsupported")
}

func (lc *localTbtcChain) ApproveDKGResult(dkgResult *tbtc.DKGChainResult) error {
	panic("unsupported")
}

func (lc *localTbtcChain) DKGParameters() (*tbtc.DKGParameters, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) GetOperatorID(operatorAddress chain.Address) (chain.OperatorID, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) PastDepositRevealedEvents(
	filter *tbtc.DepositRevealedEventFilter,
) ([]*tbtc.DepositRevealedEvent, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	eventsKey, err := buildPastDepositRevealedEventsKey(filter)
	if err != nil {
		return nil, err
	}

	events, ok := lc.pastDepositRevealedEvents[eventsKey]
	if !ok {
		return nil, fmt.Errorf("no events for given filter")
	}

	return events, nil
}

func (lc *localTbtcChain) addPastDepositRevealedEvent(
	filter *tbtc.DepositRevealedEventFilter,
	event *tbtc.DepositRevealedEvent,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	eventsKey, err := buildPastDepositRevealedEventsKey(filter)
	if err != nil {
		return err
	}

	lc.pastDepositRevealedEvents[eventsKey] = append(
		lc.pastDepositRevealedEvents[eventsKey],
		event,
	)

	return nil
}

func buildPastDepositRevealedEventsKey(
	filter *tbtc.DepositRevealedEventFilter,
) ([32]byte, error) {
	if filter == nil {
		return [32]byte{}, nil
	}

	var buffer bytes.Buffer

	startBlock := make([]byte, 8)
	binary.BigEndian.PutUint64(startBlock, filter.StartBlock)
	buffer.Write(startBlock)

	if filter.EndBlock != nil {
		endBlock := make([]byte, 8)
		binary.BigEndian.PutUint64(startBlock, *filter.EndBlock)
		buffer.Write(endBlock)
	}

	for _, depositor := range filter.Depositor {
		depositorBytes, err := hex.DecodeString(depositor.String())
		if err != nil {
			return [32]byte{}, err
		}

		buffer.Write(depositorBytes)
	}

	for _, walletPublicKeyHash := range filter.WalletPublicKeyHash {
		buffer.Write(walletPublicKeyHash[:])
	}

	return sha256.Sum256(buffer.Bytes()), nil
}

func (lc *localTbtcChain) GetDepositRequest(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) (*tbtc.DepositChainRequest, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	requestKey := buildDepositRequestKey(fundingTxHash, fundingOutputIndex)

	request, ok := lc.depositRequests[requestKey]
	if !ok {
		return nil, fmt.Errorf("no request for given key")
	}

	return request, nil
}

func (lc *localTbtcChain) setDepositRequest(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
	request *tbtc.DepositChainRequest,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	requestKey := buildDepositRequestKey(fundingTxHash, fundingOutputIndex)

	lc.depositRequests[requestKey] = request
}

func (lc *localTbtcChain) PastNewWalletRegisteredEvents(
	filter *tbtc.NewWalletRegisteredEventFilter,
) ([]*tbtc.NewWalletRegisteredEvent, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	eventsKey, err := buildPastNewWalletRegisteredEventsKey(filter)
	if err != nil {
		return nil, err
	}

	events, ok := lc.pastNewWalletRegisteredEvents[eventsKey]
	if !ok {
		return nil, fmt.Errorf("no events for given filter")
	}

	return events, nil
}

func (lc *localTbtcChain) addPastNewWalletRegisteredEvent(
	filter *tbtc.NewWalletRegisteredEventFilter,
	event *tbtc.NewWalletRegisteredEvent,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	eventsKey, err := buildPastNewWalletRegisteredEventsKey(filter)
	if err != nil {
		return err
	}

	if _, ok := lc.pastNewWalletRegisteredEvents[eventsKey]; !ok {
		lc.pastNewWalletRegisteredEvents[eventsKey] = []*tbtc.NewWalletRegisteredEvent{}
	}

	lc.pastNewWalletRegisteredEvents[eventsKey] = append(
		lc.pastNewWalletRegisteredEvents[eventsKey],
		event,
	)

	return nil
}

func buildPastNewWalletRegisteredEventsKey(
	filter *tbtc.NewWalletRegisteredEventFilter,
) ([32]byte, error) {
	if filter == nil {
		return [32]byte{}, nil
	}

	var buffer bytes.Buffer

	startBlock := make([]byte, 8)
	binary.BigEndian.PutUint64(startBlock, filter.StartBlock)
	buffer.Write(startBlock)

	if filter.EndBlock != nil {
		endBlock := make([]byte, 8)
		binary.BigEndian.PutUint64(startBlock, *filter.EndBlock)
		buffer.Write(endBlock)
	}

	for _, ecdsaWalletID := range filter.EcdsaWalletID {
		buffer.Write(ecdsaWalletID[:])
	}

	for _, walletPublicKeyHash := range filter.WalletPublicKeyHash {
		buffer.Write(walletPublicKeyHash[:])
	}

	return sha256.Sum256(buffer.Bytes()), nil
}

func (lc *localTbtcChain) GetWallet(walletPublicKeyHash [20]byte) (*tbtc.WalletChainData, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) ComputeMainUtxoHash(mainUtxo *bitcoin.UnspentTransactionOutput) [32]byte {
	panic("unsupported")
}

func (lc *localTbtcChain) BuildDepositKey(fundingTxHash bitcoin.Hash, fundingOutputIndex uint32) *big.Int {
	depositKeyBytes := buildDepositRequestKey(fundingTxHash, fundingOutputIndex)

	return new(big.Int).SetBytes(depositKeyBytes[:])
}

func buildDepositRequestKey(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) [32]byte {
	fundingOutputIndexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(fundingOutputIndexBytes, fundingOutputIndex)

	return sha256.Sum256(append(fundingTxHash[:], fundingOutputIndexBytes...))
}

func (lc *localTbtcChain) GetDepositParameters() (
	dustThreshold uint64,
	treasuryFeeDivisor uint64,
	txMaxFee uint64,
	revealAheadPeriod uint32,
	err error,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	return lc.depositParameters.dustThreshold,
		lc.depositParameters.treasuryFeeDivisor,
		lc.depositParameters.txMaxFee,
		lc.depositParameters.revealAheadPeriod,
		nil
}

func (lc *localTbtcChain) TxProofDifficultyFactor() (*big.Int, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) setDepositParameters(
	dustThreshold uint64,
	treasuryFeeDivisor uint64,
	txMaxFee uint64,
	revealAheadPeriod uint32,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.depositParameters = depositParameters{
		dustThreshold:      dustThreshold,
		treasuryFeeDivisor: treasuryFeeDivisor,
		txMaxFee:           txMaxFee,
		revealAheadPeriod:  revealAheadPeriod,
	}
}

func (lc *localTbtcChain) OnHeartbeatRequestSubmitted(
	handler func(event *tbtc.HeartbeatRequestSubmittedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localTbtcChain) OnDepositSweepProposalSubmitted(
	func(event *tbtc.DepositSweepProposalSubmittedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localTbtcChain) PastDepositSweepProposalSubmittedEvents(
	filter *tbtc.DepositSweepProposalSubmittedEventFilter,
) ([]*tbtc.DepositSweepProposalSubmittedEvent, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) GetWalletLock(
	walletPublicKeyHash [20]byte,
) (time.Time, tbtc.WalletActionType, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) ValidateDepositSweepProposal(
	proposal *tbtc.DepositSweepProposal,
	depositsExtraInfo []struct {
		*tbtc.Deposit
		FundingTx *bitcoin.Transaction
	},
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	key, err := buildDepositSweepProposalValidationKey(proposal)
	if err != nil {
		return err
	}

	result, ok := lc.depositSweepProposalValidations[key]
	if !ok {
		return fmt.Errorf("validation result unknown")
	}

	if !result {
		return fmt.Errorf("validation failed")
	}

	return nil
}

func (lc *localTbtcChain) setDepositSweepProposalValidationResult(
	proposal *tbtc.DepositSweepProposal,
	depositsExtraInfo []struct {
		*tbtc.Deposit
		FundingTx *bitcoin.Transaction
	},
	result bool,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	key, err := buildDepositSweepProposalValidationKey(proposal)
	if err != nil {
		return err
	}

	lc.depositSweepProposalValidations[key] = result

	return nil
}

func buildDepositSweepProposalValidationKey(
	proposal *tbtc.DepositSweepProposal,
) ([32]byte, error) {
	var buffer bytes.Buffer

	buffer.Write(proposal.WalletPublicKeyHash[:])

	for _, deposit := range proposal.DepositsKeys {
		buffer.Write(deposit.FundingTxHash[:])

		fundingOutputIndex := make([]byte, 4)
		binary.BigEndian.PutUint32(fundingOutputIndex, deposit.FundingOutputIndex)
		buffer.Write(fundingOutputIndex)
	}

	buffer.Write(proposal.SweepTxFee.Bytes())

	return sha256.Sum256(buffer.Bytes()), nil
}

func (lc *localTbtcChain) SubmitDepositSweepProposalWithReimbursement(
	proposal *tbtc.DepositSweepProposal,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.depositSweepProposals = append(lc.depositSweepProposals, proposal)

	return nil
}

func (lc *localTbtcChain) GetDepositSweepMaxSize() (uint16, error) {
	panic("unsupported")
}
