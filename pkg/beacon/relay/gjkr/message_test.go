package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/member"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func TestAddAndDecryptSharesFromMessage(t *testing.T) {
	sender := member.MemberIndex(4181)
	receiver := member.MemberIndex(1231)
	shareS := big.NewInt(1381319)
	shareT := big.NewInt(1010212)

	peerSharesMessage, key, err := newTestPeerSharesMessage(
		sender,
		receiver,
		shareS,
		shareT,
	)
	if err != nil {
		t.Fatal(err)
	}

	decryptedS, err := peerSharesMessage.decryptShareS(receiver, key)
	if err != nil {
		t.Fatal(err)
	}

	decryptedT, err := peerSharesMessage.decryptShareT(receiver, key)
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

func TestNoSharesForReceiver(t *testing.T) {
	sender := member.MemberIndex(4181)
	receiver := member.MemberIndex(1231)
	shareS := big.NewInt(1381319)
	shareT := big.NewInt(1010212)

	peerSharesMessage, key, err := newTestPeerSharesMessage(
		sender,
		receiver,
		shareS,
		shareT,
	)
	if err != nil {
		t.Fatal(err)
	}

	expectedError := fmt.Errorf("no shares for receiver 9")

	_, err = peerSharesMessage.decryptShareS(member.MemberIndex(9), key)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v",
			expectedError,
			err,
		)
	}

	_, err = peerSharesMessage.decryptShareT(member.MemberIndex(9), key)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v",
			expectedError,
			err,
		)
	}
}

func TestCanDecrypt(t *testing.T) {
	sender := member.MemberIndex(4181)
	receiver := member.MemberIndex(1231)

	var tests = map[string]struct {
		modifyMessage  func(msg *PeerSharesMessage)
		expectedResult bool
	}{
		"decryption possible": {
			expectedResult: true,
		},
		"decryption not possible - invalid S": {
			modifyMessage: func(msg *PeerSharesMessage) {
				msg.shares[receiver].encryptedShareS = []byte{0x01, 0x02, 0x03}
			},
			expectedResult: false,
		},
		"decryption not possible - invalid T": {
			modifyMessage: func(msg *PeerSharesMessage) {
				msg.shares[receiver].encryptedShareT = []byte{0x04, 0x05, 0x06}
			},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			shareS := big.NewInt(90787123)
			shareT := big.NewInt(62829113)

			message, key, err := newTestPeerSharesMessage(
				sender,
				receiver,
				shareS,
				shareT,
			)
			if err != nil {
				t.Fatal(err)
			}

			if test.modifyMessage != nil {
				test.modifyMessage(message)
			}

			canDecrypt := message.CanDecrypt(receiver, key)

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

func newTestPeerSharesMessage(senderID, receiverID member.MemberIndex, shareS, shareT *big.Int) (
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

	msg := newPeerSharesMessage(senderID)
	if err := msg.addShares(receiverID, shareS, shareT, key); err != nil {
		return nil, nil, err
	}

	return msg, key, nil
}
