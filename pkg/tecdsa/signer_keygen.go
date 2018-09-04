// Package tecdsa contains the code that implements Threshold ECDSA signatures.
// The approach is based on [GGN 16].
//
//     [GGN 16]: Gennaro R., Goldfeder S., Narayanan A. (2016) Threshold-Optimal
//          DSA/ECDSA Signatures and an Application to Bitcoin Wallet Security.
//          In: Manulis M., Sadeghi AR., Schneider S. (eds) Applied Cryptography
//          and Network Security. ACNS 2016. Lecture Notes in Computer Science,
//          vol 9696. Springer, Cham
package tecdsa

import (
	"crypto/rand"
	"errors"
	"fmt"
	mathrand "math/rand"

	"github.com/keep-network/keep-core/pkg/tecdsa/commitment"
	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/keep-core/pkg/tecdsa/zkp"
	"github.com/keep-network/paillier"
)

// LocalSigner represents T-ECDSA group member during the initialization
// phase. It is responsible for constructing a broadcast
// PublicKeyShareCommitmentMessage containing public DSA key share commitment
// and a KeyShareRevealMessage revealing in a Paillier-encrypted way generated
// secret DSA key share and an unencrypted public key share.
type LocalSigner struct {
	signerCore

	dsaKeyShare *dsaKeyShare

	// Intermediate value stored between the first and the second round of
	// key generation.
	//
	// In the first round, `LocalSigner` commits to the chosen public key share.
	// In the second round, it reveals the public key share along with the
	// decommitment key.
	//
	// Since a separate commitment is produced for each peer signer in the
	// group, decommitment keys must be stored separately for each peer.
	// The map's key is the peer signer's ID.
	publicDsaKeyShareDecommitmentKeys map[string]*commitment.DecommitmentKey
}

// Signer represents T-ECDSA group member in a fully initialized state,
// ready for signing. Each Signer has a reference to a ThresholdDsaKey used
// in a signing process. It represents a (t, n) threshold sharing of the
// underlying DSA key.
type Signer struct {
	signerCore

	dsaKey *ThresholdDsaKey
}

// NewLocalSigner creates a fully initialized `LocalSigner` instance for the
// provided Paillier `ThresholdPrivateKey`, group and ZKP parameters.
// Please keep in mind there should never be created two `LocalSigner`s
// for the same instance of a `ThresholdPrivateKey`.
func NewLocalSigner(
	paillierKey *paillier.ThresholdPrivateKey,
	publicParameters *PublicParameters,
	zkpParameters *zkp.PublicParameters,
	signerGroup *signerGroup,
) *LocalSigner {
	return &LocalSigner{
		signerCore: signerCore{
			ID:               generateMemberID(),
			paillierKey:      paillierKey,
			publicParameters: publicParameters,
			zkpParameters:    zkpParameters,
			signerGroup:      signerGroup,
		},
	}
}

func generateMemberID() string {
	memberID := "0"
	for memberID = fmt.Sprintf("%v", mathrand.Int31()); memberID == "0"; {
	}
	return memberID
}

// generateDsaKeyShare generates a DSA public and secret key shares and puts
// them into `dsaKeyShare`. Secret key share is a random integer from Z_q where
// `q` is the cardinality of Elliptic Curve and public key share is a point
// on the Curve g^secretKeyShare.
func (ls *LocalSigner) generateDsaKeyShare() (*dsaKeyShare, error) {
	curveParams := ls.publicParameters.Curve.Params()

	secretKeyShare, err := rand.Int(rand.Reader, curveParams.N)
	if err != nil {
		return nil, fmt.Errorf("could not generate DSA key share [%v]", err)
	}

	publicKeyShare := curve.NewPoint(
		ls.publicParameters.Curve.ScalarBaseMult(secretKeyShare.Bytes()),
	)

	return &dsaKeyShare{
		secretKeyShare: secretKeyShare,
		publicKeyShare: publicKeyShare,
	}, nil
}

// InitializeDsaKeyShares initializes key generation process by generating DSA
// key shares and publishing `PublicKeyShareCommitmentMessage` for each peer
// signer in the group. The message contains signer's public DSA key share
// commitment.
func (ls *LocalSigner) InitializeDsaKeyShares() (
	[]*PublicKeyShareCommitmentMessage,
	error,
) {
	// Generate and store signer's DSA key share
	keyShare, err := ls.generateDsaKeyShare()
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate DSA key shares [%v]", err,
		)
	}

	ls.dsaKeyShare = keyShare

	// Initialize map holding decommitment keys for each peer signer.
	ls.publicDsaKeyShareDecommitmentKeys = make(
		map[string]*commitment.DecommitmentKey,
	)

	// Generate a separate `PublicKeyShareCommitmentMessage` for each peer
	// signer. Use peer signer's commitment master public key for that.
	messages := make(
		[]*PublicKeyShareCommitmentMessage,
		ls.signerGroup.PeerSignerCount(),
	)
	for i, peerSignerID := range ls.peerSignerIDs() {
		peerProtocolParameters := ls.peerProtocolParameters[peerSignerID]
		commitment, decommitmentKey, err := commitment.Generate(
			peerProtocolParameters.commitmentMasterPublicKey,
			keyShare.publicKeyShare.Bytes(),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"could not generate DSA public key commitment [%v]", err,
			)
		}

		ls.publicDsaKeyShareDecommitmentKeys[peerSignerID] = decommitmentKey
		messages[i] = &PublicKeyShareCommitmentMessage{
			senderID:                 ls.ID,
			receiverID:               peerSignerID,
			publicKeyShareCommitment: commitment,
		}
	}

	return messages, nil
}

// RevealDsaKeyShares produces `KeyShareRevealMessage`s and should be called
// when `PublicKeyShareCommitmentMessage`s from all group members have been
// received.
//
// `KeyShareRevealMessage` contains signer's public DSA key share, decommitment
// key for this share (used to validate the commitment published in the previous
// `PublicKeyShareCommitmentMessage` message), encrypted secret DSA key share
// and ZKP for the secret key share correctness. Bear in mind the decommitment
// key is different for each peer signer. That's why, `KeyShareRevealMessage` is
// produced individually for each peer signer.
//
// Secret key share is encrypted with an additively homomorphic encryption
// scheme and sent to all other signers in the group along with the public key
// share.
func (ls *LocalSigner) RevealDsaKeyShares() ([]*KeyShareRevealMessage, error) {
	paillierRandomness, err := paillier.GetRandomNumberInMultiplicativeGroup(
		ls.paillierKey.N, rand.Reader,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate random r for Paillier [%v]", err,
		)
	}

	encryptedSecretKeyShare, err := ls.paillierKey.EncryptWithR(
		ls.dsaKeyShare.secretKeyShare, paillierRandomness,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not encrypt secret key share [%v]", err,
		)
	}

	rangeProof, err := zkp.CommitDsaPaillierKeyRange(
		ls.dsaKeyShare.secretKeyShare,
		ls.dsaKeyShare.publicKeyShare,
		encryptedSecretKeyShare,
		paillierRandomness,
		ls.zkpParameters,
		rand.Reader,
	)

	messages := make(
		[]*KeyShareRevealMessage,
		ls.signerGroup.PeerSignerCount(),
	)
	for i, peerSignerID := range ls.peerSignerIDs() {
		messages[i] = &KeyShareRevealMessage{
			senderID:                      ls.ID,
			receiverID:                    peerSignerID,
			secretKeyShare:                encryptedSecretKeyShare,
			publicKeyShare:                ls.dsaKeyShare.publicKeyShare,
			publicKeyShareDecommitmentKey: ls.publicDsaKeyShareDecommitmentKeys[peerSignerID],
			secretKeyProof:                rangeProof,
		}
	}
	return messages, nil
}

// CombineDsaKeyShares combines all group `PublicKeyShareCommitmentMessage`s and
// `KeyShareRevealMessage`s into a `ThresholdDsaKey` which is a (t, n) threshold
// sharing of an underlying secret DSA key. Secret and public
// DSA key shares are combined in the following way:
//
// E(secretKey) = E(secretKeyShare_1) + E(secretKeyShare_2) + ... + E(secretKeyShare_n)
// publicKey = publicKeyShare_1 + publicKeyShare_2 + ... + publicKeyShare_n
//
// E is an additively homomorphic encryption scheme, hence `+` operation is
// possible. Each key share share comes from the `KeyShareRevealMessage` that
// was sent by each `LocalSigner` of the signing group.
//
// Before shares are combined, messages are validated - we check whether
// the published public key share is what the signer originally committed to
// as well as we check validity of the secret key share using the provided ZKP.
//
// Every `PublicKeyShareCommitmentMessage` should have a corresponding
// `KeyShareRevealMessage`. They are matched by a signer ID contained in
// each of the messages.
//
// This function accepts all `PublicKeyShareCommitmentMessage`s and
// `KeyShareRevealMessage`s that were generated for the current signer.
// It's expected all peer signers in the group delivered both types of messages
// to the signer.
func (ls *LocalSigner) CombineDsaKeyShares(
	shareCommitments []*PublicKeyShareCommitmentMessage,
	revealedShares []*KeyShareRevealMessage,
) (*ThresholdDsaKey, error) {
	peerSignerCount := ls.signerGroup.PeerSignerCount()

	if len(shareCommitments) != peerSignerCount {
		return nil, fmt.Errorf(
			"commitments required from all group peer members; got %v, expected %v",
			len(shareCommitments),
			peerSignerCount,
		)
	}

	if len(revealedShares) != ls.signerGroup.PeerSignerCount() {
		return nil, fmt.Errorf(
			"all group peer members should reveal shares; Got %v, expected %v",
			len(revealedShares),
			peerSignerCount,
		)
	}

	// Combine secret and public key shares from peer signers
	secretKeyShares := make([]*paillier.Cypher, peerSignerCount)
	publicKeyShares := make([]*curve.Point, peerSignerCount)

	for i, commitmentMsg := range shareCommitments {
		foundMatchingRevealMessage := false

		for _, revealedSharesMsg := range revealedShares {

			if commitmentMsg.senderID == revealedSharesMsg.senderID {
				foundMatchingRevealMessage = true

				if revealedSharesMsg.isValid(
					ls.selfProtocolParameters().commitmentMasterPublicKey,
					commitmentMsg.publicKeyShareCommitment,
					ls.zkpParameters,
				) {
					secretKeyShares[i] = revealedSharesMsg.secretKeyShare
					publicKeyShares[i] = revealedSharesMsg.publicKeyShare
				} else {
					return nil, errors.New("KeyShareRevealMessage rejected")
				}
			}
		}

		if !foundMatchingRevealMessage {
			return nil, fmt.Errorf(
				"no matching share reveal message for signer with ID=%v",
				commitmentMsg.senderID,
			)
		}
	}

	// Add signer's own secret and public key share
	encryptedSecretKeyShare, err := ls.paillierKey.Encrypt(
		ls.dsaKeyShare.secretKeyShare, rand.Reader,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not encrypt secret key share [%v]", err,
		)
	}

	secretKeyShares = append(secretKeyShares, encryptedSecretKeyShare)
	publicKeyShares = append(publicKeyShares, ls.dsaKeyShare.publicKeyShare)

	// Combine signer's own and peer signers' shares together
	secretKey := ls.paillierKey.Add(secretKeyShares...)
	publicKey := publicKeyShares[0]
	for _, share := range publicKeyShares[1:] {
		publicKey = curve.NewPoint(ls.publicParameters.Curve.Add(
			publicKey.X, publicKey.Y, share.X, share.Y,
		))
	}

	return &ThresholdDsaKey{secretKey, publicKey}, nil
}

// WithDsaKey transforms `LocalSigner` into a `Signer` when the key generation
// process completes and `ThresholdDsaKey` is ready.
// There is a one instance of `ThresholdDsaKey` for all `Signer`s.
func (ls *LocalSigner) WithDsaKey(dsaKey *ThresholdDsaKey) *Signer {
	return &Signer{
		dsaKey:     dsaKey,
		signerCore: ls.signerCore,
	}
}

// PublicKey returns a public ECDSA key of the `Signer`.
// The public key is expected to be identical for all signers in
// a signing group.
func (s *Signer) PublicKey() *curve.Point {
	return s.dsaKey.PublicKey
}
