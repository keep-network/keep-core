package tbtc

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/tecdsa"
)

func TestHeartbeatAction_HappyPath(t *testing.T) {
	walletPublicKeyHex, err := hex.DecodeString(
		"0471e30bca60f6548d7b42582a478ea37ada63b402af7b3ddd57f0c95bb6843175" +
			"aa0d2053a91a050a6797d85c38f2909cb7027f2344a01986aa2f9f8ca7a0c289",
	)
	if err != nil {
		t.Fatal(err)
	}

	startBlock := uint64(10)
	expiryBlock := startBlock + heartbeatProposalValidityBlocks

	proposal := &HeartbeatProposal{
		Message: [16]byte{
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		},
	}

	heartbeatFailureCounter := newHeartbeatFailureCounter()

	// sha256(sha256(messageToSign))
	sha256d, err := hex.DecodeString("38d30dacec5083c902952ce99fc0287659ad0b1ca2086827a8e78b0bef2c8bc1")
	if err != nil {
		t.Fatal(err)
	}

	hostChain := Connect()
	hostChain.setHeartbeatProposalValidationResult(proposal, true)

	mockExecutor := &mockHeartbeatSigningExecutor{}
	inactivityClaimExecutor := &inactivityClaimExecutor{}
	action := newHeartbeatAction(
		logger,
		hostChain,
		wallet{
			publicKey: unmarshalPublicKey(walletPublicKeyHex),
		},
		mockExecutor,
		proposal,
		heartbeatFailureCounter,
		inactivityClaimExecutor,
		startBlock,
		expiryBlock,
		func(ctx context.Context, blockHeight uint64) error {
			return nil
		},
	)

	err = action.execute()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertBigIntsEqual(
		t,
		"message to sign",
		new(big.Int).SetBytes(sha256d),
		mockExecutor.requestedMessage,
	)
	testutils.AssertUintsEqual(
		t,
		"start block",
		startBlock,
		mockExecutor.requestedStartBlock,
	)
}

func TestHeartbeatAction_SigningError(t *testing.T) {
	walletPublicKeyHex, err := hex.DecodeString(
		"0471e30bca60f6548d7b42582a478ea37ada63b402af7b3ddd57f0c95bb6843175" +
			"aa0d2053a91a050a6797d85c38f2909cb7027f2344a01986aa2f9f8ca7a0c289",
	)
	if err != nil {
		t.Fatal(err)
	}

	startBlock := uint64(10)
	expiryBlock := startBlock + heartbeatProposalValidityBlocks

	proposal := &HeartbeatProposal{
		Message: [16]byte{
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		},
	}

	heartbeatFailureCounter := newHeartbeatFailureCounter()

	hostChain := Connect()
	hostChain.setHeartbeatProposalValidationResult(proposal, true)

	mockExecutor := &mockHeartbeatSigningExecutor{}
	mockExecutor.shouldFail = true

	inactivityClaimExecutor := &inactivityClaimExecutor{}

	action := newHeartbeatAction(
		logger,
		hostChain,
		wallet{
			publicKey: unmarshalPublicKey(walletPublicKeyHex),
		},
		mockExecutor,
		proposal,
		heartbeatFailureCounter,
		inactivityClaimExecutor,
		startBlock,
		expiryBlock,
		func(ctx context.Context, blockHeight uint64) error {
			return nil
		},
	)

	action.execute()
	// TODO: Uncomment
	// err = action.execute()
	// if err == nil {
	// 	t.Fatal("expected error to be returned")
	// }
	// testutils.AssertStringsEqual(
	// 	t,
	// 	"error message",
	// 	"cannot sign heartbeat message: [oofta]",
	// 	err.Error(),
	// )
}

func TestHeartbeatFailureCounter_Increment(t *testing.T) {
	walletPublicKey := createMockSigner(t).wallet.publicKey
	walletPublicKeyBytes, err := marshalPublicKey(walletPublicKey)
	if err != nil {
		t.Fatal(t)
	}

	heartbeatFailureCounter := newHeartbeatFailureCounter()

	counterKey := hex.EncodeToString(walletPublicKeyBytes)

	// Check first increment.
	heartbeatFailureCounter.increment(counterKey)
	count := heartbeatFailureCounter.get(counterKey)
	testutils.AssertUintsEqual(
		t,
		"counter value",
		1,
		uint64(count),
	)

	// Check second increment.
	heartbeatFailureCounter.increment(counterKey)
	count = heartbeatFailureCounter.get(counterKey)
	testutils.AssertUintsEqual(
		t,
		"counter value",
		2,
		uint64(count),
	)
}

func TestHeartbeatFailureCounter_Reset(t *testing.T) {
	walletPublicKey := createMockSigner(t).wallet.publicKey
	walletPublicKeyBytes, err := marshalPublicKey(walletPublicKey)
	if err != nil {
		t.Fatal(t)
	}

	heartbeatFailureCounter := newHeartbeatFailureCounter()

	counterKey := hex.EncodeToString(walletPublicKeyBytes)

	// Check reset works as the first operation.
	heartbeatFailureCounter.reset(counterKey)
	count := heartbeatFailureCounter.get(counterKey)
	testutils.AssertUintsEqual(
		t,
		"counter value",
		0,
		uint64(count),
	)

	// Check reset works after an increment.
	heartbeatFailureCounter.increment(counterKey)
	heartbeatFailureCounter.reset(counterKey)

	count = heartbeatFailureCounter.get(counterKey)
	testutils.AssertUintsEqual(
		t,
		"counter value",
		0,
		uint64(count),
	)
}

func TestHeartbeatFailureCounter_Get(t *testing.T) {
	walletPublicKey := createMockSigner(t).wallet.publicKey
	walletPublicKeyBytes, err := marshalPublicKey(walletPublicKey)
	if err != nil {
		t.Fatal(t)
	}

	heartbeatFailureCounter := newHeartbeatFailureCounter()

	counterKey := hex.EncodeToString(walletPublicKeyBytes)

	// Check get works as the first operation.
	count := heartbeatFailureCounter.get(counterKey)
	testutils.AssertUintsEqual(
		t,
		"counter value",
		0,
		uint64(count),
	)

	// Check get works after an increment.
	heartbeatFailureCounter.increment(counterKey)
	count = heartbeatFailureCounter.get(counterKey)
	testutils.AssertUintsEqual(
		t,
		"counter value",
		1,
		uint64(count),
	)

	// Construct an arbitrary public key representing a different wallet.
	x, y := walletPublicKey.Curve.Double(walletPublicKey.X, walletPublicKey.Y)
	anotherWalletPublicKey := &ecdsa.PublicKey{
		Curve: walletPublicKey.Curve,
		X:     x,
		Y:     y,
	}
	anotherWalletPublicKeyBytes, err := marshalPublicKey(anotherWalletPublicKey)
	if err != nil {
		t.Fatal(t)
	}
	anotherCounterKey := hex.EncodeToString(anotherWalletPublicKeyBytes)

	// Check get works on another wallet.
	count = heartbeatFailureCounter.get(anotherCounterKey)
	testutils.AssertUintsEqual(
		t,
		"counter value",
		0,
		uint64(count),
	)
}

type mockHeartbeatSigningExecutor struct {
	shouldFail bool

	requestedMessage    *big.Int
	requestedStartBlock uint64
}

func (mhse *mockHeartbeatSigningExecutor) sign(
	ctx context.Context,
	message *big.Int,
	startBlock uint64,
) (*tecdsa.Signature, uint32, uint64, error) {
	mhse.requestedMessage = message
	mhse.requestedStartBlock = startBlock

	if mhse.shouldFail {
		return nil, 0, 0, fmt.Errorf("oofta")
	}

	// TODO: Return the active members count and use it in unit tests.
	return &tecdsa.Signature{}, 0, startBlock + 1, nil
}
