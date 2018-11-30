package gjkr

import (
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func TestCreateNewPeerSharesMessage(t *testing.T) {
	shareS := big.NewInt(1381319)
	shareT := big.NewInt(1010212)

	peerSharesMessage, key, err := newTestPeerSharesMessage(shareS, shareT)

	decryptedS, err := peerSharesMessage.decryptShareS(key)
	if err != nil {
		t.Fatal(err)
	}

	decryptedT, err := peerSharesMessage.decryptShareT(key)
	if err != nil {
		t.Fatal(err)
	}

	if shareS.Cmp(decryptedS) != 0 {
		t.Fatalf(
			"unexpected S share\nexpected: %v\nactual:   %v",
			shareS,
			decryptedS,
		)
	}

	if shareT.Cmp(decryptedT) != 0 {
		t.Fatalf(
			"unexpected T share\nexpected: %v\nactual:   %v",
			shareT,
			decryptedT,
		)
	}
}

func TestCanDecrypt(t *testing.T) {
	var tests = map[string]struct {
		modifyMessage  func(msg *PeerSharesMessage)
		expectedResult bool
	}{
		"decryption possible": {
			expectedResult: true,
		},
		"decryption not possible - invalid S": {
			modifyMessage: func(msg *PeerSharesMessage) {
				msg.encryptedShareS = []byte{0x01, 0x02, 0x03}
			},
			expectedResult: false,
		},
		"decryption not possible - invalid T": {
			modifyMessage: func(msg *PeerSharesMessage) {
				msg.encryptedShareT = []byte{0x04, 0x05, 0x06}
			},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			shareS := big.NewInt(90787123)
			shareT := big.NewInt(62829113)

			message, key, err := newTestPeerSharesMessage(shareS, shareT)
			if err != nil {
				t.Fatal(err)
			}

			if test.modifyMessage != nil {
				test.modifyMessage(message)
			}

			canDecrypt := message.CanDecrypt(key)

			if test.expectedResult != canDecrypt {
				t.Fatalf(
					"unexpected CanDecrypt result\nexpected: %v\nactual:   %v",
					test.expectedResult,
					canDecrypt,
				)
			}
		})
	}
}

func newTestPeerSharesMessage(shareS, shareT *big.Int) (
	*PeerSharesMessage,
	ephemeral.SymmetricKey,
	error,
) {
	keyPair1, err := ephemeral.GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}

	keyPair2, err := ephemeral.GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}

	key := keyPair1.PrivateKey.Ecdh(keyPair2.PublicKey)

	peerSharesMessage, err := newPeerSharesMessage(1, 2, shareS, shareT, key)
	if err != nil {
		return nil, nil, err
	}

	return peerSharesMessage, key, nil
}
