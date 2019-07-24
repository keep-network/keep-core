package operator

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"reflect"
	"testing"
)

func TestMarshalRoundTrip(t *testing.T) {
	_, operatorPublicKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	marshalled := Marshal(operatorPublicKey)
	unmarshalled, err := Unmarshal(marshalled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(unmarshalled, operatorPublicKey) {
		t.Fatalf(
			"Unexpected unmarshalled public key\nExpected: %v\nActual:   %v",
			operatorPublicKey,
			unmarshalled,
		)
	}
}

func TestUnmarshalIncorrectKey(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	// Does not conform to the correct curve, but is a valid ecdsa and operator key
	incorrectKey := (*PublicKey)(&privateKey.PublicKey)

	unmarshalled, err := Unmarshal(Marshal(incorrectKey))

	if unmarshalled != nil {
		t.Errorf("Expected nil unmarshalled key")
	}

	expectedError := fmt.Errorf("incorrect public key bytes")
	if !reflect.DeepEqual(err, expectedError) {
		t.Fatalf(
			"Unexpected unmarshalling error\nExpected: %v\nActual:   %v",
			expectedError,
			err,
		)
	}
}
