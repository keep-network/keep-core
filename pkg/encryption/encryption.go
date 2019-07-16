package encryption

// Box is a general interface to encrypt and decrypt an array of bytes.
type Box interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}
