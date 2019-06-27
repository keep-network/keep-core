package persistence

import (
	"fmt"
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/testutils"

	"crypto/sha256"

	"github.com/keep-network/keep-core/pkg/encryption"
)

const accountPassword = "grzeski"

var (
	delegateMock   = &delegatePersistenceMock{}
	dataToEncrypt1 = []byte{'b', 'o', 'l', 'e', 'k'}
	dataToEncrypt2 = []byte{'l', 'o', 'l', 'e', 'k'}
	dataToEncrypt  = [][]byte{dataToEncrypt1, dataToEncrypt2}
)

func TestSaveReadAndDecryptData(t *testing.T) {
	encryptedPersistence := NewEncryptedPersistence(delegateMock, accountPassword)

	err := encryptedPersistence.Save(dataToEncrypt1, "dir1", "name1")
	if err != nil {
		t.Fatalf("Error occured while saving data [%v]", err)
	}
	encryptedPersistence.Save(dataToEncrypt2, "dir2", "name2")
	if err != nil {
		t.Fatalf("Error occured while saving data [%v]", err)
	}

	decryptedChan, errChan := encryptedPersistence.ReadAll()
	if err != nil {
		t.Fatalf("Error occured while reading data [%v]", err)
	}

	for err := range errChan {
		t.Error(err)
	}

	decrypted := make([][]byte, 0)
	for d := range decryptedChan {
		content, err := d.Content()
		if err != nil {
			t.Fatal(err)
		}

		decrypted = append(decrypted, content)
	}

	if len(decrypted) != len(dataToEncrypt) {
		t.Fatalf(
			"Unexpected number of decrypted items\nExpected: [%v]\nActual:   [%v]",
			len(dataToEncrypt),
			len(decrypted),
		)
	}

	for i := 0; i < len(dataToEncrypt); i++ {
		testutils.AssertBytesEqual(t, dataToEncrypt[i], decrypted[i])
	}
}

type delegatePersistenceMock struct{}

func (dpm *delegatePersistenceMock) Save(data []byte, directory string, name string) error {
	// noop
	return nil
}

func (dpm *delegatePersistenceMock) ReadAll() (<-chan DataDescriptor, <-chan error) {
	encrypted := encryptData()

	outputData := make(chan DataDescriptor, 2)
	outputErrors := make(chan error)

	outputData <- &testDataDescriptor{"1", "dir", encrypted[0]}
	outputData <- &testDataDescriptor{"2", "dir", encrypted[1]}

	close(outputData)
	close(outputErrors)

	return outputData, outputErrors
}

func (dpm *delegatePersistenceMock) Archive(directory string) error {
	// noop
	return nil
}

type testDataDescriptor struct {
	name      string
	directory string
	content   []byte
}

func (tdd *testDataDescriptor) Name() string {
	return tdd.name
}

func (tdd *testDataDescriptor) Directory() string {
	return tdd.directory
}

func (tdd *testDataDescriptor) Content() ([]byte, error) {
	return tdd.content, nil
}

func encryptData() [][]byte {
	passwordBytes := []byte(accountPassword)
	box := encryption.NewBox(sha256.Sum256(passwordBytes))

	encryptedData1, err := box.Encrypt(dataToEncrypt1)
	encryptedData2, err := box.Encrypt(dataToEncrypt2)
	if err != nil {
		fmt.Println("Error occured while encrypting data")
	}

	return [][]byte{encryptedData1, encryptedData2}
}
