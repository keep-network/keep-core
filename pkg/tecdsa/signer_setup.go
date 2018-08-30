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

	// Information about the signing group. Holds information about all the members,
	// including the signer itself.
	// Initially empty, populated as each other signer announces its presence.
	// Signers are removed from the group if they misbehave or do not reply.
	signerGroup *signerGroup

	peerProtocolParameters map[string]*protocolParameters
}

type protocolParameters struct {
	commitmentPublicKey *bn256.G2
}

func (sc *signerCore) commitmentMasterPublicKey() *bn256.G2 {
	return sc.peerProtocolParameters[sc.ID].commitmentPublicKey
}

func (sc *signerCore) commitmentVerificationMasterPublicKey(
	signerID string,
) *bn256.G2 {
	return sc.peerProtocolParameters[signerID].commitmentPublicKey
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

	sc.peerProtocolParameters = make(map[string]*protocolParameters)
	sc.peerProtocolParameters[sc.ID] = &protocolParameters{
		commitmentPublicKey: publicKey,
	}

	return &CommitmentMasterPublicKeyMessage{
		signerID:        sc.ID,
		masterPublicKey: publicKey.Marshal(),
	}, nil
}

// ReceiveCommitmentMasterPublicKeys takes all the received
// `CommitmentMasterPublicKeyMessage`s and saves the commitment master public
// key value specific for the signer. This value is used later to validate
// commitments from the given signer.
func (sc *signerCore) ReceiveCommitmentMasterPublicKeys(
	messages []*CommitmentMasterPublicKeyMessage,
) error {
	if len(messages) != sc.signerGroup.InitialGroupSize {
		return fmt.Errorf(
			"master public key messages required from all group members; got %v, expected %v",
			len(messages),
			sc.signerGroup.InitialGroupSize,
		)
	}

	for _, message := range messages {
		if message.signerID != sc.ID {
			masterPublicKey := new(bn256.G2)
			masterPublicKey.Unmarshal(
				message.masterPublicKey,
			)

			sc.peerProtocolParameters[message.signerID] = &protocolParameters{
				commitmentPublicKey: masterPublicKey,
			}
		}
	}

	return nil
}
