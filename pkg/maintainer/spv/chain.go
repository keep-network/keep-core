package spv

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/bitcoin"
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
	// funding transaction hash and output index. Returns an error if the
	// deposit was not found.
	GetDepositRequest(
		fundingTxHash bitcoin.Hash,
		fundingOutputIndex uint32,
	) (*tbtc.DepositChainRequest, error)

	// TxProofDifficultyFactor returns the number of confirmations on the
	// Bitcoin chain required to successfully evaluate an SPV proof.
	TxProofDifficultyFactor() (*big.Int, error)

	// PastDepositSweepProposalSubmittedEvents returns past
	// `DepositSweepProposalSubmitted` events.
	PastDepositSweepProposalSubmittedEvents(
		filter *tbtc.DepositSweepProposalSubmittedEventFilter,
	) ([]*tbtc.DepositSweepProposalSubmittedEvent, error)
}
