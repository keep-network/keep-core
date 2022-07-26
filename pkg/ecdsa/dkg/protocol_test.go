package dkg

import (
	"fmt"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"reflect"
	"testing"
)

// TODO: This file contains unit tests that stress each protocol phase
//       separately. We should also develop integration tests just like
//       we did for Random Beacon's DKG protocol. This will require
//       to refactor the `pkg/internal/dkgtest` package in the way that
//       it supports the ECDSA DKG as well.

func TestGenerateEphemeralKeyPair(t *testing.T) {
	groupSize := 3
	dishonestThreshold := 0

	members := initializeEphemeralKeyPairGeneratingMembersGroup(
		dishonestThreshold,
		groupSize,
	)

	// Generate ephemeral key pairs for each group member.
	messages := make(map[group.MemberIndex]*ephemeralPublicKeyMessage)
	for _, member := range members {
		message, err := member.generateEphemeralKeyPair()
		if err != nil {
			t.Fatal(err)
		}
		messages[member.id] = message
	}

	// Assert that each member has a correct state.
	for _, member := range members {
		// Assert the right key pairs count is stored in the member's state.
		expectedKeyPairsCount := groupSize - 1
		actualKeyPairsCount := len(member.ephemeralKeyPairs)
		testutils.AssertIntsEqual(
			t,
			fmt.Sprintf(
				"number of stored ephemeral key pairs for member [%v]",
				member.id,
			),
			expectedKeyPairsCount,
			actualKeyPairsCount,
		)

		// Assert the member does not hold a key pair with itself.
		_, ok := member.ephemeralKeyPairs[member.id]
		if ok {
			t.Fatalf(
				"[member:%v] found ephemeral key pair generated to self",
				member.id,
			)
		}

		// Assert key pairs are non-nil.
		for otherMemberID, keyPair := range member.ephemeralKeyPairs {
			if keyPair.PrivateKey == nil {
				t.Fatalf(
					"[member:%v] key pair's private key not set for member [%v]",
					member.id,
					otherMemberID,
				)
			}

			if keyPair.PublicKey == nil {
				t.Fatalf(
					"[member:%v] key pair's public key not set for member [%v]",
					member.id,
					otherMemberID,
				)
			}
		}
	}

	// Assert that each message is formed correctly.
	for memberID, message := range messages {
		// We should always be the sender of our own messages.
		if memberID != message.senderID {
			t.Fatalf(
				"message from incorrect sender\n"+
					"expected: [%v]\n"+
					"actual:   [%v]",
				memberID,
				message.senderID,
			)
		}

		// We should not generate an ephemeral key for ourselves.
		_, ok := message.ephemeralPublicKeys[memberID]
		if ok {
			t.Fatal("found ephemeral key generated to self")
		}
	}
}

func TestGenerateSymmetricKeys(t *testing.T) {
	groupSize := 3
	dishonestThreshold := 0

	members, messages, err := initializeSymmetricKeyGeneratingMembersGroup(
		dishonestThreshold,
		groupSize,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Generate symmetric keys for each group member.
	for _, member := range members {
		var receivedMessages []*ephemeralPublicKeyMessage
		for _, message := range messages {
			if message.senderID != member.id {
				receivedMessages = append(receivedMessages, message)
			}
		}

		err := member.generateSymmetricKeys(receivedMessages)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Assert that each member has a correct state.
	for _, member := range members {
		// Assert the right keys count is stored in the member's state.
		expectedKeysCount := groupSize - 1
		actualKeysCount := len(member.symmetricKeys)
		testutils.AssertIntsEqual(
			t,
			fmt.Sprintf(
				"number of stored symmetric keys for member [%v]",
				member.id,
			),
			expectedKeysCount,
			actualKeysCount,
		)

		// Assert all symmetric keys stored by this member are correct.
		for otherMemberID, actualKey := range member.symmetricKeys {
			var otherMemberEphemeralPublicKey *ephemeral.PublicKey
			for _, message := range messages {
				if message.senderID == otherMemberID {
					if ephemeralPublicKey, ok := message.ephemeralPublicKeys[member.id]; ok {
						otherMemberEphemeralPublicKey = ephemeralPublicKey
					}
				}
			}

			if otherMemberEphemeralPublicKey == nil {
				t.Fatalf(
					"[member:%v] no ephemeral public key from member [%v]",
					member.id,
					otherMemberID,
				)
			}

			expectedKey := ephemeral.SymmetricKey(
				member.ephemeralKeyPairs[otherMemberID].PrivateKey.Ecdh(
					otherMemberEphemeralPublicKey,
				),
			)

			if !reflect.DeepEqual(
				expectedKey,
				actualKey,
			) {
				t.Fatalf(
					"[member:%v] wrong symmetric key for member [%v]",
					member.id,
					otherMemberID,
				)
			}
		}
	}
}

func TestGenerateSymmetricKeys_InvalidEphemeralPublicKeyMessage(t *testing.T) {
	groupSize := 3
	dishonestThreshold := 0

	members, messages, err := initializeSymmetricKeyGeneratingMembersGroup(
		dishonestThreshold,
		groupSize,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Corrupt the message sent by member 2 by removing the ephemeral
	// public key generated for member 3.
	misbehavingMemberID := group.MemberIndex(2)
	delete(messages[misbehavingMemberID-1].ephemeralPublicKeys, 3)

	// Generate symmetric keys for each group member.
	for _, member := range members {
		var receivedMessages []*ephemeralPublicKeyMessage
		for _, message := range messages {
			if message.senderID != member.id {
				receivedMessages = append(receivedMessages, message)
			}
		}

		err := member.generateSymmetricKeys(receivedMessages)

		var expectedErr error
		// The misbehaved member should not get an error.
		if member.id != misbehavingMemberID {
			expectedErr = fmt.Errorf(
				"member [%v] sent invalid ephemeral "+
					"public key message",
				misbehavingMemberID,
			)
		}

		testutils.AssertErrorsEqual(t, expectedErr, err)
	}
}

func initializeEphemeralKeyPairGeneratingMembersGroup(
	dishonestThreshold int,
	groupSize int,
) []*ephemeralKeyPairGeneratingMember {
	dkgGroup := group.NewGroup(dishonestThreshold, groupSize)

	var members []*ephemeralKeyPairGeneratingMember
	for i := 1; i <= groupSize; i++ {
		id := group.MemberIndex(i)
		members = append(members, &ephemeralKeyPairGeneratingMember{
			member: &member{
				id:    id,
				group: dkgGroup,
			},
			ephemeralKeyPairs: make(map[group.MemberIndex]*ephemeral.KeyPair),
		})
	}

	return members
}
func initializeSymmetricKeyGeneratingMembersGroup(
	dishonestThreshold int,
	groupSize int,
) (
	[]*symmetricKeyGeneratingMember,
	[]*ephemeralPublicKeyMessage,
	error,
) {
	var members []*symmetricKeyGeneratingMember
	var messages []*ephemeralPublicKeyMessage

	for _, member := range initializeEphemeralKeyPairGeneratingMembersGroup(
		dishonestThreshold,
		groupSize,
	) {
		message, err := member.generateEphemeralKeyPair()
		if err != nil {
			return nil, nil, fmt.Errorf(
				"cannot generate ephemeral key pair for member [%v]: [%v]",
				member.id,
				err,
			)
		}

		members = append(members, member.initializeSymmetricKeyGeneration())
		messages = append(messages, message)
	}

	return members, messages, nil
}
