package tbtc

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"math/big"
	"testing"
)

func TestHeartbeatAction_HappyPath(t *testing.T) {
	startBlock := uint64(10)
	expiryBlock := startBlock + heartbeatProposalValidityBlocks
	messageToSign, err := hex.DecodeString("FFFFFFFFFFFFFFFF0000000000000001")
	if err != nil {
		t.Fatal(err)
	}
	// sha256(sha256(messageToSign))
	sha256d, err := hex.DecodeString("38d30dacec5083c902952ce99fc0287659ad0b1ca2086827a8e78b0bef2c8bc1")
	if err != nil {
		t.Fatal(err)
	}

	mockExecutor := &mockHeartbeatSigningExecutor{}
	action := newHeartbeatAction(
		logger,
		wallet{},
		mockExecutor,
		messageToSign,
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
	startBlock := uint64(10)
	expiryBlock := startBlock + heartbeatProposalValidityBlocks
	messageToSign, err := hex.DecodeString("FFFFFFFFFFFFFFFF0000000000000001")
	if err != nil {
		t.Fatal(err)
	}

	mockExecutor := &mockHeartbeatSigningExecutor{}
	mockExecutor.shouldFail = true

	action := newHeartbeatAction(
		logger,
		wallet{},
		mockExecutor,
		messageToSign,
		startBlock,
		expiryBlock,
		func(ctx context.Context, blockHeight uint64) error {
			return nil
		},
	)

	err = action.execute()
	if err == nil {
		t.Fatal("expected error to be returned")
	}
	testutils.AssertStringsEqual(
		t,
		"error message",
		"cannot sign heartbeat message: [oofta]",
		err.Error(),
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
) (*tecdsa.Signature, uint64, error) {
	mhse.requestedMessage = message
	mhse.requestedStartBlock = startBlock

	if mhse.shouldFail {
		return nil, 0, fmt.Errorf("oofta")
	}

	return &tecdsa.Signature{}, startBlock + 1, nil
}
