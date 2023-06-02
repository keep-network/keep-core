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

const blockDelta = 2

type testConfig struct {
	clientConfig  Config
	errorMessages expectedErrorMessages
}

// Servers details were taken from a public Electrum servers list published
// at https://1209k.com/bitcoin-eye/ele.php?chain=tbtc.
var testConfigs = map[string]testConfig{
	"electrs-esplora tcp": {
		clientConfig: Config{
			URL:                 "tcp://electrum.blockstream.info:60001",
			RequestTimeout:      timeout * 3,
			RequestRetryTimeout: timeout * 10,
		},
	},
	"electrs-esplora ssl": {
		clientConfig: Config{
			URL:                 "ssl://electrum.blockstream.info:60002",
			RequestTimeout:      timeout * 3,
			RequestRetryTimeout: timeout * 10,
		},
	},
	"electrumx tcp": {
		clientConfig: Config{
			URL:                 "tcp://tn.not.fyi:55001",
			RequestTimeout:      timeout,
			RequestRetryTimeout: timeout * 2,
		},
	},
	"electrumx wss": {
		clientConfig: Config{
			URL:                 "wss://electrumx-server.test.tbtc.network:8443",
			RequestTimeout:      timeout,
			RequestRetryTimeout: timeout * 2,
		},
	},
	"fulcrum tcp": {
		clientConfig: Config{
			URL:                 "tcp://testnet.aranguren.org:51001",
			RequestTimeout:      timeout,
			RequestRetryTimeout: timeout * 2,
		},
	},
}

var invalidTxID bitcoin.Hash

func init() {
	invalidTxID, err = bitcoin.NewHashFromString(
		"9489457dc2c5a461a0b86394741ef57731605f2c628102de9f4d90afee9ac794",
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		panic(err)
	}
}

func TestConnect_Integration(t *testing.T) {
	for testName, testConfig := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			_, cancelCtx := newTestConnection(t, testConfig.clientConfig)
			defer cancelCtx()
		})
	}
}

func TestGetTransaction_Integration(t *testing.T) {
	for testName, testConfig := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			for txName, tx := range testData.Transactions {
			electrum, cancelCtx := newTestConnection(t, testConfig.clientConfig)
			defer cancelCtx()
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
			electrum, cancelCtx := newTestConnection(t, testConfig.clientConfig)
			defer cancelCtx()

			_, err := electrum.GetTransaction(invalidTxID)

			assertMissingTransactionError(
				t,
				testConfig.clientConfig,
				fmt.Sprintf(
					"failed to get raw transaction with ID [%s]",
					invalidTxID.Hex(bitcoin.ReversedByteOrder),
				),
				err,
			)
		})
	}
}

func TestGetTransactionConfirmations_Integration(t *testing.T) {
	for testName, testConfig := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum, cancelCtx := newTestConnection(t, testConfig.clientConfig)
			defer cancelCtx()

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

					assertNumberCloseTo(t, expectedConfirmations, result, blockDelta)
				})
			}
		})
	}
}

func TestGetTransactionConfirmations_Negative_Integration(t *testing.T) {
	for testName, testConfig := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum, cancelCtx := newTestConnection(t, testConfig.clientConfig)
			defer cancelCtx()

			_, err := electrum.GetTransactionConfirmations(invalidTxID)

			assertMissingTransactionError(
				t,
				testConfig.clientConfig,
				fmt.Sprintf(
					"failed to get raw transaction with ID [%s]",
					invalidTxID.Hex(bitcoin.ReversedByteOrder),
				),
				err,
			)
		})
	}
}

func TestGetLatestBlockHeight_Integration(t *testing.T) {
	expectedResult := uint(2404094)

	for testName, testConfig := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum, cancelCtx := newTestConnection(t, testConfig.clientConfig)
			defer cancelCtx()

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

	for testName, testConfig := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum, cancelCtx := newTestConnection(t, testConfig.clientConfig)
			defer cancelCtx()

			_, err := electrum.GetBlockHeader(blockHeight)

			assertMissingBlockHeaderError(
				t,
				testConfig.clientConfig,
				"failed to get block header",
				err,
			)
		})
	}
}

func TestGetTransactionMerkleProof_Integration(t *testing.T) {
	transactionHash := testData.TxMerkleProof.TxHash
	blockHeight := testData.TxMerkleProof.BlockHeigh

	expectedResult := &testData.TxMerkleProof.MerkleProof

	for testName, testConfig := range testConfigs {
		t.Run(testName, func(t *testing.T) {
		electrum, cancelCtx := newTestConnection(t, testConfig.clientConfig)
			defer cancelCtx()

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
	blockHeight := uint(123456)

	for testName, testConfig := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum, cancelCtx := newTestConnection(t, testConfig.clientConfig)
			defer cancelCtx()

			_, err := electrum.GetTransactionMerkleProof(
				invalidTxID,
				blockHeight,
			)

			assertMissingTransactionInBlockError(
				t,
				testConfig.clientConfig,
				"failed to get merkle proof",
				err,
			)
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
			electrum, cancelCtx := newTestConnection(t, testConfig.clientConfig)
			defer cancelCtx()

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

func TestEstimateSatPerVByteFee_Integration(t *testing.T) {
	for testName, testConfig := range testConfigs {
		t.Run(testName, func(t *testing.T) {
			electrum, cancelCtx := newTestConnection(t, testConfig.clientConfig)
			defer cancelCtx()

			satPerVByteFee, err := electrum.EstimateSatPerVByteFee(1)
			if err != nil {
				t.Fatal(err)
			}

			// We expect the fee is always at least 1.
			if satPerVByteFee < 1 {
				t.Errorf("returned fee is below 1")
			}
		})
	}
}

func newTestConnection(t *testing.T, config electrum.Config) (bitcoin.Chain, context.CancelFunc) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	electrum, err := electrum.Connect(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	return electrum, cancelCtx
}

func assertNumberCloseTo(t *testing.T, expected uint, actual uint, delta uint) {
	min := expected - delta
	max := expected + delta

	if min > actual || actual > max {
		t.Errorf(
			"value %d is out of expected range: [%d,%d]",
			actual,
			min,
			max,
		)
	}
}

type expectedErrorMessages struct {
	missingTransaction        []string
	missingBlockHeader        []string
	missingTransactionInBlock []string
}

var expectedServerErrorMessages = expectedErrorMessages{
	missingTransaction: []string{
		"errNo: 0, errMsg: missing transaction",
		"errNo: 2, errMsg: daemon error: DaemonError({'code': -5, 'message': 'No such mempool or blockchain transaction. Use gettransaction for wallet transactions.'})",
		"errNo: 2, errMsg: daemon error: DaemonError({'message': 'Transaction not found.', 'code': -1})",
	},
	missingBlockHeader: []string{
		"errNo: 0, errMsg: missing header",
		"errNo: 1, errMsg: height 4,294,967,295 out of range",
		"errNo: 1, errMsg: Invalid height",
	},
	missingTransactionInBlock: []string{
		"errNo: 0, errMsg: tx not found or is unconfirmed",
		"errNo: 1, errMsg: tx 9489457dc2c5a461a0b86394741ef57731605f2c628102de9f4d90afee9ac794 not in block at height 123,456",
		"errNo: 1, errMsg: No transaction matching the requested hash found at height 123456"},
}

func assertMissingTransactionError(
	t *testing.T,
	clientConfig electrum.Config,
	clientErrorPrefix string,

	actualError error,
) {
	assertServerError(
		t,
		clientConfig,
		clientErrorPrefix,
		expectedServerErrorMessages.missingTransaction,
		actualError,
	)
}

func assertMissingBlockHeaderError(
	t *testing.T,
	clientConfig electrum.Config,
	clientErrorPrefix string,
	actualError error,
) {
	assertServerError(
		t,
		clientConfig,
		clientErrorPrefix,
		expectedServerErrorMessages.missingBlockHeader,
		actualError,
	)
}

func assertMissingTransactionInBlockError(
	t *testing.T,
	clientConfig electrum.Config,
	clientErrorPrefix string,
	actualError error,
) {
	assertServerError(
		t,
		clientConfig,
		clientErrorPrefix,
		expectedServerErrorMessages.missingTransactionInBlock,
		actualError,
	)
}

func assertServerError(
	t *testing.T,
	clientConfig electrum.Config,
	clientErrorPrefix string,
	expectedServerErrors []string,
	actualError error,
) {
	expectedErrorMsgFormat := fmt.Sprintf(
		"%s: [retry timeout [%s] exceeded; most recent error: [request failed: [%%s]]]",
		clientErrorPrefix,
		clientConfig.RequestRetryTimeout,
	)

	expectedErrorMsgStrings := make([]string, len(expectedServerErrors))
	for i, serverError := range expectedServerErrors {
		expectedErrorMsgStrings[i] = fmt.Sprintf(expectedErrorMsgFormat, serverError)
	}

	if actualError == nil {
		t.Errorf("expected error, but actual error is nil")
		return
	}

	if !slices.Contains(expectedErrorMsgStrings, actualError.Error()) {
		t.Errorf(
			"unexpected error message\nactual:\n\t%v\nexpected one of:\n\t%s",
			actualError,
			strings.Join(expectedErrorMsgStrings, "\n\t"),
		)
		return
	}
}
