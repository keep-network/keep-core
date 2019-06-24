package persistence

import (
	"crypto/sha256"

	"github.com/keep-network/keep-core/pkg/encryption"
)

// KeyLength represents the byte size of the key.
const KeyLength = encryption.KeyLength

type encryptedPersistence struct {
	delegate Handle
	box      encryption.Box
}

// NewEncryptedPersistence creates an adapter for the disk persistence to store data
// in an encrypted format
func NewEncryptedPersistence(handle Handle, password string) Handle {
	passwordBytes := []byte(password)

	return &encryptedPersistence{
		delegate: handle,
		box:      encryption.NewBox(sha256.Sum256(passwordBytes)),
	}
}

func (ep *encryptedPersistence) Save(data []byte, directory string, name string) error {
	encrypted, err := ep.box.Encrypt(data)
	if err != nil {
		return err
	}

	return ep.delegate.Save(encrypted, directory, name)
}

func (ep *encryptedPersistence) ReadAll() ([][]byte, error) {
	encryptedMemberships, err := ep.delegate.ReadAll()
	decryptedMemberships := [][]byte{}

	for _, encryptedMembership := range encryptedMemberships {
		decryptedMembership, err := ep.box.Decrypt(encryptedMembership)
		if err != nil {
			return nil, err
		}

		decryptedMemberships = append(decryptedMemberships, decryptedMembership)
	}

	return decryptedMemberships, err
}

func (ep *encryptedPersistence) Archive(directory string) error {
	return ep.delegate.Archive(directory)
}
