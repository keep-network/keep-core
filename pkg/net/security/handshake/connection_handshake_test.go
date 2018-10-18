package handshake

import (
	"errors"
	"math/rand"
	"reflect"
	"testing"
)

func TestInitiateHanshakeWithUniqueNonce(t *testing.T) {
	initiator1 := initiateHandshake()
	initiator2 := initiateHandshake()

	if initiator1.nonce1 == 0 {
		t.Fatalf("Nonce not initialized")
	}

	if initiator1.nonce1 == initiator2.nonce1 {
		t.Fatalf("nonce1 should be unique for each handshake")
	}

}

func TestAnswerHandshakeWithChallenge(t *testing.T) {
	//
	// Act 1
	//

	// initiator station
	initiator := initiateHandshake()
	act1Msg := initiator.message()

	// responder station
	responder := answerHandshake(act1Msg)

	//
	// Act 2
	//

	// responder station
	act2Msg := responder.message()

	// assert if challenge sent by responder in Act 2 is properly
	// created from `nonce1` and `nonce2`
	expectedChallenge := hashToChallenge(initiator.nonce1, responder.nonce2)
	if act2Msg.challenge != expectedChallenge {
		t.Fatalf(
			"Unexpected challenge\nExpected: %v\nActual: %v",
			expectedChallenge,
			act2Msg.challenge,
		)
	}
}

func TestRepeatChallengeToFinalize(t *testing.T) {
	nonce1 := rand.Uint64()
	nonce2 := rand.Uint64()
	expectedChallenge := hashToChallenge(nonce1, nonce2)

	//
	// Act 2
	//

	// responder station
	act2Msg := &act2Message{nonce2, expectedChallenge}

	// initiator station
	initiatorAct2 := &initiatorAct2{nonce1}
	initiatorAct3, err := initiatorAct2.next(act2Msg)
	if err != nil {
		t.Fatal(err)
	}

	//
	// Act 3
	//

	// initiator station
	act3Msg := initiatorAct3.message()

	// assert if challenge sent by initiator in Act3 is the
	// same challenge as the one received from responder in Act2
	if act3Msg.challenge != expectedChallenge {
		t.Fatalf(
			"Unexpected challenge\nExpected: %v\nActual: %v",
			expectedChallenge,
			act3Msg.challenge,
		)
	}
}

func TestFailAct2ForInvalidChallenge(t *testing.T) {
	nonce1 := rand.Uint64()
	nonce2 := rand.Uint64()

	//
	// Act 2
	//

	// responder station
	invalidChallenge := [32]byte{0xff, 0xfa}
	act2Msg := &act2Message{nonce2, invalidChallenge}

	// initiator station
	initiatorAct2 := &initiatorAct2{nonce1}
	_, err := initiatorAct2.next(act2Msg)

	// assert if initiator detects invalid challenge sent by responder
	expectedError := errors.New("unexpected responder's challenge")
	if !reflect.DeepEqual(err, expectedError) {
		t.Fatalf(
			"Unexpected error\nExpected: %v\nActual: %v",
			expectedError,
			err,
		)
	}
}

func TestFailAct3ForInvalidChallenge(t *testing.T) {
	expectedChallenge := hashToChallenge(rand.Uint64(), rand.Uint64())
	responderAct3 := &responderAct3{expectedChallenge}

	invalidChallenge := hashToChallenge(rand.Uint64(), rand.Uint64())
	initiatorAct3 := &initiatorAct3{invalidChallenge}

	//
	// Act 3
	//
	act3Msg := initiatorAct3.message()
	err := responderAct3.finalizeHandshake(act3Msg)

	// assert if responder detects invalid challenge sent by initiator
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
	//
	// Act 1
	//

	// initiator station
	initiatorAct1 := initiateHandshake()
	act1Message := initiatorAct1.message()
	initiatorAct2 := initiatorAct1.next()

	// responder station
	responderAct2 := answerHandshake(act1Message)

	//
	// Act 2
	//

	// responder station
	act2Message := responderAct2.message()
	responderAct3 := responderAct2.next()

	// initiator station
	initiatorAct3, err := initiatorAct2.next(act2Message)
	if err != nil {
		t.Fatal(err)
	}

	//
	// Act 3
	//

	// initiator station
	act3Message := initiatorAct3.message()

	// responder station
	err = responderAct3.finalizeHandshake(act3Message)
	if err != nil {
		t.Fatal(err)
	}
}
