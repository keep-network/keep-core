// Package entrytest provides a full roundtrip relay entry signing test engine
// including all the signing phases. It is executed against local chain and
// broadcast channel.
package entrytest

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/internal/interception"
	"github.com/keep-network/keep-core/pkg/net/key"
	"github.com/keep-network/keep-core/pkg/operator"

	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/entry"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"

	chainLocal "github.com/keep-network/keep-core/pkg/chain/local"
	netLocal "github.com/keep-network/keep-core/pkg/net/local"
)

var minimumStake = big.NewInt(20)

// Result of the relay entry signing protocol execution.
type Result struct {
	Entry          *event.Entry
	signerFailures []error
}

// RunTest executes the full relay entry signing roundtrip test for the provided
// group of signers and threshold. Note that the group public key and private
// key shares used by signers had to be created for the same threshold value.
// The provided interception rules are applied in the broadcast channel for
// the time of the protocol execution.
// Previous entry and seed together form a value to be signed, just like in the
// real random beacon.
func RunTest(
	signers []*dkg.ThresholdSigner,
	threshold int,
	rules interception.Rules,
	previousEntry *big.Int,
	seed *big.Int,
) (*Result, error) {
	privateKey, publicKey, err := operator.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	_, networkPublicKey := key.OperatorKeyToNetworkKey(privateKey, publicKey)

	network := interception.NewNetwork(
		netLocal.ConnectWithKey(networkPublicKey),
		rules,
	)

	chain := chainLocal.ConnectWithKey(len(signers), threshold, minimumStake, privateKey)

	return executeSigning(signers, threshold, chain, network, previousEntry, seed)
}

func executeSigning(
	signers []*dkg.ThresholdSigner,
	threshold int,
	chain chainLocal.Chain,
	network interception.Network,
	previousEntry *big.Int,
	seed *big.Int,
) (*Result, error) {
	blockCounter, err := chain.BlockCounter()
	if err != nil {
		return nil, err
	}

	broadcastChannel, err := network.ChannelFor(
		fmt.Sprintf("entry-test-%v-%v", previousEntry, seed),
	)
	if err != nil {
		return nil, err
	}

	entrySubmissionChan := make(chan *event.Entry)
	chain.ThresholdRelay().OnSignatureSubmitted(
		func(event *event.Entry) {
			entrySubmissionChan <- event
		},
	)

	var signerFailuresMutex sync.Mutex
	var signerFailures []error

	var wg sync.WaitGroup
	wg.Add(len(signers))

	currentBlockHeight, err := blockCounter.CurrentBlock()
	if err != nil {
		return nil, err
	}

	// Wait for 3 blocks before starting signing to
	// make sure all signers are ready
	startBlockHeight := currentBlockHeight + 3

	for _, signer := range signers {
		go func(signer *dkg.ThresholdSigner) {
			err := entry.SignAndSubmit(
				blockCounter,
				broadcastChannel,
				chain.ThresholdRelay(),
				previousEntry,
				seed,
				threshold,
				signer,
				startBlockHeight,
			)
			if err != nil {
				fmt.Printf("[signer:%v %v] failed with: [%v]\n", signer.MemberID(), previousEntry, err)
				signerFailuresMutex.Lock()
				signerFailures = append(signerFailures, err)
				signerFailuresMutex.Unlock()
			}
			wg.Done()
		}(signer)
	}
	wg.Wait()

	// We give 5 seconds so that OnSignatureSubmitted async handler
	// is fired. If it's not, it means no entry was published to
	// the chain.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case entry := <-entrySubmissionChan:
		return &Result{
			entry,
			signerFailures,
		}, nil

	case <-ctx.Done():
		// no entry published to the chain
		return &Result{
			nil,
			signerFailures,
		}, nil
	}
}
