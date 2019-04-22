package ephemeral

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/nacl/secretbox"
)

// SymmetricKeyLength represents the byte size of the key.
const SymmetricKeyLength = 32

// NonceSize represents the byte size of nonce for XSalsa20 cipher used for
// SymmetricKey encryption.
const NonceSize = 24

// SymmetricEcdhKey is an ephemeral Elliptic Curve key created with
// Diffie-Hellman key exchange and implementing `SymmetricKey` interface.
type SymmetricEcdhKey struct {
	key [SymmetricKeyLength]byte
}

// Ecdh performs Elliptic Curve Diffie-Hellman operation between public and
// private key. The returned value is `SymmetricEcdhKey` that can be used
// for encryption and decryption.
func (pk *PrivateKey) Ecdh(publicKey *PublicKey) *SymmetricEcdhKey {
	shared := btcec.GenerateSharedSecret(
		(*btcec.PrivateKey)(pk),
		(*btcec.PublicKey)(publicKey),
	)

	return &SymmetricEcdhKey{sha256.Sum256(shared)}
}

// Encrypt takes the input plaintext and uses XSalsa20 and Poly1305 to encrypt
// and authenticate the message with the symmetric key.
func (sek *SymmetricEcdhKey) Encrypt(plaintext []byte) ([]byte, error) {
	// The nonce needs to be unique, but not secure. Therefore we include it
	// at the beginning of the ciphertext.
	var nonce [NonceSize]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, fmt.Errorf("symmetric key encryption failed [%v]", err)
	}

	return secretbox.Seal(nonce[:], plaintext, &nonce, &sek.key), nil
}

// Decrypt takes the input ciphertext and authenticates and decrypts it.
func (sek *SymmetricEcdhKey) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	defer func() {
		// secretbox Open panics for invalid input
		if recover() != nil {
			err = errors.New("symmetric key decryption failed")
		}
	}()

	var nonce [NonceSize]byte
	copy(nonce[:], ciphertext[:NonceSize])

	plaintext, ok := secretbox.Open(nil, ciphertext[NonceSize:], &nonce, &sek.key)
	if !ok {
		err = fmt.Errorf("symmetric key decryption failed")
	}

	return
}
