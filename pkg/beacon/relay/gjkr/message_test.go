package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
)

func TestAddAndDecryptSharesFromMessage(t *testing.T) {
	sender := group.MemberIndex(41)
	receiver := group.MemberIndex(11)
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
	sender := group.MemberIndex(81)
	receiver := group.MemberIndex(11)
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

	_, err = peerSharesMessage.decryptShareS(group.MemberIndex(9), key)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v",
			expectedError,
			err,
		)
	}

	_, err = peerSharesMessage.decryptShareT(group.MemberIndex(9), key)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v",
			expectedError,
			err,
		)
	}
}

func newTestPeerSharesMessage(senderID, receiverID group.MemberIndex, shareS, shareT *big.Int) (
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
