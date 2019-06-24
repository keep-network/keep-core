package persistence

import (
	"fmt"
	"reflect"
	"testing"

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

	decrypted, err := encryptedPersistence.ReadAll()
	if err != nil {
		t.Fatalf("Error occured while reading data [%v]", err)
	}

	if !reflect.DeepEqual(
		dataToEncrypt,
		decrypted,
	) {
		t.Fatalf("invalid decrypted results: \nexpected: %v\nactual:   %v\n",
			dataToEncrypt,
			decrypted,
		)
	}
}

type delegatePersistenceMock struct{}

func (dpm *delegatePersistenceMock) Save(data []byte, directory string, name string) error {
	// noop
	return nil
}

func (dpm *delegatePersistenceMock) ReadAll() ([][]byte, error) {
	encrypted := encryptData()

	return [][]byte{encrypted[0], encrypted[1]}, nil
}

func (dpm *delegatePersistenceMock) Archive(directory string) error {
	// noop
	return nil
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
