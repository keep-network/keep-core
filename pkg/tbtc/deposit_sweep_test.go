package tbtc

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/tecdsa"

	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc/internal/test"
)

// TODO: Think about covering unhappy paths for specific steps of the deposit sweep action.
func TestDepositSweepAction_Execute(t *testing.T) {
	scenarios, err := test.LoadDepositSweepTestScenarios()
	if err != nil {
		t.Fatal(err)
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Title, func(t *testing.T) {
			now := time.Now()

			hostChain := Connect()
			bitcoinChain := newLocalBitcoinChain()

			wallet := wallet{
				// Set only relevant fields.
				publicKey: scenario.WalletPublicKey,
			}
			walletPublicKeyHash := bitcoin.PublicKeyHash(wallet.publicKey)

			// Record the transactions that will serve as sweep transaction's
			// input in the Bitcoin local chain.
			for _, transaction := range scenario.InputTransactions {
				err := bitcoinChain.BroadcastTransaction(transaction)
				if err != nil {
					t.Fatal(err)
				}
			}

			// depositsKeys will be needed to build the proposal instance.
			depositsKeys := make([]struct {
				FundingTxHash      bitcoin.Hash
				FundingOutputIndex uint32
			}, len(scenario.Deposits))

			// depositsExtraInfo will be needed to perform on-chain proposal
			// validation.
			depositsExtraInfo := make([]struct {
				*Deposit
				FundingTx *bitcoin.Transaction
			}, len(scenario.Deposits))

			depositsRevealBlocks := make([]*big.Int, len(scenario.Deposits))

			// Record all necessary deposits' data on the local host chain.
			for i, deposit := range scenario.Deposits {
				fundingTxHash := deposit.Utxo.Outpoint.TransactionHash
				fundingOutputIndex := deposit.Utxo.Outpoint.OutputIndex

				fundingTx, err := bitcoinChain.GetTransaction(fundingTxHash)
				if err != nil {
					t.Fatal(err)
				}

				depositsKeys[i] = struct {
					FundingTxHash      bitcoin.Hash
					FundingOutputIndex uint32
				}{
					FundingTxHash:      fundingTxHash,
					FundingOutputIndex: fundingOutputIndex,
				}

				depositsExtraInfo[i] = struct {
					*Deposit
					FundingTx *bitcoin.Transaction
				}{
					Deposit:   (*Deposit)(deposit),
					FundingTx: fundingTx,
				}

				// Build the deposit reveal block based on the deposit index.
				// This field can be an arbitrary value, but it is good to keep
				// it consistent.
				depositRevealBlock := uint64(100 * i)
				depositsRevealBlocks[i] = big.NewInt(int64(depositRevealBlock))

				// The deposit sweep action will look for past deposit
				// revealed events using a specific filter. We need to make
				// sure the local host chain will return the expected result
				// for that filter. We need to build the startBlock and endBlock
				// filter's parameters using the depositRevealBlock value
				// as done within the depositSweepAction.execute function.
				// We also need to use the correct wallet PKH.
				err = hostChain.setPastDepositRevealedEvents(
					&DepositRevealedEventFilter{
						StartBlock:          depositRevealBlock,
						EndBlock:            &depositRevealBlock,
						WalletPublicKeyHash: [][20]byte{walletPublicKeyHash},
					},
					[]*DepositRevealedEvent{
						{
							FundingTxHash:       fundingTxHash,
							FundingOutputIndex:  fundingOutputIndex,
							Depositor:           deposit.Depositor,
							Amount:              uint64(deposit.Utxo.Value),
							BlindingFactor:      deposit.BlindingFactor,
							WalletPublicKeyHash: deposit.WalletPublicKeyHash,
							RefundPublicKeyHash: deposit.RefundPublicKeyHash,
							RefundLocktime:      deposit.RefundLocktime,
							Vault:               deposit.Vault,
							BlockNumber:         depositRevealBlock,
						},
					},
				)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Build the sweep proposal based on the scenario data.
			proposal := &DepositSweepProposal{
				WalletPublicKeyHash:  walletPublicKeyHash,
				DepositsKeys:         depositsKeys,
				SweepTxFee:           big.NewInt(scenario.Fee),
				DepositsRevealBlocks: depositsRevealBlocks,
			}

			// Choose an arbitrary start block and expiration time.
			proposalProcessingStartBlock := uint64(100)
			proposalExpiresAt := now.Add(4 * time.Hour)

			// Simulate the on-chain proposal validation passes with success.
			err = hostChain.setDepositSweepProposalValidationResult(
				proposal,
				depositsExtraInfo,
				true,
			)
			if err != nil {
				t.Fatal(err)
			}

			// Record the wallet main UTXO hash in the local host chain so
			// the deposit action can detect it.
			var walletMainUtxoHash [32]byte
			if scenario.WalletMainUtxo != nil {
				walletMainUtxoHash = hostChain.ComputeMainUtxoHash(
					scenario.WalletMainUtxo,
				)
			}
			hostChain.setWallet(walletPublicKeyHash, &WalletChainData{
				MainUtxoHash: walletMainUtxoHash,
			})

			// Create a signing executor mock instance.
			signingExecutor := newMockWalletSigningExecutor()

			// The signatures within the scenario fixture are in the format
			// suitable for applying them directly to a Bitcoin transaction.
			// However, the signing executor operates on raw tECDSA signatures
			// so, we need to unpack them first.
			rawSignatures := make([]*tecdsa.Signature, len(scenario.Signatures))
			for i, signature := range scenario.Signatures {
				rawSignatures[i] = &tecdsa.Signature{
					R: signature.R,
					S: signature.S,
				}
			}

			// Set up the signing executor mock to return the signatures from
			// the test fixture when called with the expected parameters.
			// Note that the start block is set based on the proposal
			// processing start block as done within the action.
			signingExecutor.setSignatures(
				scenario.ExpectedSigHashes,
				proposalProcessingStartBlock,
				rawSignatures,
			)

			action := newDepositSweepAction(
				logger.With(),
				hostChain,
				bitcoinChain,
				wallet,
				signingExecutor,
				proposal,
				proposalProcessingStartBlock,
				proposalExpiresAt,
			)

			// Modify the default parameters of the action to make
			// it possible to execute in the current test environment.
			action.requiredFundingTxConfirmations = 1
			action.broadcastCheckDelay = 1 * time.Second

			err := action.execute()
			if err != nil {
				t.Fatal(err)
			}

			// Action execution that completes without an error is a sign of
			// success. However, just in case, make an additional check that
			// the expected sweep transaction was actually broadcasted on the
			// local Bitcoin chain.
			broadcastedSweepTransaction, err := bitcoinChain.GetTransaction(
				scenario.ExpectedSweepTransactionHash,
			)
			if err != nil {
				t.Fatal(err)
			}

			testutils.AssertBytesEqual(
				t,
				scenario.ExpectedSweepTransaction.Serialize(),
				broadcastedSweepTransaction.Serialize(),
			)
		})
	}
}

func TestAssembleDepositSweepTransaction(t *testing.T) {
	scenarios, err := test.LoadDepositSweepTestScenarios()
	if err != nil {
		t.Fatal(err)
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Title, func(t *testing.T) {
			bitcoinChain := newLocalBitcoinChain()

			for _, transaction := range scenario.InputTransactions {
				err := bitcoinChain.BroadcastTransaction(transaction)
				if err != nil {
					t.Fatal(err)
				}
			}

			deposits := make([]*Deposit, len(scenario.Deposits))
			for i, d := range scenario.Deposits {
				deposits[i] = &Deposit{
					Utxo:                d.Utxo,
					Depositor:           d.Depositor,
					BlindingFactor:      d.BlindingFactor,
					WalletPublicKeyHash: d.WalletPublicKeyHash,
					RefundPublicKeyHash: d.RefundPublicKeyHash,
					RefundLocktime:      d.RefundLocktime,
					Vault:               d.Vault,
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
				scenario.ExpectedSweepTransactionHash.Hex(bitcoin.InternalByteOrder),
				transaction.Hash().Hex(bitcoin.InternalByteOrder),
			)
			testutils.AssertStringsEqual(
				t,
				"sweep transaction witness hash",
				scenario.ExpectedSweepTransactionWitnessHash.Hex(bitcoin.InternalByteOrder),
				transaction.WitnessHash().Hex(bitcoin.InternalByteOrder),
			)
		})
	}
}
