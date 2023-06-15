package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/internal/testutils"

	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

func TestSaveEphemeralKeyMessagesForEvidence(t *testing.T) {
	groupSize := 2
	dishonestThreshold := 0

	// Create a group of 2 members
	ephemeralGeneratingMembers := initializeEphemeralKeyPairMembersGroup(
		dishonestThreshold,
		groupSize,
	)

	member1 := ephemeralGeneratingMembers[0]
	member2 := ephemeralGeneratingMembers[1]

	message1, err := member1.GenerateEphemeralKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := member2.GenerateEphemeralKeyPair(); err != nil {
		t.Fatal(err)
	}

	symmetricKeyMember2 := member2.InitializeSymmetricKeyGeneration()
	if err := symmetricKeyMember2.GenerateSymmetricKeys(
		[]*EphemeralPublicKeyMessage{message1},
	); err != nil {
		t.Fatal(err)
	}

	evidenceMsg := symmetricKeyMember2.evidenceLog.ephemeralPublicKeyMessage(
		member1.ID,
	)

	if !reflect.DeepEqual(message1, evidenceMsg) {
		t.Fatalf(
			"unexpected message in evidence log\nexpected: %v\n actual:   %v",
			message1,
			evidenceMsg,
		)
	}
}

func TestGenerateEphemeralKeys(t *testing.T) {
	groupSize := 3
	dishonestThreshold := 0

	// Create a group of 3 members
	ephemeralGeneratingMembers := initializeEphemeralKeyPairMembersGroup(
		dishonestThreshold,
		groupSize,
	)

	// generate ephemeral key pairs for each group member; prepare messages
	broadcastedPubKeyMessages := make(map[group.MemberIndex]*EphemeralPublicKeyMessage)
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
	receivedPubKeyMessages := make(map[group.MemberIndex][]*EphemeralPublicKeyMessage)
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
	dishonestThreshold int,
	groupSize int,
) []*EphemeralKeyPairGeneratingMember {
	dkgGroup := group.NewGroup(dishonestThreshold, groupSize)

	protocolParameters := newProtocolParameters(big.NewInt(18313131145))

	var members []*EphemeralKeyPairGeneratingMember
	for i := 1; i <= groupSize; i++ {
		id := group.MemberIndex(i)
		members = append(members, &EphemeralKeyPairGeneratingMember{
			LocalMember: &LocalMember{
				memberCore: &memberCore{
					logger:             &testutils.MockLogger{},
					ID:                 id,
					group:              dkgGroup,
					evidenceLog:        newDkgEvidenceLog(),
					protocolParameters: protocolParameters,
					sessionID:          "session-1",
				},
			},
			ephemeralKeyPairs: make(map[group.MemberIndex]*ephemeral.KeyPair),
		})
	}

	return members
}

func initializeSymmetricKeyMembersGroup(
	dishonestThreshold int,
	groupSize int,
) ([]*SymmetricKeyGeneratingMember, error) {
	keyPairMembers := initializeEphemeralKeyPairMembersGroup(dishonestThreshold, groupSize)

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
	dishonestThreshold int,
	groupSize int,
) ([]*SymmetricKeyGeneratingMember, error) {
	symmetricKeyMembers, err := initializeSymmetricKeyMembersGroup(
		dishonestThreshold,
		groupSize,
	)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%v]", err)
	}

	// generate symmetric keys with all other members of the group
	for _, member1 := range symmetricKeyMembers {
		ephemeralKeys := make(map[group.MemberIndex]*ephemeral.PublicKey)

		for _, member2 := range symmetricKeyMembers {
			if member1.ID != member2.ID {
				privKey := member1.ephemeralKeyPairs[member2.ID].PrivateKey
				pubKey := member2.ephemeralKeyPairs[member1.ID].PublicKey
				member1.symmetricKeys[member2.ID] = privKey.Ecdh(pubKey)

				ephemeralKeys[member2.ID] = member1.ephemeralKeyPairs[member2.ID].PublicKey
			}
		}

		// simulating message broadcast in the group
		for _, member := range symmetricKeyMembers {
			member.evidenceLog.PutEphemeralMessage(
				&EphemeralPublicKeyMessage{member1.ID, ephemeralKeys, "session-1"},
			)
		}
	}

	return symmetricKeyMembers, nil
}
