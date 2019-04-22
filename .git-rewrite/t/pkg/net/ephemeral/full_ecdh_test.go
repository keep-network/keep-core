package ephemeral

import "testing"

func TestFullEcdh(t *testing.T) {
	//
	// players generate ephemeral keypair
	//

	// player 1
	keyPair1, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	// player 2
	keyPair2, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	//
	// players exchange public keys and perform ECDH
	//

	// player 1:
	symmetricKey1 := keyPair1.PrivateKey.Ecdh(keyPair2.PublicKey)

	// player 2:
	symmetricKey2 := keyPair2.PrivateKey.Ecdh(keyPair1.PublicKey)

	//
	// players use symmetric key for encryption/decryption
	//

	msg := "People say nothing is impossible, but I do nothing every day"

	// player 1:
	encrypted, err := symmetricKey1.Encrypt([]byte(msg))
	if err != nil {
		t.Fatal(err)
	}

	//player 2:
	decrypted, err := symmetricKey2.Decrypt(encrypted)
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
