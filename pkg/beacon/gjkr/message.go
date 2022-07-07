package gjkr

import (
	"fmt"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/group"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
)

// EphemeralPublicKeyMessage is a message payload that carries the sender's
// ephemeral public keys generated for all other group members.
//
// The receiver performs ECDH on a sender's ephemeral public key intended for
// the receiver and on the receiver's private ephemeral key, creating a symmetric
// key used for encrypting a conversation between the sender and the receiver.
// In case of an accusation for malicious behavior, the accusing party reveals
// its private ephemeral key so that all the other group members can resolve the
// accusation looking at messages exchanged between accuser and accused party.
// To validate correctness of accuser's private ephemeral key, all group members
// must know its ephemeral public key prior to exchanging any messages. Hence,
// this message contains all the generated public keys and it is broadcast
// within the group.
type EphemeralPublicKeyMessage struct {
	senderID group.MemberIndex // i

	ephemeralPublicKeys map[group.MemberIndex]*ephemeral.PublicKey // j -> Y_ij
}

// MemberCommitmentsMessage is a message payload that carries the sender's
// commitments to coefficients of the secret shares polynomial generated
// by member in the third phase of the protocol.
//
// It is expected to be broadcast.
type MemberCommitmentsMessage struct {
	senderID group.MemberIndex

	commitments []*bn256.G1 // slice of C_ik
}

// PeerSharesMessage is a message payload that carries shares `s_ij` and `t_ij`
// calculated by the sender `i` for all other group members individually.
//
// It is expected to be broadcast within the group.
type PeerSharesMessage struct {
	senderID group.MemberIndex // i

	shares map[group.MemberIndex]*peerShares // j -> (s_ij, t_ij)
}

type peerShares struct {
	encryptedShareS []byte // s_ij
	encryptedShareT []byte // t_ij
}

// SecretSharesAccusationsMessage is a message payload that carries all of the
// sender's accusations against other members of the threshold group.
// If all other members behaved honestly from the sender's point of view, this
// message should be broadcast but with an empty map of `accusedMembersKeys`.
//
// It is expected to be broadcast.
type SecretSharesAccusationsMessage struct {
	senderID group.MemberIndex

	accusedMembersKeys map[group.MemberIndex]*ephemeral.PrivateKey
}

// MemberPublicKeySharePointsMessage is a message payload that carries the
// sender's public key share points.
//
// It is expected to be broadcast.
type MemberPublicKeySharePointsMessage struct {
	senderID group.MemberIndex

	publicKeySharePoints []*bn256.G2 // A_ik = g^{a_ik} mod p
}

// PointsAccusationsMessage is a message payload that carries all of the sender's
// accusations against other members of the threshold group after public key share
// points validation.
// If all other members behaved honestly from the sender's point of view, this
// message should be broadcast but with an empty map of `accusedMembersKeys`.
// It is expected to be broadcast.
type PointsAccusationsMessage struct {
	senderID group.MemberIndex

	accusedMembersKeys map[group.MemberIndex]*ephemeral.PrivateKey
}

// MisbehavedEphemeralKeysMessage is a message payload that carries sender's
// ephemeral private keys used to generate ephemeral symmetric keys to encrypt
// communication with members from QUAL set which were marked as disqualified
// or inactive. It is expected to be broadcast.
type MisbehavedEphemeralKeysMessage struct {
	senderID group.MemberIndex

	privateKeys map[group.MemberIndex]*ephemeral.PrivateKey
}

// SenderID returns protocol-level identifier of the message sender.
func (epkm *EphemeralPublicKeyMessage) SenderID() group.MemberIndex {
	return epkm.senderID
}

// SenderID returns protocol-level identifier of the message sender.
func (mcm *MemberCommitmentsMessage) SenderID() group.MemberIndex {
	return mcm.senderID
}

// SenderID returns protocol-level identifier of the message sender.
func (psm *PeerSharesMessage) SenderID() group.MemberIndex {
	return psm.senderID
}

// SenderID returns protocol-level identifier of the message sender.
func (ssam *SecretSharesAccusationsMessage) SenderID() group.MemberIndex {
	return ssam.senderID
}

// SenderID returns protocol-level identifier of the message sender.
func (mpkspm *MemberPublicKeySharePointsMessage) SenderID() group.MemberIndex {
	return mpkspm.senderID
}

// SenderID returns protocol-level identifier of the message sender.
func (pam *PointsAccusationsMessage) SenderID() group.MemberIndex {
	return pam.senderID
}

// SenderID returns protocol-level identifier of the message sender.
func (mekm *MisbehavedEphemeralKeysMessage) SenderID() group.MemberIndex {
	return mekm.senderID
}

func newPeerSharesMessage(senderID group.MemberIndex) *PeerSharesMessage {
	return &PeerSharesMessage{
		senderID: senderID,
		shares:   make(map[group.MemberIndex]*peerShares),
	}
}

func (psm *PeerSharesMessage) addShares(
	receiverID group.MemberIndex,
	shareS *big.Int,
	shareT *big.Int,
	symmetricKey ephemeral.SymmetricKey,
) error {
	encryptedS, err := symmetricKey.Encrypt(shareS.Bytes())
	if err != nil {
		return fmt.Errorf("could not encrypt S share [%v]", err)
	}

	encryptedT, err := symmetricKey.Encrypt(shareT.Bytes())
	if err != nil {
		return fmt.Errorf("could not encrypt T share [%v]", err)
	}

	psm.shares[receiverID] = &peerShares{encryptedS, encryptedT}

	return nil
}

func (psm *PeerSharesMessage) decryptShareS(
	receiverID group.MemberIndex,
	key ephemeral.SymmetricKey,
) (*big.Int, error) {
	shares, ok := psm.shares[receiverID]
	if !ok {
		return nil, fmt.Errorf("no shares for receiver %v", receiverID)
	}

	decryptedS, err := key.Decrypt(shares.encryptedShareS)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt S share [%v]", err)
	}

	return new(big.Int).SetBytes(decryptedS), nil
}

func (psm *PeerSharesMessage) decryptShareT(
	receiverID group.MemberIndex,
	key ephemeral.SymmetricKey,
) (*big.Int, error) {
	shares, ok := psm.shares[receiverID]
	if !ok {
		return nil, fmt.Errorf("no shares for receiver %v", receiverID)
	}

	decryptedT, err := key.Decrypt(shares.encryptedShareT)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt T share [%v]", err)
	}

	return new(big.Int).SetBytes(decryptedT), nil
}

func (psm *PeerSharesMessage) decryptShares(
	receiverID group.MemberIndex,
	key ephemeral.SymmetricKey,
) (*big.Int, *big.Int, error) {
	shareS, err := psm.decryptShareS(receiverID, key) // s_mj
	if err != nil {
		return nil, nil, fmt.Errorf("cannot decrypt share S [%v]", err)
	}
	shareT, err := psm.decryptShareT(receiverID, key) // t_mj
	if err != nil {
		return nil, nil, fmt.Errorf("cannot decrypt share T [%v]", err)
	}

	return shareS, shareT, nil
}
