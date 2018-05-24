package ethereum

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gen "github.com/keep-network/keep-core/pkg/chain/gen"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

type KeepRandomBeacon struct {
	Provider        *provider
	Caller          *gen.KeepRandomBeaconImplV1Caller
	CallerOpts      *bind.CallOpts
	Transactor      *gen.KeepRandomBeaconImplV1Transactor
	TransactorOpts  *bind.TransactOpts
	Contract        *gen.KeepRandomBeaconImplV1
	ABI             string
	ABIparsed       abi.ABI
	ContractAddress common.Address
	Name            string
}

func NewKeepRandomBeacon(pv *provider) (rv *KeepRandomBeacon, err error) {

	ContractAddressHex := pv.Config.ContractAddresses["KeepRandomBeacon"] // Proxy Address
	ContractAddress := common.HexToAddress(ContractAddressHex)

	krbTransactor, err := gen.NewKeepRandomBeaconImplV1Transactor(ContractAddress, pv.Client)
	if err != nil {
		log.Printf("Failed to instantiate a KeepRelayBeaconTranactor contract: %s", err)
		return
	}

	file, err := os.Open(pv.Config.Account.KeyFile)
	if err != nil {
		log.Printf("Failed to open keyfile: %v, %s", err, pv.Config.Account.KeyFile)
		return
	}

	optsTransactor, err := bind.NewTransactor(bufio.NewReader(file), pv.Config.Account.KeyFilePassword)
	if err != nil {
		log.Printf("Failed to read keyfile: %v, %s", err, pv.Config.Account.KeyFile)
		return
	}

	krbCaller, err := gen.NewKeepRandomBeaconImplV1Caller(ContractAddress, pv.Client)
	if err != nil {
		log.Printf("Failed to instantiate a KeepRelayBeaconCaller contract: %s", err)
		return
	}

	optsCaller := &bind.CallOpts{
		Pending: false,
		From:    ContractAddress,
		Context: nil,
	}

	parsed, err := abi.JSON(strings.NewReader(gen.KeepRandomBeaconImplV1ABI))
	if err != nil {
		log.Printf("Failed to parse ABI, error:%s", err)
		return
	}

	krbContract, err := gen.NewKeepRandomBeaconImplV1(ContractAddress, pv.Client)
	if err != nil {
		log.Printf("Failed to instantiate contract object: %v at address: %s", err, ContractAddressHex)
		return
	}

	return &KeepRandomBeacon{
		Name:            "KeepRandomBeacon", // "KeepRandomBeaconImplV1",
		Provider:        pv,
		Transactor:      krbTransactor,
		TransactorOpts:  optsTransactor,
		Caller:          krbCaller,
		CallerOpts:      optsCaller,
		Contract:        krbContract,
		ABI:             gen.KeepRandomBeaconImplV1ABI,
		ABIparsed:       parsed,
		ContractAddress: ContractAddress,
	}, nil
}

func (krb *KeepRandomBeacon) Initialized() (isInitialize bool, err error) {
	isInitialize, err = krb.Caller.Initialized(krb.CallerOpts)
	return
}

func (krb *KeepRandomBeacon) HasMinimumStake(address common.Address) (hasMinimum bool, err error) {
	hasMinimum, err = krb.Caller.HasMinimumStake(krb.CallerOpts, address)
	return
}

func (krb *KeepRandomBeacon) RequestRelayEntry(blockReward *big.Int, rawseed []byte) (tx *types.Transaction, err error) {
	seed := big.NewInt(0).SetBytes(rawseed)
	tx, err = krb.Transactor.RequestRelayEntry(krb.TransactorOpts, blockReward, seed)
	return
}

func (krb *KeepRandomBeacon) RelayEntry(requestID *big.Int, groupSignature *big.Int, groupID *big.Int, previousEntry *big.Int) (tx *types.Transaction, err error) {
	tx, err = krb.Transactor.RelayEntry(krb.TransactorOpts, requestID, groupSignature, groupID, previousEntry)
	return
}

func (krb *KeepRandomBeacon) SubmitGroupPublicKey(groupPublicKey []byte, requestID *big.Int) (tx *types.Transaction, err error) {
	gpk := ByteSliceToSliceOf1Byte(groupPublicKey)
	tx, err = krb.Transactor.SubmitGroupPublicKey(krb.TransactorOpts, gpk, requestID)
	return
}

type FxRelayEntryRequested func(requestID *big.Int, payment *big.Int, blockReward *big.Int, seed *big.Int, blockNumber *big.Int)

func (krb *KeepRandomBeacon) WatchRelayEntryRequested(success FxRelayEntryRequested, fail FxError) (err error) {
	name := "RelayEntryRequested"
	sink := make(chan *gen.KeepRandomBeaconImplV1RelayEntryRequested, 10)
	if db1 {
		fmt.Printf("Calling Watch for %s, %s\n", name, godebug.LF())
	}
	event, err := krb.Contract.WatchRelayEntryRequested(nil, sink)
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
	event, err := krb.Contract.WatchRelayEntryGenerated(nil, sink)
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
	event, err := krb.Contract.WatchRelayResetEvent(nil, sink)
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
	event, err := krb.Contract.WatchSubmitGroupPublicKeyEvent(nil, sink)
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
