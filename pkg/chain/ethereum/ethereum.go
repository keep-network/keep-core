package ethereum

import (
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/async"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
)

// ThresholdRelay converts from ethereumChain to beacon.ChainInterface.
func (ec *ethereumChain) ThresholdRelay() relay.ChainInterface {
	return ec
}

// GetConfig get the GroupSize and Threshold for groups and returns it
// in a becaon.Config struct.
func (ec *ethereumChain) GetConfig() (relay.Config, error) {
	size, err := ec.keepGroupContract.GroupSize()
	if err != nil {
		return relay.Config{}, fmt.Errorf("error calling GroupSize: [%v]", err)
	}

	threshold, err := ec.keepGroupContract.GroupThreshold()
	if err != nil {
		return relay.Config{}, fmt.Errorf("error calling GroupThreshold: [%v]", err)
	}

	return relay.Config{
		GroupSize: size,
		Threshold: threshold,
	}, nil
}

// SubmitGroupPublicKey sets up the callback functions for the submission of a
// public key for the group.
func (ec *ethereumChain) SubmitGroupPublicKey(
	groupID string,
	key [96]byte,
) *async.KeepRandomBeaconSubmitGroupPublicKeyEventPromise {

	aPromise, err := ec.keepRandomBeaconContract.WatchSubmitGroupPublicKeyEvent()
	if err != nil {
		aPromise.Fail(err)
	}

	_, err = ec.keepRandomBeaconContract.SubmitGroupPublicKey(key[:], big.NewInt(1))
	if err != nil {
		aPromise.Fail(err)
	}

	return aPromise
}
