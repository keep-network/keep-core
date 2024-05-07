package tbtc

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"golang.org/x/crypto/sha3"

	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/inactivity"
	"github.com/keep-network/keep-core/pkg/tecdsa"
)

func TestSignClaim_SigningSuccessful(t *testing.T) {
	chain := Connect()
	inactivityClaimSigner := newInactivityClaimSigner(chain)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}
	privateKeyShare := tecdsa.NewPrivateKeyShare(testData[0])

	claim := &inactivity.Claim{
		Nonce:                  big.NewInt(5),
		WalletPublicKey:        privateKeyShare.PublicKey(),
		InactiveMembersIndexes: []group.MemberIndex{11, 22, 33},
		HeartbeatFailed:        true,
	}

	signedClaim, err := inactivityClaimSigner.SignClaim(claim)
	if err != nil {
		t.Fatal(err)
	}

	expectedPublicKey := chain.Signing().PublicKey()
	if !reflect.DeepEqual(
		expectedPublicKey,
		signedClaim.PublicKey,
	) {
		t.Errorf(
			"unexpected public key\n"+
				"expected: %v\n"+
				"actual:   %v\n",
			expectedPublicKey,
			signedClaim.PublicKey,
		)
	}

	expectedInactivityClaimHash := inactivity.ClaimSignatureHash(
		sha3.Sum256(
			[]byte(fmt.Sprint(
				claim.Nonce,
				claim.WalletPublicKey,
				claim.InactiveMembersIndexes,
				claim.HeartbeatFailed,
			)),
		),
	)
	if expectedInactivityClaimHash != signedClaim.ClaimHash {
		t.Errorf(
			"unexpected claim hash\n"+
				"expected: %v\n"+
				"actual:   %v\n",
			expectedInactivityClaimHash,
			signedClaim.ClaimHash,
		)
	}

	// Since signature is different on every run (even if the same private key
	// and claim hash are used), simply verify if it's correct
	signatureVerification, err := chain.Signing().Verify(
		signedClaim.ClaimHash[:],
		signedClaim.Signature,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !signatureVerification {
		t.Errorf(
			"Signature [0x%x] was not generated properly for the claim hash "+
				"[0x%x]",
			signedClaim.Signature,
			signedClaim.ClaimHash,
		)
	}
}

func TestSignClaim_ErrorDuringInactivityClaimHashCalculation(t *testing.T) {
	chain := Connect()
	inactivityClaimSigner := newInactivityClaimSigner(chain)

	// Use nil as the claim to cause hash calculation error.
	_, err := inactivityClaimSigner.SignClaim(nil)

	expectedError := fmt.Errorf("claim is nil")
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestVerifySignature_VerifySuccessful(t *testing.T) {
	chain := Connect()
	inactivityClaimSigner := newInactivityClaimSigner(chain)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}
	privateKeyShare := tecdsa.NewPrivateKeyShare(testData[0])

	claim := &inactivity.Claim{
		Nonce:                  big.NewInt(5),
		WalletPublicKey:        privateKeyShare.PublicKey(),
		InactiveMembersIndexes: []group.MemberIndex{11, 22, 33},
		HeartbeatFailed:        true,
	}

	signedClaim, err := inactivityClaimSigner.SignClaim(claim)
	if err != nil {
		t.Fatal(err)
	}

	verificationSuccessful, err := inactivityClaimSigner.VerifySignature(
		signedClaim,
	)
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

func TestVerifySignature_VerifyFailure(t *testing.T) {
	chain := Connect()
	inactivityClaimSigner := newInactivityClaimSigner(chain)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}
	privateKeyShare := tecdsa.NewPrivateKeyShare(testData[0])

	claim := &inactivity.Claim{
		Nonce:                  big.NewInt(5),
		WalletPublicKey:        privateKeyShare.PublicKey(),
		InactiveMembersIndexes: []group.MemberIndex{11, 22, 33},
		HeartbeatFailed:        true,
	}

	signedClaim, err := inactivityClaimSigner.SignClaim(claim)
	if err != nil {
		t.Fatal(err)
	}

	anotherClaim := &inactivity.Claim{
		Nonce:                  big.NewInt(6),
		WalletPublicKey:        privateKeyShare.PublicKey(),
		InactiveMembersIndexes: []group.MemberIndex{11, 22, 33},
		HeartbeatFailed:        true,
	}

	anotherSignedClaim, err := inactivityClaimSigner.SignClaim(anotherClaim)
	if err != nil {
		t.Fatal(err)
	}

	// Assign signature from another claim to cause a signature verification
	// failure.
	signedClaim.Signature = anotherSignedClaim.Signature

	verificationSuccessful, err := inactivityClaimSigner.VerifySignature(
		signedClaim,
	)
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

func TestVerifySignature_VerifyError(t *testing.T) {
	chain := Connect()
	inactivityClaimSigner := newInactivityClaimSigner(chain)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}
	privateKeyShare := tecdsa.NewPrivateKeyShare(testData[0])

	claim := &inactivity.Claim{
		Nonce:                  big.NewInt(5),
		WalletPublicKey:        privateKeyShare.PublicKey(),
		InactiveMembersIndexes: []group.MemberIndex{11, 22, 33},
		HeartbeatFailed:        true,
	}

	signedClaim, err := inactivityClaimSigner.SignClaim(claim)
	if err != nil {
		t.Fatal(err)
	}

	// Drop the last byte of the signature to cause an error during signature
	// verification.
	signedClaim.Signature = signedClaim.Signature[:len(signedClaim.Signature)-1]

	_, err = inactivityClaimSigner.VerifySignature(signedClaim)

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

// TODO: Continue with unit tests.
