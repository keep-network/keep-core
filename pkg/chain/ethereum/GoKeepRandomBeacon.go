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
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

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

func NewKeepRandomBeacon(pv *ethereumChain) (rv *KeepRandomBeacon, err error) {

	ContractAddressHex := pv.config.ContractAddresses["KeepRandomBeacon"] // Proxy Address
	contractAddress := common.HexToAddress(ContractAddressHex)

	krbTransactor, err := gen.NewKeepRandomBeaconImplV1Transactor(contractAddress, pv.client)
	if err != nil {
		log.Printf("Failed to instantiate a KeepRelayBeaconTranactor contract: %s", err)
		return
	}

	file, err := os.Open(pv.config.Account.KeyFile)
	if err != nil {
		log.Printf("Failed to open keyfile: %v, %s", err, pv.config.Account.KeyFile)
		return
	}

	optsTransactor, err := bind.NewTransactor(bufio.NewReader(file), pv.config.Account.KeyFilePassword)
	if err != nil {
		log.Printf("Failed to read keyfile: %v, %s", err, pv.config.Account.KeyFile)
		return
	}

	krbCaller, err := gen.NewKeepRandomBeaconImplV1Caller(contractAddress, pv.client)
	if err != nil {
		log.Printf("Failed to instantiate a KeepRelayBeaconCaller contract: %s", err)
		return
	}

	optsCaller := &bind.CallOpts{
		Pending: false,
		From:    contractAddress,
		Context: nil,
	}

	krbContract, err := gen.NewKeepRandomBeaconImplV1(contractAddress, pv.client)
	if err != nil {
		log.Printf("Failed to instantiate contract object: %v at address: %s", err, ContractAddressHex)
		return
	}

	return &KeepRandomBeacon{
		name:            "KeepRandomBeacon", // "KeepRandomBeaconImplV1",
		provider:        pv,
		transactor:      krbTransactor,
		transactorOpts:  optsTransactor,
		caller:          krbCaller,
		callerOpts:      optsCaller,
		contract:        krbContract,
		contractAddress: contractAddress,
	}, nil
}

func (krb *KeepRandomBeacon) Initialized() (bool, error) {
	return krb.caller.Initialized(krb.callerOpts)
}

func (krb *KeepRandomBeacon) HasMinimumStake(address common.Address) (bool, error) {
	return krb.caller.HasMinimumStake(krb.callerOpts, address)
}

func (krb *KeepRandomBeacon) RequestRelayEntry(blockReward *big.Int, rawseed []byte) (*types.Transaction, error) {
	seed := big.NewInt(0).SetBytes(rawseed)
	return krb.transactor.RequestRelayEntry(krb.transactorOpts, blockReward, seed)
}

func (krb *KeepRandomBeacon) RelayEntry(requestID *big.Int, groupSignature *big.Int, groupID *big.Int, previousEntry *big.Int) (*types.Transaction, error) {
	return krb.transactor.RelayEntry(krb.transactorOpts, requestID, groupSignature, groupID, previousEntry)
}

func (krb *KeepRandomBeacon) SubmitGroupPublicKey(groupPublicKey []byte, requestID *big.Int) (*types.Transaction, error) {
	gpk := ByteSliceToSliceOf1Byte(groupPublicKey)
	return krb.transactor.SubmitGroupPublicKey(krb.transactorOpts, gpk, requestID)
}

type FxRelayEntryRequested func(requestID *big.Int, payment *big.Int, blockReward *big.Int, seed *big.Int, blockNumber *big.Int)

func (krb *KeepRandomBeacon) WatchRelayEntryRequested(success FxRelayEntryRequested, fail FxError) (err error) {
	name := "RelayEntryRequested"
	sink := make(chan *gen.KeepRandomBeaconImplV1RelayEntryRequested, 10)
	if db1 {
		fmt.Printf("Calling Watch for %s, %s\n", name, godebug.LF())
	}
	event, err := krb.contract.WatchRelayEntryRequested(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s events: %s", name, err)
		return
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				if db1 {
					fmt.Printf("%sGot a [%s] event! Yea!%s, %+v\n", MiscLib.ColorGreen, name, MiscLib.ColorReset, rn)
					fmt.Printf("%s        Decoded into JSON data!%s, %s\n", MiscLib.ColorGreen, MiscLib.ColorReset, godebug.SVarI(rn))
				}
				success(rn.RequestID, rn.Payment, rn.BlockReward, rn.Seed, rn.BlockNumber)

			case ee := <-event.Err():
				if db1 {
					fmt.Printf("%sGot an error: %s%s\n", MiscLib.ColorYellow, ee, MiscLib.ColorReset)
				}
				fail(ee)
			}
		}
	}()
	return
}

type FxRelayEntryGenerated func(requestID *big.Int, RequestResponse *big.Int, RequestGroupID *big.Int, PreviousEntry *big.Int, blockNumber *big.Int)

func (krb *KeepRandomBeacon) WatchRelayEntryGenerated(success FxRelayEntryGenerated, fail FxError) (err error) {
	name := "RelayEntryGenerated"
	sink := make(chan *gen.KeepRandomBeaconImplV1RelayEntryGenerated, 10)
	if db1 {
		fmt.Printf("Calling Watch for %s, %s\n", name, godebug.LF())
	}
	event, err := krb.contract.WatchRelayEntryGenerated(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s event: %s", name, err)
		return
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				if db1 {
					fmt.Printf("%sGot a [%s] event! Yea!%s, %+v\n", MiscLib.ColorGreen, name, MiscLib.ColorReset, rn)
					fmt.Printf("%s        Decoded into JSON data!%s, %s\n", MiscLib.ColorGreen, MiscLib.ColorReset, godebug.SVarI(rn))
				}
				success(rn.RequestID, rn.RequestResponse, rn.RequestGroupID, rn.PreviousEntry, rn.BlockNumber)

			case ee := <-event.Err():
				if db1 {
					fmt.Printf("%sGot an error, event %s: %s%s\n", MiscLib.ColorYellow, name, ee, MiscLib.ColorReset)
				}
				fail(ee)
			}
		}
	}()
	return
}

type FxRelayResetEvent func(LastValidRelayEntry *big.Int, LastValidRelayTxHash *big.Int, LastValidRelayBlock *big.Int)

func (krb *KeepRandomBeacon) WatchRelayResetEvent(success FxRelayResetEvent, fail FxError) (err error) {
	name := "RelayResetEvent"
	sink := make(chan *gen.KeepRandomBeaconImplV1RelayResetEvent, 10)
	if db1 {
		fmt.Printf("Calling Watch for %s, %s\n", name, godebug.LF())
	}
	event, err := krb.contract.WatchRelayResetEvent(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s event: %s", name, err)
		return
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				if db1 {
					fmt.Printf("%sGot a [%s] event! Yea!%s, %+v\n", MiscLib.ColorGreen, name, MiscLib.ColorReset, rn)
					fmt.Printf("%s        Decoded into JSON data!%s, %s\n", MiscLib.ColorGreen, MiscLib.ColorReset, godebug.SVarI(rn))
				}
				success(rn.LastValidRelayEntry, rn.LastValidRelayTxHash, rn.LastValidRelayBlock)

			case ee := <-event.Err():
				if db1 {
					fmt.Printf("%sGot an error, event %s: %s%s\n", MiscLib.ColorYellow, name, ee, MiscLib.ColorReset)
				}
				fail(ee)
			}
		}
	}()
	return
}

type FxSubmitGroupPublicKeyEvent func(GroupPublicKey []byte, RequestID *big.Int, ActivationBlockHeight *big.Int)

func (krb *KeepRandomBeacon) WatchSubmitGroupPublicKeyEvent(success FxSubmitGroupPublicKeyEvent, fail FxError) (err error) {
	name := "SubmitGroupPublicKeyEvent"
	sink := make(chan *gen.KeepRandomBeaconImplV1SubmitGroupPublicKeyEvent, 10)
	if db1 {
		fmt.Printf("Calling Watch for %s, %s\n", name, godebug.LF())
	}
	event, err := krb.contract.WatchSubmitGroupPublicKeyEvent(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s event: %s", name, err)
		return
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				if db1 {
					fmt.Printf("%sGot a [%s] event! Yea!%s, %+v\n", MiscLib.ColorGreen, name, MiscLib.ColorReset, rn)
					fmt.Printf("%s        Decoded into JSON data!%s, %s\n", MiscLib.ColorGreen, MiscLib.ColorReset, godebug.SVarI(rn))
				}
				gpk := SliceOf1ByteToByteSlice(rn.GroupPublicKey)
				success(gpk, rn.RequestID, rn.ActivationBlockHeight)

			case ee := <-event.Err():
				if db1 {
					fmt.Printf("%sGot an error, event %s: %s%s\n", MiscLib.ColorYellow, name, ee, MiscLib.ColorReset)
				}
				fail(ee)
			}
		}
	}()
	return
}

/*
GoKeepRandomBeacon.go:20:6: exported type KeepRandomBeacon should have comment or be unexported
GoKeepRandomBeacon.go:33:1: exported function NewKeepRandomBeacon should have comment or be unexported
GoKeepRandomBeacon.go:94:1: exported method KeepRandomBeacon.Initialized should have comment or be unexported
GoKeepRandomBeacon.go:98:1: exported method KeepRandomBeacon.HasMinimumStake should have comment or be unexported
GoKeepRandomBeacon.go:102:1: exported method KeepRandomBeacon.RequestRelayEntry should have comment or be unexported
GoKeepRandomBeacon.go:107:1: exported method KeepRandomBeacon.RelayEntry should have comment or be unexported
GoKeepRandomBeacon.go:111:1: exported method KeepRandomBeacon.SubmitGroupPublicKey should have comment or be unexported
GoKeepRandomBeacon.go:116:6: exported type FxRelayEntryRequested should have comment or be unexported
GoKeepRandomBeacon.go:118:1: exported method KeepRandomBeacon.WatchRelayEntryRequested should have comment or be unexported
GoKeepRandomBeacon.go:150:6: exported type FxRelayEntryGenerated should have comment or be unexported
GoKeepRandomBeacon.go:152:1: exported method KeepRandomBeacon.WatchRelayEntryGenerated should have comment or be unexported
GoKeepRandomBeacon.go:184:6: exported type FxRelayResetEvent should have comment or be unexported
GoKeepRandomBeacon.go:186:1: exported method KeepRandomBeacon.WatchRelayResetEvent should have comment or be unexported
GoKeepRandomBeacon.go:218:6: exported type FxSubmitGroupPublicKeyEvent should have comment or be unexported
GoKeepRandomBeacon.go:220:1: exported method KeepRandomBeacon.WatchSubmitGroupPublicKeyEvent should have comment or be unexported
*/
