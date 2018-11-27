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

	pubKey := keyPair.PublicKey

	marshalled := pubKey.Marshal()

	unmarshalled, err := UnmarshalPublicKey(marshalled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(unmarshalled, pubKey) {
		t.Fatal("unmarshalled public key does not match the original one")
	}
}

func TestMarshalUnmarshalPrivateKey(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	privKey := keyPair.PrivateKey

	marshalled := privKey.Marshal()

	unmarshalled := UnmarshalPrivateKey(marshalled)

	if !reflect.DeepEqual(unmarshalled, privKey) {
		t.Fatal("unmarshalled private key does not match the original one")
	}
}
