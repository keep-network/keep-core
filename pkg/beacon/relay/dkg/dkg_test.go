package dkg

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/net"

	chainLocal "github.com/keep-network/keep-core/pkg/chain/local"
)

func TestExecute_HappyPath(t *testing.T) {
	groupSize := 5
	threshold := 3

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		return msg
	}
	network := testutils.NewInterceptingNetwork(interceptor)

	result, signers, err := executeDKG(groupSize, threshold, network)
	if err != nil {
		t.Fatal(err)
	}

	assertSignersCount(t, signers, groupSize)
	assertSamePublicKey(t, result, signers)
	// TODO: assert no DQ
	// TODO: assert no IA
	// TODO: assert key is valid
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
	network := testutils.NewInterceptingNetwork(interceptor)

	result, signers, err := executeDKG(groupSize, threshold, network)
	if err != nil {
		t.Fatal(err)
	}

	assertSignersCount(t, signers, groupSize)
	honestSigners := filterOutMisbehavingSigners(signers, group.MemberIndex(1))
	assertSamePublicKey(t, result, honestSigners)
	// TODO: assert no DQ
	// TODO: assert member 1 is IA
	// TODO: assert key is valid
}

func assertSignersCount(
	t *testing.T,
	signers []*ThresholdSigner,
	expectedCount int,
) {
	if len(signers) != expectedCount {
		t.Errorf(
			"Unexpected number of signers\nExpected: [%v]\nActual:   [%v]",
			expectedCount,
			len(signers),
		)
	}
}

func assertSamePublicKey(
	t *testing.T,
	result *relaychain.DKGResult,
	signers []*ThresholdSigner,
) {
	for _, signer := range signers {
		testutils.AssertBytesEqual(
			t,
			result.GroupPublicKey,
			signer.GroupPublicKeyBytes(),
		)
	}
}

func filterOutMisbehavingSigners(
	signers []*ThresholdSigner,
	misbehavingSignersIDs ...group.MemberIndex,
) []*ThresholdSigner {
	var honestSigners []*ThresholdSigner
	for _, signer := range signers {
		isMisbehaving := false
		for _, misbehavingID := range misbehavingSignersIDs {
			if signer.MemberID() == misbehavingID {
				isMisbehaving = true
				break
			}
		}
		if !isMisbehaving {
			honestSigners = append(honestSigners, signer)
		}
	}
	return honestSigners
}

func executeDKG(
	groupSize int,
	threshold int,
	network testutils.InterceptingNetwork,
) (*relaychain.DKGResult, []*ThresholdSigner, error) {
	minimumStake, requestID, seed, startBlockHeight, err := getDKGParameters()
	if err != nil {
		return nil, nil, err
	}

	chainHandle := chainLocal.Connect(groupSize, threshold, minimumStake)
	blockCounter, err := chainHandle.BlockCounter()
	if err != nil {
		return nil, nil, err
	}
	broadcastChannel, err := network.ChannelFor("dkg_test")
	if err != nil {
		return nil, nil, err
	}

	signers := make([]*ThresholdSigner, groupSize)
	resultChan := make(chan *event.DKGResultSubmission)
	chainHandle.ThresholdRelay().OnDKGResultSubmitted(
		func(event *event.DKGResultSubmission) {
			resultChan <- event
		},
	)

	var wg sync.WaitGroup
	wg.Add(groupSize)
	for i := 0; i < groupSize; i++ {
		i := i // capture for goroutine
		go func() {
			signer, _ := ExecuteDKG(
				requestID,
				seed,
				i,
				groupSize,
				threshold,
				startBlockHeight,
				blockCounter,
				chainHandle.ThresholdRelay(),
				broadcastChannel,
			)
			signers[i] = signer
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
		return chainHandle.GetDKGResult(requestID), signers, nil
	case <-ctx.Done():
		return nil, signers, fmt.Errorf("No result published to the chain")
	}
}

func getDKGParameters() (
	minimumStake *big.Int,
	requestID *big.Int,
	seed *big.Int,
	startBlockHeight uint64,
	err error,
) {
	startBlockHeight = uint64(1)
	minimumStake = big.NewInt(20)

	requestID, err = rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return
	}

	seed, err = rand.Int(rand.Reader, big.NewInt(100000))

	return
}
