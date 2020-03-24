package dkg

import (
	"fmt"
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local"
)

var (
	playerIndex                 group.MemberIndex
	groupPublicKey              *bn256.G2
	gjkrResult                  *gjkr.Result
	dkgResultChannel            chan *event.DKGResultSubmission
	startPublicationBlockHeight uint64
	localChain                  chain.Handle
	blockCounter                chain.BlockCounter
)

func setup() {
	playerIndex = group.MemberIndex(1)
	groupPublicKey = new(bn256.G2).ScalarBaseMult(big.NewInt(10))
	gjkrResult = &gjkr.Result{GroupPublicKey: groupPublicKey}
	dkgResultChannel = make(chan *event.DKGResultSubmission)
	startPublicationBlockHeight = uint64(0)
	localChain = local.Connect(5, 3, big.NewInt(10))
	blockCounter, _ = localChain.BlockCounter()
}

func TestDecideMemberFate_HappyPath(t *testing.T) {
	setup()

	go func() {
		// DKG result must arrive before the timeout block.
		// Timeout block equals to 21 because:
		// - startPublicationBlockHeight = 0
		// - signingStateBlocks = 6
		// - groupSize = 5
		// - resultPublicationBlockStep = 3
		// So, 0 + 6 + (5 * 3) = 21
		_ = blockCounter.WaitForBlockHeight(2)
		dkgResultChannel <- &event.DKGResultSubmission{
			GroupPublicKey: groupPublicKey.Marshal(),
			Misbehaved:     []byte{},
		}
	}()

	err := decideMemberFate(
		playerIndex,
		gjkrResult,
		dkgResultChannel,
		startPublicationBlockHeight,
		localChain.ThresholdRelay(),
		blockCounter,
	)

	if err != nil {
		t.Errorf("Error [%v] should be nil", err)
	}
}

func TestDecideMemberFate_NotSameGroupPublicKey(t *testing.T) {
	setup()

	go func() {
		_ = blockCounter.WaitForBlockHeight(2)
		otherGroupPublicKey := new(bn256.G2).ScalarBaseMult(big.NewInt(11))
		dkgResultChannel <- &event.DKGResultSubmission{
			GroupPublicKey: otherGroupPublicKey.Marshal(),
			Misbehaved:     []byte{},
		}
	}()

	err := decideMemberFate(
		playerIndex,
		gjkrResult,
		dkgResultChannel,
		startPublicationBlockHeight,
		localChain.ThresholdRelay(),
		blockCounter,
	)

	if err == nil {
		t.Errorf("Error should not be nil")
	}

	expectedError := fmt.Sprintf(
		"[member:%v] could not stay in the group because "+
			"member do not support the same group public key",
		playerIndex,
	)
	if err != nil && err.Error() != expectedError {
		t.Errorf("Error [%v] should be [%v]", err, expectedError)
	}
}

func TestDecideMemberFate_MemberIsMisbehaved(t *testing.T) {
	setup()

	go func() {
		_ = blockCounter.WaitForBlockHeight(2)
		dkgResultChannel <- &event.DKGResultSubmission{
			GroupPublicKey: groupPublicKey.Marshal(),
			Misbehaved:     []byte{playerIndex},
		}
	}()

	err := decideMemberFate(
		playerIndex,
		gjkrResult,
		dkgResultChannel,
		startPublicationBlockHeight,
		localChain.ThresholdRelay(),
		blockCounter,
	)

	if err == nil {
		t.Errorf("Error should not be nil")
	}

	expectedError := fmt.Sprintf(
		"[member:%v] could not stay in the group because "+
			"member is considered as misbehaved",
		playerIndex,
	)
	if err != nil && err.Error() != expectedError {
		t.Errorf("Error [%v] should be [%v]", err, expectedError)
	}
}

func TestDecideMemberFate_Timeout(t *testing.T) {
	setup()

	err := decideMemberFate(
		playerIndex,
		gjkrResult,
		dkgResultChannel,
		startPublicationBlockHeight,
		localChain.ThresholdRelay(),
		blockCounter,
	)

	if err == nil {
		t.Errorf("Error should not be nil")
	}

	expectedError := "DKG result publication timed out"
	if err != nil && err.Error() != expectedError {
		t.Errorf("Error [%v] should be [%v]", err, expectedError)
	}
}
