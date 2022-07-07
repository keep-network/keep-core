package chain

import "github.com/keep-network/keep-core/pkg/operator"

// Signing is an interface that provides ability to sign and verify
// signatures using operator's key associated with the chain.
type Signing interface {
	// PublicKey returns operator's public key in a serialized format.
	// The returned public key is used to Sign messages and can be later used
	// for verification.
	PublicKey() []byte

	// Sign the provided message with operator's private key. Returns the
	// signature or error in case signing failed.
	Sign(message []byte) ([]byte, error)

	// Verify the provided message against the signature using operator's
	// public key. Returns true if signature is valid and false otherwise.
	// If signature verification failed for some reason, an error is returned.
	Verify(message []byte, signature []byte) (bool, error)

	// VerifyWithPublicKey verifies the provided message against the signature
	// using the provided operator's public key. Returns true if signature is
	// valid and false otherwise. If signature verification failed for some
	// reason, an error is returned.
	VerifyWithPublicKey(
		message []byte,
		signature []byte,
		publicKey []byte,
	) (bool, error)

	// PublicKeyToAddress converts operator's public key to an address
	// associated with the chain.
	// TODO: Refactor to return Address type
	PublicKeyToAddress(publicKey *operator.PublicKey) ([]byte, error)

	// PublicKeyBytesToAddress converts operator's public key bytes to an address
	// associated with the chain.
	// TODO: Refactor to return Address type
	PublicKeyBytesToAddress(publicKey []byte) []byte
}
