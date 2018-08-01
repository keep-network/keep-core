package tecdsa

import (
	"crypto/rand"
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
				round2Messages[2].secretKeyRandomFactorShare.C = big.NewInt(1337)
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

			_, _, err = round2Signers[0].CombineRound2Messages(
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

func TestSignRound3And4(t *testing.T) {
	var tests = map[string]struct {
		modifyRound3Messages func(msgs []*SignRound3Message) []*SignRound3Message
		modifyRound4Messages func(msgs []*SignRound4Message) []*SignRound4Message
		expectedError        error
	}{
		"positive validation": {
			expectedError: nil,
		},
		"negative validation - too few round 3 messages": {
			modifyRound3Messages: func(
				round3Messages []*SignRound3Message,
			) []*SignRound3Message {
				return []*SignRound3Message{round3Messages[0]}
			},
			expectedError: errors.New(
				"round 3 messages required from all group members; got 1, expected 10",
			),
		},
		"negative validation - too few round 4 messages": {
			modifyRound4Messages: func(
				round4Messages []*SignRound4Message,
			) []*SignRound4Message {
				return []*SignRound4Message{round4Messages[0]}
			},
			expectedError: errors.New(
				"round 4 messages required from all group members; got 1, expected 10",
			),
		},
		"negative validation - missing round 3 message for signer": {
			modifyRound3Messages: func(
				round3Messages []*SignRound3Message,
			) []*SignRound3Message {
				round3Messages[3].signerID = "evil"
				return round3Messages
			},
			expectedError: errors.New(
				"no matching round 4 message for signer with ID = evil",
			),
		},
		"negative validation - invalid ZKP": {
			modifyRound4Messages: func(
				round4Messages []*SignRound4Message,
			) []*SignRound4Message {
				round4Messages[5].signatureUnmaskShare.C = big.NewInt(1337)
				return round4Messages
			},
			expectedError: errors.New(
				"round 4 message rejected",
			),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			round2Signers, secretKeyRandomFactor, secretKeyMultiple, err :=
				initializeNewRound2SignerGroup()
			if err != nil {
				t.Fatal(err)
			}

			round3Messages := make([]*SignRound3Message, len(round2Signers))
			round4Messages := make([]*SignRound4Message, len(round2Signers))
			round4Signers := make([]*Round4Signer, len(round2Signers))

			for i, signer := range round2Signers {
				round3Signer, round3Message, err := signer.SignRound3(
					secretKeyRandomFactor, secretKeyMultiple,
				)
				if err != nil {
					t.Fatal(err)
				}

				round4Signer, round4Message, err := round3Signer.SignRound4()
				if err != nil {
					t.Fatal(err)
				}

				round3Messages[i] = round3Message
				round4Messages[i] = round4Message
				round4Signers[i] = round4Signer
			}

			if test.modifyRound3Messages != nil {
				round3Messages = test.modifyRound3Messages(round3Messages)
			}

			if test.modifyRound4Messages != nil {
				round4Messages = test.modifyRound4Messages(round4Messages)
			}

			_, _, err = round4Signers[0].CombineRound4Messages(
				round3Messages,
				round4Messages,
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

func initializeNewRound2SignerGroup() (
	round2Signers []*Round2Signer,
	secretKeyRandomFactor *paillier.Cypher,
	secretKeyMultiple *paillier.Cypher,
	err error,
) {
	signers, err := initializeNewSignerGroup()
	if err != nil {
		return nil, nil, nil, err
	}

	round2Signers = make([]*Round2Signer, len(signers))
	for i, signer := range signers {
		round2Signers[i] = &Round2Signer{
			Signer: *signer,
		}
	}

	paillierKey := signers[0].paillierKey

	secretKeyRandomFactor, err = paillierKey.Encrypt(
		big.NewInt(1337), rand.Reader,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	secretKeyMultiple = paillierKey.Mul(
		signers[0].dsaKey.secretKey, big.NewInt(1337),
	)

	return
}
