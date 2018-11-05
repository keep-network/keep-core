package ethereum

/*
Some of this code is MIT License - Look at https://github.com/pschlump/GCall

Copyright (c) 2018 Philip Schlump

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	goeth "github.com/ethereum/go-ethereum"               // ethereum "github.com/ethereum/go-ethereum"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi" //
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/keep-network/keep-core/pkg/chain/gen/abi"
)

// --------------------------------------------------------------------------------------------------------
//
// doWatch watches a contract for events.
//
// Input:
// 		contractName, eventName		-- Watch to watch - if eventName == "" then watch all events on contract
//		gCfg 						-- Config
//			gCfg.GetNameForTopic(log.Topics[0].String())
//
// Uses:
// 		Bind2Contract(...)
// 		ReturnTypeConverter(marshalledValues)
//		TypeOfSlice(marshalledValues)				((debug only))
//
// --------------------------------------------------------------------------------------------------------
func doWatch(contractName, eventName string, assumeProxy bool) (err error) {

	if db10001 {
		fmt.Printf("contractName [%s] eventName [%s]\n", contractName, eventName)
		fmt.Printf("Found contract [before overload check], %s\n", contractName)
	}

	var ABIraw string
	switch contractName {
	case "KeepGroup":
		if assumeProxy {
			ABIraw = abi.KeepGroupImplV1ABI
		} else {
			ABIraw = abi.KeepGroupABI
		}
	case "KeepGroupImplV1":
		ABIraw = abi.KeepGroupImplV1ABI
	default:
		return fmt.Errorf("contrct %s invalid - incorrect contract name", contractName)
	}

	contractAddressStr, ok := TestConfig.ContractAddress[contractName]
	if !ok {
		fmt.Printf("invalid contract address name: [%s] address: [%s].\n", contractName, contractAddressStr)
		return err
	}
	contractAddress := common.HexToAddress(contractAddressStr)

	/* Contract - parse into the go-eth format */
	// conn, err := ethclient.Dial(gCfg.GethURL_ws)
	_, parsedABI, err := Bind2Contract(ABIraw, contractAddress, EthConn.client, EthConn.client, EthConn.client)
	if err != nil {
		fmt.Printf("Error on Bind2Contract: %s\n", err)
		return err
	}

	query := goeth.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	var ch = make(chan types.Log)
	ctx := context.Background()

	sub, err := EthConn.client.SubscribeFilterLogs(ctx, query, ch)
	if err != nil {
		log.Println("Subscribe:", err) // xyzzy  - fix
		return err
	}

	/*
		// list out the current watched events! -- capture current events in list
		if watching, ok := CurrentWatchMap[CurrentWatchType{ContractName: contractName, EventName: eventName}]; !ok || !watching {
			CurrentWatchMap[CurrentWatchType{ContractName: contractName, EventName: eventName}] = true
			CurrentWatch = append(CurrentWatch, CurrentWatchType{ContractName: contractName, EventName: eventName})
		} else {
			fmt.Printf("Already watching %s.%s\n", contractName, eventName)
			return err
		}
	*/

	go func() {
		for {
			if db10001 {
				fmt.Printf("Waiting for event at 'select'n")
			}
			select {
			case log := <-ch:
				if len(log.Topics) > 0 {
					// PJS - xyzzy xyzzy - name := gCfg.GetNameForTopic(log.Topics[0].String())
					name := GetNameForTopic(log.Topics[0].String())
					if db10001 {
						fmt.Printf("name [%s] eventName [%s]\n", name, eventName)
					}
					if eventName == "" || name == eventName {
						fmt.Printf("Caught Event Name:%s Log:%s\n", name, PrintAsJson(log))

						if event, ok := parsedABI.Events[name]; ok {
							arguments := event.Inputs                                 // get the inputs to the event - these will determine the unpack.
							marshalledValues, err := arguments.UnpackValues(log.Data) // marshalledValues is an array of interface{}
							if err != nil {
								fmt.Printf("Error on unmarshalling event data: %s eventName:%s\n", err, name)
							} else {
								// 1. Output of watch "bytes32" data - display better as a hex string
								// 0xBBbbBB... for 32 bytes instead of an array of byte.
								typeModified := ReturnTypeConverter(marshalledValues)
								fmt.Printf("Event Data: %s\n", PrintAsJson(typeModified))

								if db10003 {
									TypeOfSlice(marshalledValues)
								}
							}
						} else {
							fmt.Printf("Error failed to lookup event [%s] in ABI\n", name)
						}
					} else {
						if db10001 {
							fmt.Printf(db10001, "%s.%s - event ignored; not watched\n", contractName, name)
						}
					}
				}
			case err := <-sub.Err():
				fmt.Printf("error=%s\n", err)
				return
			}
		}
	}()

	return nil

}

// Bind2Contract binds a generic wrapper to an already deployed contract.
func Bind2Contract(ABI string, address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, *ethabi.ABI, error) {
	parsed, err := ethabi.JSON(strings.NewReader(ABI))
	if err != nil {
		return nil, nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), &parsed, nil
}

// TypeOfSlice print out slice types.  Used in debuging.
func TypeOfSlice(t interface{}) {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)

		for i := 0; i < s.Len(); i++ {
			fmt.Printf("i=%d: type=%T\n", i, s.Index(i))
		}
	}
}

// ReturnTypeConverter will Convert return type to have correct datay types so that JSON marshal/unmarshal
// will display it correclty.
func ReturnTypeConverter(rt []interface{}) (rv []interface{}) {
	for ii := 0; ii < len(rt); ii++ {
		t := rt[ii]
		tT := fmt.Sprintf("%T", t)
		if tT == "[32]uint8" {
			uu, ok := t.([32]uint8)
			if !ok {
				panic("Should have conveted")
			}
			var ft EthBytes32
			for jj := 0; jj < 32; jj++ {
				ft[jj] = uu[jj]
			}
			rv = append(rv, ft)
		} else {
			rv = append(rv, t)
		}
		/*
			switch reflect.TypeOf(t).Kind() {
			case reflect.Slice:
				s := reflect.ValueOf(t)

				for i := 0; i < s.Len(); i++ {
					fmt.Printf("i=%d: type=%T", i, s.Index(i))
				}
			default:
				rv = append(rv, rt[ii])
			}
		*/
	}
	return
}

//			gCfg.GetNameForTopic(log.Topics[0].String())
func GetNameForTopic(topic string) string {
	return topic
}

// EthBytes32 is setup to meat the interface{} specification for JSON.
type EthBytes32 [32]uint8

const db10001 = true
const db10002 = false
const db10003 = false
