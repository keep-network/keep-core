package tecdsa

import (
	"crypto/rand"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/tecdsa/zkp"
	"github.com/keep-network/paillier"
)

type signerCore struct {
	ID string

	paillierKey *paillier.ThresholdPrivateKey

	publicParameters *PublicParameters
	zkpParameters    *zkp.PublicParameters

	// Information about the signing group. Holds information about all the
	// members, including the signer itself.
	// Initially empty, populated as each other signer announces its presence.
	// Signers are removed from the group if they misbehave or do not reply.
	signerGroup *signerGroup

	protocolParameters map[string]*protocolParameters
}

type protocolParameters struct {
	commitmentMasterPublicKey *bn256.G2
}

func (sc *signerCore) selfProtocolParameters() *protocolParameters {
	return sc.protocolParameters[sc.ID]
}

func (sc *signerCore) peerSignerIDs() []string {
	peerIDs := make([]string, 0)
	for _, peerID := range sc.signerGroup.signerIDs {
		if peerID != sc.ID {
			peerIDs = append(peerIDs, peerID)
		}
	}

	return peerIDs
}

// GenerateCommitmentMasterPublicKey produces a CommitmentMasterPublicKeyMessage
// and should be called by all members of the group on very early stage, during
// the group setup, prior to generating any commitments.
//
// `CommitmentMasterPublicKeyMessage` contains signer-specific multi-trapdoor
// commitment master public key. For security reasons, each signer should
// produce its own key.
func (sc *signerCore) GenerateCommitmentMasterPublicKey() (
	*CommitmentMasterPublicKeyMessage,
	error,
) {
	_, publicKey, err := bn256.RandomG2(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate multi-trapdoor commitment master public key [%v]",
			err,
		)
	}

	sc.protocolParameters = make(map[string]*protocolParameters)
	sc.protocolParameters[sc.ID] = &protocolParameters{
		commitmentMasterPublicKey: publicKey,
	}

	return &CommitmentMasterPublicKeyMessage{
		senderID:        sc.ID,
		masterPublicKey: publicKey.Marshal(),
	}, nil
}

// ReceiveCommitmentMasterPublicKeys takes all the received
// `CommitmentMasterPublicKeyMessage`s and saves the commitment master public
// key value specific for the signer. This value is used later to validate
// commitments from the given signer.
// It's expected to receive messages from peer signers only.
func (sc *signerCore) ReceiveCommitmentMasterPublicKeys(
	messages []*CommitmentMasterPublicKeyMessage,
) error {
	if len(messages) != sc.signerGroup.PeerSignerCount() {
		return fmt.Errorf(
			"master public key messages required from all group peer members; got %v, expected %v",
			len(messages),
			sc.signerGroup.PeerSignerCount(),
		)
	}

	for _, message := range messages {
		if message.senderID != sc.ID {
			masterPublicKey := new(bn256.G2)
			masterPublicKey.Unmarshal(message.masterPublicKey)

			sc.protocolParameters[message.senderID] = &protocolParameters{
				commitmentMasterPublicKey: masterPublicKey,
			}
		}
	}

	return nil
}
