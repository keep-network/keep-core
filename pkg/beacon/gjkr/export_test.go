package gjkr

import (
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/group"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
)

func (epkm *EphemeralPublicKeyMessage) SetSenderID(
	senderID group.MemberIndex,
) {
	epkm.senderID = senderID
}

func (epkm *EphemeralPublicKeyMessage) SetPublicKey(
	memberIndex group.MemberIndex,
	publicKey *ephemeral.PublicKey,
) {
	epkm.ephemeralPublicKeys[memberIndex] = publicKey
}

func (epkm *EphemeralPublicKeyMessage) GetPublicKey(
	memberIndex group.MemberIndex,
) *ephemeral.PublicKey {
	return epkm.ephemeralPublicKeys[memberIndex]
}

func (epkm *EphemeralPublicKeyMessage) RemovePublicKey(
	memberIndex group.MemberIndex,
) {
	delete(epkm.ephemeralPublicKeys, memberIndex)
}

func (mcm *MemberCommitmentsMessage) SetCommitment(
	index int,
	commitment *bn256.G1,
) {
	mcm.commitments[index] = commitment
}

func (mcm *MemberCommitmentsMessage) RemoveCommitment(
	index int,
) {
	mcm.commitments = append(
		mcm.commitments[:index],
		mcm.commitments[index+1:]...,
	)
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

func (psm *PeerSharesMessage) AddShares(
	receiverID group.MemberIndex,
	shareS *big.Int,
	shareT *big.Int,
	symmetricKey ephemeral.SymmetricKey,
) error {
	return psm.addShares(receiverID, shareS, shareT, symmetricKey)
}

func (psm *PeerSharesMessage) RemoveShares(memberIndex group.MemberIndex) {
	delete(psm.shares, memberIndex)
}

func (ssam *SecretSharesAccusationsMessage) SetAccusedMemberKey(
	memberIndex group.MemberIndex,
	privateKey *ephemeral.PrivateKey,
) {
	ssam.accusedMembersKeys[memberIndex] = privateKey
}

func (ssam *SecretSharesAccusationsMessage) SetAccusedMemberKeys(
	accusedMembersKeys map[group.MemberIndex]*ephemeral.PrivateKey,
) {
	ssam.accusedMembersKeys = accusedMembersKeys
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

func (pam *PointsAccusationsMessage) SetAccusedMemberKey(
	memberIndex group.MemberIndex,
	privateKey *ephemeral.PrivateKey,
) {
	pam.accusedMembersKeys[memberIndex] = privateKey
}

func (pam *PointsAccusationsMessage) SetAccusedMemberKeys(
	accusedMembersKeys map[group.MemberIndex]*ephemeral.PrivateKey,
) {
	pam.accusedMembersKeys = accusedMembersKeys
}

func (mekm *MisbehavedEphemeralKeysMessage) SetPrivateKey(
	memberIndex group.MemberIndex,
	privateKey *ephemeral.PrivateKey,
) {
	mekm.privateKeys[memberIndex] = privateKey
}

func (mekm *MisbehavedEphemeralKeysMessage) SetPrivateKeys(
	privateKeys map[group.MemberIndex]*ephemeral.PrivateKey,
) {
	mekm.privateKeys = privateKeys
}

func (mekm *MisbehavedEphemeralKeysMessage) RemovePrivateKey(
	memberIndex group.MemberIndex,
) {
	delete(mekm.privateKeys, memberIndex)
}

func GeneratePolynomial(degree int) ([]*big.Int, error) {
	return generatePolynomial(degree)
}

func EvaluateMemberShare(
	memberID group.MemberIndex,
	coefficients []*big.Int,
) *big.Int {
	// evaluateMemberShare method is stateless so we use it via
	// a fake CommittingMember in order to avoid code duplication
	cm := &CommittingMember{}
	return cm.evaluateMemberShare(memberID, coefficients)
}

func (ms *messageStorage) removeMessage(sender group.MemberIndex) {
	ms.cacheLock.Lock()
	defer ms.cacheLock.Unlock()

	delete(ms.cache, sender)
}
