package dkg2

import (
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/altbn128"
	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/states"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

// ExecuteDKG runs the full distributed key generation lifecycle.
func ExecuteDKG(
	requestID *big.Int,
	seed *big.Int,
	index int, // starts with 0
	groupSize int,
	threshold int,
	blockCounter chain.BlockCounter,
	relayChain relayChain.Interface,
	channel net.BroadcastChannel,
) (*ThresholdSigner, error) {
	// The staker index should begin with 1
	playerIndex := index + 1
	if playerIndex < 1 {
		return nil, fmt.Errorf("[member:%v] player index must be >= 1", playerIndex)
	}

	gjkrResult, signer, err := executeGJKR(playerIndex, blockCounter, channel, threshold, seed)
	if err != nil {
		return nil, fmt.Errorf("[member:%v] GJKR execution failed [%v]", playerIndex, err)
	}

	// TODO Consider removing this print after Phase 14 is implemented and
	// replace it with print at the end of DKG execution.
	fmt.Printf("[member:%v] GJKR Result: %+v\n", playerIndex, gjkrResult)

	err = executePublishing(
		requestID,
		playerIndex,
		relayChain,
		blockCounter,
		convertResult(gjkrResult, groupSize),
	)
	if err != nil {
		return nil, fmt.Errorf("publishing failed [%v]", err)
	}

	return signer, nil
}

// executeGJKR runs the GJKR distributed key generation  protocol, given a
// broadcast channel to mediate it, a block counter used for time tracking,
// a player index to use in the group, and a group size and threshold. If
// generation is successful, it returns a threshold group member who can
// participate in the group; if generation fails, it returns an error
// representing what went wrong.
func executeGJKR(
	playerIndex int,
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	threshold int,
	seed *big.Int,
) (*gjkr.Result, *ThresholdSigner, error) {
	memberID := gjkr.MemberID(playerIndex)
	fmt.Printf("[member:0x%010v] Initializing member\n", memberID)

	member, err := gjkr.NewMember(memberID, make([]gjkr.MemberID, 0), threshold, seed)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create a new member [%v]", err)
	}
	initialState := gjkr.InitializationState(channel, member)

	finalState, err := states.Execute(
		channel,
		gjkr.Init(channel),
		initialState,
		blockCounter,
	)
	if err != nil {
		return nil, nil, err
	}

	// TODO: Rename states to be exported
	if finalState, ok := finalState.(*gjkr.FinalizationState); ok {
		gjkrResult := finalState.Result()
		return gjkrResult,
			&ThresholdSigner{
				memberID:             gjkr.MemberID(finalState.MemberID()),
				groupPublicKey:       gjkrResult.GroupPublicKey,
				groupPrivateKeyShare: finalState.GroupPrivateKeyShare(),
			},
			nil
	}

	return nil, nil, fmt.Errorf("invalid state at the end")
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
	convertToByteSlice := func(memberIDsSlice []gjkr.MemberID) []byte {
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
		Success:        gjkrResult.Success,
		GroupPublicKey: groupPublicKey,
		Inactive:       convertToByteSlice(gjkrResult.Inactive),
		Disqualified:   convertToByteSlice(gjkrResult.Disqualified),
	}
}
