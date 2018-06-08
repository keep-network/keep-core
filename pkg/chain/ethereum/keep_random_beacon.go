package ethereum

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gen "github.com/keep-network/keep-core/pkg/chain/gen"
)

// KeepRandomBeacon connection information for interface to the contract
type KeepRandomBeacon struct {
	provider        *ethereumChain
	caller          *gen.KeepRandomBeaconImplV1Caller
	callerOpts      *bind.CallOpts
	transactor      *gen.KeepRandomBeaconImplV1Transactor
	transactorOpts  *bind.TransactOpts
	contract        *gen.KeepRandomBeaconImplV1
	contractAddress common.Address
	name            string
}

// NewKeepRandomBeacon creates the necessary connections and configurations for
// accessing the contract.
func newKeepRandomBeacon(pv *ethereumChain) (*KeepRandomBeacon, error) {
	// Proxy Address
	ContractAddressHex, ok := pv.config.ContractAddresses["KeepRandomBeaconImplV1"]
	if !ok {
		return nil, fmt.Errorf(
			"Error: no address information for 'KeepRandomBeacon' in configuration\n")
	}
	contractAddress := common.HexToAddress(ContractAddressHex)

	beaconTransactor, err := gen.NewKeepRandomBeaconImplV1Transactor(
		contractAddress,
		pv.client)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to instantiate a KeepRelayBeaconTranactor contract: %s",
			err)
	}

	file, err := os.Open(pv.config.Account.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to open keyfile: %v, %s",
			err, pv.config.Account.KeyFile)
	}

	optsTransactor, err := bind.NewTransactor(
		bufio.NewReader(file),
		pv.config.Account.KeyFilePassword)
	if err != nil {
		return nil, fmt.Errorf("Failed to read keyfile: %v, %s", err,
			pv.config.Account.KeyFile)
	}

	beaconCaller, err := gen.NewKeepRandomBeaconImplV1Caller(
		contractAddress,
		pv.client)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to instantiate a KeepRelayBeaconCaller contract: %s", err)
	}

	optsCaller := &bind.CallOpts{
		Pending: false,
		From:    contractAddress,
		Context: nil,
	}

	krbContract, err := gen.NewKeepRandomBeaconImplV1(contractAddress, pv.client)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to instantiate contract object: %v at address: %s",
			err, ContractAddressHex)
	}

	return &KeepRandomBeacon{
		name:            "KeepRandomBeacon", // "KeepRandomBeaconImplV1",
		provider:        pv,
		transactor:      beaconTransactor,
		transactorOpts:  optsTransactor,
		caller:          beaconCaller,
		callerOpts:      optsCaller,
		contract:        krbContract,
		contractAddress: contractAddress,
	}, nil
}

// Initialized calls the contract and returns true if the contract has
// had its Initialize method called.
func (krb *KeepRandomBeacon) Initialized() (bool, error) {
	return krb.caller.Initialized(krb.callerOpts)
}

// HasMinimumStake returns true if the specified address has sufficient
// state to participate.
func (krb *KeepRandomBeacon) HasMinimumStake(
	address common.Address,
) (bool, error) {
	return krb.caller.HasMinimumStake(krb.callerOpts, address)
}

// RequestRelayEntry start the process of generating a signature.
func (krb *KeepRandomBeacon) RequestRelayEntry(
	blockReward *big.Int,
	rawseed []byte,
) (*types.Transaction, error) {
	seed := big.NewInt(0).SetBytes(rawseed)
	return krb.transactor.RequestRelayEntry(krb.transactorOpts, blockReward, seed)
}

// SubmitGroupPublicKey upon completion of a sgiagure make the contract
// call to put it on chain.
func (krb *KeepRandomBeacon) SubmitGroupPublicKey(
	groupPublicKey []byte,
	requestID *big.Int,
) (*types.Transaction, error) {
	gpk := ByteSliceToSliceOf1Byte(groupPublicKey)
	return krb.transactor.SubmitGroupPublicKey(krb.transactorOpts, gpk, requestID)
}

// relayEntryRequested type of fucntion called for RelayEntryRequested event.
type relayEntryRequested func(requestID *big.Int, payment *big.Int,
	blockReward *big.Int, seed *big.Int, blockNumber *big.Int)

// WatchRelayEntryRequested watches for event RelayEntryRequested.
func (krb *KeepRandomBeacon) WatchRelayEntryRequested(
	success relayEntryRequested,
	fail errorCallback,
) error {
	name := "RelayEntryRequested"
	sink := make(chan *gen.KeepRandomBeaconImplV1RelayEntryRequested, 10)
	event, err := krb.contract.WatchRelayEntryRequested(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s events: %s", name, err)
		return err
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				success(rn.RequestID, rn.Payment, rn.BlockReward, rn.Seed, rn.BlockNumber)

			case ee := <-event.Err():
				fail(ee)
			}
		}
	}()
	return nil
}

// relayEntryGenerated type of fucntion called for RelayEntryGenerated event.
type relayEntryGenerated func(requestID *big.Int, RequestResponse *big.Int,
	RequestGroupID *big.Int, PreviousEntry *big.Int, blockNumber *big.Int)

// WatchRelayEntryGenerated watches for event
func (krb *KeepRandomBeacon) WatchRelayEntryGenerated(
	success relayEntryGenerated,
	fail errorCallback,
) error {
	name := "RelayEntryGenerated"
	sink := make(chan *gen.KeepRandomBeaconImplV1RelayEntryGenerated, 10)
	event, err := krb.contract.WatchRelayEntryGenerated(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s event: %s", name, err)
		return err
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				success(rn.RequestID, rn.RequestResponse,
					rn.RequestGroupID, rn.PreviousEntry, rn.BlockNumber)

			case ee := <-event.Err():
				fail(ee)
			}
		}
	}()
	return nil
}

// relayResetEvent type of fucntion called for ResetEvent event.
type relayResetEvent func(LastValidRelayEntry *big.Int,
	LastValidRelayTxHash *big.Int, LastValidRelayBlock *big.Int)

// WatchRelayResetEvent watches for event WatchRelayResetEvent
func (krb *KeepRandomBeacon) WatchRelayResetEvent(
	success relayResetEvent,
	fail errorCallback,
) error {
	name := "RelayResetEvent"
	sink := make(chan *gen.KeepRandomBeaconImplV1RelayResetEvent, 10)
	event, err := krb.contract.WatchRelayResetEvent(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s event: %s", name, err)
		return err
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				success(rn.LastValidRelayEntry, rn.LastValidRelayTxHash,
					rn.LastValidRelayBlock)

			case ee := <-event.Err():
				fail(ee)
			}
		}
	}()
	return nil
}

// submitGroupPublicKeyEvent type of fucntion called for
// SubmitGroupPublicKeyEvent event.
type submitGroupPublicKeyEvent func(GroupPublicKey []byte,
	RequestID *big.Int, ActivationBlockHeight *big.Int)

// WatchSubmitGroupPublicKeyEvent watches for event SubmitGroupPublicKeyEvent
func (krb *KeepRandomBeacon) WatchSubmitGroupPublicKeyEvent(
	success submitGroupPublicKeyEvent,
	fail errorCallback,
) error {
	name := "SubmitGroupPublicKeyEvent"
	sink := make(chan *gen.KeepRandomBeaconImplV1SubmitGroupPublicKeyEvent, 10)
	event, err := krb.contract.WatchSubmitGroupPublicKeyEvent(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s event: %s", name, err)
		return err
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				gpk := SliceOf1ByteToByteSlice(rn.GroupPublicKey)
				success(gpk, rn.RequestID, rn.ActivationBlockHeight)

			case ee := <-event.Err():
				fail(ee)
			}
		}
	}()
	return nil
}
