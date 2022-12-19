//go:build integration
// +build integration

package electrum

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/go-test/deep"
	"golang.org/x/exp/slices"

	"github.com/keep-network/keep-core/pkg/bitcoin"

	testData "github.com/keep-network/keep-core/internal/testdata/bitcoin"
)

// TODO: Include integration test in the CI.
// To run the tests execute `go test -v -tags=integration ./...`

const timeout = 2 * time.Second

// Servers details were taken from a public Electrum servers list published
// at https://1209k.com/bitcoin-eye/ele.php?chain=tbtc.
var configs = map[string]Config{
	"electrs-esplora tcp": {
		URL:                 "electrum.blockstream.info:60001",
		Protocol:            TCP,
		RequestTimeout:      timeout,
		RequestRetryTimeout: timeout * 2,
	},
	"electrs-esplora ssl": {
		URL:                 "electrum.blockstream.info:60002",
		Protocol:            SSL,
		RequestTimeout:      timeout,
		RequestRetryTimeout: timeout * 3,
	},
	"electrumx ssl": {
		URL:                 "testnet.hsmiths.com:53012",
		Protocol:            SSL,
		RequestTimeout:      timeout,
		RequestRetryTimeout: timeout * 2,
	},
	"fulcrum ssl": {
		URL:                 "blackie.c3-soft.com:57006",
		Protocol:            SSL,
		RequestTimeout:      timeout,
		RequestRetryTimeout: timeout * 2,
	},
	// TODO: Add Keep's electrum server
}

func TestGetTransaction_Integration(t *testing.T) {
	for testName, config := range configs {
		electrum := newTestConnection(t, config)

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

	replaceErrorMsgForTests := []string{"electrumx ssl", "fulcrum ssl"}

	for testName, config := range configs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config)

			expectedErrorMsg := fmt.Sprintf(
				"failed to get raw transaction with ID [%s]: [retry timeout [%s] exceeded; most recent error: [request failed: [missing transaction]]]",
				invalidTxID.Hex(bitcoin.ReversedByteOrder),
				config.RequestRetryTimeout,
			)

			// As a workaround for the problem described in https://github.com/checksum0/go-electrum/issues/5
			// we use an alternative expected error message for servers
			// that are not correctly supported by the electrum client.
			if slices.Contains(replaceErrorMsgForTests, testName) {
				expectedErrorMsg = fmt.Sprintf(
					"failed to get raw transaction with ID [%s]: [retry timeout [%s] exceeded; most recent error: [request failed: [Unmarshal received message failed: json: cannot unmarshal object into Go struct field response.error of type string]]]",
					invalidTxID.Hex(bitcoin.ReversedByteOrder),
					config.RequestRetryTimeout,
				)
			}

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
	for testName, config := range configs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config)

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

	replaceErrorMsgForTests := []string{"electrumx ssl", "fulcrum ssl"}

	for testName, config := range configs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config)

			expectedErrorMsg := fmt.Sprintf(
				"failed to get raw transaction with ID [%s]: [retry timeout [%s] exceeded; most recent error: [request failed: [missing transaction]]]",
				invalidTxID.Hex(bitcoin.ReversedByteOrder),
				config.RequestRetryTimeout,
			)

			// As a workaround for the problem described in https://github.com/checksum0/go-electrum/issues/5
			// we use an alternative expected error message for servers
			// that are not correctly supported by the electrum client.
			if slices.Contains(replaceErrorMsgForTests, testName) {
				expectedErrorMsg = fmt.Sprintf(
					"failed to get raw transaction with ID [%s]: [retry timeout [%s] exceeded; most recent error: [request failed: [Unmarshal received message failed: json: cannot unmarshal object into Go struct field response.error of type string]]]",
					invalidTxID.Hex(bitcoin.ReversedByteOrder),
					config.RequestRetryTimeout,
				)
			}

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

	for testName, config := range configs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config)

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

	for testName, config := range configs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config)

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

	replaceErrorMsgForTests := []string{"electrumx ssl", "fulcrum ssl"}

	for testName, config := range configs {
		t.Run(testName, func(t *testing.T) {
			electrum := newTestConnection(t, config)

			expectedErrorMsg := fmt.Sprintf(
				"failed to get block header: [retry timeout [%s] exceeded; most recent error: [request failed: [missing header]]]",
				config.RequestRetryTimeout,
			)

			// As a workaround for the problem described in https://github.com/checksum0/go-electrum/issues/5
			// we use an alternative expected error message for servers
			// that are not correctly supported by the electrum client.
			if slices.Contains(replaceErrorMsgForTests, testName) {
				expectedErrorMsg = fmt.Sprintf(
					"failed to get block header: [retry timeout [%s] exceeded; most recent error: [request failed: [Unmarshal received message failed: json: cannot unmarshal object into Go struct field response.error of type string]]]",
					config.RequestRetryTimeout,
				)
			}

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
