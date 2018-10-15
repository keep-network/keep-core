package security

import (
	"errors"
	"reflect"
	"testing"
)

func TestGenAndRecvActOne(t *testing.T) {
	initiatorAct1 := &initiatorAct1{}
	responderAct1 := responderAct2{}

	// execute 1st act
	// initiator -> responder
	initiatorAct2, act1Msg := initiatorAct1.genActOne()
	responderAct2 := responderAct1.recvActOne(act1Msg)

	if initiatorAct2.nonce1 != responderAct2.nonce1 {
		t.Fatalf(
			"Unexpected nonce 1\nExpected: %v\nActual: %v",
			initiatorAct2.nonce1,
			responderAct2.nonce1,
		)
	}
}

func TestGenAndRectActTwo(t *testing.T) {
	nonce1 := uint64(1337)

	initiatorAct2 := &initiatorAct2{nonce1: nonce1}
	responderAct2 := &responderAct2{nonce1: nonce1}

	// execute 2nd act
	// responder -> initiator
	responderAct3, act2Msg := responderAct2.genActTwo()
	initiatorAct3, err := initiatorAct2.recvActTwo(act2Msg)
	if err != nil {
		t.Fatal(err)
	}

	if initiatorAct3.nonce2 != responderAct3.nonce2 {
		t.Fatalf(
			"Unexpected nonce 2\nExpected: %v\nActual: %v",
			initiatorAct3.nonce1,
			responderAct3.nonce1,
		)
	}

	if initiatorAct3.challenge != responderAct3.challenge {
		t.Fatalf(
			"Unexpected challenge\nExpected: %v\nActual: %v",
			initiatorAct3.nonce1,
			responderAct3.nonce1,
		)
	}
}

func TestRecvActTwoUnexpectedChallenge(t *testing.T) {
	nonce1 := uint64(1337)

	initiatorAct2 := &initiatorAct2{nonce1: nonce1}
	responderAct2 := &responderAct2{nonce1: nonce1}

	// execute 2nd act
	// responder -> initiator
	_, act2Msg := responderAct2.genActTwo()
	act2Msg.challenge[2] = 0xff
	_, err := initiatorAct2.recvActTwo(act2Msg)

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
	nonce1 := uint64(1337)
	nonce2 := uint64(1410)

	initiatorAct3 := &initiatorAct3{
		initiatorAct2: initiatorAct2{nonce1: nonce1},
		nonce2:        nonce2,
		challenge:     hashToChallenge(nonce1, nonce2),
	}

	responderAct3 := &responderAct3{
		responderAct2: responderAct2{nonce1: nonce1},
		nonce2:        nonce2,
		challenge:     hashToChallenge(nonce1, nonce2),
	}

	// execute 3rd act
	// initiator -> responder
	act3Msg := initiatorAct3.genActThree()
	err := responderAct3.recvActThree(act3Msg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenAndRecvActThreeUnexpectedChallenge(t *testing.T) {
	nonce1 := uint64(1337)
	nonce2 := uint64(1410)

	initiatorAct3 := &initiatorAct3{
		initiatorAct2: initiatorAct2{nonce1: nonce1},
		nonce2:        nonce2,
		challenge:     hashToChallenge(nonce1, nonce2),
	}

	responderAct3 := &responderAct3{
		responderAct2: responderAct2{nonce1: nonce1},
		nonce2:        nonce2,
		challenge:     hashToChallenge(nonce1, nonce2),
	}

	// execute 3rd act
	// initiator -> responder
	act3Msg := initiatorAct3.genActThree()
	act3Msg.challenge = [32]byte{0xff, 0xfa}
	err := responderAct3.recvActThree(act3Msg)

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
	initiatorAct1 := &initiatorAct1{}
	responderAct1 := responderAct2{}

	// execute 1st act
	// initiator -> responder
	initiatorAct2, act1Msg := initiatorAct1.genActOne()
	responderAct2 := responderAct1.recvActOne(act1Msg)

	// execute 2nd act
	// responder -> initiator
	responderAct3, act2Msg := responderAct2.genActTwo()
	initiatorAct3, err := initiatorAct2.recvActTwo(act2Msg)
	if err != nil {
		t.Fatal(err)
	}

	// execute 3rd act
	// initiator -> responder
	act3Msg := initiatorAct3.genActThree()
	err = responderAct3.recvActThree(act3Msg)
	if err != nil {
		t.Fatal(err)
	}

	if initiatorAct3.nonce1 != responderAct3.nonce1 {
		t.Fatalf(
			"Unexpected nonce 1\nExpected: %v\nActual: %v",
			initiatorAct3.nonce1,
			responderAct3.nonce1,
		)
	}

	if initiatorAct3.nonce2 != responderAct3.nonce2 {
		t.Fatalf(
			"Unexpected nonce 2\nExpected: %v\nActual: %v",
			initiatorAct3.nonce1,
			responderAct3.nonce1,
		)
	}

	if initiatorAct3.challenge != responderAct3.challenge {
		t.Fatalf(
			"Unexpected challenge\nExpected: %v\nActual: %v",
			initiatorAct3.nonce1,
			responderAct3.nonce1,
		)
	}
}
