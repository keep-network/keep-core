package net

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math/rand"
)

type connectionHandshake struct {
	nonce1    uint64
	nonce2    uint64
	challenge [challengeSize]byte
}

const (
	nonceSize     = 8 // uint64 size
	challengeSize = sha256.Size

	ActOneSize   = nonceSize
	ActTwoSize   = challengeSize + nonceSize
	ActThreeSize = challengeSize
)

func (ch *connectionHandshake) GenActOne() [ActOneSize]byte {
	ch.nonce1 = rand.Uint64()

	var actOneMessage [ActOneSize]byte
	binary.LittleEndian.PutUint64(actOneMessage[:], ch.nonce1)

	return actOneMessage
}

func (ch *connectionHandshake) RecvActOne(actOne [ActOneSize]byte) {
	ch.nonce1 = binary.LittleEndian.Uint64(actOne[:])
}

func (ch *connectionHandshake) GenActTwo() [ActTwoSize]byte {
	ch.nonce2 = rand.Uint64()
	ch.challenge = hashToChallenge(ch.nonce1, ch.nonce2)

	var actTwoMessage [ActTwoSize]byte
	copy(actTwoMessage[0:], ch.challenge[:])
	binary.LittleEndian.PutUint64(actTwoMessage[challengeSize:], ch.nonce2)

	return actTwoMessage
}

func (ch *connectionHandshake) RecvActTwo(actTwo [ActTwoSize]byte) error {
	ch.nonce2 = binary.LittleEndian.Uint64(actTwo[challengeSize:])
	ch.challenge = hashToChallenge(ch.nonce1, ch.nonce2)

	var responderChallenge [challengeSize]byte
	copy(responderChallenge[:], actTwo[0:challengeSize])

	if ch.challenge != responderChallenge {
		return errors.New("unexpected responder's challenge")
	}

	return nil
}

func (ch *connectionHandshake) GenActThree() [ActThreeSize]byte {
	return ch.challenge
}

func (ch *connectionHandshake) RecvActThree(actThree [ActThreeSize]byte) error {
	if ch.challenge != actThree {
		return errors.New("unexpected initiator's challenge")
	}

	return nil
}

func hashToChallenge(nonce1 uint64, nonce2 uint64) [challengeSize]byte {
	var inputBytes [challengeSize]byte
	binary.LittleEndian.PutUint64(inputBytes[0:], nonce1)
	binary.LittleEndian.PutUint64(inputBytes[nonceSize:], nonce2)

	return sha256.Sum256(inputBytes[:])
}
