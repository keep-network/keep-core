package ethereum

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	erpc "github.com/ethereum/go-ethereum/rpc"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

// ThresholdRelay converts from ethereumChain to beacon.ChainInterface.
func (ec *ethereumChain) ThresholdRelay() relaychain.Interface {
	return ec
}

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

func (ec *ethereumChain) SubmitGroupPublicKey(
	groupID string,
	key [96]byte,
) *async.GroupRegistrationPromise {
	groupRegistrationPromise := &async.GroupRegistrationPromise{}

	err := ec.keepRandomBeaconContract.WatchSubmitGroupPublicKeyEvent(
		func(
			groupPublicKey []byte,
			requestID *big.Int,
			activationBlockHeight *big.Int,
		) {
			err := groupRegistrationPromise.Fulfill(&event.GroupRegistration{
				GroupPublicKey:        groupPublicKey,
				RequestID:             requestID,
				ActivationBlockHeight: activationBlockHeight,
			})
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"fulfilling promise failed with: [%v].\n",
					err,
				)
			}
		},
		func(err error) error {
			return groupRegistrationPromise.Fail(
				fmt.Errorf(
					"entry of group key failed with: [%v]",
					err,
				),
			)
		},
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"watch group public key event failed with: [%v].\n",
			err,
		)
		return groupRegistrationPromise
	}

	_, err = ec.keepRandomBeaconContract.SubmitGroupPublicKey(key[:], big.NewInt(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to submit GroupPublicKey [%v].\n", err)
		return groupRegistrationPromise
	}

	return groupRegistrationPromise
}

func (ec *ethereumChain) SubmitRelayEntry(
	newEntry *event.Entry,
) *async.RelayEntryPromise {
	relayEntryPromise := &async.RelayEntryPromise{}

	err := ec.keepRandomBeaconContract.WatchRelayEntryGenerated(
		func(
			requestID *big.Int,
			requestResponse *big.Int,
			requestGroupID *big.Int,
			previousEntry *big.Int,
			blockNumber *big.Int,
		) {
			var value [32]byte
			copy(value[:], requestResponse.Bytes()[:32])

			err := relayEntryPromise.Fulfill(&event.Entry{
				RequestID:     requestID,
				Value:         value,
				GroupID:       requestGroupID,
				PreviousEntry: previousEntry,
				Timestamp:     time.Now().UTC(),
			})
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"execution of fulfilling promise failed with: [%v]",
					err,
				)
			}
		},
		func(err error) error {
			return relayEntryPromise.Fail(
				fmt.Errorf(
					"entry of relay submission failed with: [%v]",
					err,
				),
			)
		},
	)
	if err != nil {
		promiseErr := relayEntryPromise.Fail(
			fmt.Errorf(
				"watch relay entry failed with: [%v]",
				err,
			),
		)
		if promiseErr != nil {
			fmt.Fprintf(
				os.Stderr,
				"execution of failing promise failed with: [%v]",
				promiseErr,
			)
		}
		return relayEntryPromise
	}

	groupSignature := big.NewInt(int64(0)).SetBytes(newEntry.Value[:])
	_, err = ec.keepRandomBeaconContract.SubmitRelayEntry(
		newEntry.RequestID,
		newEntry.GroupID,
		newEntry.PreviousEntry,
		groupSignature,
	)
	if err != nil {
		promiseErr := relayEntryPromise.Fail(
			fmt.Errorf(
				"submitting relay entry to chain failed with: [%v]",
				err,
			),
		)
		if promiseErr != nil {
			fmt.Printf(
				"execution of failing promise failed with: [%v]",
				promiseErr,
			)
		}
		return relayEntryPromise
	}

	return relayEntryPromise
}

func (ec *ethereumChain) OnRelayEntryGenerated(handle func(entry *event.Entry)) {
	err := ec.keepRandomBeaconContract.WatchRelayEntryGenerated(
		func(
			requestID *big.Int,
			requestResponse *big.Int,
			requestGroupID *big.Int,
			previousEntry *big.Int,
			blockNumber *big.Int,
		) {
			var value [32]byte
			copy(value[:], requestResponse.Bytes()[:32])

			handle(&event.Entry{
				RequestID:     requestID,
				Value:         value,
				GroupID:       requestGroupID,
				PreviousEntry: previousEntry,
				Timestamp:     time.Now().UTC(),
			})
		},
		func(err error) error {
			fmt.Fprintf(
				os.Stderr,
				"watch relay entry failed with: [%v]",
				err,
			)
			return err
		},
	)
	if err != nil {
		fmt.Printf(
			"watch relay entry failed with: [%v]",
			err,
		)
	}
}

func (ec *ethereumChain) OnRelayEntryRequested(
	handle func(request *event.Request),
) {
	err := ec.keepRandomBeaconContract.WatchRelayEntryRequested(
		func(
			requestID *big.Int,
			payment *big.Int,
			blockReward *big.Int,
			seed *big.Int,
			blockNumber *big.Int,
		) {
			handle(&event.Request{
				RequestID:   requestID,
				Payment:     payment,
				BlockReward: blockReward,
				Seed:        seed,
			})
		},
		func(err error) error {
			return fmt.Errorf("relay request event failed with %v", err)
		},
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"watch relay request failed with: [%v]",
			err,
		)
	}
}

func (ec *ethereumChain) AddStaker(
	groupMemberID string,
) *async.StakerRegistrationPromise {
	onStakerAddedPromise := &async.StakerRegistrationPromise{}

	if len(groupMemberID) != 32 {
		err := onStakerAddedPromise.Fail(
			fmt.Errorf(
				"groupMemberID wrong length, need 32, got %d",
				len(groupMemberID),
			),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Promise Fail failed [%v].\n", err)
		}
		return onStakerAddedPromise
	}

	err := ec.keepGroupContract.WatchOnStakerAdded(
		func(index int, groupMemberID []byte) {
			err := onStakerAddedPromise.Fulfill(&event.StakerRegistration{
				Index:         index,
				GroupMemberID: string(groupMemberID),
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Promise Fulfill failed [%v].\n", err)
			}
		},
		func(err error) error {
			failErr := onStakerAddedPromise.Fail(
				fmt.Errorf(
					"adding new staker failed with: [%v]",
					err,
				),
			)

			if failErr != nil {
				fmt.Fprintf(
					os.Stderr,
					"Staker added promise failing failed: [%v]\n",
					err,
				)
			}

			return failErr
		},
	)
	if err != nil {
		err = onStakerAddedPromise.Fail(
			fmt.Errorf(
				"on staker added failed with: [%v]",
				err,
			),
		)
		if err != nil {
			fmt.Printf("Staker added promise failing failed: [%v]\n", err)
		}

		return onStakerAddedPromise
	}

	_, err = ec.keepGroupContract.AddStaker(groupMemberID)
	if err != nil {
		err = onStakerAddedPromise.Fail(
			fmt.Errorf(
				"on staker added failed with: [%v]",
				err,
			),
		)
		if err != nil {
			fmt.Printf("Staker added promise failing failed: [%v]\n", err)
		}

		return onStakerAddedPromise
	}

	return onStakerAddedPromise
}

func (ec *ethereumChain) OnStakerAdded(
	handle func(staker *event.StakerRegistration),
) {
	err := ec.keepGroupContract.WatchOnStakerAdded(
		func(index int, groupMemberID []byte) {
			handle(&event.StakerRegistration{
				Index:         index,
				GroupMemberID: string(groupMemberID),
			})
		},
		func(err error) error {
			return fmt.Errorf("staker event failed with %v", err)
		},
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"failed to watch OnStakerAdded [%v].\n",
			err,
		)
	}
}

func (ec *ethereumChain) GetStakerList() ([]string, error) {
	count, err := ec.keepGroupContract.GetNStaker()
	if err != nil {
		err = fmt.Errorf("failed on call to GetNStaker: [%v]", err)
		return nil, err
	}

	listOfStakers := make([]string, 0, count)
	for i := 0; i < count; i++ {
		staker, err := ec.keepGroupContract.GetStaker(i)
		if err != nil {
			return nil,
				fmt.Errorf("at postion %d out of %d error: [%v]", i, count, err)
		}
		listOfStakers = append(listOfStakers, string(staker))
	}

	return listOfStakers, nil
}

func (ec *ethereumChain) OnGroupRegistered(
	handle func(registration *event.GroupRegistration),
) {
	err := ec.keepRandomBeaconContract.WatchSubmitGroupPublicKeyEvent(
		func(
			groupPublicKey []byte,
			requestID *big.Int,
			activationBlockHeight *big.Int,
		) {
			handle(&event.GroupRegistration{
				GroupPublicKey:        groupPublicKey,
				RequestID:             requestID,
				ActivationBlockHeight: activationBlockHeight,
			})
		},
		func(err error) error {
			return fmt.Errorf("entry of group key failed with: [%v]", err)
		},
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"watch group public key event failed with: [%v].\n",
			err,
		)
	}
}

func (ec *ethereumChain) RequestRelayEntry(
	blockReward, seed *big.Int,
) *async.RelayRequestPromise {
	promise := &async.RelayRequestPromise{}
	err := ec.keepRandomBeaconContract.WatchRelayEntryRequested(
		func(
			requestID *big.Int,
			payment *big.Int,
			blockReward *big.Int,
			seed *big.Int,
			blockNumber *big.Int,
		) {
			promise.Fulfill(&event.Request{
				RequestID:   requestID,
				Payment:     payment,
				BlockReward: blockReward,
				Seed:        seed,
			})
		},
		func(err error) error {
			return fmt.Errorf("relay request failed with %v", err)
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "relay request failed with: [%v]", err)
		return promise
	}
	_, err = ec.keepRandomBeaconContract.RequestRelayEntry(blockReward, seed.Bytes())
	if err != nil {
		promise.Fail(fmt.Errorf("failure to request a relay entry: [%v]", err))
		return promise
	}
	return promise
}

func (ec *ethereumChain) ResetStaker() (*types.Transaction, error) {
	return ec.keepGroupContract.ResetStaker()
}

// ConnectTest makes the network connection to the Ethereum network.  Note: for
// other things to work correctly the configuration will need to reference a
// websocket, "ws://", or local IPC connection.
func ConnectTest(cfg Config) (*ethereumChain, error) {
	client, err := ethclient.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf(
			"error Connecting to Geth Server: %s [%v]",
			cfg.URL,
			err,
		)
	}

	clientws, err := erpc.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf(
			"error Connecting to Geth Server: %s [%v]",
			cfg.URL,
			err,
		)
	}

	clientrpc, err := erpc.Dial(cfg.URLRPC)
	if err != nil {
		return nil, fmt.Errorf(
			"error Connecting to Geth Server: %s [%v]",
			cfg.URL,
			err,
		)
	}

	pv := &ethereumChain{
		config:    cfg,
		client:    client,
		clientRPC: clientrpc,
		clientWS:  clientws,
	}

	keepRandomBeaconContract, err := newKeepRandomBeacon(pv)
	if err != nil {
		return nil, fmt.Errorf(
			"error attaching to KeepRandomBeacon contract: [%v]",
			err,
		)
	}
	pv.keepRandomBeaconContract = keepRandomBeaconContract

	keepGroupContract, err := newKeepGroup(pv)
	if err != nil {
		return nil, fmt.Errorf("error attaching to KeepGroup contract: [%v]", err)
	}
	pv.keepGroupContract = keepGroupContract

	return pv, nil
}
