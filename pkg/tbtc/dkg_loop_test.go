package tbtc

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
)

func TestDkgRetryLoop(t *testing.T) {
	seed := big.NewInt(100)

	groupParameters := &GroupParameters{
		GroupSize:       10,
		GroupQuorum:     8,
		HonestThreshold: 6,
	}

	selectedOperators := chain.Addresses{
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

	membersIndexes := make([]group.MemberIndex, 0)
	for i := range selectedOperators {
		membersIndexes = append(
			membersIndexes,
			group.MemberIndex(i+1),
		)
	}

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	testResult := &dkg.Result{
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
	}

	var tests = map[string]struct {
		memberIndex             group.MemberIndex
		ctxFn                   func() (context.Context, context.CancelFunc)
		incomingAnnouncementsFn func(sessionID string) ([]group.MemberIndex, error)
		dkgAttemptFn            dkgAttemptFn
		expectedErr             error
		expectedResult          *dkg.Result
		expectedLastAttempt     *dkgAttemptParams
	}{
		"success on initial attempt": {
			memberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(sessionID string) ([]group.MemberIndex, error) {
				return membersIndexes, nil
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, error) {
				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			expectedLastAttempt: &dkgAttemptParams{
				number:                 1,
				startBlock:             206,
				timeoutBlock:           406, // start block + 200
				excludedMembersIndexes: []group.MemberIndex{},
			},
		},
		"success on initial attempt with missing announcements and quorum": {
			memberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(sessionID string) ([]group.MemberIndex, error) {
				// Quorum of members announced their readiness.
				return []group.MemberIndex{1, 2, 3, 4, 5, 6, 7, 8}, nil
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, error) {
				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// As only 8 members (group quorum) announced their readiness,
			// we don't have any other option than select them for the attempt.
			// Not ready members are excluded.
			expectedLastAttempt: &dkgAttemptParams{
				number:                 1,
				startBlock:             206,
				timeoutBlock:           406, // start block + 200
				excludedMembersIndexes: []group.MemberIndex{9, 10},
			},
		},
		"missing announcements without quorum on initial attempt": {
			memberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(sessionID string) ([]group.MemberIndex, error) {
				if sessionID == fmt.Sprintf("%v-%v", seed, 1) {
					// Non-quorum of members announced their readiness.
					return []group.MemberIndex{1, 2, 3, 4, 5, 6, 7}, nil
				}

				return membersIndexes, nil
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, error) {
				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// First attempt fails because the group quorum did not announce
			// readiness.
			expectedLastAttempt: &dkgAttemptParams{
				number:                 2,
				startBlock:             417, // 206 + 1 * (6 + 200 + 5)
				timeoutBlock:           617, // start block + 200
				excludedMembersIndexes: []group.MemberIndex{2, 5},
			},
		},
		"announcement error on initial attempt": {
			memberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(sessionID string) ([]group.MemberIndex, error) {
				if sessionID == fmt.Sprintf("%v-%v", seed, 1) {
					return nil, fmt.Errorf("unexpected error")
				}

				return membersIndexes, nil
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, error) {
				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// First attempt fails due to the announcer error.
			expectedLastAttempt: &dkgAttemptParams{
				number:                 2,
				startBlock:             417, // 206 + 1 * (6 + 200 + 5)
				timeoutBlock:           617, // start block + 200
				excludedMembersIndexes: []group.MemberIndex{2, 5},
			},
		},
		"DKG error on initial attempt": {
			memberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(sessionID string) ([]group.MemberIndex, error) {
				return membersIndexes, nil
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, error) {
				if attempt.number == 1 {
					return nil, fmt.Errorf("invalid data")
				}

				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// The DKG error occurs on attempt 1. For attempt 2, all members are
			// ready, and we use the random algorithm to exclude some members.
			// The algorithm excludes members 2 and 5 (same operator) for the
			// given seed.
			expectedLastAttempt: &dkgAttemptParams{
				number:                 2,
				startBlock:             417, // 206 + 1 * (6 + 200 + 5)
				timeoutBlock:           617, // start block + 150
				excludedMembersIndexes: []group.MemberIndex{2, 5},
			},
		},
		"executing member excluded": {
			memberIndex: 5,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(sessionID string) ([]group.MemberIndex, error) {
				return membersIndexes, nil
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, error) {
				if attempt.number <= 2 {
					return nil, fmt.Errorf("invalid data")
				}

				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// Member 5 is the executing one. First attempt fails and is
			// retried using the random algorithm. The 2nd attempt does not
			// return an error but member 5 is excluded for this attempt so,
			// member 5 skips attempt 2 and succeeds on attempt 3.
			expectedLastAttempt: &dkgAttemptParams{
				number:                 3,
				startBlock:             628, // 206 + 2 * (6 + 200 + 5)
				timeoutBlock:           828, // start block + 200
				excludedMembersIndexes: []group.MemberIndex{9},
			},
		},
		"loop context done": {
			memberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				ctx, cancelCtx := context.WithCancel(context.Background())
				// Cancel the context deliberately.
				cancelCtx()
				return ctx, cancelCtx
			},
			incomingAnnouncementsFn: func(sessionID string) ([]group.MemberIndex, error) {
				return membersIndexes, nil
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, error) {
				return nil, fmt.Errorf("invalid data")
			},
			expectedErr:         context.Canceled,
			expectedResult:      nil,
			expectedLastAttempt: nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			announcer := &mockDkgAnnouncer{
				outgoingAnnouncements:   make(map[string]group.MemberIndex),
				incomingAnnouncementsFn: test.incomingAnnouncementsFn,
			}

			retryLoop := newDkgRetryLoop(
				&testutils.MockLogger{},
				seed,
				200,
				test.memberIndex,
				selectedOperators,
				groupParameters,
				announcer,
			)

			ctx, cancelCtx := test.ctxFn()
			defer cancelCtx()

			var lastAttempt *dkgAttemptParams

			result, err := retryLoop.start(
				ctx,
				func(ctx context.Context, attemptStartBlock uint64) error {
					return nil
				},
				func(params *dkgAttemptParams) (*dkg.Result, error) {
					lastAttempt = params
					return test.dkgAttemptFn(params)
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
				testutils.AssertIntsEqual(
					t,
					"outgoing announcements count",
					int(test.expectedLastAttempt.number),
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
						int(test.memberIndex),
						int(memberIndex),
					)
				}
			}
		})
	}
}

type mockDkgAnnouncer struct {
	// outgoingAnnouncements holds all announcements that are sent by the
	// announcer.
	outgoingAnnouncements map[string]group.MemberIndex

	// incomingAnnouncementsFn returns all announcements that are received
	// by the announcer for the given attempt.
	incomingAnnouncementsFn func(
		sessionID string,
	) ([]group.MemberIndex, error)
}

func (mda *mockDkgAnnouncer) Announce(
	ctx context.Context,
	memberIndex group.MemberIndex,
	sessionID string,
) ([]group.MemberIndex, error) {
	mda.outgoingAnnouncements[sessionID] = memberIndex

	return mda.incomingAnnouncementsFn(sessionID)
}
