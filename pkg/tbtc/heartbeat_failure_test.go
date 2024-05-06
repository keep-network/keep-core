package tbtc

import (
	"crypto/ecdsa"
	"encoding/hex"
	"testing"

	"github.com/keep-network/keep-core/internal/testutils"
)

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
