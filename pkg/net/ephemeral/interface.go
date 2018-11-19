package ephemeral

// SymmetricKey is an ephemeral key shared between two parties that was
// established with Diffie-Hellman key exchange over a channel that does
// not need to be secure.
type SymmetricKey interface{}

type Ephemeral interface {
	// Performs ECDH with this private key and a given public key
	ECDH(remoteIdentity EphemeralPublicKey) SymmetricKey
}
