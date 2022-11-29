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
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
	"golang.org/x/crypto/sha3"
)

func TestDkgRetryLoop(t *testing.T) {
	seed := big.NewInt(100)

	chainConfig := &ChainConfig{
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
				startBlock:             367, // 206 + 1 * (6 + 150 + 5)
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
				startBlock:             367, // 206 + 1 * (6 + 150 + 5)
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
				startBlock:             367, // 206 + 1 * (6 + 150 + 5)
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
				startBlock:             528, // 206 + 2 * (6 + 150 + 5)
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
			expectedErr:         nil,
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
				chainConfig,
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

func TestDecideSigningGroupMemberFate(t *testing.T) {
	chainConfig := &ChainConfig{
		GroupSize:       10,
		GroupQuorum:     8,
		HonestThreshold: 6,
	}

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	memberIndex := group.MemberIndex(1)

	result := &dkg.Result{
		Group: group.NewGroup(
			chainConfig.GroupSize-chainConfig.HonestThreshold,
			chainConfig.GroupSize,
		),
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
	}

	resultGroupPublicKeyBytes, err := result.GroupPublicKeyBytes()
	if err != nil {
		t.Fatal(err)
	}

	var tests = map[string]struct {
		ctxFn                          func() (context.Context, context.CancelFunc)
		resultSubmittedEvent           *DKGResultSubmittedEvent
		expectedOperatingMemberIndexes []group.MemberIndex
		expectedError                  error
	}{
		"member supports the published result": {
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			resultSubmittedEvent: &DKGResultSubmittedEvent{
				MemberIndex:         2,
				GroupPublicKeyBytes: resultGroupPublicKeyBytes,
				Misbehaved:          []byte{7, 10},
				BlockNumber:         5,
			},
			// should return operating members according to the published result
			expectedOperatingMemberIndexes: []group.MemberIndex{1, 2, 3, 4, 5, 6, 8, 9},
		},
		"member supports different public key": {
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			resultSubmittedEvent: &DKGResultSubmittedEvent{
				MemberIndex:         2,
				GroupPublicKeyBytes: []byte{0x00, 0x01}, // different result
				Misbehaved:          []byte{7, 10},
				BlockNumber:         5,
			},
			expectedError: fmt.Errorf(
				"[member:%v] could not stay in the group because the "+
					"member does not support the same group public key",
				memberIndex,
			),
		},
		"member considered as misbehaved": {
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			resultSubmittedEvent: &DKGResultSubmittedEvent{
				MemberIndex:         2,
				GroupPublicKeyBytes: resultGroupPublicKeyBytes,
				Misbehaved:          []byte{memberIndex}, // member considered as misbehaved
				BlockNumber:         5,
			},
			expectedError: fmt.Errorf(
				"[member:%v] could not stay in the group because the "+
					"member is considered as misbehaving",
				memberIndex,
			),
		},
		"publication timeout exceeded": {
			ctxFn: func() (context.Context, context.CancelFunc) {
				ctx, cancelCtx := context.WithCancel(context.Background())
				// Cancel the context deliberately.
				cancelCtx()
				return ctx, cancelCtx
			},
			resultSubmittedEvent: nil, // the result is not published at all
			expectedError:        fmt.Errorf("result publication timed out"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			dkgResultChannel := make(chan *DKGResultSubmittedEvent, 1)

			if test.resultSubmittedEvent != nil {
				dkgResultChannel <- test.resultSubmittedEvent
			}

			ctx, cancelCtx := test.ctxFn()
			defer cancelCtx()

			operatingMemberIndexes, err := decideSigningGroupMemberFate(
				ctx,
				memberIndex,
				dkgResultChannel,
				result,
			)

			if !reflect.DeepEqual(
				test.expectedOperatingMemberIndexes,
				operatingMemberIndexes,
			) {
				t.Errorf(
					"unexpected operating member indexes\n"+
						"expected: [%v]\n"+
						"actual:   [%v]",
					test.expectedOperatingMemberIndexes,
					operatingMemberIndexes,
				)
			}

			if !reflect.DeepEqual(
				test.expectedError,
				err,
			) {
				t.Errorf(
					"unexpected error\n"+
						"expected: [%v]\n"+
						"actual:   [%v]",
					test.expectedError,
					err,
				)
			}
		})
	}
}

func TestFinalSigningGroup(t *testing.T) {
	chainConfig := &ChainConfig{
		GroupSize:       5,
		GroupQuorum:     3,
		HonestThreshold: 2,
	}

	selectedOperators := []chain.Address{
		"0xAA",
		"0xBB",
		"0xCC",
		"0xDD",
		"0xEE",
	}

	var tests = map[string]struct {
		selectedOperators           []chain.Address
		operatingMembersIndexes     []group.MemberIndex
		expectedFinalOperators      []chain.Address
		expectedFinalMembersIndexes map[group.MemberIndex]group.MemberIndex
		expectedError               error
	}{
		"selected operators count not equal to the group size": {
			selectedOperators:       selectedOperators[:4],
			operatingMembersIndexes: []group.MemberIndex{1, 2, 3, 4, 5},
			expectedError:           fmt.Errorf("invalid input parameters"),
		},
		"all selected operators are operating": {
			selectedOperators:           selectedOperators,
			operatingMembersIndexes:     []group.MemberIndex{5, 4, 3, 2, 1},
			expectedFinalOperators:      selectedOperators,
			expectedFinalMembersIndexes: map[group.MemberIndex]group.MemberIndex{1: 1, 2: 2, 3: 3, 4: 4, 5: 5},
		},
		"honest majority of selected operators are operating": {
			selectedOperators:           selectedOperators,
			operatingMembersIndexes:     []group.MemberIndex{5, 1, 3},
			expectedFinalOperators:      []chain.Address{"0xAA", "0xCC", "0xEE"},
			expectedFinalMembersIndexes: map[group.MemberIndex]group.MemberIndex{1: 1, 3: 2, 5: 3},
		},
		"less than honest majority of selected operators are operating": {
			selectedOperators:       selectedOperators,
			operatingMembersIndexes: []group.MemberIndex{5, 1},
			expectedError:           fmt.Errorf("invalid input parameters"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualFinalOperators, actualFinalMembersIndexes, err :=
				finalSigningGroup(
					test.selectedOperators,
					test.operatingMembersIndexes,
					chainConfig,
				)

			if !reflect.DeepEqual(test.expectedError, err) {
				t.Errorf(
					"unexpected error\n"+
						"expected: %v\n"+
						"actual:   %v\n",
					test.expectedError,
					err,
				)
			}

			if !reflect.DeepEqual(
				test.expectedFinalOperators,
				actualFinalOperators,
			) {
				t.Errorf(
					"unexpected final operators\n"+
						"expected: %v\n"+
						"actual:   %v\n",
					test.expectedFinalOperators,
					actualFinalOperators,
				)
			}

			if !reflect.DeepEqual(
				test.expectedFinalMembersIndexes,
				actualFinalMembersIndexes,
			) {
				t.Errorf(
					"unexpected final members indexes\n"+
						"expected: %v\n"+
						"actual:   %v\n",
					test.expectedFinalMembersIndexes,
					actualFinalMembersIndexes,
				)
			}
		})
	}
}

func TestSignResult_SigningSuccessful(t *testing.T) {
	chain := Connect(5, 4, 3)
	dkgResultSigner := newDkgResultSigner(chain)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}
	result := &dkg.Result{
		Group:           group.NewGroup(32, 64),
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
	}

	signedResult, err := dkgResultSigner.SignResult(result)
	if err != nil {
		t.Fatal(err)
	}

	expectedPublicKey := chain.Signing().PublicKey()
	if !reflect.DeepEqual(
		expectedPublicKey,
		signedResult.PublicKey,
	) {
		t.Errorf(
			"unexpected public key\n"+
				"expected: %v\n"+
				"actual:   %v\n",
			expectedPublicKey,
			signedResult.PublicKey,
		)
	}

	expectedDKGResultHash := dkg.ResultHash(
		sha3.Sum256([]byte(fmt.Sprint(result))),
	)
	if expectedDKGResultHash != signedResult.ResultHash {
		t.Errorf(
			"unexpected result hash\n"+
				"expected: %v\n"+
				"actual:   %v\n",
			expectedDKGResultHash,
			signedResult.ResultHash,
		)
	}

	// Since signature is different on every run (even if the same private key
	// and result hash are used), simply verify if it's correct
	signatureVerification, err := chain.Signing().Verify(
		signedResult.ResultHash[:],
		signedResult.Signature,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !signatureVerification {
		t.Errorf(
			"Signature [0x%x] was not generated properly for the result hash "+
				"[0x%x]",
			signedResult.Signature,
			signedResult.ResultHash,
		)
	}
}

func TestSignResult_ErrorDuringDkgResultHashCalculation(t *testing.T) {
	chain := Connect(5, 4, 3)
	dkgResultSigner := newDkgResultSigner(chain)

	// Use nil as the DKG result to cause hash calculation error
	_, err := dkgResultSigner.SignResult(nil)

	expectedError := fmt.Errorf(
		"dkg result hash calculation failed [%w]",
		errNilDKGResult,
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestVerifySignature_VerificationSuccessful(t *testing.T) {
	chain := Connect(5, 4, 3)
	dkgResultSigner := newDkgResultSigner(chain)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}
	result := &dkg.Result{
		Group:           group.NewGroup(32, 64),
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
	}

	signedResult, err := dkgResultSigner.SignResult(result)
	if err != nil {
		t.Fatal(err)
	}

	verificationSuccessful, err := dkgResultSigner.VerifySignature(signedResult)
	if err != nil {
		t.Fatal(err)
	}

	if !verificationSuccessful {
		t.Fatal(
			"Expected successful verification of signature, but it was " +
				"unsuccessful",
		)
	}
}

func TestVerifySignature_VerificationFailure(t *testing.T) {
	chain := Connect(5, 4, 3)
	dkgResultSigner := newDkgResultSigner(chain)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	result := &dkg.Result{
		Group:           group.NewGroup(32, 64),
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
	}
	signedResult, err := dkgResultSigner.SignResult(result)
	if err != nil {
		t.Fatal(err)
	}

	anotherResult := &dkg.Result{
		Group:           group.NewGroup(30, 64),
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
	}
	anotherSignedResult, err := dkgResultSigner.SignResult(anotherResult)
	if err != nil {
		t.Fatal(err)
	}

	// Assign signature from another result to cause a signature verification
	// failure
	signedResult.Signature = anotherSignedResult.Signature

	verificationSuccessful, err := dkgResultSigner.VerifySignature(signedResult)
	if err != nil {
		t.Fatal(err)
	}

	if verificationSuccessful {
		t.Fatal(
			"Expected unsuccessful verification of signature, but it was " +
				"successful",
		)
	}
}

func TestVerifySignature_VerificationError(t *testing.T) {
	chain := Connect(5, 4, 3)
	dkgResultSigner := newDkgResultSigner(chain)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	result := &dkg.Result{
		Group:           group.NewGroup(32, 64),
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
	}
	signedResult, err := dkgResultSigner.SignResult(result)
	if err != nil {
		t.Fatal(err)
	}

	// Drop the last byte of the signature to cause an error during signature
	// verification
	signedResult.Signature = signedResult.Signature[:len(signedResult.Signature)-1]

	_, err = dkgResultSigner.VerifySignature(signedResult)

	expectedError := fmt.Errorf(
		"failed to unmarshal signature: [asn1: syntax error: data truncated]",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\n"+
				"expected: [%+v]\n"+
				"actual:   [%+v]",
			expectedError,
			err,
		)
	}
}

func TestSubmitResult_MemberSubmitsResult(t *testing.T) {
	chain := Connect(5, 4, 3)
	dkgResultSubmitter := newDkgResultSubmitter(&testutils.MockLogger{}, chain)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}
	result := &dkg.Result{
		Group:           group.NewGroup(32, 64),
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
	}

	memberIndex := group.MemberIndex(1)
	signatures := map[group.MemberIndex][]byte{
		1: []byte("signature 1"),
		2: []byte("signature 2"),
		3: []byte("signature 3"),
		4: []byte("signature 4"),
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	err = dkgResultSubmitter.SubmitResult(
		ctx,
		memberIndex,
		result,
		signatures,
	)
	if err != nil {
		t.Fatal(err)
	}

	expectedActiveWallet, err := result.GroupPublicKeyBytes()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedActiveWallet, chain.activeWallet) {
		t.Errorf(
			"unexpected active wallet bytes \nexpected: [0x%x]\nactual:   [0x%x]\n",
			expectedActiveWallet,
			chain.activeWallet,
		)
	}
}

func TestSubmitResult_MemberDoesNotSubmitsResult(t *testing.T) {
	chain := Connect(5, 4, 3)
	dkgResultSubmitter := newDkgResultSubmitter(&testutils.MockLogger{}, chain)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}
	result := &dkg.Result{
		Group:           group.NewGroup(32, 64),
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
	}
	signatures := map[group.MemberIndex][]byte{
		1: []byte("signature 1"),
		2: []byte("signature 2"),
		3: []byte("signature 3"),
		4: []byte("signature 4"),
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	secondMemberSubmissionChannel := make(chan error)

	// Attempt to submit result for the second member on a separate goroutine.
	go func() {
		secondMemberIndex := group.MemberIndex(2)
		secondMemberErr := dkgResultSubmitter.SubmitResult(
			ctx,
			secondMemberIndex,
			result,
			signatures,
		)
		secondMemberSubmissionChannel <- secondMemberErr
	}()

	// This sleep is needed to give enough time for the second member to
	// register their result submission event handler and act properly on
	// the result submitted by the first member.
	time.Sleep(1 * time.Second)

	// While the second member is waiting for submission eligibility, submit the
	// result with the first member.
	firstMemberIndex := group.MemberIndex(1)
	firstMemberErr := dkgResultSubmitter.SubmitResult(
		ctx,
		firstMemberIndex,
		result,
		signatures,
	)
	if firstMemberErr != nil {
		t.Fatal(firstMemberErr)
	}

	// Check that the second member returned without errors
	secondMemberErr := <-secondMemberSubmissionChannel
	if err != nil {
		t.Fatal(secondMemberErr)
	}

	if chain.resultSubmitterIndex != firstMemberIndex {
		t.Errorf(
			"unexpected result submitter index \nexpected: %v\nactual:   %v\n",
			firstMemberIndex,
			chain.resultSubmitterIndex,
		)
	}

	expectedActiveWallet, err := result.GroupPublicKeyBytes()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedActiveWallet, chain.activeWallet) {
		t.Errorf(
			"unexpected active wallet bytes \nexpected: [0x%x]\nactual:   [0x%x]\n",
			expectedActiveWallet,
			chain.activeWallet,
		)
	}
}

func TestSubmitResult_TooFewSignatures(t *testing.T) {
	chain := Connect(5, 4, 3)
	dkgResultSubmitter := newDkgResultSubmitter(&testutils.MockLogger{}, chain)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}
	result := &dkg.Result{
		Group:           group.NewGroup(32, 64),
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
	}

	memberIndex := group.MemberIndex(1)
	signatures := map[group.MemberIndex][]byte{
		1: []byte("signature 1"),
		2: []byte("signature 2"),
		3: []byte("signature 3"),
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	err = dkgResultSubmitter.SubmitResult(
		ctx,
		memberIndex,
		result,
		signatures,
	)

	expectedError := fmt.Errorf(
		"could not submit result with [3] signatures for group quorum [4]",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\n"+
				"expected: [%+v]\n"+
				"actual:   [%+v]",
			expectedError,
			err,
		)
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
