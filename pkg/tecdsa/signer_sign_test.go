package tecdsa

import (
	"crypto/rand"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/paillier"
)

// In the second round, signer reveals values for which he committed to in the
// first round. After the second round, messages are validated and
// combined together. Here, we simulate the whole process.
func TestSignAndCombineRound1And2(t *testing.T) {
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
				round2Messages[2].secretKeyFactorShare.C = big.NewInt(1337)
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
			round1Signers := make([]*Round1Signer, len(signers))
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
				round1Signers[i] = round1Signer
				round2Signers[i] = round2Signer
			}

			if test.modifyRound1Messages != nil {
				round1Messages = test.modifyRound1Messages(round1Messages)
			}

			if test.modifyRound2Messages != nil {
				round2Messages = test.modifyRound2Messages(round2Messages)
			}

			paillierKey := round2Signers[0].paillierKey
			expectedSecretKeyFactor := round1Signers[0].encryptedSecretKeyFactorShare
			expectedSecretKeyMultiple := round1Signers[0].secretKeyMultipleShare
			for _, signer := range round1Signers[1:] {
				expectedSecretKeyFactor = paillierKey.Add(
					expectedSecretKeyFactor, signer.encryptedSecretKeyFactorShare,
				)
				expectedSecretKeyMultiple = paillierKey.Add(
					expectedSecretKeyMultiple, signer.secretKeyMultipleShare,
				)
			}

			secretKeyFactor, secretKeyMultiple, err :=
				round2Signers[0].CombineRound2Messages(
					round1Messages,
					round2Messages,
				)

			if !reflect.DeepEqual(test.expectedError, err) {
				t.Fatalf(
					"unexpected error\nexpected: %v\nactual: %v",
					test.expectedError,
					err,
				)
			}

			if test.expectedError == nil {
				if !reflect.DeepEqual(expectedSecretKeyFactor, secretKeyFactor) {
					t.Fatalf(
						"unexpected secret key factor\nexpected: %v\nactual: %v",
						expectedSecretKeyFactor,
						secretKeyFactor,
					)
				}

				if !reflect.DeepEqual(expectedSecretKeyMultiple, secretKeyMultiple) {
					t.Fatalf(
						"unexpected secret key multiple\nexpected: %v\nactual: %v",
						expectedSecretKeyMultiple,
						secretKeyMultiple,
					)
				}
			}
		})
	}
}

// In the fourth round, signer reveals values for which he committed to in the
// third round. After the fourth round, messages are validated and
// combined together. Here, we simulate the whole process.
func TestSignAndCombineRound3And4(t *testing.T) {
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
			round3Signers := make([]*Round3Signer, len(round2Signers))
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
				round3Signers[i] = round3Signer
				round4Signers[i] = round4Signer
			}

			paillierKey := round3Signers[0].paillierKey
			expectedSignatureUnmask := round3Signers[0].signatureUnmaskShare
			for _, signer := range round3Signers[1:] {
				expectedSignatureUnmask = paillierKey.Add(
					expectedSignatureUnmask, signer.signatureUnmaskShare,
				)
			}

			ellipticCurve := round3Signers[0].groupParameters.curve
			expectedSignatureFactorPublic := round3Signers[0].signatureFactorPublicShare
			for _, signer := range round3Signers[1:] {
				expectedSignatureFactorPublic = curve.NewPoint(ellipticCurve.Add(
					expectedSignatureFactorPublic.X,
					expectedSignatureFactorPublic.Y,
					signer.signatureFactorPublicShare.X,
					signer.signatureFactorPublicShare.Y,
				))
			}

			if test.modifyRound3Messages != nil {
				round3Messages = test.modifyRound3Messages(round3Messages)
			}

			if test.modifyRound4Messages != nil {
				round4Messages = test.modifyRound4Messages(round4Messages)
			}

			signatureUnmask, signatureFactorPublic, err :=
				round4Signers[0].CombineRound4Messages(
					round3Messages,
					round4Messages,
				)

			if !reflect.DeepEqual(test.expectedError, err) {
				t.Fatalf(
					"unexpected error\nexpected: %v\nactual: %v",
					test.expectedError,
					err,
				)
			}

			if test.expectedError == nil {
				if !reflect.DeepEqual(expectedSignatureUnmask, signatureUnmask) {
					t.Fatalf(
						"unexpected signature unmask\n expected: %v\n actual: %v",
						expectedSignatureUnmask,
						signatureUnmask,
					)
				}

				if !reflect.DeepEqual(
					expectedSignatureFactorPublic, signatureFactorPublic,
				) {
					t.Fatalf(
						"unexpected signature factor public\nexpected: %v\nactual: %v",
						expectedSignatureFactorPublic,
						signatureFactorPublic,
					)
				}
			}
		})
	}
}

// In the fifth round, signers jointly decrypt signature unmask as well as
// compute hash of the signature factor public parameter.
// Here we test the hash computation process. Threshold decryption is tested
// separately in another test.
func TestSignRound5(t *testing.T) {
	round4Signers, err := initializeNewRound4SignerGroup()
	if err != nil {
		t.Fatal(err)
	}

	signer := round4Signers[0]

	signatureUnmaskCypher, err := signer.paillierKey.Encrypt(
		big.NewInt(911), rand.Reader,
	)
	if err != nil {
		t.Fatal(err)
	}

	signatureFactorPublic := curve.NewPoint(
		publicParameters.curve.ScalarBaseMult(big.NewInt(411).Bytes()),
	)

	round5Signer, _, err := signer.SignRound5(
		signatureUnmaskCypher,
		signatureFactorPublic,
	)
	if err != nil {
		t.Fatal(err)
	}

	expectedSignatureFactorPublicHash := new(big.Int).Mod(
		signatureFactorPublic.X,
		publicParameters.curve.Params().N,
	)

	if round5Signer.signatureFactorPublicHash.Cmp(
		expectedSignatureFactorPublicHash,
	) != 0 {
		t.Fatalf(
			"unexpected signature random multiple public hash\nexpected: %v\nactual: %v",
			expectedSignatureFactorPublicHash,
			round5Signer.signatureFactorPublicHash,
		)
	}

}

// In the fifth round, signers jointly decrypt signature unmask as well as
// compute hash of the signature factor public parameter.
// After the fifth round, partial signature unmask decryptions are combined
// together. Here we test the decryption process.
func TestSignAndCombineRound5(t *testing.T) {
	signatureUnmask := big.NewInt(712)

	signatureFactorPublic := curve.NewPoint(
		publicParameters.curve.ScalarBaseMult(big.NewInt(411).Bytes()),
	)

	var tests = map[string]struct {
		modifyRound5Messages    func(msgs []*SignRound5Message) []*SignRound5Message
		expectedSignatureUnmask *big.Int
		expectedError           error
	}{
		"successful signature unmask decryption": {
			expectedSignatureUnmask: big.NewInt(712),
		},
		"negative validation - too few round 5 messages": {
			modifyRound5Messages: func(msgs []*SignRound5Message) []*SignRound5Message {
				return []*SignRound5Message{msgs[0], msgs[1]}
			},
			expectedError: errors.New(
				"round 5 messages required from all group members; got 2, expected 10",
			),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			round4Signers, err := initializeNewRound4SignerGroup()
			if err != nil {
				t.Fatal(err)
			}

			signatureUnmaskCypher, err := round4Signers[0].paillierKey.Encrypt(
				signatureUnmask, rand.Reader,
			)
			if err != nil {
				t.Fatal(err)
			}

			round5Signers := make([]*Round5Signer, len(round4Signers))
			round5Messages := make([]*SignRound5Message, len(round4Signers))

			for i, signer := range round4Signers {
				signer, message, err := signer.SignRound5(
					signatureUnmaskCypher,
					signatureFactorPublic,
				)
				if err != nil {
					t.Fatal(err)
				}

				round5Signers[i] = signer
				round5Messages[i] = message
			}

			if test.modifyRound5Messages != nil {
				round5Messages = test.modifyRound5Messages(round5Messages)
			}

			actualSignatureUnmask, err := round5Signers[0].CombineRound5Messages(
				round5Messages,
			)
			if !reflect.DeepEqual(test.expectedError, err) {
				t.Fatalf(
					"unexpected error\nexpected: %v\nactual: %v",
					test.expectedError,
					err,
				)
			}
			if test.expectedError == nil {
				if test.expectedSignatureUnmask.Cmp(actualSignatureUnmask) != 0 {
					t.Fatalf(
						"unexpected signature unmask\nexpected: %v\n actual: %v",
						test.expectedSignatureUnmask,
						actualSignatureUnmask,
					)
				}
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

// Crates and initializes a new group of `Round2Signer`s with T-ECDSA key and
// all other parameters set and ready for round 3 signing.
func initializeNewRound2SignerGroup() (
	round2Signers []*Round2Signer,
	secretKeyFactor *paillier.Cypher,
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
	secretKeyFactorPlaintext := big.NewInt(1337)

	secretKeyFactor, err = paillierKey.Encrypt(
		secretKeyFactorPlaintext, rand.Reader,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	secretKeyMultiple = paillierKey.Mul(
		signers[0].dsaKey.secretKey, secretKeyFactorPlaintext,
	)

	return
}

// Crates and initializes a new group of `Round4Signer`s with T-ECDSA key and
// all other parameters set and ready for round 5 signing.
func initializeNewRound4SignerGroup() ([]*Round4Signer, error) {
	signers, err := initializeNewSignerGroup()
	if err != nil {
		return nil, err
	}

	secretKeyFactor, err := signers[0].paillierKey.Encrypt(
		big.NewInt(7331), rand.Reader,
	)
	if err != nil {
		return nil, err
	}

	round4Signers := make([]*Round4Signer, len(signers))
	for i, signer := range signers {
		round4Signers[i] = &Round4Signer{
			Signer:          *signer,
			secretKeyFactor: secretKeyFactor,
		}
	}

	return round4Signers, nil
}
