// Package dkgtest provides a full roundtrip DKG test engine including all
// the phases. It is executed against local chain and broadcast channel.
package dkgtest

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"math"
	"math/big"
	"sync"
	"testing"
	"time"

	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
	"github.com/keep-network/keep-core/pkg/beacon/dkg"
	dkgResult "github.com/keep-network/keep-core/pkg/beacon/dkg/result"
	"github.com/keep-network/keep-core/pkg/beacon/event"
	"github.com/keep-network/keep-core/pkg/beacon/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/group"
	"github.com/keep-network/keep-core/pkg/internal/interception"
	netLocal "github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/operator"
)

// Result of a DKG test execution.
type Result struct {
	dkgResult           *beaconchain.DKGResult
	dkgResultSignatures map[group.MemberIndex][]byte
	signers             []*dkg.ThresholdSigner
	memberFailures      []error
}

// GetSigners returns all signers created from DKG protocol execution.
// If no signers were created because of protocol failures, empty slice
// is returned.
func (r *Result) GetSigners() []*dkg.ThresholdSigner {
	return r.signers
}

// RandomSeed generates a random DKG seed value. It is important to do not
// reuse the same seed value between integration tests run in parallel.
// Broadcast channel name contains a seed to avoid mixing up channel messages
// between two or more tests executed in parallel.
func RandomSeed(t *testing.T) *big.Int {
	seed, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		t.Fatal(err)
	}
	return seed
}

// RunTest executes the full DKG roundrip test for the provided group size,
// seed, and honest threshold. The provided interception rules are applied in
// the broadcast channel for the time of DKG execution.
func RunTest(
	groupSize int,
	honestThreshold int,
	seed *big.Int,
	rules interception.Rules,
) (*Result, error) {
	operatorPrivateKey, operatorPublicKey, err := operator.GenerateKeyPair(local_v1.DefaultCurve)
	if err != nil {
		return nil, err
	}

	network := interception.NewNetwork(
		netLocal.ConnectWithKey(operatorPublicKey),
		rules,
	)

	localChain := local_v1.ConnectWithKey(
		groupSize,
		honestThreshold,
		operatorPrivateKey,
	)

	blockCounter, err := localChain.BlockCounter()
	if err != nil {
		return nil, err
	}

	address, err := localChain.Signing().PublicKeyToAddress(operatorPublicKey)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot convert operator public key to chain address: [%v]",
			err,
		)
	}

	selectedStakers := make([]chain.Address, groupSize)
	for i := range selectedStakers {
		selectedStakers[i] = address
	}

	return executeDKG(
		seed,
		localChain,
		blockCounter,
		localChain.GetLastDKGResult,
		network,
		selectedStakers,
	)
}

func executeDKG(
	seed *big.Int,
	beaconChain beaconchain.Interface,
	blockCounter chain.BlockCounter,
	lastDKGResultGetter func() (
		*beaconchain.DKGResult,
		map[beaconchain.GroupMemberIndex][]byte,
	),
	network interception.Network,
	selectedStakers []chain.Address,
) (*Result, error) {
	relayConfig := beaconChain.GetConfig()

	broadcastChannel, err := network.BroadcastChannelFor(fmt.Sprintf("dkg-test-%v", seed))
	if err != nil {
		return nil, err
	}

	resultSubmissionChan := make(chan *event.DKGResultSubmission)
	_ = beaconChain.OnDKGResultSubmitted(
		func(event *event.DKGResultSubmission) {
			resultSubmissionChan <- event
		},
	)

	var signersMutex sync.Mutex
	var signers []*dkg.ThresholdSigner

	var memberFailures []error

	var wg sync.WaitGroup
	wg.Add(relayConfig.GroupSize)

	currentBlockHeight, err := blockCounter.CurrentBlock()
	if err != nil {
		return nil, err
	}

	// Wait for 3 blocks before starting DKG to
	// make sure all members are up.
	startBlockHeight := currentBlockHeight + 3

	gjkr.RegisterUnmarshallers(broadcastChannel)
	dkgResult.RegisterUnmarshallers(broadcastChannel)

	membershipValidator := group.NewStakersMembershipValidator(
		selectedStakers,
		beaconChain.Signing(),
	)

	for i := 0; i < relayConfig.GroupSize; i++ {
		i := i // capture for goroutine
		go func() {
			signer, err := dkg.ExecuteDKG(
				seed,
				uint8(i),
				relayConfig.GroupSize,
				relayConfig.DishonestThreshold(),
				membershipValidator,
				startBlockHeight,
				blockCounter,
				beaconChain,
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
		dkgResult, dkgResultSignatures := lastDKGResultGetter()
		return &Result{
			dkgResult,
			dkgResultSignatures,
			signers,
			memberFailures,
		}, nil

	case <-ctx.Done():
		// no result published to the chain
		return &Result{
			nil,
			nil,
			signers,
			memberFailures,
		}, nil
	}
}
