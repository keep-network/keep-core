package ephemeral

import (
	"reflect"
	"testing"
)

func TestMarshalUnmarshalPublicKey(t *testing.T) {
	keyPair, err := GenerateKeypair()
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
	keyPair, err := GenerateKeypair()
	if err != nil {
		t.Fatal(err)
	}

	marshalled := keyPair.PrivateKey.Marshal()

	unmarshalled := UnmarshalPrivateKey(marshalled)

	if !reflect.DeepEqual(unmarshalled, keyPair.PrivateKey) {
		t.Fatal("unmarshalled private key does not match the original one")
	}
}
