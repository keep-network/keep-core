package gjkr

import (
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func (mcm *MemberCommitmentsMessage) SetCommitment(
	index int,
	commitment *bn256.G1,
) {
	mcm.commitments[index] = commitment
}

func (psm *PeerSharesMessage) SetShares(
	memberIndex group.MemberIndex,
	encryptedShareS, encryptedShareT []byte,
) {
	psm.shares[memberIndex] = &peerShares{
		encryptedShareS: encryptedShareS,
		encryptedShareT: encryptedShareT,
	}
}

func (ssam *SecretSharesAccusationsMessage) SetAccusedMemberKey(
	memberIndex group.MemberIndex,
	privateKey *ephemeral.PrivateKey,
) {
	ssam.accusedMembersKeys[memberIndex] = privateKey
}

func (mpkspm *MemberPublicKeySharePointsMessage) SetPublicKeyShare(
	index int,
	publicKeyShare *bn256.G2,
) {
	mpkspm.publicKeySharePoints[index] = publicKeyShare
}

func (mpkspm *MemberPublicKeySharePointsMessage) RemovePublicKeyShare(
	index int,
) {
	mpkspm.publicKeySharePoints = append(
		mpkspm.publicKeySharePoints[:index],
		mpkspm.publicKeySharePoints[index+1:]...,
	)
}

func (dekm *DisqualifiedEphemeralKeysMessage) SetPrivateKey(
	memberIndex group.MemberIndex,
	privateKey *ephemeral.PrivateKey,
) {
	dekm.privateKeys[memberIndex] = privateKey
}
