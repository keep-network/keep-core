package dkg2

import (
	"math/big"
	"sync"
	"testing"

	"github.com/keep-network/keep-core/pkg/altbn128"
	"github.com/keep-network/keep-core/pkg/internal/testutils"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	chainLocal "github.com/keep-network/keep-core/pkg/chain/local"
	netLocal "github.com/keep-network/keep-core/pkg/net/local"
)

type result struct {
	signer *ThresholdSigner
	err    error
}

func TestExecuteDKGLocal(t *testing.T) {
	groupSize := 5
	threshold := 3

	requestID := big.NewInt(13)
	seed := big.NewInt(8)

	networkProvider := netLocal.Connect()
	chainHandle := chainLocal.Connect(groupSize, threshold, big.NewInt(200))

	blockCounter, err := chainHandle.BlockCounter()
	if err != nil {
		t.Fatal(err)
	}

	executeDKG := func(playerIndex int) (*ThresholdSigner, error) {
		broadcastChannel, err := networkProvider.ChannelFor("testing_channel")
		if err != nil {
			t.Fatalf("cannot generate broadcast channel [%v]", err)
		}

		signer, err := ExecuteDKG(
			requestID,
			seed,
			playerIndex,
			groupSize,
			threshold,
			blockCounter,
			chainHandle.ThresholdRelay(),
			broadcastChannel,
		)

		return signer, err
	}

	resultsChannel := make(chan *result, groupSize)

	var wg sync.WaitGroup
	wg.Add(groupSize)
	for i := 0; i < groupSize; i++ {
		memberID := i
		go func() {
			signer, err := executeDKG(memberID)
			resultsChannel <- &result{signer, err}
			wg.Done()
		}()
	}
	wg.Wait()
	close(resultsChannel)

	// read all results from the channel and put them into a slice
	var resultsSlice []*result
	for result := range resultsChannel {
		resultsSlice = append(resultsSlice, result)
	}

	// assert no errors occured
	for _, result := range resultsSlice {
		if result.err != nil {
			t.Errorf("unexpected error: [%v]", result.err)
		}
	}

	// assert all signers share the same group key
	for i := 1; i < groupSize; i++ {
		key0 := resultsSlice[0].signer.GroupPublicKeyBytes()
		keyi := resultsSlice[i].signer.GroupPublicKeyBytes()

		testutils.AssertBytesEqual(t, key0, keyi)
	}

	// TODO: Add verification of result submission when new Phase 14 is added to
	// states machine.
}

func TestConvertResult(t *testing.T) {
	groupSize := 5

	publicKey := new(bn256.G2).ScalarBaseMult(big.NewInt(2))
	compressedPublicKey := altbn128.G2Point{G2: publicKey}.Compress()

	var tests = map[string]struct {
		gjkrResult     *gjkr.Result
		expectedResult *relayChain.DKGResult
	}{
		"success: false, group public key: nil, DQ and IA: empty": {
			gjkrResult: &gjkr.Result{
				GroupPublicKey: nil,
				Disqualified:   []group.MemberIndex{},
				Inactive:       []group.MemberIndex{},
			},
			expectedResult: &relayChain.DKGResult{
				GroupPublicKey: []byte{},
				Disqualified:   []byte{0x00, 0x00, 0x00, 0x00, 0x00},
				Inactive:       []byte{0x00, 0x00, 0x00, 0x00, 0x00},
			},
		},
		"success: true, group public key: provided, DQ and IA: provided": {
			gjkrResult: &gjkr.Result{
				GroupPublicKey: publicKey,
				Disqualified:   []group.MemberIndex{1, 3, 4},
				Inactive:       []group.MemberIndex{5},
			},
			expectedResult: &relayChain.DKGResult{
				GroupPublicKey: compressedPublicKey,
				Disqualified:   []byte{0x01, 0x00, 0x01, 0x01, 0x00},
				Inactive:       []byte{0x00, 0x00, 0x00, 0x00, 0x01},
			},
		},
	}
	for _, test := range tests {
		convertedResult := convertResult(test.gjkrResult, groupSize)

		if !test.expectedResult.Equals(convertedResult) {
			t.Fatalf("\nexpected: %v\nactual:   %v\n", test.expectedResult, convertedResult)
		}
	}
}
