package gjkr

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func TestPutEphemeralPubKeyEvidenceLog(t *testing.T) {
	var tests = map[string]struct {
		sender                      MemberID
		receiver                    MemberID
		modifyPubKeyMessageLogState func(
			sender, receiver MemberID,
			log *dkgEvidenceLog,
		) error
		expectedError error
	}{
		"EphemeralPubKeyMessage successfully stored for sender, receiver": {
			sender:   MemberID(uint32(1)),
			receiver: MemberID(uint32(2)),
			modifyPubKeyMessageLogState: func(
				sender, receiver MemberID,
				log *dkgEvidenceLog,
			) error {
				return nil
			},
			expectedError: nil,
		},
		"EphemeralPubKeyMessage already exists for sender, receiver": {
			sender:   MemberID(uint32(1)),
			receiver: MemberID(uint32(2)),
			modifyPubKeyMessageLogState: func(
				sender, receiver MemberID,
				log *dkgEvidenceLog,
			) error {
				msg := &EphemeralPublicKeyMessage{
					senderID:   sender,
					receiverID: receiver,
				}
				err := log.PutEphemeralMessage(msg)
				if err != nil {
					return err
				}
				return nil
			},
			expectedError: fmt.Errorf(
				"message exists for sender 1 and receiver 2",
			),
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// set up the initial state
			dkgEvidenceLog := NewDkgEvidenceLog()

			// modify the state of the log
			err := test.modifyPubKeyMessageLogState(
				test.sender,
				test.receiver,
				dkgEvidenceLog,
			)
			if err != nil {
				t.Fatal(err)
			}

			ephemeralKey, err := ephemeral.GenerateKeyPair()
			if err != nil {
				t.Fatal(err)
			}
			// simulate adding a message to the store
			message := &EphemeralPublicKeyMessage{
				senderID:           test.sender,
				receiverID:         test.receiver,
				ephemeralPublicKey: ephemeralKey.PublicKey,
			}
			err = dkgEvidenceLog.PutEphemeralMessage(message)

			if !reflect.DeepEqual(err, test.expectedError) {
				t.Fatalf(
					"\nexpected: %s\nactual:   %s\n",
					test.expectedError,
					err,
				)
			}
		})
	}
}

func TestGetEphemeralPubKeyEvidenceLog(t *testing.T) {
	var tests = map[string]struct {
		sender         MemberID
		receiver       MemberID
		expectedResult *EphemeralPublicKeyMessage
	}{
		"valid EphemeralPubKeyMessage returned for sender, receiver": {
			sender:   MemberID(uint32(1)),
			receiver: MemberID(uint32(2)),
			expectedResult: &EphemeralPublicKeyMessage{
				senderID:   MemberID(uint32(1)),
				receiverID: MemberID(uint32(2)),
			},
		},
		"no EphemeralPubKeyMessage for sender": {
			sender:         MemberID(uint32(1)),
			receiver:       MemberID(uint32(2)),
			expectedResult: nil,
		},
		"no EphemeralPubKeyMessage for receiver": {
			sender:         MemberID(uint32(1)),
			receiver:       MemberID(uint32(2)),
			expectedResult: nil,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// set up the initial state
			dkgEvidenceLog := NewDkgEvidenceLog()

			// simulate adding a message to the store
			if test.expectedResult != nil {
				if err := dkgEvidenceLog.PutEphemeralMessage(
					test.expectedResult,
				); err != nil {
					t.Fatal(err)
				}
			}

			result := dkgEvidenceLog.ephemeralPublicKeyMessage(
				test.sender, test.receiver,
			)

			if result != test.expectedResult {
				t.Fatalf(
					"\nexpected: %d\nactual:   %d\n",
					test.expectedResult,
					result,
				)
			}
		})
	}
}
