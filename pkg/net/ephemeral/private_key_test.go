package ephemeral

import (
	"reflect"
	"testing"
)

func TestMarshalUnmarshalPublicKey(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	marshalled := keyPair.PublicKey.Marshal()

	unmarshalled, err := UnmarshalPublicKey(marshalled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(unmarshalled, keyPair.PublicKey) {
		t.Fatal("unmarshalled public key does not match the original one")
	}
}

func TestMarshalUnmarshalPrivateKey(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	marshalled := keyPair.PrivateKey.Marshal()
	unmarshalled := UnmarshalPrivateKey(marshalled)

	if !reflect.DeepEqual(unmarshalled, keyPair.PrivateKey) {
		t.Fatal("unmarshalled private key does not match the original one")
	}
}

func TestIsKeyMatching(t *testing.T) {
	keyPair1, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	keyPair2, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	ok := keyPair1.PublicKey.IsKeyMatching(keyPair1.PrivateKey)
	if !ok {
		t.Fatal("private key does not match the public key")
	}

	ok = keyPair1.PublicKey.IsKeyMatching(keyPair2.PrivateKey)
	if ok {
		t.Fatal("private key matches wrong public key")
	}
}
