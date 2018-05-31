package ethereum

import (
	"github.com/keep-network/keep-core/pkg/beacon"
)

func (ec *ethereumChain) GetBeaconConfig() beacon.Config {
	// xyzzy - pull config from chain?
	// TODO: make calls to Ethereum to pull out the configuration
	return beacon.Config{
		GroupSize: 10,
		Threshold: 4,
	}
}
