package net

import (
	"errors"
	"reflect"
	"testing"
)

func TestGenAndRecvActOne(t *testing.T) {
	initiator := &connectionHandshake{}
	responder := &connectionHandshake{}

	actOneMsg := initiator.GenActOne()
	responder.RecvActOne(actOneMsg)

	if initiator.nonce1 != responder.nonce1 {
		t.Fatalf(
			"Unexpected nonce 1\nExpected: %v\nActual: %v",
			initiator.nonce1,
			responder.nonce1,
		)
	}
}

func TestGenAndRecvActTwo(t *testing.T) {
	initiator := &connectionHandshake{}
	responder := &connectionHandshake{}

	initiator.nonce1 = uint64(1337)
	responder.nonce1 = initiator.nonce1

	actTwoMsg := responder.GenActTwo()
	err := initiator.RecvActTwo(actTwoMsg)
	if err != nil {
		t.Fatal(err)
	}

	if initiator.nonce2 != responder.nonce2 {
		t.Fatalf(
			"Unexpected nonce 2\nExpected: %v\nActual: %v",
			initiator.nonce1,
			responder.nonce1,
		)
	}

	if initiator.challenge != responder.challenge {
		t.Fatalf(
			"Unexpected challenge\nExpected: %v\nActual: %v",
			initiator.nonce1,
			responder.nonce1,
		)
	}
}

func TestRecvActTwoUnexpectedChallenge(t *testing.T) {
	initiator := &connectionHandshake{}
	responder := &connectionHandshake{}

	initiator.nonce1 = uint64(1337)
	responder.nonce1 = initiator.nonce1

	actTwoMsg := responder.GenActTwo()
	actTwoMsg[2] = 0xff
	err := initiator.RecvActTwo(actTwoMsg)

	expectedError := errors.New("unexpected responder's challenge")
	if !reflect.DeepEqual(err, expectedError) {
		t.Fatalf(
			"Unexpected error\nExpected: %v\nActual: %v",
			expectedError,
			err,
		)
	}
}

func TestGenAndRecvActThree(t *testing.T) {
	initiator := &connectionHandshake{}
	responder := &connectionHandshake{}

	initiator.nonce1 = uint64(1337)
	responder.nonce1 = initiator.nonce1

	initiator.nonce2 = uint64(1410)
	responder.nonce2 = initiator.nonce2
	initiator.challenge = hashToChallenge(responder.nonce1, responder.nonce2)
	responder.challenge = initiator.challenge

	actThreeMsg := initiator.GenActThree()
	err := responder.RecvActThree(actThreeMsg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenAndRecvActThreeUnexpectedChallenge(t *testing.T) {
	initiator := &connectionHandshake{}
	responder := &connectionHandshake{}

	initiator.nonce1 = uint64(1337)
	responder.nonce1 = initiator.nonce1

	initiator.nonce2 = uint64(1410)
	responder.nonce2 = initiator.nonce2
	initiator.challenge = hashToChallenge(responder.nonce1, responder.nonce2)
	responder.challenge = [32]byte{0xff, 0xfa}

	actThreeMsg := initiator.GenActThree()
	err := responder.RecvActThree(actThreeMsg)

	expectedError := errors.New("unexpected initiator's challenge")
	if !reflect.DeepEqual(err, expectedError) {
		t.Fatalf(
			"Unexpected error\nExpected: %v\nActual: %v",
			expectedError,
			err,
		)
	}
}

func TestFullHandshake(t *testing.T) {
	initiator := &connectionHandshake{}
	responder := &connectionHandshake{}

	actOneMsg := initiator.GenActOne()
	responder.RecvActOne(actOneMsg)

	actTwoMsg := responder.GenActTwo()
	err := initiator.RecvActTwo(actTwoMsg)
	if err != nil {
		t.Fatal(err)
	}

	actThreeMsg := initiator.GenActThree()
	err = responder.RecvActThree(actThreeMsg)
	if err != nil {
		t.Fatal(err)
	}

	if initiator.nonce1 != responder.nonce1 {
		t.Fatalf(
			"Unexpected nonce 1\nExpected: %v\nActual: %v",
			initiator.nonce1,
			responder.nonce1,
		)
	}

	if initiator.nonce2 != responder.nonce2 {
		t.Fatalf(
			"Unexpected nonce 2\nExpected: %v\nActual: %v",
			initiator.nonce1,
			responder.nonce1,
		)
	}

	if initiator.challenge != responder.challenge {
		t.Fatalf(
			"Unexpected challenge\nExpected: %v\nActual: %v",
			initiator.nonce1,
			responder.nonce1,
		)
	}
}
