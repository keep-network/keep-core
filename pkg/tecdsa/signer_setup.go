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

	// IDs of all signers in active signer's group, including the signer itself.
	// Initially empty, populated as each other signer announces its presence.
	signerIDs []string
}

// RegisterSignerID adds a signer to the list of signers the local signer
// knows about.
func (ls *LocalSigner) RegisterSignerID(ID string) {
	ls.signerIDs = append(ls.signerIDs, ID)
}

// RemoveSignerID removes a signer from the list of signers the local signer
// knows about.
func (sc *signerCore) RemoveSignerID(ID string) {
	for i := 0; i < len(sc.signerIDs); i++ {
		if sc.signerIDs[i] == ID {
			sc.signerIDs = append(sc.signerIDs[:i], sc.signerIDs[i+1:]...)
		}
	}
}

// IsActiveSigner checks if a signer with given ID is one of the signers the local
// signer knows about.
func (sc *signerCore) IsActiveSigner(ID string) bool {
	for i := 0; i < len(sc.signerIDs); i++ {
		if sc.signerIDs[i] == ID {
			return true
		}
	}
	return false
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
	if len(masterPublicKeySharesMessages) != sc.groupParameters.GroupSize {
		return nil, fmt.Errorf(
			"master public key share required from all group members; got %v, expected %v",
			len(masterPublicKeySharesMessages),
			sc.groupParameters.GroupSize,
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
