package ephemeral

import "testing"

func TestFullEcdh(t *testing.T) {
	//
	// players generate ephemeral keypair
	//

	// player 1
	player1KeyPair, err := GenerateKeypair()
	if err != nil {
		t.Fatal(err)
	}
	privKey1 := player1KeyPair.PrivateKey
	pubKey1 := player1KeyPair.PublicKey

	// player 2
	player2KeyPair, err := GenerateKeypair()
	if err != nil {
		t.Fatal(err)
	}
	privKey2 := player2KeyPair.PrivateKey
	pubKey2 := player2KeyPair.PublicKey

	//
	// players exchange public keys and perform ECDH
	//

	// player 1:
	symmetricKey1 := privKey1.Ecdh(pubKey2)

	// player 2:
	symmetricKey2 := privKey2.Ecdh(pubKey1)

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
