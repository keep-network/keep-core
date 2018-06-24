package ethereum

import (
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/chaintype"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

// ThresholdRelay converts from ethereumChain to beacon.ChainInterface.
func (ec *ethereumChain) ThresholdRelay() relaychain.Interface {
	return ec
}

// GetConfig get the GroupSize and Threshold for groups and returns it
// in a becaon.Config struct.
func (ec *ethereumChain) GetConfig() (relayconfig.Chain, error) {
	size, err := ec.keepGroupContract.GroupSize()
	if err != nil {
		return relayconfig.Chain{}, fmt.Errorf("error calling GroupSize: [%v]", err)
	}

	threshold, err := ec.keepGroupContract.GroupThreshold()
	if err != nil {
		return relayconfig.Chain{}, fmt.Errorf("error calling GroupThreshold: [%v]", err)
	}

	return relayconfig.Chain{
		GroupSize: size,
		Threshold: threshold,
	}, nil
}

// SubmitGroupPublicKey sets up the callback functions for the submission of a
// public key for the group.
func (ec *ethereumChain) SubmitGroupPublicKey(
	groupID string,
	key [96]byte,
) *async.GroupPublicKeyPromise {

	aPromise := &async.GroupPublicKeyPromise{}

	success := func(
		GroupPublicKey []byte,
		RequestID *big.Int,
		ActivationBlockHeight *big.Int,
	) {
		aPromise.Fulfill(&chaintype.GroupPublicKey{
			GroupPublicKey:        GroupPublicKey,
			RequestID:             RequestID,
			ActivationBlockHeight: ActivationBlockHeight,
		})
	}

	fail := func(err error) error {
		aPromise.Fail(err)
		return nil
	}

	err := ec.keepRandomBeaconContract.WatchSubmitGroupPublicKeyEvent(
		success,
		fail,
	)
	if err != nil {
		aPromise.Fail(fmt.Errorf("error creating event watch for request: [%v]", err))
		return nil
	}

	_, err = ec.keepRandomBeaconContract.SubmitGroupPublicKey(key[:], big.NewInt(1))
	if err != nil {
		aPromise.Fail(err)
		return nil
	}

	return aPromise
}
