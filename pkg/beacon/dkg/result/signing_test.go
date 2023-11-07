package result

import (
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"

	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// TestResultSigningAndVerificationRoundTrip simulates Phase 13 execution when
// a group of members produces messages containing signatures and than one member
// verifies messages received from other group members.
func TestResultSigningAndVerificationRoundTrip(t *testing.T) {
	groupSize := 10

	dkgResult := &beaconchain.DKGResult{
		GroupPublicKey: []byte{10},
	}

	members, beaconChains, err := initializeSigningMembers(groupSize)
	if err != nil {
		t.Fatal(err)
	}

	currentMember := members[0]
	currentSigning := beaconChains[0].Signing()

	messages := make([]*DKGResultHashSignatureMessage, 0)

	for i, member := range members {
		message, err := member.SignDKGResult(
			dkgResult,
			beaconChains[i],
		)
		if err != nil {
			t.Fatal(err)
		}

		// Don't register message from self.
		if member.index != currentMember.index {
			messages = append(messages, message)
		}
	}

	receivedValidSignatures, err := currentMember.VerifyDKGResultSignatures(
		messages,
		currentSigning,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(receivedValidSignatures) != groupSize {
		t.Errorf(
			"unexpected number of registered signatures\nexpected: %v\nactual:   %v\n",
			groupSize,
			len(receivedValidSignatures),
		)
	}
}

func TestVerifyDKGResultSignatures(t *testing.T) {
	groupSize := 10

	dkgResultHash1 := beaconchain.DKGResultHash{10}
	dkgResultHash2 := beaconchain.DKGResultHash{20}

	members, beaconChains, err := initializeSigningMembers(groupSize)
	if err != nil {
		t.Fatal(err)
	}

	verifyingMember, verifyingMemberSigning := members[0], beaconChains[0].Signing()
	verifyingMember.preferredDKGResultHash = dkgResultHash1

	selfSignature, _ := verifyingMemberSigning.Sign(dkgResultHash1[:])
	verifyingMember.selfDKGResultSignature = selfSignature

	member2, signing2 := members[1], beaconChains[1].Signing()
	member3, signing3 := members[2], beaconChains[2].Signing()
	member4, signing4 := members[3], beaconChains[3].Signing()
	member5, signing5 := members[4], beaconChains[4].Signing()

	signature21, _ := signing2.Sign(dkgResultHash1[:])

	signature311, _ := signing3.Sign(dkgResultHash1[:])
	signature312, _ := signing3.Sign(dkgResultHash1[:])

	signature411, _ := signing4.Sign(dkgResultHash1[:])
	signature421, _ := signing4.Sign(dkgResultHash2[:])

	signature52, _ := signing5.Sign(dkgResultHash2[:])

	var tests = map[string]struct {
		messages []*DKGResultHashSignatureMessage

		expectedReceivedValidSignatures map[group.MemberIndex][]byte
		expectedError                   error
	}{
		"received valid messages with signatures for the preferred result": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member2.index,
					resultHash:  dkgResultHash1,
					signature:   signature21,
					publicKey:   signing2.PublicKey(),
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   signing3.PublicKey(),
				},
			},
			expectedReceivedValidSignatures: map[group.MemberIndex][]byte{
				verifyingMember.index: selfSignature,
				member2.index:         signature21,
				member3.index:         signature311,
			},
		},
		"received messages from other member with duplicated different signatures for the preferred result": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   signing3.PublicKey(),
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature312,
					publicKey:   signing3.PublicKey(),
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   signing3.PublicKey(),
				},
			},
			expectedReceivedValidSignatures: map[group.MemberIndex][]byte{
				verifyingMember.index: selfSignature,
			},
		},
		"received messages from other member with the same signatures for the preferred result": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   signing3.PublicKey(),
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   signing3.PublicKey(),
				},
			},
			expectedReceivedValidSignatures: map[group.MemberIndex][]byte{
				verifyingMember.index: selfSignature,
			},
		},
		"received messages from other member with signatures for two different results": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member4.index,
					resultHash:  dkgResultHash1,
					signature:   signature411,
					publicKey:   signing4.PublicKey(),
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member4.index,
					resultHash:  dkgResultHash2,
					signature:   signature421,
					publicKey:   signing4.PublicKey(),
				},
			},
			expectedReceivedValidSignatures: map[group.MemberIndex][]byte{
				verifyingMember.index: selfSignature,
			},
		},
		"received a message from other member with signature for result different than preferred": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member5.index,
					resultHash:  dkgResultHash2,
					signature:   signature52,
					publicKey:   signing5.PublicKey(),
				},
			},
			expectedReceivedValidSignatures: map[group.MemberIndex][]byte{
				verifyingMember.index: selfSignature,
			},
		},
		"received a message from other member with invalid signature": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member2.index,
					resultHash:  dkgResultHash1,
					signature:   []byte{99},
					publicKey:   signing2.PublicKey(),
				},
			},
			expectedReceivedValidSignatures: map[group.MemberIndex][]byte{
				verifyingMember.index: selfSignature,
			},
		},
		"received a message from other member with invalid public key": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member2.index,
					resultHash:  dkgResultHash1,
					signature:   signature21,
					publicKey:   signing5.PublicKey(),
				},
			},
			expectedReceivedValidSignatures: map[group.MemberIndex][]byte{
				verifyingMember.index: selfSignature,
			},
		},
		"mixed cases with received valid signatures and duplicated signatures": {
			messages: []*DKGResultHashSignatureMessage{
				// Valid signature supporting the same result as preferred.
				&DKGResultHashSignatureMessage{
					senderIndex: member2.index,
					resultHash:  dkgResultHash1,
					signature:   signature21,
					publicKey:   signing2.PublicKey(),
				},
				// Multiple signatures from the same member supporting the same result as preferred.
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   signing3.PublicKey(),
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature312,
					publicKey:   signing3.PublicKey(),
				},
				// Multiple signatures from the same member supporting two different results.
				&DKGResultHashSignatureMessage{
					senderIndex: member4.index,
					resultHash:  dkgResultHash1,
					signature:   signature411,
					publicKey:   signing4.PublicKey(),
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member4.index,
					resultHash:  dkgResultHash2,
					signature:   signature421,
					publicKey:   signing4.PublicKey(),
				},
				// Member supporting different result than preferred.
				&DKGResultHashSignatureMessage{
					senderIndex: member5.index,
					resultHash:  dkgResultHash2,
					signature:   signature52,
					publicKey:   signing5.PublicKey(),
				},
			},
			expectedReceivedValidSignatures: map[group.MemberIndex][]byte{
				verifyingMember.index: selfSignature,
				member2.index:         signature21,
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			receivedValidSignatures, err := verifyingMember.VerifyDKGResultSignatures(
				test.messages,
				verifyingMemberSigning,
			)

			if !reflect.DeepEqual(err, test.expectedError) {
				t.Errorf(
					"unexpected error\nexpected: %v\nactual:   %v\n",
					test.expectedError,
					err,
				)
			}

			if !reflect.DeepEqual(
				receivedValidSignatures,
				test.expectedReceivedValidSignatures,
			) {
				t.Errorf(
					"unexpected registered received valid signatures\nexpected: %v\nactual:   %v\n",
					test.expectedReceivedValidSignatures,
					receivedValidSignatures,
				)
			}
		})
	}
}

func TestShouldAcceptMessage(t *testing.T) {
	groupSize := 5
	honestThreshold := 3

	localChain := local_v1.Connect(groupSize, honestThreshold)

	operatorsAddresses := make([]chain.Address, groupSize)
	operatorsPublicKeys := make([][]byte, groupSize)
	for i := range operatorsAddresses {
		_, operatorPublicKey, err := operator.GenerateKeyPair(
			local_v1.DefaultCurve,
		)
		if err != nil {
			t.Fatal(err)
		}

		operatorAddress, err := localChain.Signing().PublicKeyToAddress(
			operatorPublicKey,
		)
		if err != nil {
			t.Fatal(err)
		}

		operatorsAddresses[i] = operatorAddress
		operatorsPublicKeys[i] = operator.MarshalUncompressed(operatorPublicKey)
	}

	tests := map[string]struct {
		senderID           group.MemberIndex
		senderPublicKey    []byte
		inactiveMembersIDs []group.MemberIndex
		expectedResult     bool
	}{
		"message from another valid and operating member": {
			senderID:           group.MemberIndex(2),
			senderPublicKey:    operatorsPublicKeys[1],
			inactiveMembersIDs: []group.MemberIndex{},
			expectedResult:     true,
		},
		"message from another valid but non-operating member": {
			senderID:           group.MemberIndex(2),
			senderPublicKey:    operatorsPublicKeys[1],
			inactiveMembersIDs: []group.MemberIndex{2}, // 2 is inactive
			expectedResult:     false,
		},
		"message from self": {
			senderID:           group.MemberIndex(1),
			senderPublicKey:    operatorsPublicKeys[0],
			inactiveMembersIDs: []group.MemberIndex{},
			expectedResult:     false,
		},
		"message from another invalid member": {
			senderID:           group.MemberIndex(2),
			senderPublicKey:    operatorsPublicKeys[3],
			inactiveMembersIDs: []group.MemberIndex{},
			expectedResult:     false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			dkgGroup := group.NewGroup(groupSize-honestThreshold, groupSize)

			membershipValdator := group.NewMembershipValidator(
				&testutils.MockLogger{},
				operatorsAddresses,
				localChain.Signing(),
			)

			member := NewSigningMember(
				&testutils.MockLogger{},
				group.MemberIndex(1),
				dkgGroup,
				membershipValdator,
				"session-1",
			)

			for _, inactiveMemberID := range test.inactiveMembersIDs {
				dkgGroup.MarkMemberAsInactive(inactiveMemberID)
			}

			result := member.shouldAcceptMessage(test.senderID, test.senderPublicKey)

			testutils.AssertBoolsEqual(
				t,
				"result from message validator",
				test.expectedResult,
				result,
			)
		})
	}
}

func initializeSigningMembers(groupSize int) (
	[]*SigningMember,
	[]beaconchain.Interface,
	error,
) {
	honestThreshold := groupSize/2 + 1
	dishonestThreshold := groupSize - honestThreshold

	dkgGroup := group.NewGroup(dishonestThreshold, groupSize)

	members := make([]*SigningMember, groupSize)
	beaconChains := make([]beaconchain.Interface, groupSize)

	operatorsAddresses := make([]chain.Address, groupSize)

	for i := 0; i < groupSize; i++ {
		operatorPrivateKey, _, err := operator.GenerateKeyPair(local_v1.DefaultCurve)
		if err != nil {
			return nil, nil, err
		}

		localChain := local_v1.ConnectWithKey(
			groupSize,
			honestThreshold,
			operatorPrivateKey,
		)

		operatorAddress, err := localChain.Signing().PublicKeyToAddress(&operatorPrivateKey.PublicKey)
		if err != nil {
			return nil, nil, err
		}

		beaconChains[i] = localChain
		operatorsAddresses[i] = operatorAddress
	}

	membershipValidator := group.NewMembershipValidator(
		&testutils.MockLogger{},
		operatorsAddresses,
		beaconChains[0].Signing(),
	)

	for i := 0; i < groupSize; i++ {
		memberIndex := group.MemberIndex(i + 1)

		members[i] = NewSigningMember(
			&testutils.MockLogger{},
			memberIndex,
			dkgGroup,
			membershipValidator,
			"session-1",
		)
	}

	return members, beaconChains, nil
}
