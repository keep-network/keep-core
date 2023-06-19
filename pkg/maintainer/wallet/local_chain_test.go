package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/subscription"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
)

type walletLock struct {
	lockExpiration time.Time
	walletAction   tbtc.WalletActionType
}

type localChain struct {
	mutex sync.Mutex

	walletLocks map[[20]byte]*walletLock
}

func newLocalChain() *localChain {
	return &localChain{
		walletLocks: make(map[[20]byte]*walletLock),
	}
}

func (lc *localChain) BlockCounter() (chain.BlockCounter, error) {
	panic("unsupported")
}

func (lc *localChain) Signing() chain.Signing {
	panic("unsupported")
}

func (lc *localChain) OperatorKeyPair() (*operator.PrivateKey, *operator.PublicKey, error) {
	panic("unsupported")
}

func (lc *localChain) GetBlockNumberByTimestamp(timestamp uint64) (uint64, error) {
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

func (lc *localChain) SelectGroup() (*tbtc.GroupSelectionResult, error) {
	panic("unsupported")
}

func (lc *localChain) OnDKGStarted(
	func(event *tbtc.DKGStartedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localChain) PastDKGStartedEvents(
	filter *tbtc.DKGStartedEventFilter,
) ([]*tbtc.DKGStartedEvent, error) {
	panic("unsupported")
}

func (lc *localChain) OnDKGResultSubmitted(
	func(event *tbtc.DKGResultSubmittedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localChain) OnDKGResultChallenged(
	func(event *tbtc.DKGResultChallengedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localChain) OnDKGResultApproved(
	func(event *tbtc.DKGResultApprovedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localChain) AssembleDKGResult(
	submitterMemberIndex group.MemberIndex,
	groupPublicKey *ecdsa.PublicKey,
	operatingMembersIndexes []group.MemberIndex,
	misbehavedMembersIndexes []group.MemberIndex,
	signatures map[group.MemberIndex][]byte,
	groupSelectionResult *tbtc.GroupSelectionResult,
) (*tbtc.DKGChainResult, error) {
	panic("unsupported")
}

func (lc *localChain) SubmitDKGResult(dkgResult *tbtc.DKGChainResult) error {
	panic("unsupported")
}

func (lc *localChain) GetDKGState() (tbtc.DKGState, error) {
	panic("unsupported")
}

func (lc *localChain) CalculateDKGResultSignatureHash(
	groupPublicKey *ecdsa.PublicKey,
	misbehavedMembersIndexes []group.MemberIndex,
	startBlock uint64,
) (dkg.ResultSignatureHash, error) {
	panic("unsupported")
}

func (lc *localChain) IsDKGResultValid(dkgResult *tbtc.DKGChainResult) (bool, error) {
	panic("unsupported")
}

func (lc *localChain) ChallengeDKGResult(dkgResult *tbtc.DKGChainResult) error {
	panic("unsupported")
}

func (lc *localChain) ApproveDKGResult(dkgResult *tbtc.DKGChainResult) error {
	panic("unsupported")
}

func (lc *localChain) DKGParameters() (*tbtc.DKGParameters, error) {
	panic("unsupported")
}

func (lc *localChain) GetOperatorID(operatorAddress chain.Address) (chain.OperatorID, error) {
	panic("unsupported")
}

func (lc *localChain) PastDepositRevealedEvents(
	filter *tbtc.DepositRevealedEventFilter,
) ([]*tbtc.DepositRevealedEvent, error) {
	panic("unsupported")
}

func (lc *localChain) GetDepositRequest(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) (*tbtc.DepositChainRequest, error) {
	panic("unsupported")
}

func (lc *localChain) PastNewWalletRegisteredEvents(
	filter *tbtc.NewWalletRegisteredEventFilter,
) ([]*tbtc.NewWalletRegisteredEvent, error) {
	panic("unsupported")
}

func (lc *localChain) GetWallet(walletPublicKeyHash [20]byte) (*tbtc.WalletChainData, error) {
	panic("unsupported")
}

func (lc *localChain) ComputeMainUtxoHash(mainUtxo *bitcoin.UnspentTransactionOutput) [32]byte {
	panic("unsupported")
}

func (lc *localChain) BuildDepositKey(fundingTxHash bitcoin.Hash, fundingOutputIndex uint32) *big.Int {
	panic("unsupported")
}

func (lc *localChain) GetDepositParameters() (
	dustThreshold uint64,
	treasuryFeeDivisor uint64,
	txMaxFee uint64,
	revealAheadPeriod uint32,
	err error,
) {
	panic("unsupported")
}

func (lc *localChain) TxProofDifficultyFactor() (*big.Int, error) {
	panic("unsupported")
}

func (lc *localChain) OnHeartbeatRequestSubmitted(
	handler func(event *tbtc.HeartbeatRequestSubmittedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localChain) OnDepositSweepProposalSubmitted(
	func(event *tbtc.DepositSweepProposalSubmittedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localChain) PastDepositSweepProposalSubmittedEvents(
	filter *tbtc.DepositSweepProposalSubmittedEventFilter,
) ([]*tbtc.DepositSweepProposalSubmittedEvent, error) {
	panic("unsupported")
}

func (lc *localChain) GetWalletLock(
	walletPublicKeyHash [20]byte,
) (time.Time, tbtc.WalletActionType, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	walletLock, ok := lc.walletLocks[walletPublicKeyHash]
	if !ok {
		return time.Time{}, tbtc.Noop, fmt.Errorf("no lock configured for given wallet")
	}

	return walletLock.lockExpiration, walletLock.walletAction, nil
}

func (lc *localChain) setWalletLock(
	walletPublicKeyHash [20]byte,
	lockExpiration time.Time,
	walletAction tbtc.WalletActionType,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.walletLocks[walletPublicKeyHash] = &walletLock{
		lockExpiration: lockExpiration,
		walletAction:   walletAction,
	}
}

func (lc *localChain) resetWalletLock(
	walletPublicKeyHash [20]byte,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.walletLocks[walletPublicKeyHash] = &walletLock{
		lockExpiration: time.Unix(0, 0),
		walletAction:   tbtc.Noop,
	}
}

func (lc *localChain) ValidateDepositSweepProposal(
	proposal *tbtc.DepositSweepProposal,
	depositsExtraInfo []struct {
		*tbtc.Deposit
		FundingTx *bitcoin.Transaction
	},
) error {
	panic("unsupported")
}

func (lc *localChain) SubmitDepositSweepProposalWithReimbursement(
	proposal *tbtc.DepositSweepProposal,
) error {
	panic("unsupported")
}

func (lc *localChain) GetDepositSweepMaxSize() (uint16, error) {
	panic("unsupported")
}

func (lc *localChain) SubmitDepositSweepProofWithReimbursement(
	transaction *bitcoin.Transaction,
	proof *bitcoin.SpvProof,
	mainUTXO bitcoin.UnspentTransactionOutput,
	vault common.Address,
) error {
	panic("unsupported")
}
