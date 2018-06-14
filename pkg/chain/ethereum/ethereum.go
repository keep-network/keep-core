package ethereum

import (
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
	promise "github.com/keep-network/keep-core/pkg/callback"
	"github.com/keep-network/keep-core/pkg/chain/gen"
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
) error {

	err := ec.keepRandomBeaconContract.WatchSubmitGroupPublicKeyEvent(
		promise.NewPromise().OnSuccess(func(data interface{}) {
			data2, ok := data.(gen.KeepRandomBeaconSubmitGroupPublicKeyEvent)
			if !ok {
				panic("flame out crash and burn!")
			}
			ec.handlerMutex.Lock()
			for _, handler := range ec.groupPublicKeySubmissionHandlers {
				handler(groupID, data2.ActivationBlockHeight)
			}
			ec.handlerMutex.Unlock()
		}).OnFailure(func(err error) {
			msg := fmt.Sprintf("error: [%v]", err)
			ec.handlerMutex.Lock()
			for _, handler := range ec.groupPublicKeyFailureHandlers {
				handler(groupID, msg)
			}
			ec.handlerMutex.Unlock()
			return
		}),
	)
	if err != nil {
		msg := fmt.Sprintf("error creating event watch for request: [%v]", err)
		ec.handlerMutex.Lock()
		for _, handler := range ec.groupPublicKeyFailureHandlers {
			handler(groupID, msg)
		}
		ec.handlerMutex.Unlock()
		return err
	}

	_, err = ec.keepRandomBeaconContract.SubmitGroupPublicKey(key[:], big.NewInt(1))
	if err != nil {
		msg := fmt.Sprintf("error submitting request: [%v]", err)
		ec.handlerMutex.Lock()
		for _, handler := range ec.groupPublicKeyFailureHandlers {
			handler(groupID, msg)
		}
		ec.handlerMutex.Unlock()
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
