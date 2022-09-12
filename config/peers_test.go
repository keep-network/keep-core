package config

import (
	"reflect"
	"testing"

	commonEthereum "github.com/keep-network/keep-common/pkg/chain/ethereum"
)

func TestResolvePeers(t *testing.T) {
	var tests = map[string]struct {
		network       commonEthereum.Network
		expectedPeers []string
		expectedError error
	}{
		// TODO: Add mainnet support
		// "mainnet network": {},
		"goerli network": {
			network: commonEthereum.Goerli,
			expectedPeers: []string{
				"/dns4/bootstrap-0.test.keep.network/tcp/3919/ipfs/16Uiu2HAmCcfVpHwfBKNFbQuhvGuFXHVLQ65gB4sJm7HyrcZuLttH",
				"/dns4/bootstrap-1.test.keep.network/tcp/3919/ipfs/16Uiu2HAm3eJtyFKAttzJ85NLMromHuRg4yyum3CREMf6CHBBV6KY",
				"/dns4/bst-a01.test.keep.boar.network/tcp/4001/ipfs/16Uiu2HAmMosdpAuRSw1ahNhqFq8e3Y4d4c5WZkjW1FGQi5WJwWZ7",
			},
		},
		"developer network": {
			network: commonEthereum.Developer,
		},
		"unknown network": {
			network: commonEthereum.Unknown,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			cfg := &Config{}
			cfg.Ethereum.Network = test.network

			err := cfg.resolvePeers()
			if !reflect.DeepEqual(test.expectedError, err) {
				t.Errorf(
					"unexpected error\nexpected: %+v\nactual:   %+v\n",
					test.expectedError,
					err,
				)
			}

			if !reflect.DeepEqual(test.expectedPeers, cfg.LibP2P.Peers) {
				t.Errorf(
					"unexpected peers\nexpected: %v\nactual:   %v\n",
					test.expectedPeers,
					cfg.LibP2P.Peers,
				)
			}
		})
	}
}
