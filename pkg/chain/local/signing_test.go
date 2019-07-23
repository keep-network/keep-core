package local

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	signing := &localSigning{key}

	message := []byte("Two things only the greatest fools do: throw " +
		"stones at hornets' nests and threaten a witcher.")

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
