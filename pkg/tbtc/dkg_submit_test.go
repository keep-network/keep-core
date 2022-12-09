package tbtc

import (
	"context"
	"fmt"
	"github.com/keep-network/keep-core/pkg/chain"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
	"golang.org/x/crypto/sha3"
)

func TestSignResult_SigningSuccessful(t *testing.T) {
	chain := Connect(5, 4, 3)
	dkgStartBlock := uint64(2000)
	dkgResultSigner := newDkgResultSigner(chain, dkgStartBlock)

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
	dkgStartBlock := uint64(2000)
	dkgResultSigner := newDkgResultSigner(chain, dkgStartBlock)

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
	dkgStartBlock := uint64(2000)
	dkgResultSigner := newDkgResultSigner(chain, dkgStartBlock)

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
	dkgStartBlock := uint64(2000)
	dkgResultSigner := newDkgResultSigner(chain, dkgStartBlock)

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
	dkgStartBlock := uint64(2000)
	dkgResultSigner := newDkgResultSigner(chain, dkgStartBlock)

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
	const (
		groupSize       = 5
		groupQuorum     = 4
		honestThreshold = 3
	)

	localChain := Connect(groupSize, groupQuorum, honestThreshold)

	err := localChain.startDKG()
	if err != nil {
		t.Fatal(err)
	}

	operatorID, operatorAddress, err := localChain.operator()
	if err != nil {
		t.Fatal(err)
	}

	var operatorsIDs chain.OperatorIDs
	var operatorsAddresses chain.Addresses

	for memberIndex := uint8(1); memberIndex <= groupSize; memberIndex++ {
		operatorsIDs = append(operatorsIDs, operatorID)
		operatorsAddresses = append(operatorsAddresses, operatorAddress)
	}

	groupSelectionResult := &GroupSelectionResult{
		OperatorsIDs:       operatorsIDs,
		OperatorsAddresses: operatorsAddresses,
	}

	dkgResultSubmitter := newDkgResultSubmitter(
		&testutils.MockLogger{},
		localChain,
		groupSelectionResult,
	)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}
	result := &dkg.Result{
		Group:           group.NewGroup(groupSize-honestThreshold, groupSize),
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

	expectedGroupPublicKey, err := result.GroupPublicKeyBytes()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedGroupPublicKey, localChain.dkgResult.GroupPublicKey) {
		t.Errorf(
			"unexpected group public key \nexpected: [0x%x]\nactual:   [0x%x]\n",
			expectedGroupPublicKey,
			localChain.dkgResult.GroupPublicKey,
		)
	}
}

func TestSubmitResult_MemberDoesNotSubmitResult(t *testing.T) {
	const (
		groupSize       = 5
		groupQuorum     = 4
		honestThreshold = 3
	)

	localChain := Connect(groupSize, groupQuorum, honestThreshold)

	err := localChain.startDKG()
	if err != nil {
		t.Fatal(err)
	}

	operatorID, operatorAddress, err := localChain.operator()
	if err != nil {
		t.Fatal(err)
	}

	var operatorsIDs chain.OperatorIDs
	var operatorsAddresses chain.Addresses

	for memberIndex := uint8(1); memberIndex <= groupSize; memberIndex++ {
		operatorsIDs = append(operatorsIDs, operatorID)
		operatorsAddresses = append(operatorsAddresses, operatorAddress)
	}

	groupSelectionResult := &GroupSelectionResult{
		OperatorsIDs:       operatorsIDs,
		OperatorsAddresses: operatorsAddresses,
	}

	dkgResultSubmitter := newDkgResultSubmitter(
		&testutils.MockLogger{},
		localChain,
		groupSelectionResult,
	)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}
	result := &dkg.Result{
		Group:           group.NewGroup(groupSize-honestThreshold, groupSize),
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

	if localChain.dkgResult.SubmitterMemberIndex != firstMemberIndex {
		t.Errorf(
			"unexpected result submitter index \nexpected: %v\nactual:   %v\n",
			firstMemberIndex,
			localChain.dkgResult.SubmitterMemberIndex,
		)
	}

	expectedGroupPublicKey, err := result.GroupPublicKeyBytes()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedGroupPublicKey, localChain.dkgResult.GroupPublicKey) {
		t.Errorf(
			"unexpected group public key \nexpected: [0x%x]\nactual:   [0x%x]\n",
			expectedGroupPublicKey,
			localChain.dkgResult.GroupPublicKey,
		)
	}
}

func TestSubmitResult_TooFewSignatures(t *testing.T) {
	const (
		groupSize       = 5
		groupQuorum     = 4
		honestThreshold = 3
	)

	localChain := Connect(groupSize, groupQuorum, honestThreshold)

	err := localChain.startDKG()
	if err != nil {
		t.Fatal(err)
	}

	operatorID, operatorAddress, err := localChain.operator()
	if err != nil {
		t.Fatal(err)
	}

	var operatorsIDs chain.OperatorIDs
	var operatorsAddresses chain.Addresses

	for memberIndex := uint8(1); memberIndex <= groupSize; memberIndex++ {
		operatorsIDs = append(operatorsIDs, operatorID)
		operatorsAddresses = append(operatorsAddresses, operatorAddress)
	}

	groupSelectionResult := &GroupSelectionResult{
		OperatorsIDs:       operatorsIDs,
		OperatorsAddresses: operatorsAddresses,
	}

	dkgResultSubmitter := newDkgResultSubmitter(
		&testutils.MockLogger{},
		localChain,
		groupSelectionResult,
	)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}
	result := &dkg.Result{
		Group:           group.NewGroup(groupSize-honestThreshold, groupSize),
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
