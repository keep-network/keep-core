package config

import (
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

func TestResolveElectrum(t *testing.T) {
	var tests = map[string]struct {
		network       bitcoin.Network
		expectedURL   string
		expectedError error
	}{
		"mainnet network": {
			network:     bitcoin.Mainnet,
			expectedURL: "ssl://electrum.blockstream.info:50002",
		},
		"testnet network": {
			network:     bitcoin.Testnet,
			expectedURL: "ssl://electrum.blockstream.info:60002",
		},
		"regtest network": {
			network:     bitcoin.Regtest,
			expectedURL: "",
		},
		"unknown network": {
			network:     bitcoin.Unknown,
			expectedURL: "",
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
			if !reflect.DeepEqual(test.expectedURL, resolvedConfig.URL) {
				t.Errorf(
					"expected URL doesn't match resolved URL\n"+
						"expected: %+v\n"+
						"actual:   %+v\n",
					test.expectedURL,
					resolvedConfig.URL,
				)
			}
		})
	}
}
