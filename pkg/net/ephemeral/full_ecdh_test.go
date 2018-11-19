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
	_ = privKey1.ECDH(pubKey2)

	// player 2:
	_ = privKey2.ECDH(pubKey1)
}
