package local

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	message := []byte("Two things only the greatest fools do: throw " +
		"stones at hornets' nests and threaten a witcher.")

	signing, err := newSigning()
	if err != nil {
		t.Fatal(err)
	}

	signature, err := signing.Sign(message)
	if err != nil {
		t.Fatal(err)
	}

	var tests = map[string]struct {
		message                 []byte
		signature               []byte
		validSignatureExpected  bool
		validationErrorExpected bool
	}{
		"valid signature for message": {
			message:                 message,
			signature:               signature,
			validSignatureExpected:  true,
			validationErrorExpected: false,
		},
		"invalid signature for message": {
			message:                 []byte("I am sorry"),
			signature:               signature,
			validSignatureExpected:  false,
			validationErrorExpected: false,
		},
		"corrupted signature": {
			message:                 message,
			signature:               []byte("I am so sorry"),
			validSignatureExpected:  false,
			validationErrorExpected: true,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			ok, err := signing.Verify(test.message, test.signature)

			if !ok && test.validSignatureExpected {
				t.Errorf("expected valid signature but verification failed")
			}
			if ok && !test.validSignatureExpected {
				t.Errorf("expected invalid signature but verification succeeded")
			}

			if err == nil && test.validationErrorExpected {
				t.Errorf("expected signature validation error; none happened")
			}
			if err != nil && !test.validationErrorExpected {
				t.Errorf("unexpected signature validation error [%v]", err)
			}
		})
	}
}

func TestSignAndVerifyWithProvidedPublicKey(t *testing.T) {
	message := []byte("You shall not pass")

	signing1, err := newSigning()
	if err != nil {
		t.Fatal(err)
	}

	signing2, err := newSigning()
	if err != nil {
		t.Fatal(err)
	}

	publicKey := signing1.PublicKey()
	signature, err := signing1.Sign(message)
	if err != nil {
		t.Fatal(err)
	}

	var tests = map[string]struct {
		message                 []byte
		signature               []byte
		publicKey               []byte
		validSignatureExpected  bool
		validationErrorExpected bool
	}{
		"valid signature for message": {
			message:                 message,
			signature:               signature,
			publicKey:               publicKey,
			validSignatureExpected:  true,
			validationErrorExpected: false,
		},
		"invalid signature for message": {
			message:                 []byte("Fly you fools"),
			signature:               signature,
			publicKey:               publicKey,
			validSignatureExpected:  false,
			validationErrorExpected: false,
		},
		"corrupted signature": {
			message:                 message,
			signature:               []byte("My precious"),
			publicKey:               publicKey,
			validSignatureExpected:  false,
			validationErrorExpected: true,
		},
		"invalid remote public key": {
			message:                 message,
			signature:               signature,
			publicKey:               signing2.PublicKey(),
			validSignatureExpected:  false,
			validationErrorExpected: false,
		},
		"corrupted remote public key": {
			message:                 message,
			signature:               signature,
			publicKey:               []byte("A Balrog"),
			validSignatureExpected:  false,
			validationErrorExpected: true,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			ok, err := signing2.VerifyWithPublicKey(
				test.message,
				test.signature,
				test.publicKey,
			)

			if !ok && test.validSignatureExpected {
				t.Errorf("expected valid signature but verification failed")
			}
			if ok && !test.validSignatureExpected {
				t.Errorf("expected invalid signature but verification succeeded")
			}

			if err == nil && test.validationErrorExpected {
				t.Errorf("expected signature validation error; none happened")
			}
			if err != nil && !test.validationErrorExpected {
				t.Errorf("unexpected signature validation error [%v]", err)
			}
		})
	}
}

func newSigning() (*localSigning, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		return nil, err
	}

	return &localSigning{key}, nil
}
