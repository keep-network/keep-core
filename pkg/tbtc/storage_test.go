package tbtc

import (
	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"reflect"
	"testing"
)

func TestSaveSigner(t *testing.T) {
	persistenceHandle := &mockPersistenceHandle{}

	walletStorage := connectWalletStorage(persistenceHandle)

	signer := sampleSigner(t)

	err := walletStorage.saveSigner(signer)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertIntsEqual(
		t,
		"persisted items count",
		1,
		len(persistenceHandle.saved),
	)
	testutils.AssertIntsEqual(
		t,
		"cached items count",
		1,
		len(walletStorage.wallets),
	)
}

func TestLoadSigners(t *testing.T) {
	signer := sampleSigner(t)
	signerBytes, err := signer.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	persistenceHandle := &mockPersistenceHandle{
		saved: []persistence.DataDescriptor{
			&mockDescriptor{
				name:      "membership_1",
				directory: "wallet_1",
				content:   signerBytes,
			},
		},
	}

	// connectWalletStorage calls loadSigners under the hood
	walletStorage := connectWalletStorage(persistenceHandle)

	testutils.AssertIntsEqual(
		t,
		"cached items count",
		1,
		len(walletStorage.wallets),
	)

	signers := walletStorage.getSigners(signer.wallet.publicKey)

	testutils.AssertIntsEqual(
		t,
		"wallet signers count",
		1,
		len(signers),
	)

	if !reflect.DeepEqual(signer, signers[0]) {
		t.Errorf("loaded signer differs from the original one")
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
	return nil
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
	return nil
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
