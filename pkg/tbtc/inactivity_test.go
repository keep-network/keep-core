package tbtc

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/inactivity"
	"github.com/keep-network/keep-core/pkg/tecdsa"
)

func TestInactivityClaimExecutor_ClaimInactivity(t *testing.T) {
	executor, walletEcdsaID, chain := setupInactivityClaimExecutorScenario(t)

	initialNonce, err := chain.GetInactivityClaimNonce(walletEcdsaID)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	message := big.NewInt(100)
	inactiveMembersIndexes := []group.MemberIndex{1, 4}

	err = executor.claimInactivity(
		ctx,
		inactiveMembersIndexes,
		true,
		message,
	)
	if err != nil {
		t.Fatal(err)
	}

	currentNonce, err := chain.GetInactivityClaimNonce(walletEcdsaID)
	if err != nil {
		t.Fatal(err)
	}

	expectedNonceDiff := uint64(1)
	nonceDiff := currentNonce.Uint64() - initialNonce.Uint64()

	testutils.AssertUintsEqual(
		t,
		"inactivity nonce difference",
		expectedNonceDiff,
		nonceDiff,
	)
}

func TestInactivityClaimExecutor_ClaimInactivity_Busy(t *testing.T) {
	executor, _, _ := setupInactivityClaimExecutorScenario(t)

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	message := big.NewInt(100)
	inactiveMembersIndexes := []group.MemberIndex{1, 4}

	errChan := make(chan error, 1)
	go func() {
		err := executor.claimInactivity(
			ctx,
			inactiveMembersIndexes,
			true,
			message,
		)
		errChan <- err
	}()

	time.Sleep(100 * time.Millisecond)

	err := executor.claimInactivity(
		ctx,
		inactiveMembersIndexes,
		true,
		message,
	)
	testutils.AssertErrorsSame(t, errInactivityClaimExecutorBusy, err)

	err = <-errChan
	if err != nil {
		t.Errorf("unexpected error: [%v]", err)
	}
}

func setupInactivityClaimExecutorScenario(t *testing.T) (
	*inactivityClaimExecutor,
	[32]byte,
	*localChain,
) {
	groupParameters := &GroupParameters{
		GroupSize:       5,
		GroupQuorum:     4,
		HonestThreshold: 3,
	}

	operatorPrivateKey, operatorPublicKey, err := operator.GenerateKeyPair(
		local_v1.DefaultCurve,
	)
	if err != nil {
		t.Fatal(err)
	}

	localChain := ConnectWithKey(operatorPrivateKey)

	localProvider := local.ConnectWithKey(operatorPublicKey)

	operatorAddress, err := localChain.Signing().PublicKeyToAddress(
		operatorPublicKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	var operators []chain.Address
	for i := 0; i < groupParameters.GroupSize; i++ {
		operators = append(operators, operatorAddress)
	}

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(
		groupParameters.GroupSize,
	)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	signers := make([]*signer, len(testData))
	for i := range testData {
		privateKeyShare := tecdsa.NewPrivateKeyShare(testData[i])

		signers[i] = &signer{
			wallet: wallet{
				publicKey:             privateKeyShare.PublicKey(),
				signingGroupOperators: operators,
			},
			signingGroupMemberIndex: group.MemberIndex(i + 1),
			privateKeyShare:         privateKeyShare,
		}
	}

	keyStorePersistence := createMockKeyStorePersistence(t, signers...)

	walletPublicKeyHash := bitcoin.PublicKeyHash(signers[0].wallet.publicKey)
	ecdsaWalletID := [32]byte{1, 2, 3}

	localChain.setWallet(
		walletPublicKeyHash,
		&WalletChainData{
			EcdsaWalletID: ecdsaWalletID,
		},
	)

	node, err := newNode(
		groupParameters,
		localChain,
		newLocalBitcoinChain(),
		localProvider,
		keyStorePersistence,
		&mockPersistenceHandle{},
		generator.StartScheduler(),
		&mockCoordinationProposalGenerator{},
		Config{},
	)
	if err != nil {
		t.Fatal(err)
	}

	executor, ok, err := node.getInactivityClaimExecutor(
		signers[0].wallet.publicKey,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("node is supposed to control wallet signers")
	}

	return executor, ecdsaWalletID, localChain
}

func TestSignClaim_SigningSuccessful(t *testing.T) {
	chain := Connect()
	inactivityClaimSigner := newInactivityClaimSigner(chain)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}
	privateKeyShare := tecdsa.NewPrivateKeyShare(testData[0])

	claim := &inactivity.ClaimPreimage{
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

	expectedInactivityClaimHash := inactivity.ClaimHash(
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

	claim := &inactivity.ClaimPreimage{
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

	claim := &inactivity.ClaimPreimage{
		Nonce:                  big.NewInt(5),
		WalletPublicKey:        privateKeyShare.PublicKey(),
		InactiveMembersIndexes: []group.MemberIndex{11, 22, 33},
		HeartbeatFailed:        true,
	}

	signedClaim, err := inactivityClaimSigner.SignClaim(claim)
	if err != nil {
		t.Fatal(err)
	}

	anotherClaim := &inactivity.ClaimPreimage{
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

	claim := &inactivity.ClaimPreimage{
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
