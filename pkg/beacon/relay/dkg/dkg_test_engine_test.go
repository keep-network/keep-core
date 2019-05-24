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
)

type dkgTestResult struct {
	result         *relaychain.DKGResult
	signers        []*ThresholdSigner
	memberFailures []error
}

func executeDKG(
	groupSize int,
	threshold int,
	network testutils.InterceptingNetwork,
) (*dkgTestResult, error) {
	minimumStake, requestID, seed, startBlockHeight, err := getDKGParameters()
	if err != nil {
		return nil, err
	}

	chainHandle := chainLocal.Connect(groupSize, threshold, minimumStake)
	blockCounter, err := chainHandle.BlockCounter()
	if err != nil {
		return nil, err
	}
	broadcastChannel, err := network.ChannelFor("dkg_test")
	if err != nil {
		return nil, err
	}

	resultChan := make(chan *event.DKGResultSubmission)
	chainHandle.ThresholdRelay().OnDKGResultSubmitted(
		func(event *event.DKGResultSubmission) {
			resultChan <- event
		},
	)

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
				chainHandle.ThresholdRelay(),
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
			chainHandle.GetDKGResult(requestID),
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
