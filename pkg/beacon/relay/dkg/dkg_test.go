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
	"github.com/keep-network/keep-core/pkg/internal/interception"
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

func TestExecute_IA_member1_ephemeralKeyGenerationPhase(t *testing.T) {
	groupSize := 5
	threshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {

		publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := runTest(groupSize, threshold, interceptorRules)
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

func TestExecute_IA_member1and2_commitmentPhase(t *testing.T) {
	groupSize := 7
	threshold := 4

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		// drop commitment message from member 1
		commitmentMessage, ok := msg.(*gjkr.MemberCommitmentsMessage)
		if ok && commitmentMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		// drop shares message from member 2
		sharesMessage, ok := msg.(*gjkr.PeerSharesMessage)
		if ok && sharesMessage.SenderID() == group.MemberIndex(2) {
			return nil
		}

		return msg
	}

	result, err := runTest(groupSize, threshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	assertSuccessfulSignersCount(t, result, groupSize-2)
	assertMemberFailuresCount(t, result, 2)
	assertSamePublicKey(t, result)
	assertNoDisqualifiedMembers(t, result)
	assertInactiveMembers(t, result, group.MemberIndex(1), group.MemberIndex(2))
	assertValidGroupPublicKey(t, result)
}

func assertSuccessfulSignersCount(
	t *testing.T,
	testResult *dkgTestResult,
	expectedCount int,
) {
	if len(testResult.signers) != expectedCount {
		t.Errorf(
			"unexpected number of successful signers\nexpected: [%v]\nactual:   [%v]",
			expectedCount,
			len(testResult.signers),
		)
	}
}

func assertMemberFailuresCount(
	t *testing.T,
	testResult *dkgTestResult,
	expectedCount int,
) {
	if len(testResult.memberFailures) != expectedCount {
		t.Errorf(
			"unexpected number of member failures\nexpected: [%v]\nactual:   [%v]",
			expectedCount,
			len(testResult.memberFailures),
		)
	}
}

func assertNoDisqualifiedMembers(t *testing.T, testResult *dkgTestResult) {
	disqualifiedMemberByte := byte(0x01)

	for i, dq := range testResult.dkgResult.Disqualified {
		if dq == disqualifiedMemberByte {
			t.Errorf("member [%v] has been unexpectedly disqualified", i)
		}
	}
}

func assertNoInactiveMembers(t *testing.T, testResult *dkgTestResult) {
	assertInactiveMembers(t, testResult)
}

func assertInactiveMembers(
	t *testing.T,
	testResult *dkgTestResult,
	expectedInactiveMembers ...group.MemberIndex,
) {
	inactiveMemberByte := byte(0x01)
	activeMemberByte := byte(0x00)

	containsMemberIndex := func(
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

	for i, ia := range testResult.dkgResult.Inactive {
		memberIndex := i + 1 // member indexes starts from 1
		inactiveExpected := containsMemberIndex(
			group.MemberIndex(memberIndex),
			expectedInactiveMembers,
		)

		if ia == inactiveMemberByte && !inactiveExpected {
			t.Errorf(
				"member [%v] has been unexpectedly marked as inactive",
				memberIndex,
			)
		} else if ia == activeMemberByte && inactiveExpected {
			t.Errorf(
				"member [%v] has not been unexpectedly marked as inactive",
				memberIndex,
			)
		}
	}
}

func assertSamePublicKey(t *testing.T, testResult *dkgTestResult) {
	for _, signer := range testResult.signers {
		testutils.AssertBytesEqual(
			t,
			testResult.dkgResult.GroupPublicKey,
			signer.GroupPublicKeyBytes(),
		)
	}
}

func assertValidGroupPublicKey(t *testing.T, testResult *dkgTestResult) {
	_, err := altbn128.DecompressToG2(testResult.dkgResult.GroupPublicKey)
	if err != nil {
		t.Errorf("invalid group public key: [%v]", err)
	}
}

type dkgTestResult struct {
	dkgResult      *relaychain.DKGResult
	signers        []*ThresholdSigner
	memberFailures []error
}

func runTest(
	groupSize int,
	threshold int,
	rules interception.Rules,
) (*dkgTestResult, error) {
	privateKey, publicKey, err := operator.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	_, networkPublicKey := key.OperatorKeyToNetworkKey(privateKey, publicKey)

	network := interception.NewNetwork(
		netLocal.ConnectWithKey(networkPublicKey),
		rules,
	)

	chain := chainLocal.ConnectWithKey(groupSize, threshold, minimumStake, privateKey)

	return executeDKG(groupSize, threshold, chain, network)
}

func executeDKG(
	groupSize int,
	threshold int,
	chain chainLocal.Chain,
	network interception.Network,
) (*dkgTestResult, error) {
	blockCounter, err := chain.BlockCounter()
	if err != nil {
		return nil, err
	}

	seed, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		return nil, err
	}

	broadcastChannel, err := network.ChannelFor(fmt.Sprintf("dkg-test-%v", seed))
	if err != nil {
		return nil, err
	}

	resultSubmissionChan := make(chan *event.DKGResultSubmission)
	chain.ThresholdRelay().OnDKGResultSubmitted(
		func(event *event.DKGResultSubmission) {
			resultSubmissionChan <- event
		},
	)

	var signersMutex sync.Mutex
	var signers []*ThresholdSigner

	var memberFailures []error

	var wg sync.WaitGroup
	wg.Add(groupSize)

	currentBlockHeight, err := blockCounter.CurrentBlock()
	if err != nil {
		return nil, err
	}

	// Wait for 3 blocks before starting DKG to
	// make sure all members are up.
	startBlockHeight := currentBlockHeight + 3

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
				chain.Signing(),
				broadcastChannel,
			)
			if signer != nil {
				signersMutex.Lock()
				signers = append(signers, signer)
				signersMutex.Unlock()
			}
			if err != nil {
				fmt.Printf("failed with: [%v]\n", err)
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
	case <-resultSubmissionChan:
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
