package maintainer

import (
	"testing"

	testData "github.com/keep-network/keep-core/internal/testdata/bitcoin"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

func TestAssembleTransactionProof(t *testing.T) {
	// TODO: Rewrite to parametrized test once more cases are added.
	transactionHash := testData.AssembleProof["single input"].BitcoinChainData.TransactionHash
	transaction := testData.AssembleProof["single input"].BitcoinChainData.Transaction
	accumulatedConfirmations := testData.AssembleProof["single input"].BitcoinChainData.AccumulatedTxConfirmations
	requiredConfirmations := testData.AssembleProof["single input"].RequiredConfirmations
	blockHeaders := testData.AssembleProof["single input"].BitcoinChainData.HeadersChain
	transactionMerkleProof := testData.AssembleProof["single input"].BitcoinChainData.TransactionMerkleProof

	bitcoinChain := connectLocalBitcoinChain()

	var transactions = map[bitcoin.Hash]*bitcoin.Transaction{
		transactionHash: &transaction,
	}
	bitcoinChain.SetTransactions(transactions)

	var transactionConfirmations = map[bitcoin.Hash]uint{
		transactionHash: accumulatedConfirmations,
	}
	bitcoinChain.SetTransactionConfirmations(transactionConfirmations)

	bitcoinChain.SetBlockHeaders(blockHeaders)

	bitcoinChain.SetTransactionMerkleProof(transactionMerkleProof)

	_, _, err := AssembleTransactionProof(
		transactionHash,
		requiredConfirmations,
		bitcoinChain,
	)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Add check for the returned proof and transaction.
}
