//go:build integration
// +build integration

package electrum

import (
	"context"
	"encoding/hex"
	"math"
	"testing"
	"time"

	"github.com/go-test/deep"
	"golang.org/x/exp/slices"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

// TODO: Add negative tests
// TODO: Include integration test in the CI.
// To run the tests execute `go test -v -tags=integration ./...`

const TxID = "c580e0e352570d90e303d912a506055ceeb0ee06f97dce6988c69941374f5479"

var transactionHash bitcoin.Hash

// Servers details were taken from a public Electrum servers list published
// at https://1209k.com/bitcoin-eye/ele.php?chain=tbtc.
var configs = map[string]Config{
	"electrs-esplora tcp": {
		URL:                 "electrum.blockstream.info:60001",
		Protocol:            TCP,
		RequestRetryTimeout: 1 * time.Second,
	},
	"electrs-esplora ssl": {
		URL:                 "electrum.blockstream.info:60002",
		Protocol:            SSL,
		RequestRetryTimeout: 1 * time.Second,
	},
	"electrumx ssl": {
		URL:                 "testnet.hsmiths.com:53012",
		Protocol:            SSL,
		RequestRetryTimeout: 1 * time.Second,
	},
	"fulcrum ssl": {
		URL:                 "blackie.c3-soft.com:57006",
		Protocol:            SSL,
		RequestRetryTimeout: 1 * time.Second,
	},
	// TODO: Add Keep's electrum server
}

func init() {
	var err error
	transactionHash, err = bitcoin.NewHashFromString(
		TxID,
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		panic(err)
	}
}

func TestGetTransaction_Integration(t *testing.T) {
	expectedResult := bitcoinTestTx(t)

	for testName, config := range configs {
		t.Run(testName, func(t *testing.T) {
			electrs := newTestConnection(t, config)

			result, err := electrs.GetTransaction(transactionHash)
			if err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(result, expectedResult); diff != nil {
				t.Errorf("compare failed: %v", diff)
			}
		})
	}
}

func TestGetTransactionConfirmations_Integration(t *testing.T) {
	expectedResult := uint(271247)

	for testName, config := range configs {
		t.Run(testName, func(t *testing.T) {
			electrs := newTestConnection(t, config)

			result, err := electrs.GetTransactionConfirmations(transactionHash)
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

func TestGetLatestBlockHeight_Integration(t *testing.T) {
	expectedResult := uint(2404094)

	for testName, config := range configs {
		t.Run(testName, func(t *testing.T) {
			electrs := newTestConnection(t, config)

			result, err := electrs.GetLatestBlockHeight()
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
			electrs := newTestConnection(t, config)

			result, err := electrs.GetBlockHeader(blockHeight)
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
			electrs := newTestConnection(t, config)

			expectedErrorMsg := "failed to get block header: [retry timeout [1s] exceeded; most recent error: [GetBlockHeader failed: [missing header]]]"

			// As a workaround for the problem described in https://github.com/checksum0/go-electrum/issues/5
			// we use an alternative expected error message for servers
			// that are not correctly supported by the electrum client.
			if slices.Contains(replaceErrorMsgForTests, testName) {
				expectedErrorMsg = "failed to get block header: [retry timeout [1s] exceeded; most recent error: [GetBlockHeader failed: [Unmarshal received message failed: json: cannot unmarshal object into Go struct field response.error of type string]]]"
			}

			_, err := electrs.GetBlockHeader(blockHeight)
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
	electrs, err := Connect(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	return electrs
}

func bitcoinTestTx(t *testing.T) *bitcoin.Transaction {
	prevTxHash, err := bitcoin.NewHashFromString(
		"e788a344a86f7e369511fe37ebd1d74686dde694ee99d06db5db3d4a14719b1d",
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		t.Fatal(err)
	}

	txInScript, err := hex.DecodeString("47304402206f8553c07bcdc0c3b906311888103d623ca9096ca0b28b7d04650a029a01fcf9022064cda02e39e65ace712029845cfcf58d1b59617d753c3fd3556f3551b609bbb00121039d61d62dcd048d3f8550d22eb90b4af908db60231d117aeede04e7bc11907bfa")
	if err != nil {
		t.Fatal(err)
	}

	txOutScript0, err := hex.DecodeString("a9143ec459d0f3c29286ae5df5fcc421e2786024277e87")
	if err != nil {
		t.Fatal(err)
	}
	txOutScript1, err := hex.DecodeString("0014e257eccafbc07c381642ce6e7e55120fb077fbed")
	if err != nil {
		t.Fatal(err)
	}

	return &bitcoin.Transaction{
		Version: 1,
		Inputs: []*bitcoin.TransactionInput{
			{
				Outpoint: &bitcoin.TransactionOutpoint{
					TransactionHash: prevTxHash,
					OutputIndex:     1,
				},
				SignatureScript: txInScript,
				Sequence:        4294967295,
			},
		},
		Outputs: []*bitcoin.TransactionOutput{
			{
				PublicKeyScript: txOutScript0,
				Value:           20000,
			},
			{
				PublicKeyScript: txOutScript1,
				Value:           1360550,
			},
		},
		Locktime: 0,
	}
}
