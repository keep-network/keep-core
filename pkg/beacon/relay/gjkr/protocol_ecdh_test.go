package gjkr

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func TestGenerateEphemeralKeys(t *testing.T) {
	groupSize := 3

	// Create a group of 3 members
	ephemeralGeneratingMembers := initializeEphemeralKeyPairMembersGroup(
		groupSize,
		groupSize, // threshold = groupSize
		nil,
	)

	// generate ephemeral key pairs for each group member; prepare messages
	broadcastedPubKeyMessages := make(map[MemberID]*EphemeralPublicKeyMessage)
	for _, ephemeralGeneratingMember := range ephemeralGeneratingMembers {
		message, err := ephemeralGeneratingMember.GenerateEphemeralKeyPair()
		if err != nil {
			t.Fatal(err)
		}
		broadcastedPubKeyMessages[ephemeralGeneratingMember.ID] = message
	}

	for memberID, message := range broadcastedPubKeyMessages {
		// We should always be the sender of our own messages
		if message.senderID != memberID {
			t.Fatalf("message from incorrect sender got %v want %v",
				message.senderID,
				memberID,
			)
		}

		// We should not generate an ephemeral key for ourselves
		_, ok := message.ephemeralPublicKeys[memberID]
		if ok {
			t.Fatal("found ephemeral key generated to self")
		}
	}

	// Simulate the each member receiving all messages from the network
	receivedPubKeyMessages := make(map[MemberID][]*EphemeralPublicKeyMessage)
	for memberID, ephemeralPubKeyMessage := range broadcastedPubKeyMessages {
		for _, otherMember := range ephemeralGeneratingMembers {
			// We would only receive messages from the other members
			if memberID != otherMember.ID {
				receivedPubKeyMessages[otherMember.ID] = append(
					receivedPubKeyMessages[otherMember.ID],
					ephemeralPubKeyMessage,
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

func initializeEphemeralKeyPairMembersGroup(
	threshold int,
	groupSize int,
	dkg *DKG,
) []*EphemeralKeyPairGeneratingMember {
	group := &Group{
		dishonestThreshold: threshold,
	}

	var members []*EphemeralKeyPairGeneratingMember
	for i := 1; i <= groupSize; i++ {
		id := MemberID(i)
		members = append(members, &EphemeralKeyPairGeneratingMember{
			memberCore: &memberCore{
				ID:             id,
				group:          group,
				protocolConfig: dkg,
			},
			ephemeralKeyPairs: make(map[MemberID]*ephemeral.KeyPair),
		})
		group.RegisterMemberID(id)
	}

	return members
}

func initializeSymmetricKeyMembersGroup(
	threshold int,
	groupSize int,
	dkg *DKG,
) ([]*SymmetricKeyGeneratingMember, error) {
	keyPairMembers := initializeEphemeralKeyPairMembersGroup(threshold, groupSize, dkg)

	// generate ephemeral key pairs for all other members of the group
	for _, member1 := range keyPairMembers {
		for _, member2 := range keyPairMembers {
			if member1.ID != member2.ID {

				keyPair, err := ephemeral.GenerateKeyPair()
				if err != nil {
					return nil, fmt.Errorf(
						"SymmetricKeyGeneratingMember initialization failed [%v]",
						err,
					)
				}
				member1.ephemeralKeyPairs[member2.ID] = keyPair
			}
		}
	}

	symmetricKeyMembers := make([]*SymmetricKeyGeneratingMember, len(keyPairMembers))
	for i, keyPairMember := range keyPairMembers {
		symmetricKeyMembers[i] = keyPairMember.InitializeSymmetricKeyGeneration()
	}

	return symmetricKeyMembers, nil
}

// generateGroupWithEphemeralKeys executes first two phases of DKG protocol and
// returns a fully initialized group of `SymmetricKeyGeneratingMember`s with all
// ephemeral keys generated (private, public, and symmetric key).
func generateGroupWithEphemeralKeys(
	threshold int,
	groupSize int,
	dkg *DKG,
) ([]*SymmetricKeyGeneratingMember, error) {
	symmetricKeyMembers, err := initializeSymmetricKeyMembersGroup(
		threshold,
		groupSize,
		dkg,
	)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%v]", err)
	}

	// generate symmetric keys with all other members of the group
	for _, member1 := range symmetricKeyMembers {
		ephemeralKeys := make(map[MemberID]*ephemeral.PublicKey)

		for _, member2 := range symmetricKeyMembers {
			if member1.ID != member2.ID {
				privKey := member1.ephemeralKeyPairs[member2.ID].PrivateKey
				pubKey := member2.ephemeralKeyPairs[member1.ID].PublicKey
				member1.symmetricKeys[member2.ID] = privKey.Ecdh(pubKey)

				ephemeralKeys[member2.ID] = member1.ephemeralKeyPairs[member2.ID].PublicKey
			}
		}

		member1.protocolConfig.evidenceLog.PutEphemeralMessage(
			&EphemeralPublicKeyMessage{member1.ID, ephemeralKeys},
		)
	}

	return symmetricKeyMembers, nil
}
