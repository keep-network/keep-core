package ethereum

import (
	"fmt"
	"math/big"

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
) error {
	applyError := func(msg string) {
		ec.handlerMutex.Lock()
		for _, handler := range ec.groupPublicKeyFailureHandlers {
			handler(groupID, msg)
		}
		ec.handlerMutex.Unlock()
	}

	success := func(
		GroupPublicKey []byte,
		RequestID *big.Int,
		ActivationBlockHeight *big.Int,
	) {
		ec.handlerMutex.Lock()
		for _, handler := range ec.groupPublicKeySubmissionHandlers {
			handler(groupID, ActivationBlockHeight)
		}
		ec.handlerMutex.Unlock()
	}

	fail := func(err error) error {
		applyError(fmt.Sprintf("error: [%v]", err))
		return err
	}

	err := ec.keepRandomBeaconContract.WatchSubmitGroupPublicKeyEvent(
		success,
		fail,
	)
	if err != nil {
		applyError(fmt.Sprintf("error creating event watch for request: [%v]", err))
		return err
	}

	_, err = ec.keepRandomBeaconContract.SubmitGroupPublicKey(key[:], big.NewInt(1))
	if err != nil {
		applyError(fmt.Sprintf("error submitting request: [%v]", err))
	}
	return err
}

// OnGroupPublicKeySubmissionFailed associates a handler for a error event.
func (ec *ethereumChain) OnGroupPublicKeySubmissionFailed(
	handler func(groupID string, errorMessage string),
) error {
	ec.handlerMutex.Lock()
	ec.groupPublicKeyFailureHandlers = append(
		ec.groupPublicKeyFailureHandlers,
		handler,
	)
	ec.handlerMutex.Unlock()
	return nil
}

// OnGroupPublicKeySubmitted associates a handler for a success event.
func (ec *ethereumChain) OnGroupPublicKeySubmitted(
	handler func(groupID string, activationBlock *big.Int),
) error {
	ec.handlerMutex.Lock()
	ec.groupPublicKeySubmissionHandlers = append(
		ec.groupPublicKeySubmissionHandlers,
		handler,
	)
	ec.handlerMutex.Unlock()
	return nil
}
