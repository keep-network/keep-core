package tecdsa

import (
	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/tecdsa/commitment"
	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/keep-core/pkg/tecdsa/zkp"
	"github.com/keep-network/paillier"
)

// PublicKeyShareCommitmentMessage is a message payload that carries signer's
// commitment for a public DSA key share the signer generated.
// It's the very first message exchanged between signers during the T-ECDSA
// distributed key generation process. The message is expected to be broadcast
// publicly.
type PublicKeyShareCommitmentMessage struct {
	signerID string

	publicKeyShareCommitment *commitment.MultiTrapdoorCommitment // C_i
}

// KeyShareRevealMessage is a message payload that carries the sender's share of
// public and secret DSA key during T-ECDSA distributed key generation as well
// as proofs of correctness for the shares. Sender's share is encrypted with
// (t, n) Paillier threshold key. The message is expected to be broadcast
// publicly.
type KeyShareRevealMessage struct {
	signerID string

	secretKeyShare *paillier.Cypher // α_i = E(x_i)
	publicKeyShare *curve.Point     // y_i

	publicKeyShareDecommitmentKey *commitment.DecommitmentKey   // D_i
	secretKeyProof                *zkp.DsaPaillierKeyRangeProof // Π_i
}

// isValid checks secret and public key share against zero knowledge range proof
// shipped alongside them as well as commitment generated by the signer in the
// first phase of the key generation process. This function should be called
// for each received KeyShareRevealMessage before it's combined to a final key.
func (msg *KeyShareRevealMessage) isValid(
	commitmentMasterPublicKey *bn256.G2, // h
	publicKeyShareCommitment *commitment.MultiTrapdoorCommitment, // C_i
	zkpParams *zkp.PublicParameters,
) bool {
	commitmentValid := publicKeyShareCommitment.Verify(
		commitmentMasterPublicKey,
		msg.publicKeyShareDecommitmentKey,
		msg.publicKeyShare.Bytes(),
	)

	zkpValid := msg.secretKeyProof.Verify(
		msg.secretKeyShare, msg.publicKeyShare, zkpParams,
	)

	return commitmentValid && zkpValid
}
