package ethereum

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
)

func (ec *ethereumChain) RandomBeacon() beacon.ChainInterface {
	return ec
}

func (ec *ethereumChain) ThresholdRelay() relay.ChainInterface {
	return ec
}

func (ec *ethereumChain) GetConfig() beacon.Config {
	// TODO - how to pull config from ETH
	return ec.GetBeaconConfig()
}

func (ec *ethereumChain) SubmitGroupPublicKey(groupID string, key [96]byte) error {
	// TODO - translate groupID into a request id
	// TODO - how will requestID get set - on request?

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

	tx, err := ec.krb.SubmitGroupPublicKey(key[:], ec.requestID)
	ec.tx = tx
	if err != nil {
		applyError(fmt.Sprintf("error submitting request:%s", err))
	}
	return err
}

func (ec *ethereumChain) OnGroupPublicKeySubmissionFailed(handler func(groupID string, errorMessage string)) error {
	ec.handlerMutex.Lock()
	ec.groupPublicKeyFailureHandlers = append(ec.groupPublicKeyFailureHandlers, handler)
	ec.handlerMutex.Unlock()
	return nil
}

func (ec *ethereumChain) OnGroupPublicKeySubmitted(handler func(groupID string, activationBlock *big.Int)) error {
	ec.handlerMutex.Lock()
	ec.groupPublicKeySubmissionHandlers = append(ec.groupPublicKeySubmissionHandlers, handler)
	ec.handlerMutex.Unlock()
	return nil
}

func (ec *ethereumChain) GetLastTx() *types.Transaction {
	return ec.tx
}
