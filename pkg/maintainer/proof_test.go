package maintainer

import (
	"reflect"
	"testing"

	testData "github.com/keep-network/keep-core/internal/testdata/bitcoin"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

func TestAssembleTransactionProof(t *testing.T) {
	for testName, test := range testData.SpvProof {
		t.Run(testName, func(t *testing.T) {
			requiredConfirmations := test.RequiredConfirmations
			transactionHash := test.BitcoinChainData.TransactionHash
			transaction := test.BitcoinChainData.Transaction
			accumulatedConfirmations := test.BitcoinChainData.AccumulatedTxConfirmations
			blockHeaders := test.BitcoinChainData.HeadersChain
			transactionMerkleProof := test.BitcoinChainData.TransactionMerkleProof
			expectedProof := test.ExpectedProof
			expectedTx := &transaction

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

			tx, proof, err := AssembleTransactionProof(
				transactionHash,
				requiredConfirmations,
				bitcoinChain,
			)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(expectedProof, proof) {
				t.Errorf(
					"unexpected proof\nexpected: %v\nactual:   %v\n",
					expectedProof,
					proof,
				)
			}
			if !reflect.DeepEqual(expectedTx, tx) {
				t.Errorf(
					"unexpected transaction\nexpected: %v\nactual:   %v\n",
					expectedTx,
					tx,
				)
			}
		})
	}
}
