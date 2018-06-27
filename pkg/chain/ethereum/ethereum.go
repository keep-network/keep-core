package ethereum

import (
	"fmt"
	"math/big"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay"
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

func (ec *ethereumChain) SubmitRelayEntry(entry relay.Entry) *async.RelayEntryPromise {
	relayEntryPromise := &async.RelayEntryPromise{}

	err := ec.keepRandomBeaconContract.WatchRelayEntryGenerated(
		func(
			RequestID *big.Int,
			RequestResponse *big.Int,
			RequestGroupID *big.Int,
			PreviousEntry *big.Int,
			blockNumber *big.Int,
		) {
			err := relayEntryPromise.Fulfill(&relay.Entry{
				RequestID:     RequestID,
				Value:         RequestResponse.Bytes(),
				GroupID:       RequestGroupID,
				PreviousEntry: PreviousEntry,
				Timestamp:     time.Now().UTC(),
			})
			if err != nil {
				fmt.Println(err)
			}
		},
		func(err error) error { return relayEntryPromise.Fail(err) },
	)
	if err != nil {
		promiseErr := relayEntryPromise.Fail(err)
		if promiseErr != nil {
			fmt.Println(promiseErr)
		}
		return nil
	}

	groupSignature := big.NewInt(int64(0)).SetBytes(entry.Value)
	_, err = ec.keepRandomBeaconContract.SubmitRelayEntry(
		entry.RequestID,
		entry.GroupID,
		entry.PreviousEntry,
		groupSignature,
	)
	if err != nil {
		promiseErr := relayEntryPromise.Fail(err)
		if promiseErr != nil {
			fmt.Println(promiseErr)
		}
		return nil
	}

	return relayEntryPromise
}
