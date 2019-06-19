package ephemeral

import (
	"crypto/sha256"

	"github.com/btcsuite/btcd/btcec"
	"github.com/keep-network/keep-core/pkg/secret"
)

// SymmetricKeyLength represents the byte size of the key.
const SymmetricKeyLength = 32

// NonceSize represents the byte size of nonce for XSalsa20 cipher used for
// SymmetricKey encryption.
const NonceSize = 24

// SymmetricEcdhKey is an ephemeral Elliptic Curve key created with
// Diffie-Hellman key exchange and implementing `SymmetricKey` interface.
type SymmetricEcdhKey struct {
	secret *secret.Secret
}

// Ecdh performs Elliptic Curve Diffie-Hellman operation between public and
// private key. The returned value is `SymmetricEcdhKey` that can be used
// for encryption and decryption.
func (pk *PrivateKey) Ecdh(publicKey *PublicKey) *SymmetricEcdhKey {
	shared := btcec.GenerateSharedSecret(
		(*btcec.PrivateKey)(pk),
		(*btcec.PublicKey)(publicKey),
	)

	return &SymmetricEcdhKey{
		secret: secret.NewSecret(sha256.Sum256(shared)),
	}
}

// Encrypt plaintext.
func (sek *SymmetricEcdhKey) Encrypt(plaintext []byte) ([]byte, error) {
	return sek.secret.Encrypt(plaintext)
}

// Decrypt ciphertext.
func (sek *SymmetricEcdhKey) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	return sek.secret.Decrypt(ciphertext)
}
