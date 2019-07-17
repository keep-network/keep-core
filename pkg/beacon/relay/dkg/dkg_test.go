package dkg

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/altbn128"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	chainLocal "github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/internal/interceptors"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
	netLocal "github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/operator"
)

var minimumStake = big.NewInt(20)

func TestExecute_HappyPath(t *testing.T) {
	groupSize := 5
	threshold := 3

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		return msg
	}

	result, err := runTest(groupSize, threshold, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	assertSuccessfulSignersCount(t, result, groupSize)
	assertMemberFailuresCount(t, result, 0)
	assertSamePublicKey(t, result)
	assertNoDisqualifiedMembers(t, result)
	assertNoInactiveMembers(t, result)
	assertValidGroupPublicKey(t, result)
}

func TestExecute_IA_member1_commitmentPhase(t *testing.T) {
	groupSize := 5
	threshold := 3

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		// drop commitment message from member 1
		commitmentMsg, ok := msg.(*gjkr.MemberCommitmentsMessage)
		if ok && commitmentMsg.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := runTest(groupSize, threshold, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	assertSuccessfulSignersCount(t, result, groupSize-1)
	assertMemberFailuresCount(t, result, 1)
	assertSamePublicKey(t, result)
	assertNoDisqualifiedMembers(t, result)
	assertInactiveMembers(t, result, group.MemberIndex(1))
	assertValidGroupPublicKey(t, result)
}

func assertSuccessfulSignersCount(
	t *testing.T,
	result *dkgTestResult,
	expectedCount int,
) {
	if len(result.signers) != expectedCount {
		t.Errorf(
			"Unexpected number of successful signers\nExpected: [%v]\nActual:   [%v]",
			expectedCount,
			len(result.signers),
		)
	}
}

func assertMemberFailuresCount(
	t *testing.T,
	result *dkgTestResult,
	expectedCount int,
) {
	if len(result.memberFailures) != expectedCount {
		t.Errorf(
			"Unexpected number of member failures\nExpected: [%v]\nActual:   [%v]",
			expectedCount,
			len(result.memberFailures),
		)
	}
}

func assertSamePublicKey(t *testing.T, result *dkgTestResult) {
	for _, signer := range result.signers {
		testutils.AssertBytesEqual(
			t,
			result.result.GroupPublicKey,
			signer.GroupPublicKeyBytes(),
		)
	}
}

func assertNoDisqualifiedMembers(t *testing.T, result *dkgTestResult) {
	disqualifiedMemberByte := byte(0x01)

	for i, dq := range result.result.Disqualified {
		if dq == disqualifiedMemberByte {
			t.Errorf("Member [%v] has been disqualified", i)
		}
	}
}

func assertNoInactiveMembers(t *testing.T, result *dkgTestResult) {
	assertInactiveMembers(t, result)
}

func assertInactiveMembers(
	t *testing.T,
	result *dkgTestResult,
	expectedInactive ...group.MemberIndex,
) {
	inactiveMemberByte := byte(0x01)
	activeMemberByte := byte(0x00)

	containsIndex := func(
		index group.MemberIndex,
		indexes []group.MemberIndex,
	) bool {
		for _, i := range indexes {
			if i == index {
				return true
			}
		}

		return false
	}

	for i, ia := range result.result.Inactive {
		index := i + 1 // member indexes starts from 1
		inactiveExpected := containsIndex(group.MemberIndex(index), expectedInactive)

		if ia == inactiveMemberByte && !inactiveExpected {
			t.Errorf("Member [%v] has been marked as inactive", index)
		} else if ia == activeMemberByte && inactiveExpected {
			t.Errorf("Member [%v] has not been marked as inactive", index)
		}
	}
}

func assertValidGroupPublicKey(t *testing.T, result *dkgTestResult) {
	_, err := altbn128.DecompressToG2(result.result.GroupPublicKey)
	if err != nil {
		t.Errorf("Invalid group public key: [%v]", err)
	}
}

func runTest(
	groupSize int,
	threshold int,
	interceptor interceptors.NetworkMessageInterceptor,
) (*dkgTestResult, error) {
	privateKey, publicKey, err := operator.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	_, networkPublicKey := key.OperatorKeyToNetworkKey(privateKey, publicKey)

	network := interceptors.NewInterceptingNetwork(
		netLocal.ConnectWithKey(networkPublicKey),
		interceptor,
	)

	chain := chainLocal.ConnectWithKey(groupSize, threshold, minimumStake, privateKey)

	return executeDKG(groupSize, threshold, chain, network)
}

type dkgTestResult struct {
	result         *relaychain.DKGResult
	signers        []*ThresholdSigner
	memberFailures []error
}

func executeDKG(
	groupSize int,
	threshold int,
	chain chainLocal.Chain,
	network interceptors.InterceptingNetwork,
) (*dkgTestResult, error) {
	blockCounter, err := chain.BlockCounter()
	if err != nil {
		return nil, err
	}

	broadcastChannel, err := network.ChannelFor("dkg_test")
	if err != nil {
		return nil, err
	}

	resultChan := make(chan *event.DKGResultSubmission)
	chain.ThresholdRelay().OnDKGResultSubmitted(
		func(event *event.DKGResultSubmission) {
			resultChan <- event
		},
	)

	startBlockHeight := uint64(1)
	seed, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		return nil, err
	}

	var signers []*ThresholdSigner
	var memberFailures []error

	var wg sync.WaitGroup
	wg.Add(groupSize)
	for i := 0; i < groupSize; i++ {
		i := i // capture for goroutine
		go func() {
			signer, err := ExecuteDKG(
				seed,
				i,
				groupSize,
				threshold,
				startBlockHeight,
				blockCounter,
				chain.ThresholdRelay(),
				broadcastChannel,
			)
			if signer != nil {
				signers = append(signers, signer)
			}
			if err != nil {
				fmt.Printf("Failed with: [%v]\n", err)
				memberFailures = append(memberFailures, err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// We give 5 seconds so that OnDKGResultSubmitted async handler
	// is fired. If it's not, than it means no result was published
	// to the chain.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case _ = <-resultChan:
		// result was published to the chain, let's fetch it
		return &dkgTestResult{
			chain.GetLastDKGResult(),
			signers,
			memberFailures,
		}, nil

	case <-ctx.Done():
		// no result published to the chain
		return &dkgTestResult{
			nil,
			signers,
			memberFailures,
		}, nil
	}
}
