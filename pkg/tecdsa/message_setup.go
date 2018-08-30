package tecdsa

// MasterPublicKeyShareMessage is a message payload that carries a share of a
// master public key of a multi-trapdoor commitment.
// It's a to be exchanged prior to key generation and signing processess in order
// to build a master public key needed for commitments generation.
// The message is expected to be broadcast publicly.
type MasterPublicKeyShareMessage struct {
	signerID string

	masterPublicKeyShare []byte
}
