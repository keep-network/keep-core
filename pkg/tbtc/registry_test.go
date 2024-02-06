package tbtc

import (
	"crypto/ecdsa"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tecdsa"

	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/internal/testutils"
)

func TestWalletRegistry_RegisterSigner(t *testing.T) {
	persistenceHandle := &mockPersistenceHandle{}

	walletRegistry := newWalletRegistry(persistenceHandle)

	signer := createMockSigner(t)

	walletStorageKey := getWalletStorageKey(signer.wallet.publicKey)

	err := walletRegistry.registerSigner(signer)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertIntsEqual(
		t,
		"registered wallets count",
		1,
		len(walletRegistry.walletCache),
	)

	testutils.AssertIntsEqual(
		t,
		"registered wallet signers count",
		1,
		len(walletRegistry.walletCache[walletStorageKey].signers),
	)

	expectedWalletPublicKeyHash := bitcoin.PublicKeyHash(signer.wallet.publicKey)
	testutils.AssertBytesEqual(
		t,
		expectedWalletPublicKeyHash[:],
		walletRegistry.walletCache[walletStorageKey].walletPublicKeyHash[:],
	)

	if !reflect.DeepEqual(signer, walletRegistry.walletCache[walletStorageKey].signers[0]) {
		t.Errorf("registered wallet signer differs from the original one")
	}

	testutils.AssertIntsEqual(
		t,
		"persisted wallet signers count",
		1,
		len(persistenceHandle.saved),
	)
}

func TestWalletRegistry_GetSigners(t *testing.T) {
	persistenceHandle := &mockPersistenceHandle{}

	walletRegistry := newWalletRegistry(persistenceHandle)

	signer := createMockSigner(t)

	err := walletRegistry.registerSigner(signer)
	if err != nil {
		t.Fatal(err)
	}

	fetchedSigners := walletRegistry.getSigners(signer.wallet.publicKey)

	testutils.AssertIntsEqual(
		t,
		"fetched wallet signers count",
		1,
		len(fetchedSigners),
	)

	if !reflect.DeepEqual(signer, fetchedSigners[0]) {
		t.Errorf("fetched wallet signer differs from the original one")
	}
}

func TestWalletRegistry_getWalletByPublicKeyHash(t *testing.T) {
	persistenceHandle := &mockPersistenceHandle{}

	walletRegistry := newWalletRegistry(persistenceHandle)

	signer := createMockSigner(t)

	err := walletRegistry.registerSigner(signer)
	if err != nil {
		t.Fatal(err)
	}

	walletPublicKeyHash := bitcoin.PublicKeyHash(signer.wallet.publicKey)

	wallet, ok := walletRegistry.getWalletByPublicKeyHash(walletPublicKeyHash)
	if !ok {
		t.Error("should return a wallet")
	}

	testutils.AssertStringsEqual(t, "wallet", signer.wallet.String(), wallet.String())
}

func TestWalletRegistry_getWalletByPublicKeyHash_NotFound(t *testing.T) {
	persistenceHandle := &mockPersistenceHandle{}

	walletRegistry := newWalletRegistry(persistenceHandle)

	signer := createMockSigner(t)

	err := walletRegistry.registerSigner(signer)
	if err != nil {
		t.Fatal(err)
	}

	x, y := tecdsa.Curve.ScalarBaseMult(big.NewInt(100).Bytes())

	walletPublicKeyHash := bitcoin.PublicKeyHash(&ecdsa.PublicKey{
		Curve: tecdsa.Curve,
		X:     x,
		Y:     y,
	})

	_, ok := walletRegistry.getWalletByPublicKeyHash(walletPublicKeyHash)
	if ok {
		t.Error("should not return a wallet")
	}
}

func TestWalletRegistry_PrePopulateWalletCache(t *testing.T) {
	signer := createMockSigner(t)
	signerBytes, err := signer.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	walletStorageKey := getWalletStorageKey(signer.wallet.publicKey)

	persistenceHandle := &mockPersistenceHandle{
		saved: []persistence.DataDescriptor{
			&mockDescriptor{
				name:      "membership_1",
				directory: "wallet_1",
				content:   signerBytes,
			},
		},
	}

	// Cache pre-population happens within newWalletRegistry.
	walletRegistry := newWalletRegistry(persistenceHandle)

	testutils.AssertIntsEqual(
		t,
		"loaded wallets count",
		1,
		len(walletRegistry.walletCache),
	)

	expectedWalletPublicKeyHash := bitcoin.PublicKeyHash(signer.wallet.publicKey)
	testutils.AssertBytesEqual(
		t,
		expectedWalletPublicKeyHash[:],
		walletRegistry.walletCache[walletStorageKey].walletPublicKeyHash[:],
	)

	testutils.AssertIntsEqual(
		t,
		"loaded wallet signers count",
		1,
		len(walletRegistry.walletCache[walletStorageKey].signers),
	)

	if !reflect.DeepEqual(signer, walletRegistry.walletCache[walletStorageKey].signers[0]) {
		t.Errorf("loaded wallet signer differs from the original one")
	}
}

func TestWalletStorage_SaveSigner(t *testing.T) {
	persistenceHandle := &mockPersistenceHandle{}

	walletStorage := newWalletStorage(persistenceHandle)

	signer := createMockSigner(t)

	err := walletStorage.saveSigner(signer)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertIntsEqual(
		t,
		"persisted wallet signers count",
		1,
		len(persistenceHandle.saved),
	)
}

func TestWalletStorage_LoadSigners(t *testing.T) {
	signer := createMockSigner(t)
	signerBytes, err := signer.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	walletStorageKey := getWalletStorageKey(signer.wallet.publicKey)

	persistenceHandle := &mockPersistenceHandle{
		saved: []persistence.DataDescriptor{
			&mockDescriptor{
				name:      "membership_1",
				directory: "wallet_1",
				content:   signerBytes,
			},
		},
	}

	walletStorage := newWalletStorage(persistenceHandle)

	signersByWallet := walletStorage.loadSigners()

	testutils.AssertIntsEqual(
		t,
		"loaded wallets count",
		1,
		len(signersByWallet),
	)

	testutils.AssertIntsEqual(
		t,
		"loaded wallet signers count",
		1,
		len(signersByWallet[walletStorageKey]),
	)

	if !reflect.DeepEqual(signer, signersByWallet[walletStorageKey][0]) {
		t.Errorf("loaded wallet signer differs from the original one")
	}
}

type mockPersistenceHandle struct {
	saved []persistence.DataDescriptor
}

func (mph *mockPersistenceHandle) Save(
	data []byte,
	directory string,
	name string,
) error {
	mph.saved = append(mph.saved, &mockDescriptor{
		name:      name,
		directory: directory,
		content:   data,
	})

	return nil
}

func (mph *mockPersistenceHandle) Snapshot(
	data []byte,
	directory string,
	name string,
) error {
	panic("not implemented")
}

func (mph *mockPersistenceHandle) ReadAll() (
	<-chan persistence.DataDescriptor,
	<-chan error,
) {
	outputData := make(chan persistence.DataDescriptor, len(mph.saved))
	outputErrors := make(chan error)

	for _, descriptor := range mph.saved {
		outputData <- descriptor
	}

	close(outputData)
	close(outputErrors)

	return outputData, outputErrors
}

func (mph *mockPersistenceHandle) Archive(directory string) error {
	panic("not implemented")
}

func (mph *mockPersistenceHandle) Delete(directory string, name string) error {
	panic("not implemented")
}

type mockDescriptor struct {
	name      string
	directory string
	content   []byte
}

func (md *mockDescriptor) Name() string {
	return md.name
}

func (md *mockDescriptor) Directory() string {
	return md.directory
}

func (md *mockDescriptor) Content() ([]byte, error) {
	return md.content, nil
}
