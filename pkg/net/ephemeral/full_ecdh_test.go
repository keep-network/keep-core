package ephemeral

import "testing"

func TestFullECDH(t *testing.T) {
	//
	// players generate ephemeral keypair
	//

	// player 1
	privKey1, pubKey1, err := GenerateEphemeralKeypair()
	if err != nil {
		t.Fatal(err)
	}

	// player 2
	privKey2, pubKey2, err := GenerateEphemeralKeypair()
	if err != nil {
		t.Fatal(err)
	}

	//
	// players exchange public keys and perform ECDH
	//

	// player 1:
	symmetricKey1 := privKey1.ECDH(pubKey2)

	// player 2:
	symmetricKey2 := privKey2.ECDH(pubKey1)

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
