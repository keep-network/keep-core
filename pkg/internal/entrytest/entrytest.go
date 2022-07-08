// Package entrytest provides a full roundtrip relay entry signing test engine
// including all the signing phases. It is executed against local chain and
// broadcast channel.
package entrytest

import (
	"context"
	"crypto/rand"
	"fmt"
	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
	"github.com/keep-network/keep-core/pkg/chain"
	"math/big"
	"sync"
	"time"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/internal/interception"
	"github.com/keep-network/keep-core/pkg/operator"

	"github.com/keep-network/keep-core/pkg/beacon/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/entry"
	"github.com/keep-network/keep-core/pkg/beacon/event"

	chainLocal "github.com/keep-network/keep-core/pkg/chain/local"
	netLocal "github.com/keep-network/keep-core/pkg/net/local"
)

var minimumStake = big.NewInt(20)

// Result of the relay entry signing protocol execution.
type Result struct {
	entry          []byte
	signerFailures []error
}

// EntryValue returns the value of relay entry from the result as G1 or
// nil if no entry was produced because of signers failures.
// Error is returned if the entry produced by signers can not be unmarshalled
// to G1 because it is corrupted.
func (r *Result) EntryValue() (*bn256.G1, error) {
	if r.entry == nil {
		return nil, nil
	}

	g1 := new(bn256.G1)
	_, err := g1.Unmarshal(r.entry)
	if err != nil {
		return nil, fmt.Errorf("corrupted entry: [%v]", err)
	}

	return g1, nil
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
	previousEntry []byte,
) (*Result, error) {
	operatorPrivateKey, operatorPublicKey, err := operator.GenerateKeyPair(chainLocal.DefaultCurve)
	if err != nil {
		return nil, err
	}

	network := interception.NewNetwork(
		netLocal.ConnectWithKey(operatorPublicKey),
		rules,
	)

	chain := chainLocal.ConnectWithKey(len(signers), threshold, minimumStake, operatorPrivateKey)

	blockCounter, err := chain.BlockCounter()
	if err != nil {
		return nil, err
	}

	return executeSigning(
		signers,
		threshold,
		chain.ThresholdRelay(),
		blockCounter,
		chain.GetLastRelayEntry,
		network,
		previousEntry,
	)
}

func executeSigning(
	signers []*dkg.ThresholdSigner,
	threshold int,
	relayChain beaconchain.Interface,
	blockCounter chain.BlockCounter,
	lastRelayEntryGetter func() []byte,
	network interception.Network,
	previousEntry []byte,
) (*Result, error) {
	// local broadcast channel implementation is global for all tests;
	// to avoid conflicts between tests we need to randomize channel name
	// so that no channel name is shared between two tests
	randomSelector, err := rand.Int(rand.Reader, big.NewInt(10000000000))
	if err != nil {
		return nil, err
	}
	broadcastChannel, err := network.BroadcastChannelFor(
		fmt.Sprintf("entry-test-%v", randomSelector),
	)
	if err != nil {
		return nil, err
	}

	entrySubmissionChan := make(chan *event.EntrySubmitted)
	_ = relayChain.OnRelayEntrySubmitted(
		func(event *event.EntrySubmitted) {
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

	entry.RegisterUnmarshallers(broadcastChannel)

	for _, signer := range signers {
		go func(signer *dkg.ThresholdSigner) {
			err := entry.SignAndSubmit(
				blockCounter,
				broadcastChannel,
				relayChain,
				previousEntry,
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

	// We give 5 seconds so that OnRelayEntrySubmitted async handler
	// is fired. If it's not, it means no entry was published to
	// the chain.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case <-entrySubmissionChan:
		entry := lastRelayEntryGetter()
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
