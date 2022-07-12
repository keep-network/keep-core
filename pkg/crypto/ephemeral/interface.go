package ephemeral

// SymmetricKey is an ephemeral key shared between two parties that was
// established with Diffie-Hellman key exchange over a channel that does
// not need to be secure.
type SymmetricKey interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}
