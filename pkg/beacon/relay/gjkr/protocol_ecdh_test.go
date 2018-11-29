package gjkr

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func TestGenerateSymmetricKeyGroup(t *testing.T) {
	groupSize := 3

	// Create a group of 3 members
	ephemeralGeneratingMembers := createEphemeralKeyMembersGroup(
		groupSize,
		groupSize, // threshold = groupSize
		nil,
	)

	// generate ephemeral key pairs for each group member; prepare messages
	broadcastedPubKeyMessages := make(map[MemberID][]*EphemeralPublicKeyMessage)
	for _, ephemeralGeneratingMember := range ephemeralGeneratingMembers {
		messages, err := ephemeralGeneratingMember.GenerateEphemeralKeyPair()
		if err != nil {
			t.Fatal(err)
		}
		broadcastedPubKeyMessages[ephemeralGeneratingMember.ID] = messages
	}

	for memberID, ephemeralPubKeyMessages := range broadcastedPubKeyMessages {
		// We should have groupSize - 1 messages per member.
		if len(ephemeralPubKeyMessages) != groupSize-1 {
			t.Fatalf(
				"expected %d messages, got %d messages",
				groupSize-1,
				len(ephemeralPubKeyMessages),
			)
		}

		for _, message := range ephemeralPubKeyMessages {
			// We should always be the sender of our own messages
			if message.senderID != memberID {
				t.Fatalf("message from incorrect sender got %v want %v",
					message.senderID,
					memberID,
				)
			}

			// We should never have a message addressed to ourselves
			if message.receiverID == memberID {
				t.Fatal("found message addressed to self")
			}
		}
	}

	// Simulate the each member receiving all messages from the network
	receivedPubKeyMessages := make(map[MemberID][]*EphemeralPublicKeyMessage)
	for memberID, ephemeralPubKeyMessages := range broadcastedPubKeyMessages {
		for _, otherMember := range ephemeralGeneratingMembers {
			// We would only receive messages from the other members
			if memberID != otherMember.ID {
				receivedPubKeyMessages[otherMember.ID] = append(
					receivedPubKeyMessages[otherMember.ID],
					ephemeralPubKeyMessages...,
				)
			}
		}
	}

	// Move to the next phase, using the previous phase as state
	var symmetricGeneratingMembers []*SymmetricKeyGeneratingMember
	for _, member := range ephemeralGeneratingMembers {
		symmetricGeneratingMembers = append(
			symmetricGeneratingMembers,
			member.InitializeSymmetricKeyGeneration(),
		)
	}

	// For each member, attempt to generate a symmetric key
	for _, symmetricGeneratingMember := range symmetricGeneratingMembers {
		if err := symmetricGeneratingMember.GenerateSymmetricKeys(
			receivedPubKeyMessages[symmetricGeneratingMember.ID],
		); err != nil {
			t.Fatalf(
				"failed to generate symmetric key with error %v",
				err,
			)
		}
	}

	// Ensure that for each member, we generated the correct number of
	// symmetric keys (groupSize - 1 keys)
	for _, symmetricGeneratingMember := range symmetricGeneratingMembers {
		symmetricKeys := symmetricGeneratingMember.symmetricKeys
		keySlice := reflect.ValueOf(symmetricKeys).MapKeys()
		if len(keySlice) != groupSize-1 {
			t.Fatalf(
				"expected %d keys, got %d keys",
				groupSize-1,
				len(keySlice),
			)
		}
	}
}

func createEphemeralKeyMembersGroup(
	threshold int,
	groupSize int,
	dkg *DKG,
) []*EphemeralKeyGeneratingMember {
	group := &Group{
		groupSize:          groupSize,
		dishonestThreshold: threshold,
	}

	var members []*EphemeralKeyGeneratingMember
	for i := 1; i <= groupSize; i++ {
		id := MemberID(i)
		members = append(members, &EphemeralKeyGeneratingMember{
			memberCore: &memberCore{
				ID:             id,
				group:          group,
				protocolConfig: dkg,
			},
			ephemeralKeys: make(map[MemberID]*ephemeral.KeyPair),
		})
		group.RegisterMemberID(id)
	}

	return members
}

func initializeEphemeralKeyMembersGroup(
	threshold int,
	groupSize int,
	dkg *DKG,
) ([]*EphemeralKeyGeneratingMember, error) {
	members := createEphemeralKeyMembersGroup(threshold, groupSize, dkg)

	// generate ephemeral key pairs for all other members of the group
	for _, member1 := range members {
		for _, member2 := range members {
			if member1.ID != member2.ID {

				keyPair, err := ephemeral.GenerateKeyPair()
				if err != nil {
					return nil, fmt.Errorf(
						"SymmetricKeyGeneratingMember initialization failed [%v]",
						err,
					)
				}
				member1.ephemeralKeys[member2.ID] = keyPair
			}
		}
	}

	return members, nil
}

func initializeSymmetricKeyMembersGroup(
	threshold int,
	groupSize int,
	dkg *DKG,
) ([]*SymmetricKeyGeneratingMember, error) {
	ephemeralKeyMembers, err := initializeEphemeralKeyMembersGroup(
		threshold,
		groupSize,
		dkg,
	)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%v]", err)
	}

	var members []*SymmetricKeyGeneratingMember
	for _, ephemeralKeyMember := range ephemeralKeyMembers {
		members = append(
			members,
			ephemeralKeyMember.InitializeSymmetricKeyGeneration(),
		)
	}

	// generate symmetric keys with all other members of the group
	for _, member1 := range members {
		for _, member2 := range members {
			if member1.ID != member2.ID {
				privKey := member1.ephemeralKeys[member2.ID].PrivateKey
				pubKey := member2.ephemeralKeys[member1.ID].PublicKey
				member1.symmetricKeys[member2.ID] = privKey.Ecdh(pubKey)
			}
		}
	}

	return members, nil
}
