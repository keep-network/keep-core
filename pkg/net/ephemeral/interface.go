package ephemeral

// SymmetricKey is a session-scoped ECDH symmetric key.
type SymmetricKey interface {
	Encrypt([]byte) ([]byte, error)
}

// PrivateKey is a session-scoped private elliptic curve key.
type PrivateKey interface {
	Decrypt([]byte) ([]byte, error)

	Ecdh(PublicKey) SymmetricKey

	Marshal() ([]byte, error)
}

// PublicKey is a session-scoped public elliptic curve key.
type PublicKey interface {
	Marshal() ([]byte, error)
}
