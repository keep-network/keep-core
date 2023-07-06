package spv

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"math/big"
	"sync"
)

type localChain struct {
	mutex sync.Mutex
}

func newLocalChain() *localChain {
	return &localChain{}
}

func (lc *localChain) SubmitDepositSweepProofWithReimbursement(
	transaction *bitcoin.Transaction,
	proof *bitcoin.SpvProof,
	mainUTXO bitcoin.UnspentTransactionOutput,
	vault common.Address,
) error {
	panic("unsupported")
}

func (lc *localChain) GetDepositRequest(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) (*tbtc.DepositChainRequest, bool, error) {
	panic("unsupported")
}

func (lc *localChain) GetWallet(walletPublicKeyHash [20]byte) (
	*tbtc.WalletChainData,
	error,
) {
	panic("unsupported")
}

func (lc *localChain) ComputeMainUtxoHash(
	mainUtxo *bitcoin.UnspentTransactionOutput,
) [32]byte {
	panic("unsupported")
}

func (lc *localChain) TxProofDifficultyFactor() (*big.Int, error) {
	panic("unsupported")
}

func (lc *localChain) BlockCounter() (chain.BlockCounter, error) {
	panic("unsupported")
}

func (lc *localChain) PastDepositSweepProposalSubmittedEvents(
	filter *tbtc.DepositSweepProposalSubmittedEventFilter,
) (
	[]*tbtc.DepositSweepProposalSubmittedEvent,
	error,
) {
	panic("unsupported")
}

func (lc *localChain) PastRedemptionProposalSubmittedEvents(
	filter *tbtc.RedemptionProposalSubmittedEventFilter,
) (
	[]*tbtc.RedemptionProposalSubmittedEvent,
	error,
) {
	panic("unsupported")
}

func (lc *localChain) GetPendingRedemptionRequest(
	walletPublicKeyHash [20]byte,
	redeemerOutputScript bitcoin.Script,
) (*tbtc.RedemptionRequest, bool, error) {
	panic("unsupported")
}

func (lc *localChain) SubmitRedemptionProofWithReimbursement(
	transaction *bitcoin.Transaction,
	proof *bitcoin.SpvProof,
	mainUTXO bitcoin.UnspentTransactionOutput,
	walletPublicKeyHash [20]byte,
) error {
	panic("unsupported")
}

func (lc *localChain) Ready() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) IsAuthorized(address chain.Address) (bool, error) {
	panic("unsupported")
}

func (lc *localChain) IsAuthorizedForRefund(address chain.Address) (
	bool,
	error,
) {
	panic("unsupported")
}

func (lc *localChain) Signing() chain.Signing {
	panic("unsupported")
}

func (lc *localChain) Retarget(headers []*bitcoin.BlockHeader) error {
	panic("unsupported")
}

func (lc *localChain) RetargetWithRefund(headers []*bitcoin.BlockHeader) error {
	panic("unsupported")
}

func (lc *localChain) CurrentEpoch() (uint64, error) {
	panic("unsupported")
}

func (lc *localChain) ProofLength() (uint64, error) {
	panic("unsupported")
}

func (lc *localChain) GetCurrentAndPrevEpochDifficulty() (
	*big.Int,
	*big.Int,
	error,
) {
	panic("unsupported")
}
