// go:build integration

package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/bitcoin/electrum"
)

var testBitcoinNetworks = []bitcoin.Network{bitcoin.Testnet, bitcoin.Mainnet}

func TestEmbeddedElectrumServerConnect_Integration(t *testing.T) {
	for _, bitcoinNetwork := range testBitcoinNetworks {
		urls, err := readElectrumUrls(bitcoinNetwork)
		if err != nil {
			t.Error(err)
		}

		for _, url := range urls {
			t.Run(fmt.Sprintf("%s/%s", bitcoinNetwork, url), func(t *testing.T) {
				newTestConnection(t, url)
			})
		}
	}
}

func newTestConnection(t *testing.T, url string) bitcoin.Chain {
	electrum, err := electrum.Connect(context.Background(), electrum.Config{URL: url})
	if err != nil {
		t.Fatal(err)
	}

	return electrum
}
