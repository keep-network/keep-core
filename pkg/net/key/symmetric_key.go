package key

// SymmetricKey is a session-scoped ECDH symmetric key.
type SymmetricKey interface {
	Encrypt([]byte) ([]byte, error)
}

// EphemeralPrivateKey is a session-scoped private elliptic curve key.
type EphemeralPrivateKey interface {
	Decrypt([]byte) ([]byte, error)

	Ecdh(EphemeralPublicKey) SymmetricKey

	Marshal([]byte, error)
	Unmarshal([]byte) error
}

// EphemeralPublicKey is a session-scoped public elliptic curve key.
type EphemeralPublicKey interface {
	Marshal([]byte, error)
	Unmarshal([]byte) error
}
