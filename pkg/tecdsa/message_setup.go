package tecdsa

// CommitmentMasterPublicKeyMessage is a message payload that carries a
// master public key of a multi-trapdoor commitment for the specific signer.
//
// The message is expected to be broadcast publicly.
type CommitmentMasterPublicKeyMessage struct {
	senderID string

	masterPublicKey []byte
}
