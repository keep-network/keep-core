package ethereum

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
)

// RandomBeacon converts from ethereumChain to beacon.ChainInterface.
func (ec *ethereumChain) RandomBeacon() beacon.ChainInterface {
	return ec
}

// ThresholdRelay converts from ethereumChain to beacon.ChainInterface.
func (ec *ethereumChain) ThresholdRelay() relay.ChainInterface {
	return nil
}

// GetConfig get the GroupSize and Threshold for groups and returns it
// in a becaon.Config struct.
func (ec *ethereumChain) GetConfig() (beacon.Config, error) {
	size, err := ec.keepGroupContract.GroupSize()
	if err != nil {
		return beacon.Config{}, fmt.Errorf("error calling GroupSize: %s\n", err)
	}

	threshold, err := ec.keepGroupContract.GroupThreshold()
	if err != nil {
		return beacon.Config{}, fmt.Errorf("error calling GroupThreshold: %s\n", err)
	}

	return beacon.Config{
		GroupSize: size,
		Threshold: threshold,
	}, nil
}
