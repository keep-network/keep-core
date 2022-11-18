package tbtc

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/net/local"
	"reflect"
	"strings"
	"testing"
)

func TestNode_GetSigningExecutor(t *testing.T) {
	chain := Connect(10, 8, 6)
	provider := local.Connect()

	signer := sampleSigner(t)
	signerBytes, err := signer.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// Populate the mock keystore with the sample signer's data. This is
	// required to make the node controlling the signer's wallet.
	keyStorePersistance := &mockPersistenceHandle{
		saved: []persistence.DataDescriptor{
			&mockDescriptor{
				name:      "membership_1",
				directory: "wallet_1",
				content:   signerBytes,
			},
		},
	}

	node := newNode(
		chain,
		provider,
		keyStorePersistance,
		&mockPersistenceHandle{},
		generator.StartScheduler(),
		Config{},
	)

	walletPublicKey := signer.wallet.publicKey
	walletPublicKeyBytes, err := marshalPublicKey(walletPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertIntsEqual(
		t,
		"cache size",
		0,
		len(node.signingExecutors),
	)

	executor, err := node.getSigningExecutor(walletPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertIntsEqual(
		t,
		"cache size",
		1,
		len(node.signingExecutors),
	)

	testutils.AssertIntsEqual(
		t,
		"signers count",
		1,
		len(executor.signers),
	)

	if !reflect.DeepEqual(signer, executor.signers[0]) {
		t.Errorf("executor holds an unexpected signer")
	}

	expectedChannel := fmt.Sprintf(
		"%s-%s",
		ProtocolName,
		hex.EncodeToString(walletPublicKeyBytes),
	)
	testutils.AssertStringsEqual(
		t,
		"broadcast channel",
		expectedChannel,
		executor.broadcastChannel.Name(),
	)

	_, err = node.getSigningExecutor(walletPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	// The executor was already created in the previous call so cached instance
	// should be returned and no new executors should be created.
	testutils.AssertIntsEqual(
		t,
		"cache size",
		1,
		len(node.signingExecutors),
	)

	// Construct an arbitrary public key representing a wallet that is not
	// controlled by the node. We need to make sure the public key's points
	// are on the curve to avoid troubles during processing.
	x, y := walletPublicKey.Curve.Double(walletPublicKey.X, walletPublicKey.Y)
	nonControlledWalletPublicKey := &ecdsa.PublicKey{
		Curve: walletPublicKey.Curve,
		X:     x,
		Y:     y,
	}

	_, err = node.getSigningExecutor(nonControlledWalletPublicKey)
	if !strings.Contains(
		err.Error(),
		"node does not control signers of wallet with public key",
	) {
		t.Errorf("unexpected error")
	}
}
