package result

import (
	"crypto/ecdsa"
	"math/big"
	"reflect"
	"testing"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/operator"
)

// TestResultSigningAndVerificationRoundTrip simulates Phase 13 execution when
// a group of members produces messages containing signatures and than one member
// verifies messages received from other group members.
func TestResultSigningAndVerificationRoundTrip(t *testing.T) {
	groupSize := 10

	dkgResult := &relayChain.DKGResult{
		GroupPublicKey: []byte{10},
	}

	members, chainHandles, err := initializeSigningMembers(groupSize)
	if err != nil {
		t.Fatal(err)
	}

	currentMember := members[0]
	currentSigning := chainHandles[0].Signing()

	messages := make([]*DKGResultHashSignatureMessage, 0)

	for i, member := range members {
		message, err := member.SignDKGResult(
			dkgResult,
			chainHandles[i].ThresholdRelay(),
			chainHandles[i].Signing(),
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

	dkgResultHash1 := relayChain.DKGResultHash{10}
	dkgResultHash2 := relayChain.DKGResultHash{20}

	members, chainHandles, err := initializeSigningMembers(groupSize)
	if err != nil {
		t.Fatal(err)
	}

	verifyingMember, verifyingMemberSigning := members[0], chainHandles[0].Signing()
	verifyingMember.preferredDKGResultHash = dkgResultHash1

	selfSignature, _ := verifyingMemberSigning.Sign(dkgResultHash1[:])
	verifyingMember.selfDKGResultSignature = selfSignature

	member2, signing2 := members[1], chainHandles[1].Signing()
	member3, signing3 := members[2], chainHandles[2].Signing()
	member4, signing4 := members[3], chainHandles[3].Signing()
	member5, signing5 := members[4], chainHandles[4].Signing()

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
				{
					senderIndex: member2.index,
					resultHash:  dkgResultHash1,
					signature:   signature21,
					publicKey:   signing2.PublicKey(),
				},
				{
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
				{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   signing3.PublicKey(),
				},
				{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature312,
					publicKey:   signing3.PublicKey(),
				},
				{
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
				{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   signing3.PublicKey(),
				},
				{
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
				{
					senderIndex: member4.index,
					resultHash:  dkgResultHash1,
					signature:   signature411,
					publicKey:   signing4.PublicKey(),
				},
				{
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
				{
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
				{
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
				{
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
				{
					senderIndex: member2.index,
					resultHash:  dkgResultHash1,
					signature:   signature21,
					publicKey:   signing2.PublicKey(),
				},
				// Multiple signatures from the same member supporting the same result as preferred.
				{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   signing3.PublicKey(),
				},
				{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature312,
					publicKey:   signing3.PublicKey(),
				},
				// Multiple signatures from the same member supporting two different results.
				{
					senderIndex: member4.index,
					resultHash:  dkgResultHash1,
					signature:   signature411,
					publicKey:   signing4.PublicKey(),
				},
				{
					senderIndex: member4.index,
					resultHash:  dkgResultHash2,
					signature:   signature421,
					publicKey:   signing4.PublicKey(),
				},
				// Member supporting different result than preferred.
				{
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

func initializeSigningMembers(groupSize int) ([]*SigningMember, []chain.Handle, error) {
	honestThreshold := groupSize/2 + 1
	dishonestThreshold := groupSize - honestThreshold
	minimumStake := big.NewInt(200)

	dkgGroup := group.NewDkgGroup(dishonestThreshold, groupSize)

	members := make([]*SigningMember, groupSize)
	chainHandles := make([]chain.Handle, groupSize)

	for i := 0; i < groupSize; i++ {
		memberIndex := group.MemberIndex(i + 1)

		members[i] = NewSigningMember(
			memberIndex,
			dkgGroup,
			&mockMembershipValidator{},
		)

		privateKey, _, err := operator.GenerateKeyPair()
		if err != nil {
			return nil, nil, err
		}

		chainHandles[i] = local.ConnectWithKey(
			groupSize,
			honestThreshold,
			minimumStake,
			privateKey,
		)
	}

	return members, chainHandles, nil
}

type mockMembershipValidator struct{}

func (mmv *mockMembershipValidator) IsInGroup(
	publicKey *ecdsa.PublicKey,
) bool {
	return true
}

func (mmv *mockMembershipValidator) IsValidMembership(
	memberID group.MemberIndex,
	publicKey []byte,
) bool {
	return true
}
