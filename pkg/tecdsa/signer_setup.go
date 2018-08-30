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

	signatureParameters *PublicSignatureParameters
	zkpParameters       *zkp.PublicParameters

	// Information about the signing group. Holds information about all the members,
	// including the signer itself.
	// Initially empty, populated as each other signer announces its presence.
	// Signers are removed from the group if they misbehave or do not reply.
	signerGroup *signerGroup
}

// GenerateMasterPublicKeyShare produces a MasterPublicKeyShareMessage and should
// be called by all members of the group on very early stage prior to generating
// any commitments.
//
// `MasterPublicKeyShareMessage` contains signer's multi-trapdoor commitment master
// public key share.
//
// The shares should be combined and set as master public key for each signer.
func (sc *signerCore) GenerateMasterPublicKeyShare() (*MasterPublicKeyShareMessage, error) {
	_, hShare, err := bn256.RandomG2(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("could not generate multi-trapdoor commitment master trapdoor public key share [%v]", err)
	}

	return &MasterPublicKeyShareMessage{
		signerID:             sc.ID,
		masterPublicKeyShare: hShare.Marshal(),
	}, nil
}

// CombineMasterPublicKeyShares combines all group `MasterPublicKeyShareMessage`s
// into a `masterPublicKey`.
//
// The shares are expected to be points in G2 abstract cyclic group of bn256 curve.
// Shares are combined by points addition.
func (sc *signerCore) CombineMasterPublicKeyShares(
	masterPublicKeySharesMessages []*MasterPublicKeyShareMessage,
) (*bn256.G2, error) {
	if len(masterPublicKeySharesMessages) != sc.signerGroup.InitialGroupSize {
		return nil, fmt.Errorf(
			"master public key share required from all group members; got %v, expected %v",
			len(masterPublicKeySharesMessages),
			sc.signerGroup.InitialGroupSize,
		)
	}

	masterPublicKey := new(bn256.G2)
	masterPublicKey.Unmarshal(
		masterPublicKeySharesMessages[0].masterPublicKeyShare,
	)

	for _, message := range masterPublicKeySharesMessages[1:] {
		masterPublicKeyShare := new(bn256.G2)
		masterPublicKeyShare.Unmarshal(message.masterPublicKeyShare)
		masterPublicKey = new(bn256.G2).Add(masterPublicKey, masterPublicKeyShare)
	}
	return masterPublicKey, nil
}
