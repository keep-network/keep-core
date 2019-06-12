package persistence

import (
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

type encryptedPersistence struct {
	handle   Handle
	password [32]byte
}

const (
	// SymmetricKeyLength represents the byte size of the key.
	SymmetricKeyLength = 32

	// NonceSize represents the byte size of nonce for XSalsa20 cipher used for
	// SymmetricKey encryption.
	NonceSize = 24
)

// NewEncryptedPersistence creates an adapter for the disk persistence to store data
// in an encrypted format
func NewEncryptedPersistence(handle Handle, password string) Handle {
	var secret [SymmetricKeyLength]byte
	copy(secret[:], password)

	return &encryptedPersistence{
		handle:   handle,
		password: secret,
	}
}

func (ep *encryptedPersistence) Save(data []byte, directory string, name string) error {
	var nonce [NonceSize]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return fmt.Errorf("key encryption failed [%v]", err)
	}

	encrypted := secretbox.Seal(nonce[:], data, &nonce, &ep.password)

	return ep.handle.Save(encrypted, directory, name)
}

func (ep *encryptedPersistence) ReadAll() ([][]byte, error) {
	encryptedMemberships, err := ep.handle.ReadAll()

	defer func() {
		// secretbox Open panics for invalid input
		if recover() != nil {
			err = fmt.Errorf("key decryption failed")
		}
	}()

	decryptedMemberships := [][]byte{}
	var nonce [NonceSize]byte

	// TODO: you iterate over the memberships here but in storage.go::readAll() you
	// itereate one more time. Combine the two?
	for _, encryptedMembership := range encryptedMemberships {
		copy(nonce[:], encryptedMembership[:NonceSize])
		decryptedMembership, ok := secretbox.Open(nil, encryptedMembership[NonceSize:], &nonce, &ep.password)
		if !ok {
			err = fmt.Errorf("key decryption failed")
		}
		decryptedMemberships = append(decryptedMemberships, decryptedMembership)
	}

	return decryptedMemberships, err
}

func (ep *encryptedPersistence) Archive(directory string) error {
	return ep.handle.Archive(directory)
}
