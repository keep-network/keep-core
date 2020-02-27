package ethereum

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
)

func TestResolveContractByName(t *testing.T) {
	config := ethereum.Config{
		ContractAddresses: make(map[string]string),
	}

	barbarianAddress := "0x0b185C37E1C9D01437c800a8B60fA0845742c271"
	conquerorAddress := "0x1895e4A71d0956553cf80f2ccac69642a2f1bFF4"
	destroyerAddress := "0xYOLO" // invalid HEX address

	config.ContractAddresses["ConanTheBarbarian"] = barbarianAddress
	config.ContractAddresses["ConanTheConqueror"] = conquerorAddress
	config.ContractAddresses["ConanTheDestroyer"] = destroyerAddress

	toCommonAddress := func(address string) *common.Address {
		commonAddress := common.HexToAddress(address)
		return &commonAddress
	}

	var tests = map[string]struct {
		queriedContractName     string
		expectedContractAddress *common.Address
		expectedError           error
	}{
		"known contract name #1": {
			queriedContractName:     "ConanTheBarbarian",
			expectedContractAddress: toCommonAddress(barbarianAddress),
			expectedError:           nil,
		},
		"known contract name #2": {
			queriedContractName:     "ConanTheConqueror",
			expectedContractAddress: toCommonAddress(conquerorAddress),
			expectedError:           nil,
		},
		"unknown contract name": {
			queriedContractName:     "ConanTheLegend",
			expectedContractAddress: nil,
			expectedError: fmt.Errorf(
				"no address information for [ConanTheLegend] in configuration",
			),
		},
		"invalid contract address": {
			queriedContractName:     "ConanTheDestroyer",
			expectedContractAddress: nil,
			expectedError: fmt.Errorf(
				"configured address [0xYOLO] for contract [ConanTheDestroyer] is not valid hex address",
			),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			address, err := addressForContract(config, test.queriedContractName)
			if !reflect.DeepEqual(test.expectedContractAddress, address) {
				t.Fatalf(
					"Unexpected contract address\nExpected: [%v]\nActual:   [%v]\n",
					test.expectedContractAddress,
					address,
				)
			}
			if !reflect.DeepEqual(test.expectedError, err) {
				t.Fatalf(
					"Unexpected error\nExpected: [%v]\nActual:   [%v]\n",
					test.expectedError,
					err,
				)
			}
		})
	}
}
