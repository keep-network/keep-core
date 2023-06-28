package wallet

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/coordinator"
	"github.com/keep-network/keep-core/pkg/tbtc"
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

func (lc *localChain) PastDepositRevealedEvents(
	filter *tbtc.DepositRevealedEventFilter,
) ([]*tbtc.DepositRevealedEvent, error) {
	panic("unsupported")
}

func (lc *localChain) GetDepositRequest(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) (*tbtc.DepositChainRequest, bool, error) {
	panic("unsupported")
}

func (lc *localChain) PastNewWalletRegisteredEvents(
	filter *coordinator.NewWalletRegisteredEventFilter,
) ([]*coordinator.NewWalletRegisteredEvent, error) {
	panic("unsupported")
}

func (lc *localChain) PastRedemptionRequestedEvents(
	filter *tbtc.RedemptionRequestedEventFilter,
) ([]*tbtc.RedemptionRequestedEvent, error) {
	panic("unsupported")
}

func (lc *localChain) BuildDepositKey(fundingTxHash bitcoin.Hash, fundingOutputIndex uint32) *big.Int {
	panic("unsupported")
}

func (lc *localChain) BuildRedemptionKey(
	walletPublicKeyHash [20]byte,
	redeemerOutputScript bitcoin.Script,
) (*big.Int, error) {
	panic("unsupported")
}

func (lc *localChain) GetPendingRedemptionRequest(
	walletPublicKeyHash [20]byte,
	redeemerOutputScript bitcoin.Script,
) (*tbtc.RedemptionRequest, error) {
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

func (lc *localChain) GetRedemptionParameters() (
	dustThreshold uint64,
	treasuryFeeDivisor uint64,
	txMaxFee uint64,
	txMaxTotalFee uint64,
	timeout uint32,
	timeoutSlashingAmount *big.Int,
	timeoutNotifierRewardMultiplier uint32,
	err error,
) {
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

func (lc *localChain) SubmitRedemptionProposalWithReimbursement(
	proposal *tbtc.RedemptionProposal,
) error {
	panic("unsupported")
}

func (lc *localChain) GetDepositSweepMaxSize() (uint16, error) {
	panic("unsupported")
}

func (lc *localChain) ValidateRedemptionProposal(
	proposal *tbtc.RedemptionProposal,
) error {
	panic("unsupported")
}

func (lc *localChain) GetRedemptionMaxSize() (uint16, error) {
	panic("unsupported")
}

func (lc *localChain) GetRedemptionRequestMinAge() (uint32, error) {
	panic("unsupported")
}
