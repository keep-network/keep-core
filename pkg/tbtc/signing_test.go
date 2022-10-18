package tbtc

import (
	"context"
	"fmt"
	fuzz "github.com/google/gofuzz"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/operator"
	"math/big"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/keep-core/pkg/tecdsa/signing"
)

func TestSigningAnnouncementMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &signingAnnouncementMessage{
		senderID:      group.MemberIndex(38),
		message:       big.NewInt(100),
		attemptNumber: 3,
	}
	unmarshaled := &signingAnnouncementMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzSigningAnnouncementMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID      group.MemberIndex
			message       *big.Int
			attemptNumber uint64
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&message)
		f.Fuzz(&attemptNumber)

		msg := &signingAnnouncementMessage{
			senderID:      senderID,
			message:       message,
			attemptNumber: attemptNumber,
		}

		_ = pbutils.RoundTrip(msg, &signingAnnouncementMessage{})
	}
}

func TestFuzzSigningAnnouncementMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&signingAnnouncementMessage{})
}

func TestSigningRetryLoop(t *testing.T) {
	chainConfig := &ChainConfig{
		GroupSize:       10,
		HonestThreshold: 6,
	}

	signingGroupOperators := chain.Addresses{
		"address-1",
		"address-2",
		"address-8",
		"address-4",
		"address-2",
		"address-6",
		"address-7",
		"address-8",
		"address-9",
		"address-8",
	}

	signingGroupMembersIndexes := make([]group.MemberIndex, 0)
	for i := range signingGroupOperators {
		signingGroupMembersIndexes = append(
			signingGroupMembersIndexes,
			group.MemberIndex(i+1),
		)
	}

	testResult := &signing.Result{
		Signature: &tecdsa.Signature{
			R:          big.NewInt(300),
			S:          big.NewInt(400),
			RecoveryID: 2,
		},
	}

	var tests = map[string]struct {
		signingGroupMemberIndex group.MemberIndex
		ctxFn                   func() (context.Context, context.CancelFunc)
		incomingAnnouncementsFn func(attemptNumber uint64) ([]group.MemberIndex, error)
		signingAttemptFn        signingAttemptFn
		expectedErr             error
		expectedResult          *signing.Result
		expectedLastAttempt     *signingAttemptParams
	}{
		"success on initial attempt": {
			signingGroupMemberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			incomingAnnouncementsFn: func(
				attemptNumber uint64,
			) ([]group.MemberIndex, error) {
				return signingGroupMembersIndexes, nil
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// The signing random retry algorithm invoked with the test seed
			// excludes 4 members (6 is the honest threshold) from the first
			// attempt: 3, 7, 8 and 10.
			expectedLastAttempt: &signingAttemptParams{
				number:                 1,
				startBlock:             206,
				excludedMembersIndexes: []group.MemberIndex{3, 7, 8, 10},
			},
		},
		"success on initial attempt with missing announcements": {
			signingGroupMemberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			incomingAnnouncementsFn: func(
				attemptNumber uint64,
			) ([]group.MemberIndex, error) {
				// Only the following members announced their readiness.
				return []group.MemberIndex{1, 2, 3, 6, 7, 9}, nil
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// As only 6 members (honest threshold) announced their readiness,
			// we don't have any other option than select them for the attempt.
			// Not ready members are excluded.
			expectedLastAttempt: &signingAttemptParams{
				number:                 1,
				startBlock:             206,
				excludedMembersIndexes: []group.MemberIndex{4, 5, 8, 10},
			},
		},
		"announcement error on initial attempt": {
			signingGroupMemberIndex: 4,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			incomingAnnouncementsFn: func(
				attemptNumber uint64,
			) ([]group.MemberIndex, error) {
				if attemptNumber <= 1 {
					return nil, fmt.Errorf("unexpected error")
				}

				return signingGroupMembersIndexes, nil
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// Member 4 is the executing one. The first attempt fails and
			// the signing random retry algorithm invoked with the test seed
			// excludes 3 members (6 is the honest threshold) from the second
			// attempt: 1, 2 and 5. The additional exclusion round that trims
			// the included members list to the honest threshold size adds
			// member 9 to the final excluded members list.
			expectedLastAttempt: &signingAttemptParams{
				number:                 2,
				startBlock:             290, // 206 + 1 * (6 + 73 + 5)
				excludedMembersIndexes: []group.MemberIndex{1, 2, 5, 9},
			},
		},
		"signing error on initial attempt": {
			signingGroupMemberIndex: 4,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			incomingAnnouncementsFn: func(
				attemptNumber uint64,
			) ([]group.MemberIndex, error) {
				return signingGroupMembersIndexes, nil
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				if attempt.number <= 1 {
					return nil, fmt.Errorf("invalid data")
				}

				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// Member 4 is the executing one. The first attempt fails and
			// the signing random retry algorithm invoked with the test seed
			// excludes 3 members (6 is the honest threshold) from the second
			// attempt: 1, 2 and 5. The additional exclusion round that trims
			// the included members list to the honest threshold size adds
			// member 9 to the final excluded members list.
			expectedLastAttempt: &signingAttemptParams{
				number:                 2,
				startBlock:             290, // 206 + 1 * (6 + 73 + 5)
				excludedMembersIndexes: []group.MemberIndex{1, 2, 5, 9},
			},
		},
		"executing member excluded": {
			signingGroupMemberIndex: 2,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			incomingAnnouncementsFn: func(
				attemptNumber uint64,
			) ([]group.MemberIndex, error) {
				return signingGroupMembersIndexes, nil
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				if attempt.number <= 5 {
					return nil, fmt.Errorf("invalid data")
				}

				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// Member 2 is the executing one. First 5 attempts fail and are
			// retried using the random algorithm. The 6th attempt does not
			// return an error but member 2 is excluded for this attempt so,
			// member 2 skips attempt 6 and ends on attempt 7.
			expectedLastAttempt: &signingAttemptParams{
				number:                 7,
				startBlock:             710, // 206 + 6 * (6 + 73 + 5)
				excludedMembersIndexes: []group.MemberIndex{1, 5, 6, 9},
			},
		},
		"loop context done": {
			signingGroupMemberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				ctx, cancelCtx := context.WithCancel(context.Background())
				// Cancel the context deliberately.
				cancelCtx()
				return ctx, cancelCtx
			},
			incomingAnnouncementsFn: func(
				attemptNumber uint64,
			) ([]group.MemberIndex, error) {
				return signingGroupMembersIndexes, nil
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				return nil, fmt.Errorf("invalid data")
			},
			expectedErr:         nil,
			expectedResult:      nil,
			expectedLastAttempt: nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			announcer := &mockSigningAnnouncer{
				outgoingAnnouncements: make(map[uint64]struct {
					signingGroupMemberIndex group.MemberIndex
					message                 *big.Int
				}),
				incomingAnnouncementsFn: test.incomingAnnouncementsFn,
			}

			message := big.NewInt(100)

			retryLoop := newSigningRetryLoop(
				&testutils.MockLogger{},
				message,
				200,
				test.signingGroupMemberIndex,
				signingGroupOperators,
				chainConfig,
				announcer,
			)

			ctx, cancelCtx := test.ctxFn()
			defer cancelCtx()

			var lastAttempt *signingAttemptParams

			result, err := retryLoop.start(
				ctx,
				func(context.Context, uint64) error {
					return nil
				},
				func(params *signingAttemptParams) (*signing.Result, error) {
					lastAttempt = params
					return test.signingAttemptFn(params)
				},
			)

			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf(
					"unexpected error\n"+
						"expected: [%+v]\n"+
						"actual:   [%+v]",
					test.expectedErr,
					err,
				)
			}

			if !reflect.DeepEqual(test.expectedResult, result) {
				t.Errorf(
					"unexpected result\n"+
						"expected: [%+v]\n"+
						"actual:   [%+v]",
					test.expectedResult,
					result,
				)
			}

			if !reflect.DeepEqual(test.expectedLastAttempt, lastAttempt) {
				t.Errorf(
					"unexpected last attempt\n"+
						"expected: [%+v]\n"+
						"actual:   [%+v]",
					test.expectedLastAttempt,
					lastAttempt,
				)
			}

			if test.expectedLastAttempt != nil {
				expectedOutgoingAnnouncementsCount :=
					int(test.expectedLastAttempt.number)
				actualOutgoingAnnouncementsCount :=
					len(announcer.outgoingAnnouncements)
				if expectedOutgoingAnnouncementsCount !=
					actualOutgoingAnnouncementsCount {
					t.Errorf(
						"unexpected outgoing announcements count\n"+
							"expected: [%+v]\n"+
							"actual:   [%+v]",
						expectedOutgoingAnnouncementsCount,
						actualOutgoingAnnouncementsCount,
					)
				}

				for attemptNumber, outgoingAnnouncement := range announcer.outgoingAnnouncements {
					if test.signingGroupMemberIndex !=
						outgoingAnnouncement.signingGroupMemberIndex {
						t.Errorf(
							"unexpected outgoing announcement's member "+
								"index for attempt [%v]\n"+
								"expected: [%+v]\n"+
								"actual:   [%+v]",
							attemptNumber,
							test.signingGroupMemberIndex,
							outgoingAnnouncement.signingGroupMemberIndex,
						)
					}

					if message.Cmp(outgoingAnnouncement.message) != 0 {
						t.Errorf(
							"unexpected outgoing announcement's message "+
								"for attempt [%v]\n"+
								"expected: [%+v]\n"+
								"actual:   [%+v]",
							attemptNumber,
							message,
							outgoingAnnouncement.message,
						)
					}
				}
			}
		})
	}
}

func TestBroadcastSigningAnnouncer(t *testing.T) {
	chainConfig := &ChainConfig{
		GroupSize:       5,
		HonestThreshold: 3,
	}

	operatorPrivateKey, operatorPublicKey, err := operator.GenerateKeyPair(
		local_v1.DefaultCurve,
	)
	if err != nil {
		t.Fatal(err)
	}

	localChain := local_v1.ConnectWithKey(
		chainConfig.GroupSize,
		chainConfig.HonestThreshold,
		operatorPrivateKey,
	)

	operatorAddress, err := localChain.Signing().PublicKeyToAddress(
		operatorPublicKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	var signingGroupOperators []chain.Address
	for i := 0; i < chainConfig.GroupSize; i++ {
		signingGroupOperators = append(signingGroupOperators, operatorAddress)
	}

	localProvider := local.ConnectWithKey(operatorPublicKey)

	type memberResult struct {
		memberIndex group.MemberIndex
		readyMembersIndexes []group.MemberIndex
	}

	type memberError struct {
		memberIndex group.MemberIndex
		err error
	}

	var tests = map[string]struct {
		message                    *big.Int
		broadcastingMembersIndexes []group.MemberIndex
		expectedErrors             map[group.MemberIndex]error
		expectedResults            map[group.MemberIndex][]group.MemberIndex
	}{
		"all members broadcasted announcements": {
			message: big.NewInt(100),
			broadcastingMembersIndexes: []group.MemberIndex{1, 2, 3, 4, 5},
			expectedErrors: make(map[group.MemberIndex]error),
			expectedResults: map[group.MemberIndex][]group.MemberIndex{
				1: {1, 2, 3, 4, 5},
				2: {1, 2, 3, 4, 5},
				3: {1, 2, 3, 4, 5},
				4: {1, 2, 3, 4, 5},
				5: {1, 2, 3, 4, 5},
			},
		},
		"honest majority of members broadcasted announcements": {
			message: big.NewInt(200),
			broadcastingMembersIndexes: []group.MemberIndex{1, 3, 5},
			expectedErrors: make(map[group.MemberIndex]error),
			expectedResults: map[group.MemberIndex][]group.MemberIndex{
				1: {1, 3, 5},
				3: {1, 3, 5},
				5: {1, 3, 5},
			},
		},
		"minority of members broadcasted announcements": {
			message: big.NewInt(300),
			broadcastingMembersIndexes: []group.MemberIndex{1, 3},
			expectedErrors: map[group.MemberIndex]error{
				1: fmt.Errorf("ready members count is lesser than the honest threshold"),
				3: fmt.Errorf("ready members count is lesser than the honest threshold"),
			},
			expectedResults: make(map[group.MemberIndex][]group.MemberIndex),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			broadcastChannel, err := localProvider.BroadcastChannelFor(
				test.message.Text(16),
			)
			if err != nil {
				t.Fatal(err)
			}

			membershipValidator := group.NewMembershipValidator(
				&testutils.MockLogger{},
				signingGroupOperators,
				localChain.Signing(),
			)

			announcer := newBroadcastSigningAnnouncer(
				chainConfig,
				broadcastChannel,
				membershipValidator,
			)

			resultsChan := make(
				chan *memberResult,
				len(test.broadcastingMembersIndexes),
			)
			errorsChan := make(
				chan *memberError,
				len(test.broadcastingMembersIndexes),
			)

			wg := sync.WaitGroup{}
			wg.Add(len(test.broadcastingMembersIndexes))

			for _, broadcastingMemberIndex :=
				range test.broadcastingMembersIndexes {
				go func(memberIndex group.MemberIndex) {
					defer wg.Done()

					ctx, cancelCtx := context.WithTimeout(
						context.Background(),
						100 * time.Millisecond,
					)
					defer cancelCtx()

					readyMembersIndexes, err := announcer.announce(
						ctx,
						memberIndex,
						test.message,
						1,
					)
					if err != nil {
						errorsChan <- &memberError{memberIndex, err}
						return
					}

					resultsChan <- &memberResult{memberIndex, readyMembersIndexes}
				}(broadcastingMemberIndex)
			}

			wg.Wait()

			close(resultsChan)
			results := make(map[group.MemberIndex][]group.MemberIndex)
			for r := range resultsChan {
				results[r.memberIndex] = r.readyMembersIndexes
			}

			close(errorsChan)
			errors := make(map[group.MemberIndex]error)
			for e := range errorsChan {
				errors[e.memberIndex] = e.err
			}

			if !reflect.DeepEqual(test.expectedErrors, errors) {
				t.Errorf(
					"unexpected errors\n" +
						"expected: [%v]\n" +
						"actual:   [%v]",
					test.expectedErrors,
					errors,
				)
			}

			if !reflect.DeepEqual(test.expectedResults, results) {
				t.Errorf(
					"unexpected results\n" +
						"expected: [%v]\n" +
						"actual:   [%v]",
					test.expectedResults,
					results,
				)
			}
		})
	}
}

type mockSigningAnnouncer struct {
	outgoingAnnouncements map[uint64]struct {
		signingGroupMemberIndex group.MemberIndex
		message                 *big.Int
	}

	incomingAnnouncementsFn func(
		attemptNumber uint64,
	) ([]group.MemberIndex, error)
}

func (msa *mockSigningAnnouncer) announce(
	ctx context.Context,
	signingGroupMemberIndex group.MemberIndex,
	message *big.Int,
	attemptNumber uint64,
) ([]group.MemberIndex, error) {
	msa.outgoingAnnouncements[attemptNumber] = struct {
		signingGroupMemberIndex group.MemberIndex
		message                 *big.Int
	}{signingGroupMemberIndex, message}

	return msa.incomingAnnouncementsFn(attemptNumber)
}
