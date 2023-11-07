package tbtc

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
)

func TestNode_GetSigningExecutor(t *testing.T) {
	groupParameters := &GroupParameters{
		GroupSize:       5,
		GroupQuorum:     4,
		HonestThreshold: 3,
	}

	localChain := Connect()
	localProvider := local.Connect()

	signer := createMockSigner(t)

	// Populate the mock keystore with the mock signer's data. This is
	// required to make the node controlling the signer's wallet.
	keyStorePersistence := createMockKeyStorePersistence(t, signer)

	node, err := newNode(
		groupParameters,
		localChain,
		newLocalBitcoinChain(),
		localProvider,
		keyStorePersistence,
		&mockPersistenceHandle{},
		generator.StartScheduler(),
		Config{},
	)
	if err != nil {
		t.Fatal(err)
	}

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

	executor, ok, err := node.getSigningExecutor(walletPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("node is supposed to control wallet signers")
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

	_, ok, err = node.getSigningExecutor(walletPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("node is supposed to control wallet signers")
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

	_, ok, err = node.getSigningExecutor(nonControlledWalletPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("node is not supposed to control wallet signers")
	}
}

// createMockSigner creates a mock signer instance that can be used for
// test cases that needs a placeholder signer. The produced signer cannot
// be used to test actual signing scenarios.
func createMockSigner(t *testing.T) *signer {
	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	privateKeyShare := tecdsa.NewPrivateKeyShare(testData[0])

	signingGroupOperators := []chain.Address{
		"address-1",
		"address-2",
		"address-3",
		"address-3",
		"address-5",
	}

	return &signer{
		wallet: wallet{
			publicKey:             privateKeyShare.PublicKey(),
			signingGroupOperators: signingGroupOperators,
		},
		signingGroupMemberIndex: group.MemberIndex(1),
		privateKeyShare:         privateKeyShare,
	}
}

// createMockKeyStorePersistence creates a mock key store that can be used
// to create test node instances. The key store is populated with the given
// signers.
func createMockKeyStorePersistence(
	t *testing.T,
	signers ...*signer,
) *mockPersistenceHandle {
	walletToSigners := make(map[string][]*signer)
	for _, signer := range signers {
		keyBytes, err := marshalPublicKey(signer.wallet.publicKey)
		if err != nil {
			t.Fatal(err)
		}

		key := hex.EncodeToString(keyBytes)

		walletToSigners[key] = append(walletToSigners[key], signer)
	}

	descriptors := make([]persistence.DataDescriptor, 0)

	for key, signers := range walletToSigners {
		for i, signer := range signers {
			signerBytes, err := signer.Marshal()
			if err != nil {
				t.Fatal(err)
			}

			// Construct the descriptor in the same way as it happens in the
			// real world.
			descriptor := &mockDescriptor{
				name:      fmt.Sprintf("membership_%v", i+1),
				directory: key[2:], // trim the 04 prefix
				content:   signerBytes,
			}

			descriptors = append(descriptors, descriptor)
		}
	}

	return &mockPersistenceHandle{
		saved: descriptors,
	}
}
