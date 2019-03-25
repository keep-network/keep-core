package operator

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

// PrivateKey represents peer's static key associated with an on-chain
// stake. It is used to authenticate the peer and for attributability (signing).
type PrivateKey = ecdsa.PrivateKey

// PublicKey represents peer's static key associated with an on-chain
// stake. It is used to authenticate the peer and for attributability
// (verification).
type PublicKey = ecdsa.PublicKey

// Signature is the resulting slice of bytes when a PrivateKey signs a hashed
// message.
type Signature = []byte

// SignatureSize is a byte size of the calculated Signature.
const SignatureSize = 65

// GenerateKeyPair generates a new, random static key based on
// secp256k1 ethereum curve.
func GenerateKeyPair() (*PrivateKey, *PublicKey, error) {
	ecdsaKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, nil, err
	}

	return (*PrivateKey)(ecdsaKey), (*PublicKey)(&ecdsaKey.PublicKey), nil
}

// EthereumKeyToOperatorKey transforms a `go-ethereum`-based ECDSA key into the
// format supported by all packages used in keep-core.
func EthereumKeyToOperatorKey(ethereumKey *keystore.Key) (*PrivateKey, *PublicKey) {
	privKey := ethereumKey.PrivateKey
	return (*PrivateKey)(privKey), (*PublicKey)(&privKey.PublicKey)
}

// Sign calculates an ECDSA signature.
//
// This function is susceptible to chosen plaintext attacks that can leak
// information about the private key that is used for signing. Callers must
// be aware that the given hash cannot be chosen by an adversery. Common
// solution is to hash any input before calculating the signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
func Sign(hash []byte, privateKey *PrivateKey) (Signature, error) {
	sig, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, err
	}
	return Signature(sig), nil
}

// VerifySignature checks that the given pubkey created signature over message.
// The public key should be in compressed (33 bytes) or uncompressed (65 bytes) format.
// The signature should be in [R || S] format.
func VerifySignature(publicKey *PublicKey, hash []byte, sig Signature) error {
	// Convert the operator's static key into an uncompressed public key
	// which should be 65 bytes in length.
	uncompressedPubKey := crypto.FromECDSAPub(publicKey)
	// If our sig is in the [R || S || V] format, ensure we strip out
	// the Ethereum-specific recovery-id, V, if it already hasn't been done.
	if len(sig) == 65 {
		sig = sig[:len(sig)-1]
	}

	// The sig should be 64 bytes.
	if len(sig) != 64 {
		return fmt.Errorf(
			"malformed signature %+v with length %d",
			sig[:],
			len(sig),
		)
	}

	if verified := crypto.VerifySignature(
		uncompressedPubKey,
		hash,
		sig[:],
	); !verified {
		return fmt.Errorf("failed to verify signature")
	}

	return nil
}

// Marshal takes an operator's PublicKey and produces an uncompressed public key
// as a slice of bytes (as specified in ANSI X9.62).
func Marshal(publicKey *PublicKey) []byte {
	return elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
}
