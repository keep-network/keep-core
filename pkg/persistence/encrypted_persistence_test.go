package persistence

import (
	"fmt"
	"reflect"
	"testing"

	sec "github.com/keep-network/keep-core/pkg/secret"
)

const (
	accountPassword    = "grzeski"
	symmetricKeyLength = 32
)

var (
	handleMock     = &handlePersistenceMock{}
	dataToEncrypt1 = []byte{'b', 'o', 'l', 'e', 'k'}
	dataToEncrypt2 = []byte{'l', 'o', 'l', 'e', 'k'}
	dataToEncrypt  = [][]byte{dataToEncrypt1, dataToEncrypt2}
)

func TestReadAndDecryptData(t *testing.T) {
	encryptedPersistence := NewEncryptedPersistence(handleMock, accountPassword)

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

type handlePersistenceMock struct{}

func (hpm *handlePersistenceMock) Save(data []byte, directory string, name string) error {
	// noop
	return nil
}

func (hpm *handlePersistenceMock) ReadAll() ([][]byte, error) {
	encrypted := encryptData()

	return [][]byte{encrypted[0], encrypted[1]}, nil
}

func (hpm *handlePersistenceMock) Archive(directory string) error {
	// noop
	return nil
}

func encryptData() [][]byte {
	var secretBytes [32]byte
	copy(secretBytes[:], accountPassword)
	secret := sec.NewSecret(secretBytes)

	encryptedData1, err := secret.Encrypt(dataToEncrypt1)
	encryptedData2, err := secret.Encrypt(dataToEncrypt2)
	if err != nil {
		fmt.Println("Error occured while encrypting data")
	}

	return [][]byte{encryptedData1, encryptedData2}
}
