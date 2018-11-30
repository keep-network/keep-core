package gjkr

import (
	"fmt"
	"sync"
)

// For complaint resolution, group members need to have access to messages
// exchanged between the accuser and the accused party. There are two situations
// in the DKG protocol where group members generate values for every other group
// member:
//
// - Ephemeral ECDH (phase 2) - after each group member generates an ephemeral
// keypair for each other group member and broadcasts those ephemeral public keys
// in the clear (phase 1), group members must ECDH those public keys with the
// ephemeral private key for that group member to derive a symmetric key.
// In the case of an accusation, members performing compliant resolution need to
// validate the private ephemeral key revealed by the accuser. To perform the
// validation, members need to compare public ephemeral key published by the
// accuser in phase 1 with the private ephemeral key published by the accuser.
//
// - Polynomial generation (phase 3) - each group member generates two sharing
// polynomials, and calculates shares as points on these polynomials individually
// for each other group member. Shares are publicly broadcast, encrypted with a
// symmetric key established between the sender and receiver. In the case of an
// accusation, members performing compliant resolution need to look at the shares
// sent by the accused party. To do this, they read the round 3 message from the
// buffer, passing the symmetric key used between the accuser and accused so that
// the round 3 message from the accused party can be decrypted.
type evidenceLog interface {
	// ephemeralPublicKeyMessage returns the `EphemeralPublicKeyMessage`
	// broadcast in the first protocol round by the given sender for the
	// given receiver.
	ephemeralPublicKeyMessage(
		sender MemberID,
		receiver MemberID,
	) *EphemeralPublicKeyMessage

	// peerSharesMessage returns the `PeerShareMessage` broadcast in the third
	// protocol round by the given sender for the given receiver.
	peerSharesMessage(
		sender MemberID,
		receiver MemberID,
	) *PeerSharesMessage
}

// dkgEvidenceLog is an implementation of an evidenceLog, containing two map of
// maps, from sender to receiver to message.
type dkgEvidenceLog struct {
	// senderID -> receiverID -> *EphemeralPublicKeyMessage
	pubKeyMessageLog     map[MemberID]map[MemberID]*EphemeralPublicKeyMessage
	pubKeyMessageLogLock sync.Mutex

	// senderID -> receiverID -> *PeerSharesMessage
	peerSharesMessageLog     map[MemberID]map[MemberID]*PeerSharesMessage
	peerSharesMessageLogLock sync.Mutex
}

// PutEphemeralMessage is a function that takes a single EphemeralPubKeyMessage
// and stores that information as evidence for future accusation trials for a
// given (sender, receiver) pair. If a message already exists for the given pair,
// we return an error to the user.
func (d *dkgEvidenceLog) PutEphemeralMessage(
	pubKeyMessage *EphemeralPublicKeyMessage,
) error {
	d.pubKeyMessageLogLock.Lock()
	defer d.pubKeyMessageLogLock.Unlock()

	senderLog, ok := d.pubKeyMessageLog[pubKeyMessage.senderID]
	if !ok {
		senderLog = make(map[MemberID]*EphemeralPublicKeyMessage)
	}

	if message, ok := senderLog[pubKeyMessage.receiverID]; ok {
		return fmt.Errorf(
			"message %v exists for sender %v and receiver %v",
			message,
			pubKeyMessage.senderID,
			pubKeyMessage.receiverID,
		)
	}

	senderLog[pubKeyMessage.receiverID] = pubKeyMessage
	return nil
}

// PutPeerSharesMessage is a function that takes a single EphemeralPubKeyMessage
// and stores that information as evidence for future accusation trials for a
// given (sender, receiver) pair. If a message already exists for the given pair,
// we return an error to the user.
func (d *dkgEvidenceLog) PutPeerSharesMessage(
	sharesMessage *PeerSharesMessage,
) error {
	d.peerSharesMessageLogLock.Lock()
	defer d.peerSharesMessageLogLock.Unlock()

	senderLog, ok := d.peerSharesMessageLog[sharesMessage.senderID]
	if !ok {
		senderLog = make(map[MemberID]*PeerSharesMessage)
	}

	if message, ok := senderLog[sharesMessage.receiverID]; ok {
		return fmt.Errorf(
			"message %v exists for sender %v and receiver %v",
			message,
			sharesMessage.senderID,
			sharesMessage.receiverID,
		)
	}

	senderLog[sharesMessage.receiverID] = sharesMessage
	return nil
}

func (d *dkgEvidenceLog) ephemeralPublicKeyMessage(
	sender MemberID,
	receiver MemberID,
) *EphemeralPublicKeyMessage {
	d.pubKeyMessageLogLock.Lock()
	defer d.pubKeyMessageLogLock.Unlock()

	senderLog, ok := d.pubKeyMessageLog[sender]
	if !ok {
		return nil
	}

	message, ok := senderLog[receiver]
	if !ok {
		return nil
	}

	return message
}

func (d *dkgEvidenceLog) peerSharesMessage(
	sender MemberID,
	receiverID MemberID,
) *PeerSharesMessage {
	d.peerSharesMessageLogLock.Lock()
	defer d.peerSharesMessageLogLock.Unlock()

	senderLog, ok := d.peerSharesMessageLog[sender]
	if !ok {
		return nil
	}

	message, ok := senderLog[receiverID]
	if !ok {
		return nil
	}

	return message
}
