package result

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/member"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/operator"
)

// TestResultSigningAndVerificationRoundTrip simulates Phase 13 execution when
// a group of members produces messages containing signatures and than one member
// verifies messages received from other group members.
func TestResultSigningAndVerificationRoundTrip(t *testing.T) {
	groupSize := 10
	threshold := 5
	minimumStake := big.NewInt(200)

	dkgResult := &relayChain.DKGResult{
		GroupPublicKey: []byte{10},
	}

	members, err := initializeSigningMembers(groupSize, threshold, minimumStake)
	if err != nil {
		t.Fatal(err)
	}

	expectedResultHash, err := members[0].chainHandle.ThresholdRelay().CalculateDKGResultHash(dkgResult)
	if err != nil {
		t.Fatal(err)
	}

	currentMember := members[0]
	messages := make([]*DKGResultHashSignatureMessage, 0)

	for _, member := range members {
		message, err := member.SignDKGResult(dkgResult)
		if err != nil {
			t.Fatal(err)
		}

		if message.senderIndex != member.index {
			t.Errorf(
				"unexpected mesage sender index\nexpected: %+v\nactual:   %+v\n",
				member.index,
				message.senderIndex,
			)
		}
		if message.resultHash != expectedResultHash {
			t.Errorf(
				"unexpected result hash\nexpected: %+v\nactual:   %+v\n",
				expectedResultHash,
				message.resultHash,
			)
		}

		err = operator.VerifySignature(
			message.publicKey,
			expectedResultHash[:],
			message.signature,
		)
		if err != nil {
			t.Errorf(
				"invalid signature [%v]\nsignature:  %v\npublic key: %+v\n",
				message.signature,
				message.publicKey,
				err,
			)
		}

		// Don't register message from self.
		if member.index != currentMember.index {
			messages = append(messages, message)
		}

		if len(currentMember.receivedValidResultSignatures) != 1 {
			t.Errorf(
				"unexpected number of registered valid signatures\nexpected: %v\nactual:   %v\n",
				1,
				len(currentMember.receivedValidResultSignatures),
			)
		}
	}

	err = currentMember.VerifyDKGResultSignatures(messages)
	if err != nil {
		t.Fatal(err)
	}

	if len(currentMember.receivedValidResultSignatures) != groupSize {
		t.Errorf(
			"unexpected number of registered valid signatures\nexpected: %v\nactual:   %v\n",
			groupSize,
			len(currentMember.receivedValidResultSignatures),
		)
	}

	for _, message := range messages {
		if !bytes.Equal(currentMember.receivedValidResultSignatures[message.senderIndex],
			message.signature) {
			t.Errorf(
				"unexpected registered signature for member %d\nexpected: %x\nactual:   %x\n",
				message.senderIndex,
				message.signature,
				currentMember.receivedValidResultSignatures[message.senderIndex],
			)
		}
	}
}

func TestVerifyDKGResultSignatures(t *testing.T) {
	threshold := 3
	groupSize := 5
	minimumStake := big.NewInt(200)

	dkgResultHash1 := relayChain.DKGResultHash{10}
	dkgResultHash2 := relayChain.DKGResultHash{20}

	members, err := initializeSigningMembers(groupSize, threshold, minimumStake)
	if err != nil {
		t.Fatal(err)
	}

	verifyingMember := members[0]
	verifyingMember.preferredDKGResultHash = dkgResultHash1

	member2 := members[1]
	member3 := members[2]
	member4 := members[3]
	member5 := members[4]

	signature21, _ := operator.Sign(dkgResultHash1[:], member2.privateKey)

	signature311, _ := operator.Sign(dkgResultHash1[:], member3.privateKey)
	signature312, _ := operator.Sign(dkgResultHash1[:], member3.privateKey)

	signature411, _ := operator.Sign(dkgResultHash1[:], member4.privateKey)
	signature421, _ := operator.Sign(dkgResultHash2[:], member4.privateKey)

	signature52, _ := operator.Sign(dkgResultHash2[:], member5.privateKey)

	var tests = map[string]struct {
		messages []*DKGResultHashSignatureMessage

		expectedReceivedValidSignatures map[member.Index]operator.Signature
		expectedError                   error
	}{
		"received valid messages with signatures for the preferred result": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member2.index,
					resultHash:  dkgResultHash1,
					signature:   signature21,
					publicKey:   &member2.privateKey.PublicKey,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   &member3.privateKey.PublicKey,
				},
			},
			expectedReceivedValidSignatures: map[member.Index]operator.Signature{
				member2.index: signature21,
				member3.index: signature311,
			},
		},
		"received messages from other member with duplicated different signatures for the preferred result": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   &member3.privateKey.PublicKey,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature312,
					publicKey:   &member3.privateKey.PublicKey,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   &member3.privateKey.PublicKey,
				},
			},
			expectedReceivedValidSignatures: map[member.Index]operator.Signature{},
		},
		"received messages from other member with the same signatures for the preferred result": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   &member3.privateKey.PublicKey,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   &member3.privateKey.PublicKey,
				},
			},
			expectedReceivedValidSignatures: map[member.Index]operator.Signature{},
		},
		"received messages from other member with signatures for two different results": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member4.index,
					resultHash:  dkgResultHash1,
					signature:   signature411,
					publicKey:   &member4.privateKey.PublicKey,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member4.index,
					resultHash:  dkgResultHash2,
					signature:   signature421,
					publicKey:   &member4.privateKey.PublicKey,
				},
			},
			expectedReceivedValidSignatures: map[member.Index]operator.Signature{},
		},
		"received a message from other member with signature for result different than preferred": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member5.index,
					resultHash:  dkgResultHash2,
					signature:   signature52,
					publicKey:   &member5.privateKey.PublicKey,
				},
			},
			expectedReceivedValidSignatures: map[member.Index]operator.Signature{},
		},
		"received a message from other member with invalid signature": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member2.index,
					resultHash:  dkgResultHash1,
					signature:   operator.Signature{99},
					publicKey:   &member2.privateKey.PublicKey,
				},
			},
			expectedReceivedValidSignatures: map[member.Index]operator.Signature{},
		},
		"received a message from other member with invalid public key": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member2.index,
					resultHash:  dkgResultHash1,
					signature:   signature21,
					publicKey:   &members[0].privateKey.PublicKey,
				},
			},
			expectedReceivedValidSignatures: map[member.Index]operator.Signature{},
		},
		"mixed cases with received valid signatures and duplicated signatures": {
			messages: []*DKGResultHashSignatureMessage{
				// Valid signature supporting the same result as preferred.
				&DKGResultHashSignatureMessage{
					senderIndex: member2.index,
					resultHash:  dkgResultHash1,
					signature:   signature21,
					publicKey:   &member2.privateKey.PublicKey,
				},
				// Multiple signatures from the same member supporting the same result as preferred.
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
					publicKey:   &member3.privateKey.PublicKey,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature312,
					publicKey:   &member3.privateKey.PublicKey,
				},
				// Multiple signatures from the same member supporting two different results.
				&DKGResultHashSignatureMessage{
					senderIndex: member4.index,
					resultHash:  dkgResultHash1,
					signature:   signature411,
					publicKey:   &member4.privateKey.PublicKey,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member4.index,
					resultHash:  dkgResultHash2,
					signature:   signature421,
					publicKey:   &member4.privateKey.PublicKey,
				},
				// Member supporting different result than preferred.
				&DKGResultHashSignatureMessage{
					senderIndex: member5.index,
					resultHash:  dkgResultHash2,
					signature:   signature52,
					publicKey:   &member5.privateKey.PublicKey,
				},
			},
			expectedReceivedValidSignatures: map[member.Index]operator.Signature{
				member2.index: signature21,
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			verifyingMember.receivedValidResultSignatures = make(map[member.Index]operator.Signature)

			err := verifyingMember.VerifyDKGResultSignatures(test.messages)

			if !reflect.DeepEqual(err, test.expectedError) {
				t.Errorf(
					"unexpected error\nexpected: %v\nactual:   %v\n",
					test.expectedError,
					err,
				)
			}

			if !reflect.DeepEqual(
				verifyingMember.receivedValidResultSignatures,
				test.expectedReceivedValidSignatures,
			) {
				t.Errorf(
					"unexpected registered received valid signatures\nexpected: %v\nactual:   %v\n",
					test.expectedReceivedValidSignatures,
					verifyingMember.receivedValidResultSignatures,
				)
			}
		})
	}
}

func initializeSigningMembers(
	groupSize int,
	threshold int,
	minimumStake *big.Int,
) ([]*SigningMember, error) {
	chainHandle := local.Connect(groupSize, threshold, minimumStake)

	members := make([]*SigningMember, 0)
	for i := 1; i <= groupSize; i++ {
		privateKey, _, err := operator.GenerateKeyPair()
		if err != nil {
			return nil, err
		}

		members = append(members, &SigningMember{
			index:                         member.Index(i),
			chainHandle:                   chainHandle,
			privateKey:                    privateKey,
			receivedValidResultSignatures: make(map[member.Index]operator.Signature),
		})
	}

	return members, nil
}
