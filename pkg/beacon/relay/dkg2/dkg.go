package dkg2

import (
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/altbn128"
	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg2/result"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/member"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/operator"
)

// ExecuteDKG runs the full distributed key generation lifecycle.
func ExecuteDKG(
	requestID *big.Int,
	seed *big.Int,
	index int, // starts with 0
	operatorPrivateKey *operator.PrivateKey,
	groupSize int,
	threshold int,
	blockCounter chain.BlockCounter,
	relayChain relayChain.Interface,
	channel net.BroadcastChannel,
) (*ThresholdSigner, error) {
	playerIndex := member.Index(index + 1)
	err := playerIndex.Validate()
	if err != nil {
		return nil, fmt.Errorf("[member:%v] %v", playerIndex, err)
	}

	gjkrResult, err := gjkr.Execute(
		playerIndex,
		blockCounter,
		channel,
		threshold,
		seed,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"[member:%v] GJKR execution failed [%v]",
			playerIndex,
			err,
		)
	}

	err = result.SignAndSubmit(
		operatorPrivateKey,
		channel,
		relayChain,
		blockCounter,
		playerIndex,
		requestID,
		convertResult(gjkrResult, groupSize),
		gjkrResult.Disqualified,
		gjkrResult.Inactive,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"[member:%v] DKG signing and submission process failed [%v]",
			playerIndex,
			err,
		)
	}

	// TODO Consider removing this print after Phase 14 is implemented and
	// replace it with print at the end of DKG execution.
	fmt.Printf("[member:%v] DKG Result: %+v\n", playerIndex, gjkrResult)

	return &ThresholdSigner{
		memberID:             playerIndex,
		groupPublicKey:       gjkrResult.GroupPublicKey,
		groupPrivateKeyShare: gjkrResult.GroupPrivateKeyShare,
	}, nil
}

// convertResult transforms GJKR protocol execution result to a chain specific
// DKG result form. It serializes a group public key to bytes and converts
// disqualified and inactive members lists to a boolean list where each entry
// corresponds to a member in the group and true/false value indicates status of
// the member.
func convertResult(gjkrResult *gjkr.Result, groupSize int) *relayChain.DKGResult {
	groupPublicKey := make([]byte, 0)

	// We convert the point G2, to compress the point correctly
	// (ensuring we encode the parity bit).
	if gjkrResult.GroupPublicKey != nil {
		altbn128GroupPublicKey := altbn128.G2Point{G2: gjkrResult.GroupPublicKey}
		groupPublicKey = altbn128GroupPublicKey.Compress()
	}

	// convertToByteSlice converts slice containing members IDs to a slice of
	// group size length where 0x01 entry indicates the member was found on
	// passed members IDs slice. It assumes member IDs for a group starts iterating
	// from 1. E.g. for a group size of 3 with a passed members ID slice {2} the
	// resulting byte slice will be {0x00, 0x01, 0x00}.
	convertToByteSlice := func(memberIDsSlice []member.Index) []byte {
		bytes := make([]byte, groupSize)
		for index := range bytes {
			for _, memberID := range memberIDsSlice {
				if memberID.Equals(index + 1) {
					bytes[index] = 0x01
				}
			}
		}
		return bytes
	}

	return &relayChain.DKGResult{
		GroupPublicKey: groupPublicKey,
		Inactive:       convertToByteSlice(gjkrResult.Inactive),
		Disqualified:   convertToByteSlice(gjkrResult.Disqualified),
	}
}
