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
type initiatorAct1 struct{}

type responderAct1 struct{}

type act1Message struct {
	nonce1 uint64
}

func (i *initiatorAct1) genActOne() (*initiatorAct2, *act1Message) {
	nonce1 := rand.Uint64()

	initiator := &initiatorAct2{nonce1: nonce1}
	message := &act1Message{nonce1: nonce1}

	return initiator, message
}

func (r *responderAct1) recvActOne(act1Msg *act1Message) *responderAct2 {
	return &responderAct2{
		nonce1: act1Msg.nonce1,
	}
}

// In the second act, responder generates `nonce2`, an 8-byte unsigned integer
// and computes a `challenge` which is the result of calling a cryptographic
// hash function SHA256 on the concatenated bytes of `nonce1` and `nonce2`.
//
// act2Message should be signed with responder's static private key.
type initiatorAct2 struct {
	initiatorAct1
	nonce1 uint64
}

type responderAct2 struct {
	responderAct1
	nonce1 uint64
}

type act2Message struct {
	nonce2    uint64
	challenge [sha256.Size]byte
}

func (r *responderAct2) genActTwo() (*responderAct3, *act2Message) {
	nonce2 := rand.Uint64()
	challenge := hashToChallenge(r.nonce1, nonce2)

	responder := &responderAct3{
		responderAct2: *r,
		nonce2:        nonce2,
		challenge:     challenge,
	}
	message := &act2Message{
		nonce2:    responder.nonce2,
		challenge: challenge,
	}

	return responder, message
}

func (i *initiatorAct2) recvActTwo(act2Msg *act2Message) (*initiatorAct3, error) {

	challenge := hashToChallenge(i.nonce1, act2Msg.nonce2)
	if challenge != act2Msg.challenge {
		return nil, errors.New("unexpected responder's challenge")
	}

	initiator := &initiatorAct3{
		initiatorAct2: *i,
		nonce2:        act2Msg.nonce2,
		challenge:     act2Msg.challenge,
	}

	return initiator, nil
}

// In the third act, initiator recomputes the challenge from `nonce1` and
// `nonce2`; if it matches the challenge sent by the responder in the previous
// act, it answers with a message containing the `challenge`.
//
// act3Message should be signed with initiator's static private key.
type initiatorAct3 struct {
	initiatorAct2

	nonce2    uint64
	challenge [sha256.Size]byte
}

type responderAct3 struct {
	responderAct2

	nonce2    uint64
	challenge [sha256.Size]byte
}

type act3Message struct {
	challenge [sha256.Size]byte
}

func (i *initiatorAct3) genActThree() *act3Message {
	return &act3Message{challenge: i.challenge}
}

func (r *responderAct3) recvActThree(act3Msg *act3Message) error {
	if r.challenge != act3Msg.challenge {
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
