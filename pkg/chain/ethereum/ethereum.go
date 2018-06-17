package ethereum

import (
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
	callback "github.com/keep-network/keep-core/pkg/callback"
)

// RandomBeacon converts from ethereumChain to beacon.ChainInterface.
func (ec *ethereumChain) RandomBeacon() beacon.ChainInterface {
	return ec
}

// ThresholdRelay converts from ethereumChain to beacon.ChainInterface.
func (ec *ethereumChain) ThresholdRelay() relay.ChainInterface {
	return ec
}

// GetConfig get the GroupSize and Threshold for groups and returns it
// in a becaon.Config struct.
func (ec *ethereumChain) GetConfig() (beacon.Config, error) {
	size, err := ec.keepGroupContract.GroupSize()
	if err != nil {
		return beacon.Config{}, fmt.Errorf("error calling GroupSize: [%v]", err)
	}

	threshold, err := ec.keepGroupContract.GroupThreshold()
	if err != nil {
		return beacon.Config{}, fmt.Errorf("error calling GroupThreshold: [%v]", err)
	}

	return beacon.Config{
		GroupSize: size,
		Threshold: threshold,
	}, nil
}

// SubmitGroupPublicKey sets up the callback functions for the submission of a
// public key for the group.
func (ec *ethereumChain) SubmitGroupPublicKey(
	groupID string,
	key [96]byte,
) *callback.Promise {

	aPromise := &callback.Promise{}
	err := ec.keepRandomBeaconContract.WatchSubmitGroupPublicKeyEvent(aPromise)
	if err != nil {
		aPromise.Fail(err)
	}

	_, err = ec.keepRandomBeaconContract.SubmitGroupPublicKey(key[:], big.NewInt(1))
	if err != nil {
		aPromise.Fail(err)
	}

	return aPromise
}
