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
				"round 1 messages required from all group peer members; got 1, expected 9",
			),
		},
		"negative validation - too few round 2 messages": {
			modifyRound2Messages: func(
				round2Messages []*SignRound2Message,
			) []*SignRound2Message {
				return []*SignRound2Message{round2Messages[0]}
			},
			expectedError: errors.New(
				"round 2 messages required from all group peer members; got 1, expected 9",
			),
		},
		"negative validation - missing round 2 message for signer": {
			modifyRound1Messages: func(
				round1Messages []*SignRound1Message,
			) []*SignRound1Message {
				round1Messages[1].senderID = "evil"
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
			signers, _, err := initializeNewSignerGroup()
			if err != nil {
				t.Fatal(err)
			}

			round1Signers := make([]*Round1Signer, len(signers))
			round2Signers := make([]*Round2Signer, len(signers))
			var round1Messages []*SignRound1Message
			var round2Messages []*SignRound2Message

			for i, signer := range signers {
				round1Signer, signersRound1Messages, err := signer.SignRound1()
				if err != nil {
					t.Fatal(err)
				}

				round2Signer, signersRound2Messages, err := round1Signer.SignRound2()
				if err != nil {
					t.Fatal(err)
				}

				round1Signers[i] = round1Signer
				round2Signers[i] = round2Signer

				round1Messages = append(round1Messages, signersRound1Messages...)
				round2Messages = append(round2Messages, signersRound2Messages...)

			}

			receiver := round2Signers[0]

			round1Messages = signRound1MessagesForReceiver(round1Messages, receiver.ID)
			round2Messages = signRound2MessagesForReceiver(round2Messages, receiver.ID)

			if test.modifyRound1Messages != nil {
				round1Messages = test.modifyRound1Messages(round1Messages)

			}

			if test.modifyRound2Messages != nil {
				round2Messages = test.modifyRound2Messages(round2Messages)

			}

			expectedSecretKeyFactor := receiver.encryptedSecretKeyFactorShare
			expectedSecretKeyMultiple := receiver.secretKeyMultipleShare
			for _, signer := range round1Signers[1:] {
				expectedSecretKeyFactor = receiver.paillierKey.Add(
					expectedSecretKeyFactor, signer.encryptedSecretKeyFactorShare,
				)
				expectedSecretKeyMultiple = receiver.paillierKey.Add(
					expectedSecretKeyMultiple, signer.secretKeyMultipleShare,
				)
			}

			secretKeyFactor, secretKeyMultiple, err :=
				receiver.CombineRound2Messages(
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
				"round 3 messages required from all group peer members; got 1, expected 9",
			),
		},
		"negative validation - too few round 4 messages": {
			modifyRound4Messages: func(
				round4Messages []*SignRound4Message,
			) []*SignRound4Message {
				return []*SignRound4Message{round4Messages[0]}
			},
			expectedError: errors.New(
				"round 4 messages required from all group peer members; got 1, expected 9",
			),
		},
		"negative validation - missing round 3 message for signer": {
			modifyRound3Messages: func(
				round3Messages []*SignRound3Message,
			) []*SignRound3Message {
				round3Messages[3].senderID = "evil"
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
			round2Signers, _, secretKeyRandomFactor, secretKeyMultiple, err :=
				initializeNewRound2SignerGroup()
			if err != nil {
				t.Fatal(err)
			}

			round3Signers := make([]*Round3Signer, len(round2Signers))
			round4Signers := make([]*Round4Signer, len(round2Signers))
			var round3Messages []*SignRound3Message
			var round4Messages []*SignRound4Message

			for i, signer := range round2Signers {
				round3Signer, signersRound3Messages, err := signer.SignRound3(
					secretKeyRandomFactor, secretKeyMultiple,
				)
				if err != nil {
					t.Fatal(err)
				}

				round4Signer, signersRound4Messages, err := round3Signer.SignRound4()
				if err != nil {
					t.Fatal(err)
				}

				round3Signers[i] = round3Signer
				round4Signers[i] = round4Signer

				round3Messages = append(round3Messages, signersRound3Messages...)
				round4Messages = append(round4Messages, signersRound4Messages...)
			}

			receiver := round4Signers[0]

			round3Messages = signRound3MessagesForReceiver(round3Messages, receiver.ID)
			round4Messages = signRound4MessagesForReceiver(round4Messages, receiver.ID)

			expectedSignatureUnmask := receiver.signatureUnmaskShare
			for _, signer := range round3Signers[1:] {
				expectedSignatureUnmask = signer.paillierKey.Add(
					expectedSignatureUnmask, signer.signatureUnmaskShare,
				)
			}

			ellipticCurve := receiver.publicParameters.Curve
			expectedSignatureFactorPublic := receiver.signatureFactorPublicShare
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
				receiver.CombineRound4Messages(
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
	round4Signers, parameters, err := initializeNewRound4SignerGroup()
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
		parameters.Curve.ScalarBaseMult(big.NewInt(411).Bytes()),
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
		parameters.curveCardinality(),
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
			round4Signers, parameters, err := initializeNewRound4SignerGroup()
			if err != nil {
				t.Fatal(err)
			}

			signatureFactorPublic := curve.NewPoint(
				parameters.Curve.ScalarBaseMult(big.NewInt(411).Bytes()),
			)

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

// In the sixth round, each signer evaluates a final signature in an encrypted
// form using the parameters evaluated so far.
// Partial decryptions are combined together in order to present the signature
// in a decrypted form.
func TestSignAndCombineRound6(t *testing.T) {
	paillierKeys, publicParameters, zkpParameters, signerGroup, err := readTestParameters()
	if err != nil {
		t.Fatal(err)
	}

	curveCardinality := publicParameters.curveCardinality() // q

	secretKey := big.NewInt(211) // x = 211
	publicKey := curve.NewPoint( // y = g^x
		publicParameters.Curve.ScalarBaseMult(secretKey.Bytes()),
	)

	encryptedSecretKey, err := paillierKeys[0].Encrypt(secretKey, rand.Reader)
	if err != nil {
		t.Fatalf("paillier encryption failed [%v]", err)
	}

	ecdsaKey := &ThresholdEcdsaKey{encryptedSecretKey, publicKey}

	secretKeyFactor := big.NewInt(314) // ρ = 314

	// u = E(ρ)
	encryptedSecretKeyFactor, err := paillierKeys[0].Encrypt(
		secretKeyFactor, rand.Reader,
	)
	if err != nil {
		t.Fatalf("paillier encryption failed [%v]", err)
	}

	// v = E(ρx)
	secretKeyMultiple, err := paillierKeys[0].Encrypt(
		new(big.Int).Mul(secretKey, secretKeyFactor),
		rand.Reader,
	)
	if err != nil {
		t.Fatalf("paillier encryption failed [%v]", err)
	}

	signatureFactorSecret := big.NewInt(708) // k = 708
	signatureFactorPublic := curve.NewPoint( // R = g^k
		publicParameters.Curve.ScalarBaseMult(signatureFactorSecret.Bytes()),
	)
	signatureFactorMask := big.NewInt(9) // c = 9

	// TDec(w) = kρ + cq
	signatureUnmask := new(big.Int).Add(
		new(big.Int).Mul(signatureFactorSecret, secretKeyFactor),
		new(big.Int).Mul(signatureFactorMask, curveCardinality),
	)

	// r = H'(R)
	signatureFactorPublicHash := new(big.Int).Mod(
		signatureFactorPublic.X,
		curveCardinality,
	)

	signers := make([]*Round5Signer, len(paillierKeys))
	for i := 0; i < len(signers); i++ {
		signers[i] = &Round5Signer{
			Signer: *NewLocalSigner(
				&paillierKeys[i], publicParameters, zkpParameters, signerGroup,
			).WithEcdsaKey(ecdsaKey),

			secretKeyFactor:           encryptedSecretKeyFactor,
			secretKeyMultiple:         secretKeyMultiple,
			signatureFactorPublic:     signatureFactorPublic,
			signatureFactorPublicHash: signatureFactorPublicHash,
		}
	}

	messageHash := make([]byte, 32) // m
	_, err = rand.Read(messageHash)
	if err != nil {
		t.Fatal(err)
	}

	round6Messages := make([]*SignRound6Message, len(signers))
	for i, signer := range signers {
		round6Messages[i], err = signer.SignRound6(signatureUnmask, messageHash)
		if err != nil {
			t.Fatalf("could not execute round 6 of signing [%v]", err)
		}
	}

	signature, err := signers[0].CombineRound6Messages(round6Messages)
	if err != nil {
		t.Fatalf("could not combine round 6 messages [%v]", err)
	}

	// s = k^{-1} (m + xr)
	expectedS := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).ModInverse(signatureFactorSecret, curveCardinality),
			new(big.Int).Add(
				new(big.Int).SetBytes(messageHash[:]),
				new(big.Int).Mul(secretKey, signatureFactorPublicHash),
			),
		),
		curveCardinality,
	)
	if expectedS.Cmp(publicParameters.halfCurveCardinality()) == 1 {
		expectedS = new(big.Int).Sub(curveCardinality, expectedS)
	}

	if signature.S.Cmp(expectedS) != 0 {
		t.Errorf("Unexpected S\nExpected: %v\nActual: %v", expectedS, signature.S)
	}
}

// Crates and initializes a new group of `Signer`s with T-ECDSA key set and
// ready for signing.
func initializeNewSignerGroup() ([]*Signer, *PublicParameters, error) {
	localGroup, publicParameters, key, err := initializeNewLocalGroupWithFullKey()
	if err != nil {
		return nil, nil, err
	}

	signers := make([]*Signer, len(localGroup))
	for i, localSigner := range localGroup {
		signers[i] = &Signer{
			ecdsaKey:   key,
			signerCore: localSigner.signerCore,
		}
	}

	return signers, publicParameters, nil
}

// Crates and initializes a new group of `Round2Signer`s with T-ECDSA key and
// all other parameters set and ready for round 3 signing.
func initializeNewRound2SignerGroup() (
	round2Signers []*Round2Signer,
	publicParameters *PublicParameters,
	secretKeyFactor *paillier.Cypher,
	secretKeyMultiple *paillier.Cypher,
	err error,
) {
	signers, publicParameters, err := initializeNewSignerGroup()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	round2Signers = make([]*Round2Signer, len(signers))
	for i, signer := range signers {
		round2Signers[i] = &Round2Signer{
			&Round1Signer{
				Signer: *signer,
			},
		}
	}

	paillierKey := signers[0].paillierKey
	secretKeyFactorPlaintext := big.NewInt(1337)

	secretKeyFactor, err = paillierKey.Encrypt(
		secretKeyFactorPlaintext, rand.Reader,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	secretKeyMultiple = paillierKey.Mul(
		signers[0].ecdsaKey.secretKey, secretKeyFactorPlaintext,
	)

	return
}

// Crates and initializes a new group of `Round4Signer`s with T-ECDSA key and
// all other parameters set and ready for round 5 signing.
func initializeNewRound4SignerGroup() (
	[]*Round4Signer,
	*PublicParameters,
	error,
) {
	signers, parameters, err := initializeNewSignerGroup()
	if err != nil {
		return nil, nil, err
	}

	secretKeyFactor, err := signers[0].paillierKey.Encrypt(
		big.NewInt(7331), rand.Reader,
	)
	if err != nil {
		return nil, nil, err
	}

	round4Signers := make([]*Round4Signer, len(signers))
	for i, signer := range signers {
		round4Signers[i] = &Round4Signer{
			Round3Signer: &Round3Signer{
				Signer:          *signer,
				secretKeyFactor: secretKeyFactor,
			},
		}
	}

	return round4Signers, parameters, nil
}

func signRound1MessagesForReceiver(
	messages []*SignRound1Message,
	receiverID string,
) []*SignRound1Message {
	filtered := make([]*SignRound1Message, 0)
	for _, message := range messages {
		if message.receiverID == receiverID {
			filtered = append(filtered, message)
		}
	}
	return filtered
}

func signRound2MessagesForReceiver(
	messages []*SignRound2Message,
	receiverID string,
) []*SignRound2Message {
	filtered := make([]*SignRound2Message, 0)
	for _, message := range messages {
		if message.receiverID == receiverID {
			filtered = append(filtered, message)
		}
	}
	return filtered
}

func signRound3MessagesForReceiver(
	messages []*SignRound3Message,
	receiverID string,
) []*SignRound3Message {
	filtered := make([]*SignRound3Message, 0)
	for _, message := range messages {
		if message.receiverID == receiverID {
			filtered = append(filtered, message)
		}
	}
	return filtered
}

func signRound4MessagesForReceiver(
	messages []*SignRound4Message,
	receiverID string,
) []*SignRound4Message {
	filtered := make([]*SignRound4Message, 0)
	for _, message := range messages {
		if message.receiverID == receiverID {
			filtered = append(filtered, message)
		}
	}
	return filtered
}
