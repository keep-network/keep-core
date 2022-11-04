//go:build integration
// +build integration

package electrs

import (
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

// TODO: Include integration test in the CI.
// To run the tests execute `go test -v -tags=integration ./...`

func TestGetTransaction_Integration(t *testing.T) {
	electrs := newBlockstreamConnection(t)

	testGetTransaction(t, electrs)
}

func TestGetTransactionConfirmations_Integration(t *testing.T) {
	electrs := newBlockstreamConnection(t)

	expectedResult := uint(268506)

	testGetTransactionConfirmations(t, electrs, expectedResult, false)
}

func TestGetLatestBlockHeight_Integration(t *testing.T) {
	electrs := newBlockstreamConnection(t)

	expectedResult := uint(2404094)

	testGetLatestBlockHeight(t, electrs, expectedResult, false)
}

func TestGetBlockHeader_Integration(t *testing.T) {
	electrs := newBlockstreamConnection(t)

	testGetBlockHeader(t, electrs)
}

func newBlockstreamConnection(t *testing.T) bitcoin.Chain {
	config := Config{
		URL:            "https://blockstream.info/testnet/api/",
		RequestTimeout: 1 * time.Second,
		RetryTimeout:   15 * time.Second,
	}
	electrs, err := Connect(config)
	if err != nil {
		t.Fatal(err)
	}

	return electrs
}
