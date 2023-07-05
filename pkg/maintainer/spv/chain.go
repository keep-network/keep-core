package spv

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

type Chain interface {
	// SubmitDepositSweepProofWithReimbursement submits the deposit sweep proof
	// via MaintainerProxy. It is used to prove the deposit sweep Bitcoin
	// transaction and update depositors' balances. The caller is reimbursed.
	SubmitDepositSweepProofWithReimbursement(
		transaction *bitcoin.Transaction,
		proof *bitcoin.SpvProof,
		mainUTXO bitcoin.UnspentTransactionOutput,
		vault common.Address,
	) error

	// GetDepositRequest gets the on-chain deposit request for the given
	// funding transaction hash and output index.The returned values represent:
	// - deposit request which is non-nil only when the deposit request was
	//   found,
	// - boolean value which is true if the deposit request was found, false
	//   otherwise,
	// - error which is non-nil only when the function execution failed. It will
	//   be nil if the deposit request was not found, but the function execution
	//   succeeded.
	GetDepositRequest(
		fundingTxHash bitcoin.Hash,
		fundingOutputIndex uint32,
	) (*tbtc.DepositChainRequest, bool, error)

	GetWallet(
		walletPublicKeyHash [20]byte,
	) (*tbtc.WalletChainData, error)

	// ComputeMainUtxoHash computes the hash of the provided main UTXO
	// according to the on-chain Bridge rules.
	ComputeMainUtxoHash(mainUtxo *bitcoin.UnspentTransactionOutput) [32]byte

	// TxProofDifficultyFactor returns the number of confirmations on the
	// Bitcoin chain required to successfully evaluate an SPV proof.
	TxProofDifficultyFactor() (*big.Int, error)

	// BlockCounter returns the chain's block counter.
	BlockCounter() (chain.BlockCounter, error)

	// PastDepositSweepProposalSubmittedEvents returns past
	// `DepositSweepProposalSubmitted` events.
	PastDepositSweepProposalSubmittedEvents(
		filter *tbtc.DepositSweepProposalSubmittedEventFilter,
	) ([]*tbtc.DepositSweepProposalSubmittedEvent, error)

	// PastRedemptionProposalSubmittedEvents returns past
	// `RedemptionProposalSubmitted` events.
	PastRedemptionProposalSubmittedEvents(
		filter *tbtc.RedemptionProposalSubmittedEventFilter,
	) ([]*tbtc.RedemptionProposalSubmittedEvent, error)

	// GetPendingRedemptionRequest gets the on-chain pending redemption request
	// for the given wallet public key hash and redeemer output script.
	// Returns an error if the request was not found.
	GetPendingRedemptionRequest(
		walletPublicKeyHash [20]byte,
		redeemerOutputScript bitcoin.Script,
	) (*tbtc.RedemptionRequest, error)
}
