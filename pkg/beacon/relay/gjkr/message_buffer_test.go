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
		expectedError func(sender, receiver MemberID) error
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
			expectedError: func(sender, receiver MemberID) error {
				return nil
			},
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
				if err := log.PutEphemeralMessage(msg); err != nil {
					return err
				}
				return nil
			},
			expectedError: func(sender, receiver MemberID) error {
				return fmt.Errorf(
					"message exists for sender %v and receiver %v",
					sender,
					receiver,
				)
			},
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

			expectedError := test.expectedError(
				test.sender, test.receiver,
			)
			if !reflect.DeepEqual(err, expectedError) {
				t.Fatalf(
					"\nexpected: %s\nactual:   %s\n",
					expectedError,
					err,
				)
			}
		})
	}
}
