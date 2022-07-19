package dkg

import (
	"fmt"
	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"math/big"
	"reflect"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/event"
	"github.com/keep-network/keep-core/pkg/beacon/gjkr"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/group"
)

var (
	playerIndex                 group.MemberIndex
	groupPublicKey              *bn256.G2
	gjkrResult                  *gjkr.Result
	dkgResultChannel            chan *event.DKGResultSubmission
	startPublicationBlockHeight uint64
	beaconChain                 beaconchain.Interface
	blockCounter                chain.BlockCounter
)

func setup() {
	playerIndex = group.MemberIndex(1)
	groupPublicKey = new(bn256.G2).ScalarBaseMult(big.NewInt(10))
	gjkrResult = &gjkr.Result{GroupPublicKey: groupPublicKey}
	dkgResultChannel = make(chan *event.DKGResultSubmission, 1)
	startPublicationBlockHeight = uint64(0)
	localChain := local_v1.Connect(5, 3)
	beaconChain = localChain
	blockCounter, _ = beaconChain.BlockCounter()
}

func TestDecideMemberFate_HappyPath(t *testing.T) {
	setup()

	dkgResultChannel <- &event.DKGResultSubmission{
		GroupPublicKey: groupPublicKey.Marshal(),
		Misbehaved:     []byte{},
	}

	err := decideMemberFate(
		playerIndex,
		gjkrResult,
		dkgResultChannel,
		startPublicationBlockHeight,
		beaconChain,
		blockCounter,
	)

	if err != nil {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			nil,
			err,
		)
	}
}

func TestDecideMemberFate_NotSameGroupPublicKey(t *testing.T) {
	setup()

	otherGroupPublicKey := new(bn256.G2).ScalarBaseMult(big.NewInt(11))
	dkgResultChannel <- &event.DKGResultSubmission{
		GroupPublicKey: otherGroupPublicKey.Marshal(),
		Misbehaved:     []byte{},
	}

	err := decideMemberFate(
		playerIndex,
		gjkrResult,
		dkgResultChannel,
		startPublicationBlockHeight,
		beaconChain,
		blockCounter,
	)

	expectedError := fmt.Errorf(
		"[member:%v] could not stay in the group because "+
			"member do not support the same group public key",
		playerIndex,
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestDecideMemberFate_MemberIsMisbehaved(t *testing.T) {
	setup()

	dkgResultChannel <- &event.DKGResultSubmission{
		GroupPublicKey: groupPublicKey.Marshal(),
		Misbehaved:     []byte{playerIndex},
	}

	err := decideMemberFate(
		playerIndex,
		gjkrResult,
		dkgResultChannel,
		startPublicationBlockHeight,
		beaconChain,
		blockCounter,
	)

	expectedError := fmt.Errorf(
		"[member:%v] could not stay in the group because "+
			"member is considered as misbehaving",
		playerIndex,
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestDecideMemberFate_Timeout(t *testing.T) {
	setup()

	err := decideMemberFate(
		playerIndex,
		gjkrResult,
		dkgResultChannel,
		startPublicationBlockHeight,
		beaconChain,
		blockCounter,
	)

	expectedError := fmt.Errorf("DKG result publication timed out")
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}
