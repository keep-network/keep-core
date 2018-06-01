package ethereum

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
)

// RandomBeacon converts from ethereumChain to beacon.ChainInterface.
// This class supports the necessary interface.
func (ec *ethereumChain) RandomBeacon() beacon.ChainInterface {
	return ec
}

// ThresholdRelay converts from ethereumChain to beacon.ChainInterface.
// This class supports the necessary interface.
func (ec *ethereumChain) ThresholdRelay() relay.ChainInterface {
	return ec
}

// GetConfig get the GroupSize and Threshold for groups and returns it
// in a becaon.Config struct.
func (ec *ethereumChain) GetConfig() beacon.Config {
	size, err := ec.kg.GetGroupSize()
	if err != nil {
	}

	threshold, err := ec.kg.GetGroupThreshold()
	if err != nil {
	}

	rv := beacon.Config{
		GroupSize: size,      // Formerly 10,
		Threshold: threshold, // Formerly 4,
	}
	return rv
}

// SubmitGroupPublicKey sets up the callback functions for the submission of a
// public key for the group.
// I am still not understanding something about this function!!!!!!!!!!!!!!!!
func (ec *ethereumChain) SubmitGroupPublicKey(groupID string, key [96]byte) error {

	applyError := func(msg string) {
		ec.handlerMutex.Lock()
		for _, handler := range ec.groupPublicKeyFailureHandlers {
			handler(groupID, msg)
		}
		ec.handlerMutex.Unlock()
	}

	success := func(GroupPublicKey []byte, RequestID *big.Int, ActivationBlockHeight *big.Int) {
		ec.handlerMutex.Lock()
		for _, handler := range ec.groupPublicKeySubmissionHandlers {
			handler(groupID, ActivationBlockHeight)
		}
		ec.handlerMutex.Unlock()
	}

	fail := func(err error) error {
		applyError(fmt.Sprintf("%s", err))
		return err
	}

	err := ec.krb.WatchSubmitGroupPublicKeyEvent(success, fail)
	if err != nil {
		applyError(fmt.Sprintf("error creating event watch for request:%s", err))
		return err
	}

	requestID := Sum256([]byte(groupID))
	ec.requestID = big.NewInt(0).SetBytes(requestID[:])

	tx, err := ec.krb.SubmitGroupPublicKey(key[:], ec.requestID)
	ec.tx = tx
	if err != nil {
		applyError(fmt.Sprintf("error submitting request:%s", err))
	}
	return err
}

// OnGroupPublicKeySubmissionFailed associates a handler for a error event
func (ec *ethereumChain) OnGroupPublicKeySubmissionFailed(handler func(groupID string, errorMessage string)) error {
	ec.handlerMutex.Lock()
	ec.groupPublicKeyFailureHandlers = append(ec.groupPublicKeyFailureHandlers, handler)
	ec.handlerMutex.Unlock()
	return nil
}

// OnGroupPublicKeySubmitted associates a handler for a success event
func (ec *ethereumChain) OnGroupPublicKeySubmitted(handler func(groupID string, activationBlock *big.Int)) error {
	ec.handlerMutex.Lock()
	ec.groupPublicKeySubmissionHandlers = append(ec.groupPublicKeySubmissionHandlers, handler)
	ec.handlerMutex.Unlock()
	return nil
}

func (ec *ethereumChain) GetLastTx() *types.Transaction {
	return ec.tx
}

func (ec *ethereumChain) GetLastRequestID() *big.Int {
	return ec.requestID
}
