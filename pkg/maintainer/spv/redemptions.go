package spv

import "github.com/keep-network/keep-core/pkg/bitcoin"

// SubmitRedemptionProof prepares redemption proof for the given transaction
// and submits it to the on-chain contract. If the number of required
// confirmations is `0`, an error is returned.
//
// TODO: Expose this function through the maintainer-cli tool.
func SubmitRedemptionProof(
	transactionHash bitcoin.Hash,
	requiredConfirmations uint,
	btcChain bitcoin.Chain,
	spvChain Chain,
) error {
	panic("not implemented yet")
}

func getUnprovenRedemptionTransactions(
	historyDepth uint64,
	transactionLimit int,
	btcChain bitcoin.Chain,
	spvChain Chain,
) (
	[]*bitcoin.Transaction,
	error,
) {
	panic("not implemented yet")
}
