package config

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/bitcoin/electrum"
	"golang.org/x/exp/slices"
	"reflect"
	"testing"
)

func TestResolveElectrum(t *testing.T) {
	var tests = map[string]struct {
		network         bitcoin.Network
		expectedConfigs []electrum.Config
		expectedError   error
	}{
		"mainnet network": {
			network: bitcoin.Mainnet,
			expectedConfigs: []electrum.Config{
				{
					URL:      "electrum.blockstream.info:50002",
					Protocol: electrum.SSL,
				},
			},
		},
		"testnet network": {
			network: bitcoin.Testnet,
			expectedConfigs: []electrum.Config{
				{
					URL:      "electrum.blockstream.info:60002",
					Protocol: electrum.SSL,
				},
			},
		},
		"regtest network": {
			network: bitcoin.Regtest,
			expectedConfigs: []electrum.Config{
				{
					URL:      "",
					Protocol: electrum.Unknown,
				},
			},
		},
		"unknown network": {
			network: bitcoin.Unknown,
			expectedConfigs: []electrum.Config{
				{
					URL:      "",
					Protocol: electrum.Unknown,
				},
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			cfg := &Config{}
			cfg.Bitcoin.Network = test.network

			err := cfg.resolveElectrum()
			if !reflect.DeepEqual(test.expectedError, err) {
				t.Errorf(
					"unexpected error\nexpected: %+v\nactual:   %+v\n",
					test.expectedError,
					err,
				)
			}

			resolvedConfig := cfg.Bitcoin.Electrum
			if !slices.Contains(test.expectedConfigs, resolvedConfig) {
				t.Errorf(
					"expected configs set doesn't contain resolved config\n"+
						"expected: %+v\n"+
						"actual:   %+v\n",
					test.expectedConfigs,
					resolvedConfig,
				)
			}
		})
	}
}
