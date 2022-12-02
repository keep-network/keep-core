package tbtc

import (
	"encoding/hex"
	"fmt"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/internal/tbtctest"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"testing"
)

func TestAssembleDepositSweepTransaction(t *testing.T) {
	scenarios, err := tbtctest.LoadDepositSweepTestScenarios()
	if err != nil {
		t.Fatal(err)
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Title, func(t *testing.T) {
			bitcoinChain := newMockBitcoinChain()

			for _, transaction := range scenario.InputTransactions {
				err := bitcoinChain.addTransaction(transaction)
				if err != nil {
					t.Fatal(err)
				}
			}

			deposits := make([]*deposit, len(scenario.Deposits))
			for i, d := range scenario.Deposits {
				deposits[i] = &deposit{
					utxo:                d.Utxo,
					depositor:           d.Depositor,
					blindingFactor:      d.BlindingFactor,
					walletPublicKeyHash: d.WalletPublicKeyHash,
					refundPublicKeyHash: d.RefundPublicKeyHash,
					refundLocktime:      d.RefundLocktime,
					vault: chain.Address(
						hex.EncodeToString(d.Vault[:]),
					),
				}
			}

			builder, err := assembleDepositSweepTransaction(
				bitcoinChain,
				scenario.WalletPublicKey,
				scenario.WalletMainUtxo,
				deposits,
				scenario.Fee,
			)
			if err != nil {
				t.Fatal(err)
			}

			sigHashes, err := builder.ComputeSignatureHashes()
			if err != nil {
				t.Fatal(err)
			}

			for i, sigHash := range sigHashes {
				testutils.AssertBigIntsEqual(
					t,
					fmt.Sprintf("sighash for input [%v]", i),
					scenario.ExpectedSigHashes[i],
					sigHash,
				)
			}

			transaction, err := builder.AddSignatures(scenario.Signatures)
			if err != nil {
				t.Fatal(err)
			}

			testutils.AssertBytesEqual(
				t,
				scenario.ExpectedSweepTransaction.Serialize(),
				transaction.Serialize(),
			)
			testutils.AssertStringsEqual(
				t,
				"sweep transaction hash",
				scenario.ExpectedSweepTransactionHash.String(bitcoin.InternalByteOrder),
				transaction.Hash().String(bitcoin.InternalByteOrder),
			)
			testutils.AssertStringsEqual(
				t,
				"sweep transaction witness hash",
				scenario.ExpectedSweepTransactionWitnessHash.String(bitcoin.InternalByteOrder),
				transaction.WitnessHash().String(bitcoin.InternalByteOrder),
			)
		})
	}
}

type mockBitcoinChain struct {
	transactions map[bitcoin.Hash]*bitcoin.Transaction
}

func newMockBitcoinChain() *mockBitcoinChain {
	return &mockBitcoinChain{
		transactions: make(map[bitcoin.Hash]*bitcoin.Transaction),
	}
}

func (mbc *mockBitcoinChain) GetTransaction(
	transactionHash bitcoin.Hash,
) (*bitcoin.Transaction, error) {
	if transaction, exists := mbc.transactions[transactionHash]; exists {
		return transaction, nil
	}

	return nil, fmt.Errorf("transaction not found")
}

func (mbc *mockBitcoinChain) GetTransactionConfirmations(
	transactionHash bitcoin.Hash,
) (uint, error) {
	panic("not implemented")
}

func (mbc *mockBitcoinChain) BroadcastTransaction(
	transaction *bitcoin.Transaction,
) error {
	panic("not implemented")
}

func (mbc *mockBitcoinChain) GetCurrentBlockNumber() (uint, error) {
	panic("not implemented")
}

func (mbc *mockBitcoinChain) GetBlockHeader(
	blockNumber uint,
) (*bitcoin.BlockHeader, error) {
	panic("not implemented")
}

func (mbc *mockBitcoinChain) addTransaction(
	transaction *bitcoin.Transaction,
) error {
	transactionHash := transaction.Hash()

	if _, exists := mbc.transactions[transactionHash]; exists {
		return fmt.Errorf("transaction already exists")
	}

	mbc.transactions[transactionHash] = transaction

	return nil
}
