package ephemeral

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	msg := "I’m just a little black rain cloud, hovering under the honey tree."

	symmetricKey, err := newEcdhSymmetricKey()
	if err != nil {
		t.Fatal(err)
	}

	encrypted, err := symmetricKey.Encrypt([]byte(msg))
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := symmetricKey.Decrypt(encrypted)
	if err != nil {
		t.Fatal(err)
	}

	decryptedString := string(decrypted)
	if decryptedString != msg {
		t.Fatalf(
			"unexpected message\nexpected: %v\nactual: %v",
			msg,
			decryptedString,
		)
	}
}

func TestCiphertextRandomized(t *testing.T) {
	msg := `You can't stay in your corner of the forest waiting 
			 for others to come to you. You have to go to them sometimes.`

	symmetricKey, err := newEcdhSymmetricKey()
	if err != nil {
		t.Fatal(err)
	}

	encrypted1, err := symmetricKey.Encrypt([]byte(msg))
	if err != nil {
		t.Fatal(err)
	}

	encrypted2, err := symmetricKey.Encrypt([]byte(msg))
	if err != nil {
		t.Fatal(err)
	}

	if len(encrypted1) != len(encrypted2) {
		t.Fatalf(
			"expected the same length of ciphertexts (%v vs %v)",
			len(encrypted1),
			len(encrypted2),
		)
	}

	if reflect.DeepEqual(encrypted1, encrypted2) {
		t.Fatalf("expected two different ciphertexts")
	}
}

func TestGracefullyHandleBrokenCipher(t *testing.T) {
	symmetricKey, err := newEcdhSymmetricKey()
	if err != nil {
		t.Fatal(err)
	}

	brokenCipher := []byte{0x01, 0x02, 0x03}

	_, err = symmetricKey.Decrypt(brokenCipher)

	expectedError := fmt.Errorf("symmetric key decryption failed")
	if !reflect.DeepEqual(expectedError, err) {
		t.Fatalf(
			"unexpected error\nexpected: %v\nactual:   %v",
			expectedError,
			err,
		)
	}
}

func newEcdhSymmetricKey() (*SymmetricEcdhKey, error) {
	keyPair1, err := GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	keyPair2, err := GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	return keyPair1.PrivateKey.Ecdh(keyPair2.PublicKey), nil
}
