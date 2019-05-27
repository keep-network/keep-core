package dkg

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	chainLocal "github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/net/key"
	netLocal "github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/operator"
)

var minimumStake = big.NewInt(20)

func runTest(
	groupSize int,
	threshold int,
	interceptor testutils.NetworkMessageInterceptor,
) (*dkgTestResult, error) {
	privateKey, publicKey, err := operator.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	_, networkPublicKey := key.OperatorKeyToNetworkKey(privateKey, publicKey)

	network := testutils.NewInterceptingNetwork(
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
	network testutils.InterceptingNetwork,
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
	requestID, seed, err := getDKGParameters()
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
				requestID,
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
			chain.GetDKGResult(requestID),
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

func getDKGParameters() (
	requestID *big.Int,
	seed *big.Int,
	err error,
) {
	requestID, err = rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return
	}

	seed, err = rand.Int(rand.Reader, big.NewInt(100000))

	return
}
