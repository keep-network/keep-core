package security

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math/rand"
)

// In the first act, initiator sends a message containing randomly generated
// `nonce1`, an 8-byte (64-bit) unsigned integer to the responder.
//
// act1Message should be signed with initiator's static private key.
type act1Message struct {
	nonce1 uint64
}

// In the second act, responder generates `nonce2`, an 8-byte unsigned integer
// and computes a `challenge` which is the result of calling a cryptographic
// hash function SHA256 on the concatenated bytes of `nonce1` and `nonce2`.
//
// act2Message should be signed with responder's static private key.
type act2Message struct {
	nonce2    uint64
	challenge [sha256.Size]byte
}

// In the third act, initiator recomputes the challenge from `nonce1` and
// `nonce2`; if it matches the challenge sent by the responder in the previous
// act, it answers with a message containing the `challenge`.
//
// act3Message should be signed with initiator's static private key.
type act3Message struct {
	challenge [sha256.Size]byte
}

type initiatorAct1 struct {
	nonce1 uint64
}

func initiateHandshake() *initiatorAct1 {
	return &initiatorAct1{
		nonce1: rand.Uint64(),
	}
}

func (ia1 *initiatorAct1) message() *act1Message {
	return &act1Message{
		nonce1: ia1.nonce1,
	}
}

func (ia1 *initiatorAct1) next() *initiatorAct2 {
	return &initiatorAct2{
		nonce1: ia1.nonce1,
	}
}

func answerHandshake(act1Msg *act1Message) *responderAct2 {
	nonce1 := act1Msg.nonce1
	nonce2 := rand.Uint64()
	challenge := hashToChallenge(nonce1, nonce2)

	return &responderAct2{nonce2, challenge}
}

type initiatorAct2 struct {
	nonce1 uint64
}

type responderAct2 struct {
	nonce2    uint64
	challenge [sha256.Size]byte
}

func (ra2 *responderAct2) message() *act2Message {
	return &act2Message{
		nonce2:    ra2.nonce2,
		challenge: ra2.challenge,
	}
}

func (ra2 *responderAct2) next() *responderAct3 {
	return &responderAct3{
		challenge: ra2.challenge,
	}
}

func (ia2 *initiatorAct2) next(act2Msg *act2Message) (*initiatorAct3, error) {
	expectedChallenge := hashToChallenge(ia2.nonce1, act2Msg.nonce2)
	if expectedChallenge != act2Msg.challenge {
		return nil, errors.New("unexpected responder's challenge")
	}

	return &initiatorAct3{
		challenge: act2Msg.challenge,
	}, nil
}

type responderAct3 struct {
	challenge [sha256.Size]byte
}

type initiatorAct3 struct {
	challenge [sha256.Size]byte
}

func (ia3 *initiatorAct3) message() *act3Message {
	return &act3Message{
		challenge: ia3.challenge,
	}
}

func (ra3 *responderAct3) finalizeHandshake(act3Msg *act3Message) error {
	if ra3.challenge != act3Msg.challenge {
		return errors.New("unexpected initiator's challenge")
	}

	return nil
}

func hashToChallenge(nonce1 uint64, nonce2 uint64) [sha256.Size]byte {
	var inputBytes [sha256.Size]byte
	binary.LittleEndian.PutUint64(inputBytes[0:], nonce1)
	binary.LittleEndian.PutUint64(inputBytes[8:], nonce2)

	return sha256.Sum256(inputBytes[:])
}
