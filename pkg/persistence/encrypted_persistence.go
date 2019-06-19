package persistence

import (
	"github.com/keep-network/keep-core/pkg/secret"
)

// SymmetricKeyLength represents the byte size of the key.
const SymmetricKeyLength = 32

type encryptedPersistence struct {
	handle Handle
	secret *secret.Secret
}

// NewEncryptedPersistence creates an adapter for the disk persistence to store data
// in an encrypted format
func NewEncryptedPersistence(handle Handle, password string) Handle {
	var secretBytes [SymmetricKeyLength]byte
	copy(secretBytes[:], password)

	return &encryptedPersistence{
		handle: handle,
		secret: secret.NewSecret(secretBytes),
	}
}

func (ep *encryptedPersistence) Save(data []byte, directory string, name string) error {
	encrypted, err := ep.secret.Encrypt(data)
	if err != nil {
		return err
	}

	return ep.handle.Save(encrypted, directory, name)
}

func (ep *encryptedPersistence) ReadAll() ([][]byte, error) {
	encryptedMemberships, err := ep.handle.ReadAll()
	decryptedMemberships := [][]byte{}

	for _, encryptedMembership := range encryptedMemberships {
		decryptedMembership, err := ep.secret.Decrypt(encryptedMembership)
		if err != nil {
			return nil, err
		}

		decryptedMemberships = append(decryptedMemberships, decryptedMembership)
	}

	return decryptedMemberships, err
}

func (ep *encryptedPersistence) Archive(directory string) error {
	return ep.handle.Archive(directory)
}
