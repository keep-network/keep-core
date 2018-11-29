package gjkr

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func initializeSymmetricKeyMembersGroup(
	threshold int,
	groupSize int,
	dkg *DKG,
) ([]*SymmetricKeyGeneratingMember, error) {
	group := &Group{
		groupSize:          groupSize,
		dishonestThreshold: threshold,
	}

	var members []*SymmetricKeyGeneratingMember

	for i := 1; i <= groupSize; i++ {
		id := MemberID(i)
		members = append(members, &SymmetricKeyGeneratingMember{
			EphemeralKeyGeneratingMember: &EphemeralKeyGeneratingMember{
				memberCore: &memberCore{
					ID:             id,
					group:          group,
					protocolConfig: dkg,
				},
				ephemeralKeys: make(map[MemberID]*ephemeral.KeyPair),
			},
			symmetricKeys: make(map[MemberID]ephemeral.SymmetricKey),
		})
		group.RegisterMemberID(id)
	}

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

func TestGenerateSymmetricKeyGroup(t *testing.T) {
	groupSize := 3
	group := &Group{
		groupSize: groupSize,
	}

	// Create a group of 3 members
	var ephemeralGeneratingMembers []*EphemeralKeyGeneratingMember
	for i := 1; i <= groupSize; i++ {
		id := MemberID(i)
		ephemeralGeneratingMembers = append(ephemeralGeneratingMembers,
			&EphemeralKeyGeneratingMember{
				memberCore: &memberCore{
					ID:             id,
					group:          group,
					protocolConfig: nil,
				},
				ephemeralKeys: make(map[MemberID]*ephemeral.KeyPair),
			},
		)
		group.RegisterMemberID(id)
	}

	// generate ephemeral key pairs for each group member; prepare messages
	broadcastedPubKeyMessages := make(map[MemberID][]*EphemeralPublicKeyMessage)
	for _, ephemeralGeneratingMember := range ephemeralGeneratingMembers {
		messages, err := ephemeralGeneratingMember.GenerateEphemeralKeyPair()
		if err != nil {

		}
		broadcastedPubKeyMessages[ephemeralGeneratingMember.ID] = messages
	}

	// We should have groupSize members.
	members := reflect.ValueOf(broadcastedPubKeyMessages).MapKeys()
	if len(members) != groupSize {
		t.Fatalf(
			"expected messages for %d members, got %d members",
			groupSize,
			len(members),
		)
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
		symmetricKeyMap := make(map[MemberID]ephemeral.SymmetricKey)
		symmetricGeneratingMembers = append(symmetricGeneratingMembers,
			&SymmetricKeyGeneratingMember{
				EphemeralKeyGeneratingMember: member,
				symmetricKeys:                symmetricKeyMap,
			},
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
