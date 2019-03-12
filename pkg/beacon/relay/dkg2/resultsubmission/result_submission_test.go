package resultsubmission

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/chain/local"
)

func TestResultSigningAndVerificationRoundTrip(t *testing.T) {
	groupSize := 10
	threshold := 5
	minimumStake := big.NewInt(200)

	dkgResult := &relayChain.DKGResult{
		GroupPublicKey: []byte{10},
	}

	members, err := initializeResultSigningMembers(groupSize, threshold, minimumStake)
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

		member.otherMembersPublicKeys[message.senderIndex] = &member.privateKey.PublicKey

		if message.senderIndex != member.index {
			t.Errorf("\nexpected: %+v\nactual:   %+v\n", member.index, message.senderIndex)
		}
		if message.resultHash != expectedResultHash {
			t.Errorf("\nexpected: %+v\nactual:   %+v\n", expectedResultHash, message.resultHash)
		}

		if !member.verifySignature(message.senderIndex, expectedResultHash, message.signature) {
			t.Errorf("invalid signature")
		}

		if member.index != currentMember.index {
			messages = append(messages, message)
		}

		if len(currentMember.receivedValidResultSignatures) != 1 {
			t.Errorf(
				"\nexpected: %v\nactual:   %v\n",
				1,
				len(currentMember.receivedValidResultSignatures),
			)
		}
	}

	actualAccusations, err := currentMember.VerifyDKGResultSignatures(messages)
	if err != nil {
		t.Fatal(err)
	}

	if len(currentMember.receivedValidResultSignatures) != groupSize {
		t.Errorf(
			"\nexpected: %v\nactual:   %v\n",
			groupSize,
			len(currentMember.receivedValidResultSignatures),
		)
	}

	if len(actualAccusations) != 0 {
		t.Errorf("\nexpected: %v\nactual:   %v\n", 0, len(actualAccusations))
	}

	for _, message := range messages {
		if !bytes.Equal(currentMember.receivedValidResultSignatures[message.senderIndex],
			message.signature) {
			t.Errorf(
				"\nexpected: %x\nactual:   %x\n",
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

	members, err := initializeResultSigningMembers(groupSize, threshold, minimumStake)
	if err != nil {
		t.Fatal(err)
	}

	verifyingMember := members[0]
	verifyingMember.preferredDKGResultHash = dkgResultHash1

	member2 := members[1]
	member3 := members[2]
	member4 := members[3]
	member5 := members[4]

	signature21, _ := member2.sign(dkgResultHash1)

	signature311, _ := member3.sign(dkgResultHash1)
	signature312, _ := member3.sign(dkgResultHash1)

	signature411, _ := member4.sign(dkgResultHash1)
	signature421, _ := member4.sign(dkgResultHash2)

	signature52, _ := member5.sign(dkgResultHash2)

	var tests = map[string]struct {
		messages []*DKGResultHashSignatureMessage

		expectedReceivedValidSignatures map[MemberIndex]Signature
		expectedAccusations             map[MemberIndex][]Signature
		expectedError                   error
	}{
		"received valid messages with signatures for the preferred result": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member2.index,
					resultHash:  dkgResultHash1,
					signature:   signature21,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
				},
			},
			expectedReceivedValidSignatures: map[MemberIndex]Signature{
				member2.index: signature21,
				member3.index: signature311,
			},
			expectedAccusations: map[MemberIndex][]Signature{},
		},
		"received messages from other member with duplicated signatures for the preferred result": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature312,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
				},
			},
			expectedReceivedValidSignatures: map[MemberIndex]Signature{},
			expectedAccusations: map[MemberIndex][]Signature{
				member3.index: []Signature{signature311, signature312, signature311},
			},
		},
		"received messages from other member with signatures for two different results": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member4.index,
					resultHash:  dkgResultHash1,
					signature:   signature411,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member4.index,
					resultHash:  dkgResultHash2,
					signature:   signature421,
				},
			},
			expectedReceivedValidSignatures: map[MemberIndex]Signature{},
			expectedAccusations: map[MemberIndex][]Signature{
				member4.index: []Signature{signature411, signature421},
			},
		},
		"received a message from other member with signature for result different than preferred": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member5.index,
					resultHash:  dkgResultHash2,
					signature:   signature52,
				},
			},
			expectedReceivedValidSignatures: map[MemberIndex]Signature{},
			expectedAccusations:             map[MemberIndex][]Signature{},
		},
		"received a message from other member with invalid signature": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member2.index,
					resultHash:  dkgResultHash1,
					signature:   Signature{99},
				},
			},
			expectedReceivedValidSignatures: map[MemberIndex]Signature{},
			expectedAccusations:             map[MemberIndex][]Signature{},
		},
		"mixed cases with received valid signatures and duplicated signatures": {
			messages: []*DKGResultHashSignatureMessage{
				&DKGResultHashSignatureMessage{
					senderIndex: member2.index,
					resultHash:  dkgResultHash1,
					signature:   signature21,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature311,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member3.index,
					resultHash:  dkgResultHash1,
					signature:   signature312,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member4.index,
					resultHash:  dkgResultHash1,
					signature:   signature411,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member4.index,
					resultHash:  dkgResultHash2,
					signature:   signature421,
				},
				&DKGResultHashSignatureMessage{
					senderIndex: member5.index,
					resultHash:  dkgResultHash2,
					signature:   signature52,
				},
			},
			expectedReceivedValidSignatures: map[MemberIndex]Signature{
				member2.index: signature21,
			},
			expectedAccusations: map[MemberIndex][]Signature{
				member3.index: []Signature{signature311, signature312},
				member4.index: []Signature{signature411, signature421},
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			verifyingMember.receivedValidResultSignatures = make(map[MemberIndex]Signature)

			actualAccusations, err := verifyingMember.VerifyDKGResultSignatures(test.messages)

			if !reflect.DeepEqual(err, test.expectedError) {
				t.Errorf("\nexpected: %v\nactual:   %v\n", test.expectedError, err)
			}

			if !reflect.DeepEqual(actualAccusations, test.expectedAccusations) {
				t.Errorf(
					"\nexpected: %+v\nactual:   %+v\n",
					test.expectedAccusations,
					actualAccusations,
				)
			}

			if !reflect.DeepEqual(
				verifyingMember.receivedValidResultSignatures,
				test.expectedReceivedValidSignatures,
			) {
				t.Errorf(
					"\nexpected: %v\nactual:   %v\n",
					test.expectedReceivedValidSignatures,
					verifyingMember.receivedValidResultSignatures,
				)
			}
		})
	}
}

func initializeResultSigningMembers(groupSize, threshold int, minimumStake *big.Int) ([]*ResultSigningMember, error) {
	chainHandle := local.Connect(groupSize, threshold, minimumStake)

	privateKeys := make(map[int]*ecdsa.PrivateKey)
	for i := 1; i <= groupSize; i++ {
		privateKey, err := crypto.GenerateKey() // TODO: Replace with static.GenerateKey
		if err != nil {
			return nil, err
		}
		privateKeys[i] = privateKey
	}

	members := make([]*ResultSigningMember, 0)
	for i := 1; i <= groupSize; i++ {
		peerMemberPublicKeys := make(map[MemberIndex]*ecdsa.PublicKey)

		for j := 1; j <= groupSize; j++ {
			if i != j {
				peerMemberPublicKeys[MemberIndex(j)] = &privateKeys[j].PublicKey
			}
		}

		members = append(members, &ResultSigningMember{
			index:                         MemberIndex(i),
			chainHandle:                   chainHandle,
			privateKey:                    privateKeys[i],
			otherMembersPublicKeys:        peerMemberPublicKeys,
			receivedValidResultSignatures: make(map[MemberIndex]Signature),
		})
	}

	return members, nil
}
