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

func TestSigningRequestID(t *testing.T) {
	message1 := new(big.Int).SetBytes([]byte{0x03, 0x09, 0xAA, 0xFF, 0x45})
	message2 := new(big.Int).SetBytes([]byte{0x02, 0x08, 0xBB, 0xEE, 0x46})
	message3 := new(big.Int).SetBytes([]byte{0x01, 0x07, 0xCC, 0xDD, 0x47})

	requestID := newSigningRequestID([]*big.Int{message1, message2, message3})

	expectedRequestID := "e4e52e3d8044b72236e52b5c4e40dcf62afb46f7de48bc0d7630cac62ccda7a5"

	testutils.AssertStringsEqual(
		t,
		"request ID",
		expectedRequestID,
		requestID.String(),
	)
}

func TestSigningRetryLoop(t *testing.T) {
	requestID := newSigningRequestID([]*big.Int{big.NewInt(100)})

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
		incomingAnnouncementsFn func(sessionID string) ([]group.MemberIndex, error)
		signingAttemptFn        signingAttemptFn
		expectedErr             error
		expectedResult          *signing.Result
		expectedLastAttempt     *signingAttemptParams
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
		"missing announcements without honest majority on initial attempt": {
			signingGroupMemberIndex: 3,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(
				sessionID string,
			) ([]group.MemberIndex, error) {
				if sessionID == fmt.Sprintf("%s-%v", requestID, 1) {
					// Minority of members announced their readiness.
					return []group.MemberIndex{1, 2, 3, 6, 7}, nil
				}

				return signingGroupMembersIndexes, nil
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// Member 3 is the executing one. The first attempt's announcement
			// fails and the signing random retry algorithm invoked with the
			// test seed excludes 3 members (6 is the honest threshold) from the
			// second attempt: 1, 2 and 5. The additional exclusion round that
			// trims the included members list to the honest threshold size
			// adds member 9 to the final excluded members list.
			expectedLastAttempt: &signingAttemptParams{
				number:                 2,
				startBlock:             317, // 206 + 1 * (6 + 100 + 5)
				excludedMembersIndexes: []group.MemberIndex{1, 2, 5, 9},
			},
		},
		"announcement error on initial attempt": {
			signingGroupMemberIndex: 4,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			incomingAnnouncementsFn: func(
				sessionID string,
			) ([]group.MemberIndex, error) {
				if sessionID == fmt.Sprintf("%s-%v", requestID, 1) {
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
				startBlock:             317, // 206 + 1 * (6 + 100 + 5)
				excludedMembersIndexes: []group.MemberIndex{1, 2, 5, 9},
			},
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
				startBlock:             317, // 206 + 1 * (6 + 100 + 5)
				excludedMembersIndexes: []group.MemberIndex{1, 2, 5, 9},
			},
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
				startBlock:             872, // 206 + 6 * (6 + 100 + 5)
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
				sessionID string,
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
				outgoingAnnouncements:   make(map[string]group.MemberIndex),
				incomingAnnouncementsFn: test.incomingAnnouncementsFn,
			}

			retryLoop := newSigningRetryLoop(
				&testutils.MockLogger{},
				requestID,
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
						int(test.signingGroupMemberIndex),
						int(memberIndex),
					)
				}
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
