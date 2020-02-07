package ethereum

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// TestCalculateDKGResultHash validates if calculated DKG result hash matches
// expected one.
//
// Expected hashes have been calculated on-chain with:
// `keccak256(abi.encode(groupPubKey, misbehaved))`
func TestCalculateDKGResultHash(t *testing.T) {
	chain := &ethereumChain{}

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

				if len(signaturesSlice) != (SignatureSize * len(indicesSlice)) {
					t.Errorf(
						"invalid signatures slice size\nexpected: %v\nactual:   %v\n",
						(SignatureSize * len(indicesSlice)),
						len(signaturesSlice),
					)
				}
			}

			for i, actualMemberIndex := range indicesSlice {
				memberIndex := chain.GroupMemberIndex(actualMemberIndex.Uint64())

				actualSignature := signaturesSlice[SignatureSize*i : SignatureSize*(i+1)]
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
	chain := &ethereumChain{}
	toBigInt := func(number string) *big.Int {
		bigInt, _ := new(big.Int).SetString(number, 10)
		return bigInt
	}

	ticketValue := toBigInt("77475267169740498967948014258679832639111923451618263020575217281118610489031")
	stakerValue := toBigInt("471938313681866282067432403796053736964016932944")

	var tests = map[string]struct {
		ticketValue        *big.Int
		stakerValue        *big.Int
		virtualStakerIndex *big.Int
		expectedPacked     string
	}{
		"virtual staker index minimum value": {
			ticketValue:        ticketValue,
			stakerValue:        stakerValue,
			virtualStakerIndex: toBigInt("1"),
			expectedPacked:     "ab49727f1f1c661a52aa72262c904281c49765499f85a774c459885000000001",
		},
		"virtual staker index maximum value": {
			ticketValue:        ticketValue,
			stakerValue:        stakerValue,
			virtualStakerIndex: toBigInt("4294967295"),
			expectedPacked:     "ab49727f1f1c661a52aa72262c904281c49765499f85a774c4598850ffffffff",
		},
		"zero ticket value": {
			ticketValue:        toBigInt("0"),
			stakerValue:        toBigInt("640134992772870476466797915370027482254406660188"),
			virtualStakerIndex: toBigInt("12"),
			expectedPacked:     "00000000000000007020a5556ba1ce5f92c81063a13d33512cf1305c0000000c",
		},
		"low ticket value (below natural threshold)": {
			ticketValue:        toBigInt("442342886742415014920381897080165736613327114059325198266614648165032201400"),
			stakerValue:        toBigInt("640134992772870476466797915370027482254406660188"),
			virtualStakerIndex: toBigInt("12"),
			expectedPacked:     "00fa5b718feae4ee7020a5556ba1ce5f92c81063a13d33512cf1305c0000000c",
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
