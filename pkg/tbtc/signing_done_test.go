package tbtc

import (
	"context"
	"fmt"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/keep-core/pkg/tecdsa/signing"
	"golang.org/x/exp/slices"
	"math/big"
	"sync"
	"testing"
)

func TestSigningDoneCheck(t *testing.T) {
	chainConfig := &ChainConfig{
		GroupSize:       5,
		GroupQuorum:     4,
		HonestThreshold: 3,
	}

	doneCheck := setupSigningDoneCheck(t, chainConfig)

	memberIndexes := make([]group.MemberIndex, doneCheck.groupSize)
	for i := range memberIndexes {
		memberIndex := group.MemberIndex(i + 1)
		memberIndexes[i] = memberIndex
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	message := big.NewInt(100)
	attemptNumber := uint64(2)
	attemptMemberIndexes := memberIndexes[:chainConfig.HonestThreshold]
	result := &signing.Result{
		Signature: &tecdsa.Signature{
			R:          big.NewInt(200),
			S:          big.NewInt(300),
			RecoveryID: 2,
		},
	}

	type outcome struct {
		memberIndex group.MemberIndex
		signature   *tecdsa.Signature
		endBlock    uint64
		err         error
	}

	wg := sync.WaitGroup{}
	wg.Add(len(memberIndexes))
	outcomesChan := make(chan *outcome, len(memberIndexes))

	for _, memberIndex := range memberIndexes {
		go func(memberIndex group.MemberIndex) {
			defer wg.Done()

			if slices.Contains(attemptMemberIndexes, memberIndex) {
				endBlock, err := doneCheck.exchange(
					ctx,
					memberIndex,
					message,
					attemptNumber,
					attemptMemberIndexes,
					result,
					500+uint64(memberIndex),
				)
				outcomesChan <- &outcome{
					memberIndex: memberIndex,
					signature:   result.Signature,
					endBlock:    endBlock,
					err:         err,
				}
			} else {
				result, endBlock, err := doneCheck.listen(
					ctx,
					message,
					attemptNumber,
					attemptMemberIndexes,
				)
				outcomesChan <- &outcome{
					memberIndex: memberIndex,
					signature:   result.Signature,
					endBlock:    endBlock,
					err:         err,
				}
			}
		}(memberIndex)
	}

	wg.Wait()
	close(outcomesChan)

	// We exchanged `500 + uint64(memberIndex)` and latest member has index 3.
	expectedEndBlock := 503

	for outcome := range outcomesChan {
		if !result.Signature.Equals(outcome.signature) {
			t.Errorf(
				"unexpected signature in for member [%v]\n"+
					"expected: [%v]\n"+
					"actual:   [%v]",
				outcome.memberIndex,
				result.Signature,
				outcome.signature,
			)
		}

		testutils.AssertIntsEqual(
			t,
			fmt.Sprintf("end block for member [%v]", outcome.memberIndex),
			expectedEndBlock,
			int(outcome.endBlock),
		)

		if outcome.err != nil {
			t.Errorf(
				"unexpected error for member [%v]: [%v]",
				outcome.memberIndex,
				outcome.err,
			)
		}
	}
}

// setupSigningDoneCheck sets up an instance of the signing done check ready
// to perform test checks.
func setupSigningDoneCheck(
	t *testing.T,
	chainConfig *ChainConfig,
) *signingDoneCheck {
	operatorPrivateKey, operatorPublicKey, err := operator.GenerateKeyPair(
		local_v1.DefaultCurve,
	)
	if err != nil {
		t.Fatal(err)
	}

	localChain := ConnectWithKey(
		chainConfig.GroupSize,
		chainConfig.GroupQuorum,
		chainConfig.HonestThreshold,
		operatorPrivateKey,
	)

	localProvider := local.ConnectWithKey(operatorPublicKey)

	operatorAddress, err := localChain.Signing().PublicKeyToAddress(
		operatorPublicKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	var operators []chain.Address
	for i := 0; i < chainConfig.GroupSize; i++ {
		operators = append(operators, operatorAddress)
	}

	broadcastChannel, err := localProvider.BroadcastChannelFor("channel")
	if err != nil {
		t.Fatal(err)
	}

	membershipValidator := group.NewMembershipValidator(
		&testutils.MockLogger{},
		operators,
		localChain.Signing(),
	)

	return newSigningDoneCheck(
		chainConfig.GroupSize,
		broadcastChannel,
		membershipValidator,
	)
}
