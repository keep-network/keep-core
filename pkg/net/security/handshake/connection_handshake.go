// Package handshake contains the code that implements authentication handshake
// performed when a new connection between two peers is established, as
// described in the network security implementation [RFC], section 1.2.3
// and 1.2.4.
//
// Each peer wanting to join a network needs to provide a proof of ownership of
// an on-chain identity with an associated stake. As part of the network join
// handshake, peer responding to the handshake will also provide proof of its
// own stake. The same handshake is executed when two peers already being a part
// of the network establish a new connection with each other.
//
// The handshake is a 3-round procedure when two parties called initiator and
// responder exchange messages. The entire handshake procedure can be described
// with the following diagram:
//
// # INITIATOR                             RESPONDER
//
// [Act 1]
// nonce1 = random_nonce()
// act1Message{nonce1, protocol_id1} ---->
//
// --------------------------------------- [Act 2]
// --------------------------------------- nonce2 = random_nonce()
// --------------------------------------- challenge = sha256(nonce1 || nonce2)
// --------------------------------------- <---- act2Message{challenge, nonce2, protocol_id2}
//
// [Act 3]
// challenge = sha256(nonce1 || nonce2)
// act3Message{challenge} ---->
//
// act1Message, act2Message, and act3Message are messages exchanged between
// initiator and responder in acts one, two, and three of the handshake,
// respectively.
//
// InitiatorAct1, InitiatorAct2, and InitiatorAct3 represent the state of the
// initiator in rounds one, two, and three of the handshake, respectively.
//
// ResponderAct2 and ResponderAct3 represent the state of the responder in
// rounds two and three of the handshake, respectively. Since the first act of
// the handshake is initiated by the initiator and the responder has no internal
// state before receiving the first message, there is no representation for
// responder state in the act one.
//
//	[RFC]: /docs/rfc/rfc-2-network-security-implementation.adoc
package handshake

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
)

// Act1Message is sent in the first handshake act by the initiator to the
// responder. It contains randomly generated `nonce1`, an 8-byte (64-bit)
// unsigned integer as well as the protocol identifier.
//
// act1Message should be signed with initiator's static private key.
type Act1Message struct {
	nonce1    uint64
	protocol1 string
}

// Act2Message is sent in the second handshake act by the responder to the
// initiator. It contains randomly generated `nonce2`, which is an 8-byte
// unsigned integer, `challenge`, which is the result of SHA256 on the
// concatenated bytes of `nonce1` and `nonce2`, and the protocol identifier.
//
// act2Message should be signed with responder's static private key.
type Act2Message struct {
	nonce2    uint64
	challenge [sha256.Size]byte
	protocol2 string
}

// Act3Message is sent in the third handshake act by the initiator to the
// responder. It contains the challenge that has been recomputed by the
// initiator as a SHA256 of the concatenated bytes of `nonce1` and `nonce2`.
//
// act3Message should be signed with initiator's static private key.
type Act3Message struct {
	challenge [sha256.Size]byte
}

// InitiatorAct1 represents the state of the initiator in the first act of the
// handshake protocol.
type InitiatorAct1 struct {
	nonce1    uint64
	protocol1 string
}

// InitiateHandshake function allows to initiate a handshake by creating
// and initializing a state machine representing initiator in the first round
// of the handshake, ready to execute the protocol.
func InitiateHandshake(protocol string) (*InitiatorAct1, error) {
	nonce1, err := randomNonce()
	if err != nil {
		return nil, fmt.Errorf("could not initiate the handshake: [%v]", err)
	}

	return &InitiatorAct1{nonce1, protocol}, nil
}

// Message returns the message sent by initiator to the responder in the first
// act of the handshake protocol.
func (ia1 *InitiatorAct1) Message() *Act1Message {
	return &Act1Message{nonce1: ia1.nonce1, protocol1: ia1.protocol1}
}

// Next performs a state transition and returns initiator in a state ready to
// execute the second act of the handshake protocol.
func (ia1 *InitiatorAct1) Next() *InitiatorAct2 {
	return &InitiatorAct2{nonce1: ia1.nonce1, protocol1: ia1.protocol1}
}

// AnswerHandshake is used to initiate a responder as a result of receiving
// message from initiator in the first act of the handshake protocol.
// The returned responder is in a state ready to execute the second act of the
// handshake protocol.
// The function also validates if both parties run the same protocol.
func AnswerHandshake(message *Act1Message, protocol string) (*ResponderAct2, error) {
	if message.protocol1 != protocol {
		return nil, fmt.Errorf("unsupported protocol: [%v]", message.protocol1)
	}

	nonce1 := message.nonce1
	nonce2, err := randomNonce()
	if err != nil {
		return nil, fmt.Errorf("could not answer the handshake: [%v]", err)
	}
	challenge := hashToChallenge(nonce1, nonce2)

	return &ResponderAct2{nonce2, challenge, protocol}, nil
}

// InitiatorAct2 represents the state of the initiator in the second act of the
// handshake protocol.
type InitiatorAct2 struct {
	nonce1    uint64
	protocol1 string
}

// ResponderAct2 represents the state of the responder in the second act of the
// handshake protocol.
type ResponderAct2 struct {
	nonce2    uint64
	challenge [sha256.Size]byte
	protocol2 string
}

// Message returns the message sent by responder to the initiator in the second
// act of the handshake protocol.
func (ra2 *ResponderAct2) Message() *Act2Message {
	return &Act2Message{
		nonce2:    ra2.nonce2,
		challenge: ra2.challenge,
		protocol2: ra2.protocol2,
	}
}

// Next performs a state transition and returns responder in a state ready to
// execute the third act of the handshake protocol.
func (ra2 *ResponderAct2) Next() *ResponderAct3 {
	return &ResponderAct3{challenge: ra2.challenge}
}

// Next performs a state transition and returns initiator in a state ready to
// execute the third act of the handshake protocol.
//
// Function validates the challenge received from responder in the second act of
// the protocol. If the challenge is the same as expected one, new state of
// initiator is returned. Otherwise, function reports an error and handshake
// protocol should be immediately aborted.
//
// The function also validates if both parties run the same protocol.
func (ia2 *InitiatorAct2) Next(message *Act2Message) (*InitiatorAct3, error) {
	if message.protocol2 != ia2.protocol1 {
		return nil, fmt.Errorf("unsupported protocol: [%v]", message.protocol2)
	}

	expectedChallenge := hashToChallenge(ia2.nonce1, message.nonce2)
	if expectedChallenge != message.challenge {
		return nil, fmt.Errorf("unexpected responder's challenge")
	}

	return &InitiatorAct3{challenge: message.challenge}, nil
}

// InitiatorAct3 represents the state of the initiator in the third act of the
// handshake protocol.
type InitiatorAct3 struct {
	challenge [sha256.Size]byte
}

// ResponderAct3 represents the state of the responder in the third act of the
// handshake protocol.
type ResponderAct3 struct {
	challenge [sha256.Size]byte
}

// Message returns the message sent by initiator to the responder in the third
// act of the handshake protocol.
func (ia3 *InitiatorAct3) Message() *Act3Message {
	return &Act3Message{challenge: ia3.challenge}
}

// FinalizeHandshake is used in the third act of the handshake protocol to
// inform responder about a message sent by initiator. Responder validates
// the challenge in the message comparing it with the one expected.
// If both challenges are equal, handshake has completed successfully and
// function returns nil. Otherwise, if challenge is not as expected, function
// returns an error and it means the handshake protocol failed.
func (ra3 *ResponderAct3) FinalizeHandshake(message *Act3Message) error {
	if ra3.challenge != message.challenge {
		return errors.New("unexpected initiator's challenge")
	}

	return nil
}

// hashToChallenge computes a challenge as a SHA256 hash of the concatenated
// bytes of `nonce1` and `nonce2`.
func hashToChallenge(nonce1 uint64, nonce2 uint64) [sha256.Size]byte {
	var inputBytes [sha256.Size]byte
	binary.LittleEndian.PutUint64(inputBytes[0:], nonce1)
	binary.LittleEndian.PutUint64(inputBytes[8:], nonce2)

	return sha256.Sum256(inputBytes[:])
}

// randomNonce uses a cryptographically secure pseudorandom number generator
// to produce an 8-byte unsigned integer nonce.
func randomNonce() (uint64, error) {
	bytes := make([]byte, 8)
	_, err := crand.Read(bytes)
	if err != nil {
		return 0, fmt.Errorf("could not generate a new nonce: [%v]", err)
	}

	return binary.LittleEndian.Uint64(bytes), nil
}
