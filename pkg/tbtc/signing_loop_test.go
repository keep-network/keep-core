package tbtc

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/internal/testutils"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/keep-core/pkg/tecdsa/signing"
)

func TestSigningRetryLoop(t *testing.T) {
	message := big.NewInt(100)

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
		signingGroupMemberIndex     group.MemberIndex
		ctxFn                       func() (context.Context, context.CancelFunc)
		incomingAnnouncementsFn     func(sessionID string) ([]group.MemberIndex, error)
		signingAttemptFn            signingAttemptFn
		exchangeDoneChecksOutcomeFn func(attemptNumber uint64, endBlock uint64) (uint64, error)
		listenDoneChecksOutcomeFn   func(attemptNumber uint64) (*signing.Result, uint64, error)
		expectedOutgoingDoneChecks  []*signingDoneMessage
		expectedErr                 error
		expectedResult              *signingRetryLoopResult
		expectedLastExecutedAttempt *signingAttemptParams
		outgoingAnnouncementsCount  uint
	}{
		"success on initial attempt": {
			signingGroupMemberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(
				sessionID string,
			) ([]group.MemberIndex, error) {
				return signingGroupMembersIndexes, nil
			},
			signingAttemptFn: func(
				attempt *signingAttemptParams,
			) (*signing.Result, uint64, error) {
				return testResult, 215, nil // an arbitrary end block
			},
			exchangeDoneChecksOutcomeFn: func(
				attemptNumber uint64,
				endBlock uint64,
			) (uint64, error) {
				// Simulate that the done check phase determines the same
				// end block as the executing signer.
				return endBlock, nil
			},
			expectedOutgoingDoneChecks: []*signingDoneMessage{
				{
					senderID:      1,
					message:       message,
					attemptNumber: 1,
					signature:     testResult.Signature,
					endBlock:      215,
				},
			},
			expectedErr: nil,
			expectedResult: &signingRetryLoopResult{
				result:              testResult,
				latestEndBlock:      215, // the end block resolved by the done check phase
				attemptTimeoutBlock: 226, // start block of the first attempt + 20
			},
			// The signing random retry algorithm invoked with the test seed
			// excludes 4 members (6 is the honest threshold) from the first
			// attempt: 3, 7, 8 and 10.
			expectedLastExecutedAttempt: &signingAttemptParams{
				number:                 1,
				startBlock:             206,
				timeoutBlock:           226, // start block of the first attempt + 20
				excludedMembersIndexes: []group.MemberIndex{3, 7, 8, 10},
			},
			outgoingAnnouncementsCount: 1,
		},
		"success on initial attempt with missing announcements and honest majority": {
			signingGroupMemberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(
				sessionID string,
			) ([]group.MemberIndex, error) {
				// Honest majority of members announced their readiness.
				return []group.MemberIndex{1, 2, 3, 6, 7, 9}, nil
			},
			signingAttemptFn: func(
				attempt *signingAttemptParams,
			) (*signing.Result, uint64, error) {
				return testResult, 215, nil // an arbitrary end block
			},
			exchangeDoneChecksOutcomeFn: func(
				attemptNumber uint64,
				endBlock uint64,
			) (uint64, error) {
				// Simulate that the done check phase determines the same
				// end block as the executing signer.
				return endBlock, nil
			},
			expectedOutgoingDoneChecks: []*signingDoneMessage{
				{
					senderID:      1,
					message:       message,
					attemptNumber: 1,
					signature:     testResult.Signature,
					endBlock:      215,
				},
			},
			expectedErr: nil,
			expectedResult: &signingRetryLoopResult{
				result:              testResult,
				latestEndBlock:      215, // the end block resolved by the done check phase
				attemptTimeoutBlock: 226, // start block of the first attempt + 20
			},
			// As only 6 members (honest threshold) announced their readiness,
			// we don't have any other option than select them for the attempt.
			// Not ready members are excluded.
			expectedLastExecutedAttempt: &signingAttemptParams{
				number:                 1,
				startBlock:             206,
				timeoutBlock:           226, // start block of the first attempt + 20
				excludedMembersIndexes: []group.MemberIndex{4, 5, 8, 10},
			},
			outgoingAnnouncementsCount: 1,
		},
		"missing announcements without honest majority on initial attempt": {
			signingGroupMemberIndex: 3,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(
				sessionID string,
			) ([]group.MemberIndex, error) {
				if sessionID == fmt.Sprintf("%v-%v", message, 1) {
					// Minority of members announced their readiness.
					return []group.MemberIndex{1, 2, 3, 6, 7}, nil
				}

				return signingGroupMembersIndexes, nil
			},
			signingAttemptFn: func(
				attempt *signingAttemptParams,
			) (*signing.Result, uint64, error) {
				return testResult, 247, nil // an arbitrary end block
			},
			exchangeDoneChecksOutcomeFn: func(
				attemptNumber uint64,
				endBlock uint64,
			) (uint64, error) {
				// Simulate that the done check phase determines the same
				// end block as the executing signer.
				return endBlock, nil
			},
			expectedOutgoingDoneChecks: []*signingDoneMessage{
				{
					senderID:      3,
					message:       message,
					attemptNumber: 2,
					signature:     testResult.Signature,
					endBlock:      247,
				},
			},
			expectedErr: nil,
			expectedResult: &signingRetryLoopResult{
				result:              testResult,
				latestEndBlock:      247, // the end block resolved by the done check phase
				attemptTimeoutBlock: 257, // start block of the second attempt + 20
			},
			// Member 3 is the executing one. The first attempt's announcement
			// fails and the signing random retry algorithm invoked with the
			// test seed excludes 3 members (6 is the honest threshold) from the
			// second attempt: 1, 2 and 5. The additional exclusion round that
			// trims the included members list to the honest threshold size
			// adds member 9 to the final excluded members list.
			expectedLastExecutedAttempt: &signingAttemptParams{
				number:                 2,
				startBlock:             237, // 206 + 1 * (6 + 20 + 5)
				timeoutBlock:           257, // start block of the second attempt + 20
				excludedMembersIndexes: []group.MemberIndex{1, 2, 5, 9},
			},
			outgoingAnnouncementsCount: 2,
		},
		"announcement error on initial attempt": {
			signingGroupMemberIndex: 4,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(
				sessionID string,
			) ([]group.MemberIndex, error) {
				if sessionID == fmt.Sprintf("%v-%v", message, 1) {
					return nil, fmt.Errorf("unexpected error")
				}

				return signingGroupMembersIndexes, nil
			},
			signingAttemptFn: func(
				attempt *signingAttemptParams,
			) (*signing.Result, uint64, error) {
				return testResult, 247, nil // an arbitrary end block
			},
			exchangeDoneChecksOutcomeFn: func(
				attemptNumber uint64,
				endBlock uint64,
			) (uint64, error) {
				// Simulate that the done check phase determines the same
				// end block as the executing signer.
				return endBlock, nil
			},
			expectedOutgoingDoneChecks: []*signingDoneMessage{
				{
					senderID:      4,
					message:       message,
					attemptNumber: 2,
					signature:     testResult.Signature,
					endBlock:      247,
				},
			},
			expectedErr: nil,
			expectedResult: &signingRetryLoopResult{
				result:              testResult,
				latestEndBlock:      247, // the end block resolved by the done check phase
				attemptTimeoutBlock: 257, // start block of the second attempt + 20
			},
			// Member 4 is the executing one. The first attempt fails and
			// the signing random retry algorithm invoked with the test seed
			// excludes 3 members (6 is the honest threshold) from the second
			// attempt: 1, 2 and 5. The additional exclusion round that trims
			// the included members list to the honest threshold size adds
			// member 9 to the final excluded members list.
			expectedLastExecutedAttempt: &signingAttemptParams{
				number:                 2,
				startBlock:             237, // 206 + 1 * (6 + 20 + 5)
				timeoutBlock:           257, // start block of the second attempt + 20
				excludedMembersIndexes: []group.MemberIndex{1, 2, 5, 9},
			},
			outgoingAnnouncementsCount: 2,
		},
		"signing error on initial attempt": {
			signingGroupMemberIndex: 4,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(
				sessionID string,
			) ([]group.MemberIndex, error) {
				return signingGroupMembersIndexes, nil
			},
			signingAttemptFn: func(
				attempt *signingAttemptParams,
			) (*signing.Result, uint64, error) {
				if attempt.number <= 1 {
					return nil, 0, fmt.Errorf("invalid data")
				}

				return testResult, 247, nil // an arbitrary end block
			},
			exchangeDoneChecksOutcomeFn: func(
				attemptNumber uint64,
				endBlock uint64,
			) (uint64, error) {
				// Simulate that the done check phase determines the same
				// end block as the executing signer.
				return endBlock, nil
			},
			expectedOutgoingDoneChecks: []*signingDoneMessage{
				{
					senderID:      4,
					message:       message,
					attemptNumber: 2,
					signature:     testResult.Signature,
					endBlock:      247,
				},
			},
			expectedErr: nil,
			expectedResult: &signingRetryLoopResult{
				result:              testResult,
				latestEndBlock:      247, // the end block resolved by the done check phase
				attemptTimeoutBlock: 257, // start block of the second attempt + 20
			},
			// Member 4 is the executing one. The first attempt fails and
			// the signing random retry algorithm invoked with the test seed
			// excludes 3 members (6 is the honest threshold) from the second
			// attempt: 1, 2 and 5. The additional exclusion round that trims
			// the included members list to the honest threshold size adds
			// member 9 to the final excluded members list.
			expectedLastExecutedAttempt: &signingAttemptParams{
				number:                 2,
				startBlock:             237, // 206 + 1 * (6 + 20 + 5)
				timeoutBlock:           257, // start block of the second attempt + 20
				excludedMembersIndexes: []group.MemberIndex{1, 2, 5, 9},
			},
			outgoingAnnouncementsCount: 2,
		},
		"executing member excluded": {
			signingGroupMemberIndex: 2,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(
				sessionID string,
			) ([]group.MemberIndex, error) {
				return signingGroupMembersIndexes, nil
			},
			signingAttemptFn: func(
				attempt *signingAttemptParams,
			) (*signing.Result, uint64, error) {
				return nil, 0, fmt.Errorf("invalid data")
			},
			listenDoneChecksOutcomeFn: func(
				attemptNumber uint64,
			) (*signing.Result, uint64, error) {
				// Simulate the result and the end block have been determined
				// by listening for signing done checks.
				if attemptNumber == 2 {
					return testResult, 247, nil
				}

				panic("undefined behavior")
			},
			expectedOutgoingDoneChecks: nil,
			expectedErr:                nil,
			expectedResult: &signingRetryLoopResult{
				result:              testResult,
				latestEndBlock:      247, // the end block resolved by the done check phase
				attemptTimeoutBlock: 257, // start block of the second attempt + 20
			},
			// Member 2 is the executing one. The first attempt fails
			// and is the last attempt executed by this member because member
			// 2 is excluded from the second attempt that produced the signature.
			expectedLastExecutedAttempt: &signingAttemptParams{
				number:                 1,
				startBlock:             206,
				timeoutBlock:           226, // start block of the first attempt + 20
				excludedMembersIndexes: []group.MemberIndex{3, 7, 8, 10},
			},
			// The second announcement is done at the beginning of the
			// second attempt for which member 2 is eventually excluded.
			outgoingAnnouncementsCount: 2,
		},
		"done checks exchange error": {
			signingGroupMemberIndex: 4,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(
				sessionID string,
			) ([]group.MemberIndex, error) {
				return signingGroupMembersIndexes, nil
			},
			signingAttemptFn: func(
				attempt *signingAttemptParams,
			) (*signing.Result, uint64, error) {
				if attempt.number == 1 {
					return testResult, 215, nil // an arbitrary end block
				}

				if attempt.number == 2 {
					return testResult, 247, nil // an arbitrary end block
				}

				panic("undefined behavior")
			},
			exchangeDoneChecksOutcomeFn: func(
				attemptNumber uint64,
				endBlock uint64,
			) (uint64, error) {
				// Fail the done check for the first attempt.
				if attemptNumber == 1 {
					return 0, fmt.Errorf("network error")
				}

				// Simulate that the done check phase determines the same
				// end block as the executing signer.
				return endBlock, nil
			},
			expectedOutgoingDoneChecks: []*signingDoneMessage{
				{
					senderID:      4,
					message:       message,
					attemptNumber: 1,
					signature:     testResult.Signature,
					endBlock:      215,
				},
				{
					senderID:      4,
					message:       message,
					attemptNumber: 2,
					signature:     testResult.Signature,
					endBlock:      247,
				},
			},
			expectedErr: nil,
			expectedResult: &signingRetryLoopResult{
				result:              testResult,
				latestEndBlock:      247, // the end block resolved by the done check phase
				attemptTimeoutBlock: 257, // start block of the second attempt + 20
			},
			// Member 4 is the executing one. The first attempt done check
			// exchange fails and the signing random retry algorithm invoked
			// with the test seed excludes 3 members (6 is the honest threshold)
			// from the second attempt: 1, 2 and 5. The additional exclusion
			// round that trims the included members list to the honest
			// threshold size adds member 9 to the final excluded members list.
			expectedLastExecutedAttempt: &signingAttemptParams{
				number:                 2,
				startBlock:             237, // 206 + 1 * (6 + 20 + 5)
				timeoutBlock:           257, // start block of the second attempt + 20
				excludedMembersIndexes: []group.MemberIndex{1, 2, 5, 9},
			},
			outgoingAnnouncementsCount: 2,
		},
		"done checks listen error": {
			signingGroupMemberIndex: 8,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(
				sessionID string,
			) ([]group.MemberIndex, error) {
				return signingGroupMembersIndexes, nil
			},
			signingAttemptFn: func(
				attempt *signingAttemptParams,
			) (*signing.Result, uint64, error) {
				if attempt.number == 2 {
					return testResult, 247, nil // an arbitrary end block
				}

				panic("undefined behavior")
			},
			listenDoneChecksOutcomeFn: func(
				attemptNumber uint64,
			) (*signing.Result, uint64, error) {
				// Fail the done check for the first attempt.
				if attemptNumber == 1 {
					return nil, 247, fmt.Errorf("network error")
				}

				panic("undefined behavior")
			},
			exchangeDoneChecksOutcomeFn: func(
				attemptNumber uint64,
				endBlock uint64,
			) (uint64, error) {
				// Simulate that the done check phase determines the same
				// end block as the executing signer.
				return endBlock, nil
			},
			expectedOutgoingDoneChecks: []*signingDoneMessage{
				{
					senderID:      8,
					message:       message,
					attemptNumber: 2,
					signature:     testResult.Signature,
					endBlock:      247,
				},
			},
			expectedErr: nil,
			expectedResult: &signingRetryLoopResult{
				result:              testResult,
				latestEndBlock:      247, // the end block resolved by the done check phase
				attemptTimeoutBlock: 257, // start block of the second attempt + 20
			},
			// Member 8 is the executing one and is excluded from the first attempt.
			// The first attempt done check listen fails and the signing random
			// retry algorithm invoked with the test seed excludes 3
			// members (6 is the honest threshold) from the second attempt:
			// 1, 2 and 5. The additional exclusion round that trims the
			// included members list to the honest threshold size adds
			//member 9 to the final excluded members list.
			expectedLastExecutedAttempt: &signingAttemptParams{
				number:                 2,
				startBlock:             237, // 206 + 1 * (6 + 20 + 5)
				timeoutBlock:           257, // start block of the second attempt + 20
				excludedMembersIndexes: []group.MemberIndex{1, 2, 5, 9},
			},
			outgoingAnnouncementsCount: 2,
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
				sessionID string,
			) ([]group.MemberIndex, error) {
				return signingGroupMembersIndexes, nil
			},
			signingAttemptFn: func(
				attempt *signingAttemptParams,
			) (*signing.Result, uint64, error) {
				return nil, 0, fmt.Errorf("invalid data")
			},
			expectedErr:                 context.Canceled,
			expectedResult:              nil,
			expectedLastExecutedAttempt: nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			announcer := &mockSigningAnnouncer{
				outgoingAnnouncements:   make(map[string]group.MemberIndex),
				incomingAnnouncementsFn: test.incomingAnnouncementsFn,
			}

			doneCheck := &mockSigningDoneCheck{
				exchangeDoneChecksOutcomeFn: test.exchangeDoneChecksOutcomeFn,
				listenDoneChecksOutcomeFn:   test.listenDoneChecksOutcomeFn,
			}

			retryLoop := newSigningRetryLoop(
				&testutils.MockLogger{},
				message,
				200,
				test.signingGroupMemberIndex,
				signingGroupOperators,
				chainConfig,
				announcer,
				doneCheck,
			)

			ctx, cancelCtx := test.ctxFn()
			defer cancelCtx()

			var lastExecutedAttempt *signingAttemptParams

			result, err := retryLoop.start(
				ctx,
				func(context.Context, uint64) error {
					return nil
				},
				func(params *signingAttemptParams) (*signing.Result, uint64, error) {
					lastExecutedAttempt = params
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

			if !reflect.DeepEqual(test.expectedLastExecutedAttempt, lastExecutedAttempt) {
				t.Errorf(
					"unexpected last executed attempt\n"+
						"expected: [%+v]\n"+
						"actual:   [%+v]",
					test.expectedLastExecutedAttempt,
					lastExecutedAttempt,
				)
			}

			if test.expectedLastExecutedAttempt != nil {
				testutils.AssertIntsEqual(
					t,
					"outgoing announcements count",
					int(test.outgoingAnnouncementsCount),
					len(announcer.outgoingAnnouncements),
				)

				for sessionID, memberIndex := range announcer.outgoingAnnouncements {
					testutils.AssertIntsEqual(
						t,
						fmt.Sprintf(
							"outgoing announcement's member index "+
								"for session [%v]",
							sessionID,
						),
						int(test.signingGroupMemberIndex),
						int(memberIndex),
					)
				}
			}

			if !reflect.DeepEqual(
				test.expectedOutgoingDoneChecks,
				doneCheck.outgoingDoneChecks,
			) {
				t.Errorf(
					"unexpected outgoing done checks\n"+
						"expected: [%v]\n"+
						"actual:   [%v]",
					test.expectedOutgoingDoneChecks,
					doneCheck.outgoingDoneChecks,
				)
			}
		})
	}
}

type mockSigningAnnouncer struct {
	// outgoingAnnouncements holds all announcements that are sent by the
	// announcer.
	outgoingAnnouncements map[string]group.MemberIndex

	// incomingAnnouncementsFn returns all announcements that are received
	// by the announcer for the given attempt.
	incomingAnnouncementsFn func(
		sessionID string,
	) ([]group.MemberIndex, error)
}

func (msa *mockSigningAnnouncer) Announce(
	ctx context.Context,
	memberIndex group.MemberIndex,
	sessionID string,
) ([]group.MemberIndex, error) {
	msa.outgoingAnnouncements[sessionID] = memberIndex

	return msa.incomingAnnouncementsFn(sessionID)
}

type mockSigningDoneCheck struct {
	outgoingDoneChecks []*signingDoneMessage

	exchangeDoneChecksOutcomeFn func(
		attemptNumber uint64,
		endBlock uint64,
	) (uint64, error)

	listenDoneChecksOutcomeFn func(
		attemptNumber uint64,
	) (*signing.Result, uint64, error)
}

func (msdc *mockSigningDoneCheck) exchange(
	ctx context.Context,
	memberIndex group.MemberIndex,
	message *big.Int,
	attemptNumber uint64,
	attemptTimeoutBlock uint64,
	attemptMembersIndexes []group.MemberIndex,
	result *signing.Result,
	endBlock uint64,
) (uint64, error) {
	msdc.outgoingDoneChecks = append(msdc.outgoingDoneChecks, &signingDoneMessage{
		senderID:      memberIndex,
		message:       message,
		attemptNumber: attemptNumber,
		signature:     result.Signature,
		endBlock:      endBlock,
	})

	return msdc.exchangeDoneChecksOutcomeFn(attemptNumber, endBlock)
}

func (msdc *mockSigningDoneCheck) listen(
	ctx context.Context,
	message *big.Int,
	attemptNumber uint64,
	attemptTimeoutBlock uint64,
	attemptMembersIndexes []group.MemberIndex,
) (*signing.Result, uint64, error) {
	return msdc.listenDoneChecksOutcomeFn(attemptNumber)
}
