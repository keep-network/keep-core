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

	// sha256(sha256(messageToSign))
	sha256d, err := hex.DecodeString("38d30dacec5083c902952ce99fc0287659ad0b1ca2086827a8e78b0bef2c8bc1")
	if err != nil {
		t.Fatal(err)
	}

	hostChain := Connect()
	hostChain.setHeartbeatProposalValidationResult(proposal, true)

	mockExecutor := &mockHeartbeatSigningExecutor{}
	action := newHeartbeatAction(
		logger,
		hostChain,
		wallet{
			publicKey: unmarshalPublicKey(walletPublicKeyHex),
		},
		mockExecutor,
		proposal,
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

	hostChain := Connect()
	hostChain.setHeartbeatProposalValidationResult(proposal, true)

	mockExecutor := &mockHeartbeatSigningExecutor{}
	mockExecutor.shouldFail = true

	action := newHeartbeatAction(
		logger,
		hostChain,
		wallet{
			publicKey: unmarshalPublicKey(walletPublicKeyHex),
		},
		mockExecutor,
		proposal,
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
