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

	groupParameters *PublicParameters
	zkpParameters   *zkp.PublicParameters

	commitmentPublicKeys map[string]*bn256.G2
}

func (sc *signerCore) commitmentMasterPublicKey() *bn256.G2 {
	return sc.commitmentPublicKeys[sc.ID]
}

func (sc *signerCore) commitmentVerificationMasterPublicKey(
	signerID string,
) *bn256.G2 {
	return sc.commitmentPublicKeys[signerID]
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
	if len(messages) != sc.groupParameters.GroupSize {
		return fmt.Errorf(
			"master public key messages required from all group members; got %v, expected %v",
			len(messages),
			sc.groupParameters.GroupSize,
		)
	}

	sc.commitmentPublicKeys = make(map[string]*bn256.G2)
	for _, message := range messages {
		masterPublicKey := new(bn256.G2)
		masterPublicKey.Unmarshal(
			message.masterPublicKey,
		)

		sc.commitmentPublicKeys[message.signerID] = masterPublicKey
	}

	return nil
}
