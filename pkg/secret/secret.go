package secret

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	sb "golang.org/x/crypto/nacl/secretbox"
)

const (
	// SymmetricKeyLength represents the byte size of the key.
	SymmetricKeyLength = 32

	// NonceSize represents the byte size of nonce for XSalsa20 cipher used for
	// SymmetricKey encryption.
	NonceSize = 24
)

// Secret is used to encrypt and decrypt a plaintext.
type Secret struct {
	key [SymmetricKeyLength]byte
}

// NewSecret uses XSalsa20 and Poly1305 to encrypt and decrypt the plaintext
// with the key.
func NewSecret(secret [SymmetricKeyLength]byte) *Secret {
	return &Secret{
		key: secret,
	}
}

// Encrypt takes the input plaintext and uses XSalsa20 and Poly1305 to encrypt
// the plaintext with the key.
func (s *Secret) Encrypt(plaintext []byte) ([]byte, error) {
	// The nonce needs to be unique, but not secure. Therefore we include it
	// at the beginning of the ciphertext.
	var nonce [NonceSize]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, fmt.Errorf("key encryption failed [%v]", err)
	}

	return sb.Seal(nonce[:], plaintext, &nonce, &s.key), nil
}

// Decrypt takes the input ciphertext and decrypts it.
func (s *Secret) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	defer func() {
		// secretbox Open panics for invalid input
		if recover() != nil {
			err = errors.New("symmetric key decryption failed")
		}
	}()

	var nonce [NonceSize]byte
	copy(nonce[:], ciphertext[:NonceSize])

	plaintext, ok := sb.Open(nil, ciphertext[NonceSize:], &nonce, &s.key)
	if !ok {
		err = fmt.Errorf("symmetric key decryption failed")
	}

	return
}
