package ethereum

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// TestCalculateDKGResultHash validates if calculated DKG result hash matches
// expected one.
//
// Expected hashes have been calculated on-chain with:
// `keccak256(abi.encode(groupPubKey, misbehaved))`
func TestCalculateDKGResultHash(t *testing.T) {
	chain := &Chain{}

	var tests = map[string]struct {
		dkgResult    *relaychain.DKGResult
		expectedHash string
	}{

		"dkg result with no misbehaving members": {
			dkgResult: &relaychain.DKGResult{
				GroupPublicKey: []byte{0x64},
				Misbehaved:     []byte{},
			},
			expectedHash: "f1918e8562236eb17adc8502332f4c9c82bc14e19bfc0aa10ab674ff75b3d2f3",
		},
		"dkg result with misbehaving members": {
			dkgResult: &relaychain.DKGResult{
				GroupPublicKey: []byte{0x64},
				Misbehaved:     []byte{0x03, 0x05},
			},
			expectedHash: "9b84bec611298ebcd371abd418e5716f511d7ff3f086cc574a84afe01afb02ec",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			expectedHash := common.Hex2Bytes(test.expectedHash)

			actualHash, err := chain.CalculateDKGResultHash(test.dkgResult)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(expectedHash, actualHash[:]) {
				t.Errorf(
					"\nexpected: %v\nactual:   %x\n",
					test.expectedHash,
					actualHash,
				)
			}
		})
	}
}

func TestConvertSignaturesToChainFormat(t *testing.T) {
	signature1 := common.LeftPadBytes([]byte("marry"), 65)
	signature2 := common.LeftPadBytes([]byte("had"), 65)
	signature3 := common.LeftPadBytes([]byte("a"), 65)
	signature4 := common.LeftPadBytes([]byte("little"), 65)
	signature5 := common.LeftPadBytes([]byte("lamb"), 65)

	invalidSignature := common.LeftPadBytes([]byte("invalid"), 64)

	var tests = map[string]struct {
		signaturesMap map[chain.GroupMemberIndex][]byte
		expectedError error
	}{
		"one valid signature": {
			signaturesMap: map[uint8][]byte{
				1: signature1,
			},
		},
		"five valid signatures": {
			signaturesMap: map[chain.GroupMemberIndex][]byte{
				3: signature3,
				1: signature1,
				4: signature4,
				5: signature5,
				2: signature2,
			},
		},
		"invalid signature": {
			signaturesMap: map[chain.GroupMemberIndex][]byte{
				1: signature1,
				2: invalidSignature,
			},
			expectedError: fmt.Errorf("invalid signature size for member [2] got [64]-bytes but required [65]-bytes"),
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			indicesSlice, signaturesSlice, err :=
				convertSignaturesToChainFormat(test.signaturesMap)

			if !reflect.DeepEqual(err, test.expectedError) {
				t.Errorf(
					"invalid error\nexpected: %v\nactual:   %v\n",
					test.expectedError,
					err,
				)
			}

			if test.expectedError == nil {
				if len(indicesSlice) != len(test.signaturesMap) {
					t.Errorf(
						"invalid member indices slice length\nexpected: %v\nactual:   %v\n",
						len(test.signaturesMap),
						len(indicesSlice),
					)
				}

				if len(signaturesSlice) != (ethutil.SignatureSize * len(indicesSlice)) {
					t.Errorf(
						"invalid signatures slice size\nexpected: %v\nactual:   %v\n",
						(ethutil.SignatureSize * len(indicesSlice)),
						len(signaturesSlice),
					)
				}
			}

			for i, actualMemberIndex := range indicesSlice {
				memberIndex := chain.GroupMemberIndex(actualMemberIndex.Uint64())

				actualSignature := signaturesSlice[ethutil.SignatureSize*i : ethutil.SignatureSize*(i+1)]
				if !bytes.Equal(
					test.signaturesMap[memberIndex],
					actualSignature,
				) {
					t.Errorf(
						"invalid signatures for member %v\nexpected: %v\nactual:   %v\n",
						actualMemberIndex,
						test.signaturesMap[memberIndex],
						actualSignature,
					)
				}
			}
		})
	}
}

func TestPackTicket(t *testing.T) {
	chain := &Chain{}
	toBigInt := func(number string) *big.Int {
		bigInt, _ := new(big.Int).SetString(number, 10)
		return bigInt
	}

	toBigIntFromAddress := func(address string) *big.Int {
		return new(big.Int).SetBytes(common.HexToAddress(address).Bytes())
	}

	var tests = map[string]struct {
		ticketValue        [8]byte
		stakerValue        *big.Int
		virtualStakerIndex *big.Int
		expectedPacked     string
	}{
		"virtual staker index minimum value": {
			ticketValue:        [8]byte{255, 255, 255, 255, 255, 255, 255, 255},
			stakerValue:        toBigInt("471938313681866282067432403796053736964016932944"),
			virtualStakerIndex: toBigInt("1"),
			expectedPacked:     "ffffffffffffffff52aa72262c904281c49765499f85a774c459885000000001",
		},
		"virtual staker index maximum value": {
			ticketValue:        [8]byte{255, 255, 255, 255, 255, 255, 255, 255},
			stakerValue:        toBigInt("471938313681866282067432403796053736964016932944"),
			virtualStakerIndex: toBigInt("4294967295"),
			expectedPacked:     "ffffffffffffffff52aa72262c904281c49765499f85a774c4598850ffffffff",
		},
		"zero ticket value": {
			ticketValue:        [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
			stakerValue:        toBigInt("640134992772870476466797915370027482254406660188"),
			virtualStakerIndex: toBigInt("12"),
			expectedPacked:     "00000000000000007020a5556ba1ce5f92c81063a13d33512cf1305c0000000c",
		},
		"low ticket value": {
			ticketValue:        [8]byte{0, 0, 0, 0, 0, 0, 255, 255},
			stakerValue:        toBigInt("640134992772870476466797915370027482254406660188"),
			virtualStakerIndex: toBigInt("12"),
			expectedPacked:     "000000000000ffff7020a5556ba1ce5f92c81063a13d33512cf1305c0000000c",
		},
		"staker value is derived from an address without leading zeros": {
			ticketValue:        [8]byte{255, 255, 255, 255, 255, 255, 255, 255},
			stakerValue:        toBigIntFromAddress("0x13b6b8e2cb25f86aa9f3f4eb55ff92684c6efb2d"),
			virtualStakerIndex: toBigInt("1"),
			expectedPacked:     "ffffffffffffffff13b6b8e2cb25f86aa9f3f4eb55ff92684c6efb2d00000001",
		},
		"staker value is derived from an address with one leading zero": {
			ticketValue:        [8]byte{255, 255, 255, 255, 255, 255, 255, 255},
			stakerValue:        toBigIntFromAddress("0x017fe79753873f1e87085ab6972715c6c12015e6"),
			virtualStakerIndex: toBigInt("1"),
			expectedPacked:     "ffffffffffffffff017fe79753873f1e87085ab6972715c6c12015e600000001",
		},
		"staker value is derived from an address with two leading zeros": {
			ticketValue:        [8]byte{255, 255, 255, 255, 255, 255, 255, 255},
			stakerValue:        toBigIntFromAddress("0x00a1b551e309e0bf36388e549d075222a3197e0c"),
			virtualStakerIndex: toBigInt("1"),
			expectedPacked:     "ffffffffffffffff00a1b551e309e0bf36388e549d075222a3197e0c00000001",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			ticket := &relaychain.Ticket{
				Value: test.ticketValue,
				Proof: &relaychain.TicketProof{
					StakerValue:        test.stakerValue,
					VirtualStakerIndex: test.virtualStakerIndex,
				},
			}

			actualTicketBytes := chain.packTicket(ticket)

			expectedTicketBytes, _ := hex.DecodeString(test.expectedPacked)

			if !bytes.Equal(expectedTicketBytes, actualTicketBytes[:]) {
				t.Errorf(
					"\nexpected: %v\nactual:   %x\n",
					test.expectedPacked,
					actualTicketBytes,
				)
			}
		})
	}
}
