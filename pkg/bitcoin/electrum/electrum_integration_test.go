//go:build integration
// +build integration

package electrum

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/go-test/deep"

	"github.com/keep-network/keep-core/pkg/bitcoin"

	testData "github.com/keep-network/keep-core/internal/testdata/bitcoin"
)

// TODO: Include integration test in the CI.
// To run the tests execute `go test -v -tags=integration ./...`

const timeout = 2 * time.Second

type serverImplementation int

const (
	electrumX serverImplementation = iota
	fulcrum
	esploraElectrs
)

type testConfig struct {
	clientConfig         Config
	serverImplementation serverImplementation
}

// Servers details were taken from a public Electrum servers list published
// at https://1209k.com/bitcoin-eye/ele.php?chain=tbtc.
var testConfigs = map[string]testConfig{
	"electrs-esplora tcp": {
		clientConfig: Config{
			URL:                 "electrum.blockstream.info:60001",
			Protocol:            TCP,
			RequestTimeout:      timeout,
			RequestRetryTimeout: timeout * 2,
		},
		serverImplementation: esploraElectrs,
	},
	"electrs-esplora ssl": {
		clientConfig: Config{
			URL:                 "electrum.blockstream.info:60002",
			Protocol:            SSL,
			RequestTimeout:      timeout,
			RequestRetryTimeout: timeout * 3,
		},
		serverImplementation: esploraElectrs,
	},
	"electrumx ssl": {
		clientConfig: Config{
			URL:                 "testnet.qtornado.com:51002",
			Protocol:            SSL,
			RequestTimeout:      timeout,
			RequestRetryTimeout: timeout * 2,
		},
		serverImplementation: electrumX,
	},
	"fulcrum ssl": {
		clientConfig: Config{
			URL:                 "blackie.c3-soft.com:57006",
			Protocol:            SSL,
			RequestTimeout:      timeout,
			RequestRetryTimeout: timeout * 2,
		},
		serverImplementation: fulcrum,
	},
	// TODO: Add Keep's electrum server
}

var (
	missingTransactionServerMsgs = map[serverImplementation]string{
		electrumX:      "errNo: 2, errMsg: daemon error: DaemonError({'code': -5, 'message': 'No such mempool or blockchain transaction. Use gettransaction for wallet transactions.'})",
		fulcrum:        "errNo: 2, errMsg: daemon error: DaemonError({'code': -5, 'message': 'No such mempool or blockchain transaction. Use gettransaction for wallet transactions.'})",
		esploraElectrs: "errNo: 0, errMsg: missing transaction",
	}

	missingTransactionInBlockMsgs = map[serverImplementation]string{
		electrumX:      "errNo: 1, errMsg: tx aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa not in block at height 123,456",
		fulcrum:        "errNo: 1, errMsg: No transaction matching the requested hash found at height 123456",
		esploraElectrs: "errNo: 0, errMsg: tx not found or is unconfirmed",
	}

	missingHeaderServerMsgs = map[serverImplementation]string{
		electrumX:      "errNo: 1, errMsg: height 4,294,967,295 out of range",
		fulcrum:        "errNo: 1, errMsg: Invalid height",
		esploraElectrs: "errNo: 0, errMsg: missing header",
	}
)

func TestGetTransaction_Integration(t *testing.T) {
	for testName, config := range testConfigs {
		electrum := newTestConnection(t, config.clientConfig)

		t.Run(testName, func(t *testing.T) {
			for txName, tx := range testData.Transactions {
				t.Run(txName, func(t *testing.T) {
					result, err := electrum.GetTransaction(tx.TxHash)
					if err != nil {
						t.Fatal(err)
					}

					if diff := deep.Equal(result, &tx.BitcoinTx); diff != nil {
						t.Errorf("compare failed: %v", diff)
					}
				})
			}
		})
	}
}

func TestGetTransaction_Negative_Integration(t *testing.T) {
	invalidTxID, err := bitcoin.NewHashFromString(
		"ecc246ac58e682c8edccabb6476bb5482df541847b774085cdb8bfc53165cd34",
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		t.Fatal(err)
	}

	for testName, config := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config.clientConfig)

			expectedErrorMsg := fmt.Sprintf(
				"failed to get raw transaction with ID [%s]: [retry timeout [%s] exceeded; most recent error: [request failed: [%s]]]",
				invalidTxID.Hex(bitcoin.ReversedByteOrder),
				config.clientConfig.RequestRetryTimeout,
				missingTransactionServerMsgs[config.serverImplementation],
			)

			_, err := electrum.GetTransaction(invalidTxID)
			if err.Error() != expectedErrorMsg {
				t.Errorf(
					"invalid error\nexpected: %v\nactual:   %v",
					expectedErrorMsg,
					err,
				)
			}
		})
	}
}

func TestGetTransactionConfirmations_Integration(t *testing.T) {
	for testName, config := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config.clientConfig)

			for txName, tx := range testData.Transactions {
				t.Run(txName, func(t *testing.T) {
					latestBlockHeight, err := electrum.GetLatestBlockHeight()
					if err != nil {
						t.Fatalf("failed to get the latest block height: %s", err)
					}
					expectedConfirmations := latestBlockHeight - tx.BlockHeight

					result, err := electrum.GetTransactionConfirmations(tx.TxHash)
					if err != nil {
						t.Fatal(err)
					}

					assertConfirmationsCloseTo(t, expectedConfirmations, result)
				})
			}
		})
	}
}

func TestGetTransactionConfirmations_Negative_Integration(t *testing.T) {
	invalidTxID, err := bitcoin.NewHashFromString(
		"ecc246ac58e682c8edccabb6476bb5482df541847b774085cdb8bfc53165cd34",
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		t.Fatal(err)
	}

	for testName, config := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config.clientConfig)

			expectedErrorMsg := fmt.Sprintf(
				"failed to get raw transaction with ID [%s]: [retry timeout [%s] exceeded; most recent error: [request failed: [%s]]]",
				invalidTxID.Hex(bitcoin.ReversedByteOrder),
				config.clientConfig.RequestRetryTimeout,
				missingTransactionServerMsgs[config.serverImplementation],
			)

			_, err := electrum.GetTransactionConfirmations(invalidTxID)
			if err.Error() != expectedErrorMsg {
				t.Errorf(
					"invalid error\nexpected: %v\nactual:   %v",
					expectedErrorMsg,
					err,
				)
			}
		})
	}
}

func TestGetLatestBlockHeight_Integration(t *testing.T) {
	expectedResult := uint(2404094)

	for testName, config := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config.clientConfig)

			result, err := electrum.GetLatestBlockHeight()
			if err != nil {
				t.Fatal(err)
			}

			if result < expectedResult {
				t.Errorf(
					"invalid result (greater or equal match)\nexpected: %v\nactual:   %v",
					expectedResult,
					result,
				)
			}
		})
	}
}

func TestGetBlockHeader_Integration(t *testing.T) {
	blockHeight := uint(2135502)

	previousBlockHeaderHash, err := bitcoin.NewHashFromString(
		"000000000066450030efdf72f233ed2495547a32295deea1e2f3a16b1e50a3a5",
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		t.Fatal(err)
	}

	merkleRootHash, err := bitcoin.NewHashFromString(
		"1251774996b446f85462d5433f7a3e384ac1569072e617ab31e86da31c247de2",
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		t.Fatal(err)
	}

	expectedResult := &bitcoin.BlockHeader{
		Version:                 536870916,
		PreviousBlockHeaderHash: previousBlockHeaderHash,
		MerkleRootHash:          merkleRootHash,
		Time:                    1641914003,
		Bits:                    436256810,
		Nonce:                   778087099,
	}

	for testName, config := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config.clientConfig)

			result, err := electrum.GetBlockHeader(blockHeight)
			if err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(result, expectedResult); diff != nil {
				t.Errorf("compare failed: %v", diff)
			}
		})
	}
}

func TestGetBlockHeader_Negative_Integration(t *testing.T) {
	blockHeight := uint(math.MaxUint32)

	for testName, config := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config.clientConfig)

			expectedErrorMsg := fmt.Sprintf(
				"failed to get block header: [retry timeout [%s] exceeded; most recent error: [request failed: [%s]]]",
				config.clientConfig.RequestRetryTimeout,
				missingHeaderServerMsgs[config.serverImplementation],
			)

			_, err := electrum.GetBlockHeader(blockHeight)
			if err.Error() != expectedErrorMsg {
				t.Errorf(
					"invalid error\nexpected: %v\nactual:   %v",
					expectedErrorMsg,
					err,
				)
			}
		})
	}
}

func TestGetTransactionMerkleProof_Integration(t *testing.T) {
	transactionHash, err := bitcoin.NewHashFromString(
		"4e210df8041914be65ec026f2963c3ae79ff867424c40523edb1adc257fde772",
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		t.Fatal(err)
	}

	blockHeight := uint(1569342)

	expectedResult := &bitcoin.TransactionMerkleProof{
		BlockHeight: 1569342,
		MerkleNodes: []string{
			"8b5bbb5bdf6727bf70fad4f46fe4eaab04c98119ffbd2d95c29adf32d26f8452",
			"53637bacb07965e4a8220836861d1b16c6da29f10ea9ab53fc4eca73074f98b9",
			"0267e738108d094ceb05217e2942e9c2a4c6389ac47f476f572c9a319ce4dfbc",
			"34e00deec50c48d99678ca2b52b82d6d5432326159c69e7233d0dde0924874b4",
			"7a53435e6c86a3620cdbae510901f17958f0540314214379197874ed8ed7a913",
			"6315dbb7ce350ceaa16cd4c35c5a147005e8b38ca1e9531bd7320629e8d17f5b",
			"40380cdadc0206646208871e952af9dcfdff2f104305ce463aed5eeaf7725d2f",
			"5d74bae6a71fd1cff2416865460583319a40343650bd4bb89de0a6ae82097037",
			"296ddccfc659e0009aad117c8ed15fb6ff81c2bade73fbc89666a22708d233f9",
		},
		Position: 176,
	}

	for testName, config := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config.clientConfig)

			result, err := electrum.GetTransactionMerkleProof(
				transactionHash,
				blockHeight,
			)
			if err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(result, expectedResult); diff != nil {
				t.Errorf("compare failed: %v", diff)
			}
		})
	}
}

func TestGetTransactionMerkleProof_Negative_Integration(t *testing.T) {
	for testName, config := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config.clientConfig)

			expectedErrorMsg := fmt.Sprintf(
				"failed to get merkle proof: [retry timeout [%s] exceeded; most recent error: [request failed: [%s]]]",
				config.clientConfig.RequestRetryTimeout,
				missingTransactionInBlockMsgs[config.serverImplementation],
			)

			transactionHash, err := bitcoin.NewHashFromString(
				// use incorrect hash
				"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				bitcoin.InternalByteOrder,
			)
			if err != nil {
				t.Fatal(err)
			}

			blockHeight := uint(123456)
			_, err = electrum.GetTransactionMerkleProof(
				transactionHash,
				blockHeight,
			)
			if err.Error() != expectedErrorMsg {
				t.Errorf(
					"invalid error\nexpected: %v\nactual:   %v",
					expectedErrorMsg,
					err,
				)
			}
		})
	}
}

func TestGetTransactionsForPublicKeyHash_Integration(t *testing.T) {
	var publicKeyHash [20]byte
	publicKeyHashBytes, err := hex.DecodeString("e6f9d74726b19b75f16fe1e9feaec048aa4fa1d0")
	if err != nil {
		t.Fatal(err)
	}
	copy(publicKeyHash[:], publicKeyHashBytes)

	// To determine the expected five latest transactions for comparison, we
	// use a block explorer to browse the history for the two addresses the
	// e6f9d74726b19b75f16fe1e9feaec048aa4fa1d0 public key hash translates to:
	//
	// - P2WPKH testnet address: https://live.blockcypher.com/btc-testnet/address/tb1qumuaw3exkxdhtut0u85latkqfz4ylgwstkdzsx
	// - P2PKH testnet address: https://live.blockcypher.com/btc-testnet/address/n2aF1Rj6PK26quhGRo8YoRQYjwm37Zjnkb
	//
	// Then, we take all transactions for both addresses and pick the latest five.
	expectedHashes := []string{
		"f65bc5029251f0042aedb37f90dbb2bfb63a2e81694beef9cae5ec62e954c22e",
		"44863a79ce2b8fec9792403d5048506e50ffa7338191db0e6c30d3d3358ea2f6",
		"4c6b33b7c0550e0e536a5d119ac7189d71e1296fcb0c258e0c115356895bc0e6",
		"605edd75ae0b4fa7cfc7aae8f1399119e9d7ecc212e6253156b60d60f4925d44",
		"4f9affc5b418385d5aa61e23caa0b55156bf0682d5fedf2d905446f3f88aec6c",
	}

	for testName, config := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config.clientConfig)

			transactions, err := electrum.GetTransactionsForPublicKeyHash(publicKeyHash, 5)
			if err != nil {
				t.Fatal(err)
			}

			hashes := make([]string, len(transactions))
			for i, transaction := range transactions {
				hash := transaction.Hash()
				hashes[i] = hash.Hex(bitcoin.ReversedByteOrder)
			}

			if !reflect.DeepEqual(expectedHashes, hashes) {
				t.Errorf(
					"unexpected transactions\nexpected: %v\nactual:   %v",
					expectedHashes,
					hashes,
				)
			}
		})
	}
}

func newTestConnection(t *testing.T, config Config) bitcoin.Chain {
	electrum, err := Connect(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	return electrum
}

func assertConfirmationsCloseTo(t *testing.T, expected uint, actual uint) {
	delta := uint(2)

	min := expected - delta
	max := expected + delta

	if min > actual || actual > max {
		t.Errorf(
			"confirmations number %d out of expected range: [%d,%d]",
			actual,
			min,
			max,
		)
	}
}
