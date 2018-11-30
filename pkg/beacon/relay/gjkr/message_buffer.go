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
// log, and decrypt it using the symmetric key used between the accuser and
// accused party. The key is publicly revealed by the accuser.
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
	pubKeyMessageLog *messageStorage

	// senderID -> receiverID -> *PeerSharesMessage
	peerSharesMessageLog *messageStorage
}

// NewDkgEvidenceLog returns a dkgEvidenceLog with backing stores for future
// accusations against EphemeralPublicKeyMessages and PeerShareMessages.
func NewDkgEvidenceLog() *dkgEvidenceLog {
	return &dkgEvidenceLog{
		pubKeyMessageLog:     newMessageStorage(),
		peerSharesMessageLog: newMessageStorage(),
	}
}

// PutEphemeralMessage is a function that takes a single EphemeralPubKeyMessage
// and stores that information as evidence for future accusation trials for a
// given (sender, receiver) pair. If a message already exists for the given pair,
// we return an error to the user.
func (d *dkgEvidenceLog) PutEphemeralMessage(
	pubKeyMessage *EphemeralPublicKeyMessage,
) error {
	return d.pubKeyMessageLog.putMessage(
		pubKeyMessage.senderID,
		pubKeyMessage.receiverID,
		pubKeyMessage,
	)
}

// PutPeerSharesMessage is a function that takes a single PeerSharesMessage
// and stores that information as evidence for future accusation trials for a
// given (sender, receiver) pair. If a message already exists for the given pair,
// we return an error to the user.
func (d *dkgEvidenceLog) PutPeerSharesMessage(
	sharesMessage *PeerSharesMessage,
) error {
	return d.peerSharesMessageLog.putMessage(
		sharesMessage.senderID,
		sharesMessage.receiverID,
		sharesMessage,
	)
}

func (d *dkgEvidenceLog) ephemeralPublicKeyMessage(
	sender MemberID,
	receiver MemberID,
) *EphemeralPublicKeyMessage {
	storedMessage := d.pubKeyMessageLog.getMessage(sender, receiver)
	switch message := storedMessage.(type) {
	case *EphemeralPublicKeyMessage:
		return message
	}
	return nil
}

func (d *dkgEvidenceLog) peerSharesMessage(
	sender MemberID,
	receiver MemberID,
) *PeerSharesMessage {
	storedMessage := d.peerSharesMessageLog.getMessage(sender, receiver)
	switch message := storedMessage.(type) {
	case *PeerSharesMessage:
		return message
	}
	return nil
}

// messageStorage is the underlying cache used by our evidenceLog implementation
// it implements a generic get and put of messages through a mapping of a
// (sender, receiver) pair.
type messageStorage struct {
	cache     map[MemberID]map[MemberID]interface{}
	cacheLock sync.Mutex
}

func newMessageStorage() *messageStorage {
	return &messageStorage{
		cache: make(map[MemberID]map[MemberID]interface{}),
	}
}

func (ms *messageStorage) getMessage(sender, receiver MemberID) interface{} {
	ms.cacheLock.Lock()
	defer ms.cacheLock.Unlock()

	senderLog, ok := ms.cache[sender]
	if !ok {
		return nil
	}

	message, ok := senderLog[receiver]
	if !ok {
		return nil
	}

	return message
}

func (ms *messageStorage) putMessage(
	sender, receiver MemberID,
	message interface{},
) error {
	ms.cacheLock.Lock()
	defer ms.cacheLock.Unlock()

	if _, ok := ms.cache[sender]; !ok {
		ms.cache[sender] = make(map[MemberID]interface{})
	}

	if _, ok := ms.cache[sender][receiver]; ok {
		return fmt.Errorf(
			"message exists for sender %v and receiver %v",
			sender,
			receiver,
		)
	}

	ms.cache[sender][receiver] = message
	return nil
}
