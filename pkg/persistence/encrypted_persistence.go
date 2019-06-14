package persistence

import (
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	// SymmetricKeyLength represents the byte size of the key.
	SymmetricKeyLength = 32

	// NonceSize represents the byte size of nonce for XSalsa20 cipher used for
	// SymmetricKey encryption.
	NonceSize = 24
)

type encryptedPersistence struct {
	handle   Handle
	password [SymmetricKeyLength]byte
}

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
	encrypted, err := encrypt(data, ep.password)
	if err != nil {
		return err
	}

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

	for _, encryptedMembership := range encryptedMemberships {
		copy(nonce[:], encryptedMembership[:NonceSize])
		decryptedMembership, err := decrypt(encryptedMembership, ep.password)
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

func decrypt(encryptedMembership []byte, key [32]byte) ([]byte, error) {
	var nonce [NonceSize]byte
	copy(nonce[:], encryptedMembership[:NonceSize])
	decryptedMembership, ok := secretbox.Open(nil, encryptedMembership[NonceSize:], &nonce, &key)
	if !ok {
		return nil, fmt.Errorf("key decryption failed for [%v]", encryptedMembership[NonceSize:])
	}

	return decryptedMembership, nil
}

func encrypt(data []byte, key [32]byte) ([]byte, error) {
	var nonce [NonceSize]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return []byte{}, fmt.Errorf("key encryption failed [%v]", err)
	}

	return secretbox.Seal(nonce[:], data, &nonce, &key), nil
}
