package ephemeral

import (
	"crypto/sha256"

	"github.com/btcsuite/btcd/btcec"
)

// SymmetricEcdhKey is an ephemeral Elliptic Curve key created with
// Diffie-Hellman key exchange and implementing `SymmetricKey` interface.
type SymmetricEcdhKey struct {
	key [sha256.Size]byte
}

// ECDH performs Elliptic Curve Diffie-Hellman operation between public and
// private key. The returned value is `SymmetricEcdhKey` that can be used
// for encryption and decryption.
func (privk *EphemeralPrivateKey) ECDH(
	remotePublicKey *EphemeralPublicKey) *SymmetricEcdhKey {
	shared := btcec.GenerateSharedSecret(
		(*btcec.PrivateKey)(privk), (*btcec.PublicKey)(remotePublicKey),
	)

	return &SymmetricEcdhKey{sha256.Sum256(shared)}
}
