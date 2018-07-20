package tecdsa

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/paillier"
)

// In the second round, signer reveals values for which he committed to in the
// first round. At the beginning of 3rd round, messages are validated and
// combined together. Here, we test the validation process.
func TestSignRound1And2(t *testing.T) {
	var tests = map[string]struct {
		modifyRound1Messages func(msgs []*SignRound1Message) []*SignRound1Message
		modifyRound2Messages func(msgs []*SignRound2Message) []*SignRound2Message
		expectedError        error
	}{
		"positive validation": {
			expectedError: nil,
		},
		"negative validation - too few round 1 messages": {
			modifyRound1Messages: func(
				round1Messages []*SignRound1Message,
			) []*SignRound1Message {
				return []*SignRound1Message{round1Messages[0]}
			},
			expectedError: errors.New(
				"round 1 messages required from all group members; got 1, expected 10",
			),
		},
		"negative validation - too few round 2 messages": {
			modifyRound2Messages: func(
				round2Messages []*SignRound2Message,
			) []*SignRound2Message {
				return []*SignRound2Message{round2Messages[0]}
			},
			expectedError: errors.New(
				"round 2 messages required from all group members; got 1, expected 10",
			),
		},
		"negative validation - missing round 2 message for signer": {
			modifyRound1Messages: func(
				round1Messages []*SignRound1Message,
			) []*SignRound1Message {
				round1Messages[1].signerID = "evil"
				return round1Messages
			},
			expectedError: errors.New(
				"no matching round 2 message for signer with ID = evil",
			),
		},
		"negative validation - invalid ZKP": {
			modifyRound2Messages: func(
				round2Messages []*SignRound2Message,
			) []*SignRound2Message {
				round2Messages[2].randomFactorShare = &paillier.Cypher{
					C: big.NewInt(1337),
				}
				return round2Messages
			},
			expectedError: errors.New(
				"round 2 message rejected",
			),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			signers, err := initializeNewSignerGroup()
			if err != nil {
				t.Fatal(err)
			}

			round1Messages := make([]*SignRound1Message, len(signers))
			round2Messages := make([]*SignRound2Message, len(signers))
			round2Signers := make([]*Round2Signer, len(signers))

			for i, signer := range signers {
				round1Signer, round1Message, err := signer.SignRound1()
				if err != nil {
					t.Fatal(err)
				}

				round2Signer, round2Message, err := round1Signer.SignRound2()
				if err != nil {
					t.Fatal(err)
				}

				round1Messages[i] = round1Message
				round2Messages[i] = round2Message
				round2Signers[i] = round2Signer
			}

			if test.modifyRound1Messages != nil {
				round1Messages = test.modifyRound1Messages(round1Messages)
			}

			if test.modifyRound2Messages != nil {
				round2Messages = test.modifyRound2Messages(round2Messages)
			}

			_, _, err = round2Signers[0].combineMessages(
				round1Messages,
				round2Messages,
			)

			if !reflect.DeepEqual(test.expectedError, err) {
				t.Fatalf(
					"unexpected error\nexpected %v\nactual %v",
					test.expectedError,
					err,
				)
			}
		})
	}
}

// Crates and initializes a new group of `Signer`s with T-ECDSA key set and
// ready for signing.
func initializeNewSignerGroup() ([]*Signer, error) {
	localGroup, key, err := initializeNewLocalGroupWithFullKey()
	if err != nil {
		return nil, err
	}

	signers := make([]*Signer, len(localGroup))
	for i, localSigner := range localGroup {
		signers[i] = &Signer{
			dsaKey:     key,
			signerCore: localSigner.signerCore,
		}
	}

	return signers, nil
}
