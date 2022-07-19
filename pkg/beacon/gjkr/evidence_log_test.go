package gjkr

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/group"
)

func TestPutEphemeralPubKeyEvidenceLog(t *testing.T) {
	var tests = map[string]struct {
		sender                      group.MemberIndex
		modifyPubKeyMessageLogState func(
			sender group.MemberIndex,
			log *dkgEvidenceLog,
		) error
		expectedError error
	}{
		"EphemeralPubKeyMessage successfully stored for sender": {
			sender: group.MemberIndex(1),
			modifyPubKeyMessageLogState: func(
				sender group.MemberIndex,
				log *dkgEvidenceLog,
			) error {
				return nil
			},
			expectedError: nil,
		},
		"EphemeralPubKeyMessage already exists for sender": {
			sender: group.MemberIndex(1),
			modifyPubKeyMessageLogState: func(
				sender group.MemberIndex,
				log *dkgEvidenceLog,
			) error {
				msg := &EphemeralPublicKeyMessage{
					senderID: sender,
				}
				err := log.PutEphemeralMessage(msg)
				if err != nil {
					return err
				}
				return nil
			},
			expectedError: fmt.Errorf(
				"message exists for sender 1",
			),
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// set up the initial state
			dkgEvidenceLog := newDkgEvidenceLog()

			// modify the state of the log
			err := test.modifyPubKeyMessageLogState(
				test.sender,
				dkgEvidenceLog,
			)
			if err != nil {
				t.Fatal(err)
			}

			// simulate adding a message to the store
			message := &EphemeralPublicKeyMessage{
				senderID: test.sender,
			}
			err = dkgEvidenceLog.PutEphemeralMessage(message)

			if !reflect.DeepEqual(err, test.expectedError) {
				t.Fatalf(
					"\nexpected: %v\nactual:   %v",
					test.expectedError,
					err,
				)
			}
		})
	}
}

func TestGetEphemeralPubKeyEvidenceLog(t *testing.T) {
	var tests = map[string]struct {
		sender         group.MemberIndex
		expectedResult *EphemeralPublicKeyMessage
	}{
		"valid EphemeralPubKeyMessage returned for sender": {
			sender: group.MemberIndex(uint32(1)),
			expectedResult: &EphemeralPublicKeyMessage{
				senderID: group.MemberIndex(uint32(1)),
			},
		},
		"no EphemeralPubKeyMessage for sender": {
			sender:         group.MemberIndex(uint32(1)),
			expectedResult: nil,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// set up the initial state
			dkgEvidenceLog := newDkgEvidenceLog()

			// simulate adding a message to the store
			if test.expectedResult != nil {
				if err := dkgEvidenceLog.PutEphemeralMessage(
					test.expectedResult,
				); err != nil {
					t.Fatal(err)
				}
			}

			result := dkgEvidenceLog.ephemeralPublicKeyMessage(test.sender)

			if result != test.expectedResult {
				t.Fatalf(
					"\nexpected: %v\nactual:   %v\n",
					test.expectedResult,
					result,
				)
			}
		})
	}
}
