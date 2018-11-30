package gjkr

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/pedersen"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func TestCalculateSharesAndCommitments(t *testing.T) {
	threshold := 3
	groupSize := 5

	members, err := initializeCommittingMembersGroup(threshold, groupSize, nil)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}

	member := members[0]
	sharesMessages, commitmentsMessage, err := member.CalculateMembersSharesAndCommitments()
	if err != nil {
		t.Fatalf("shares and commitments calculation failed [%s]", err)
	}

	if len(member.secretCoefficients) != (threshold + 1) {
		t.Fatalf("\nexpected: %v secret coefficients\nactual:   %v\n",
			threshold+1,
			len(member.secretCoefficients),
		)
	}
	if len(sharesMessages) != (groupSize - 1) {
		t.Fatalf("\nexpected: %v other members shares messages\nactual:   %v\n",
			groupSize-1,
			len(sharesMessages),
		)
	}

	if len(commitmentsMessage.commitments) != (threshold + 1) {
		t.Fatalf("\nexpected: %v calculated commitments\nactual:   %v\n",
			threshold+1,
			len(commitmentsMessage.commitments),
		)
	}
}

func TestSharesAndCommitmentsCalculationAndVerification(t *testing.T) {
	threshold := 3
	groupSize := 5

	config, err := predefinedDKG()
	if err != nil {
		t.Fatalf("predefined config initialization failed [%s]", err)
	}

	var alterOtherMemberSharesMessage = func(
		message *OtherMemberSharesMessage, symmetricKey ephemeral.SymmetricKey,
		alterS bool,
		alterT bool,
	) *OtherMemberSharesMessage {
		oldShareS, err := message.decryptShareS(symmetricKey)
		if err != nil {
			t.Fatal(err)
		}

		oldShareT, err := message.decryptShareT(symmetricKey)
		if err != nil {
			t.Fatal(err)
		}

		var newShareS = oldShareS
		var newShareT = oldShareT

		if alterS {
			newShareS = testutils.NewRandInt(oldShareS, config.Q)
		}
		if alterT {
			newShareT = testutils.NewRandInt(oldShareT, config.Q)
		}

		msg, err := newOtherMemberSharesMessage(
			message.senderID,
			message.receiverID,
			newShareS,
			newShareT,
			symmetricKey,
		)
		if err != nil {
			t.Fatal(err)
		}

		return msg
	}

	var tests = map[string]struct {
		modifyOtherMemberShareMessages func(
			messages []*OtherMemberSharesMessage, symmetricKeys map[MemberID]ephemeral.SymmetricKey,
		)
		modifyCommitmentsMessages func(messages []*MemberCommitmentsMessage)
		expectedError             error
		expectedAccusedIDs        []MemberID
	}{
		"positive validation - no accusations": {
			expectedError: nil,
		},
		"negative validation - changed share S": {
			modifyOtherMemberShareMessages: func(
				messages []*OtherMemberSharesMessage, symmetricKeys map[MemberID]ephemeral.SymmetricKey,
			) {
				// current member ID = 1, we modify first message on the list
				// so it's a message from member with ID = 2
				messages[0] = alterOtherMemberSharesMessage(
					messages[0],
					symmetricKeys[messages[0].senderID],
					true,
					false,
				)
			},
			expectedError:      nil,
			expectedAccusedIDs: []MemberID{MemberID(2)},
		},
		"negative validation - changed two shares T": {
			modifyOtherMemberShareMessages: func(
				messages []*OtherMemberSharesMessage, symmetricKeys map[MemberID]ephemeral.SymmetricKey,
			) {
				// current member ID = 1, we modify second message on the list
				// so it's a message from member with ID = 3
				messages[1] = alterOtherMemberSharesMessage(
					messages[1],
					symmetricKeys[messages[1].senderID],
					false,
					true,
				)

				// current member ID = 1, we modify third message on the list
				// so it's a message from member with ID = 4
				messages[2] = alterOtherMemberSharesMessage(
					messages[2],
					symmetricKeys[messages[2].senderID],
					false,
					true,
				)
			},
			expectedError:      nil,
			expectedAccusedIDs: []MemberID{MemberID(3), MemberID(4)},
		},
		"negative validation - changed commitment": {
			modifyCommitmentsMessages: func(messages []*MemberCommitmentsMessage) {
				messages[3].commitments[1] = testutils.NewRandInt(
					messages[3].commitments[1], config.Q,
				)
			},
			expectedError:      nil,
			expectedAccusedIDs: []MemberID{MemberID(5)},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			members, err := initializeCommittingMembersGroup(threshold, groupSize, nil)
			if err != nil {
				t.Fatalf("group initialization failed [%s]", err)
			}
			currentMember := members[0]

			var sharesMessages []*OtherMemberSharesMessage
			var commitmentsMessages []*MemberCommitmentsMessage

			for _, member := range members {
				sharesMessage, commitmentsMessage, err := member.CalculateMembersSharesAndCommitments()
				if err != nil {
					t.Fatalf("shares and commitments calculation failed [%s]", err)
				}
				sharesMessages = append(sharesMessages, sharesMessage...)
				commitmentsMessages = append(commitmentsMessages, commitmentsMessage)
			}

			filteredSharesMessages := filterOtherMemberSharesMessage(sharesMessages, currentMember.ID)
			filteredCommitmentsMessages := filterMemberCommitmentsMessages(commitmentsMessages, currentMember.ID)

			if test.modifyOtherMemberShareMessages != nil {
				test.modifyOtherMemberShareMessages(filteredSharesMessages, currentMember.symmetricKeys)
			}
			if test.modifyCommitmentsMessages != nil {
				test.modifyCommitmentsMessages(filteredCommitmentsMessages)
			}

			// Simulate step to next Phase
			// TODO Handle by Next function
			verifyingMember := &CommitmentsVerifyingMember{
				CommittingMember:                    currentMember,
				receivedValidSharesS:                make(map[MemberID]*big.Int),
				receivedValidSharesT:                make(map[MemberID]*big.Int),
				receivedValidOtherMemberCommitments: make(map[MemberID][]*big.Int),
			}

			accusedMessage, err := verifyingMember.VerifyReceivedSharesAndCommitmentsMessages(
				filteredSharesMessages,
				filteredCommitmentsMessages,
			)

			if !reflect.DeepEqual(test.expectedError, err) {
				t.Fatalf(
					"\nexpected: %v\nactual:   %v\n",
					test.expectedError,
					err,
				)
			}

			if len(accusedMessage.accusedMembersKeys) != len(test.expectedAccusedIDs) {
				t.Fatalf("\nexpected: %v accusations\nactual:   %v\n",
					len(test.expectedAccusedIDs),
					len(accusedMessage.accusedMembersKeys),
				)
			}

			expectedAccusedMembersKeys := make(map[MemberID]*ephemeral.PrivateKey)
			for _, id := range test.expectedAccusedIDs {
				expectedAccusedMembersKeys[id] = verifyingMember.ephemeralKeyPairs[id].PrivateKey
			}

			if !reflect.DeepEqual(accusedMessage.accusedMembersKeys, expectedAccusedMembersKeys) {
				t.Fatalf("incorrect accused members IDs\nexpected: %v\nactual:   %v\n",
					expectedAccusedMembersKeys,
					accusedMessage.accusedMembersKeys,
				)
			}

			expectedReceivedSharesLength := groupSize - 1 - len(test.expectedAccusedIDs)
			if len(verifyingMember.receivedValidSharesS) != expectedReceivedSharesLength {
				t.Fatalf("\nexpected: %v received shares S\nactual:   %v\n",
					expectedReceivedSharesLength,
					len(verifyingMember.receivedValidSharesS),
				)
			}
			if len(verifyingMember.receivedValidSharesT) != expectedReceivedSharesLength {
				t.Fatalf("\nexpected: %v received shares T\nactual:   %v\n",
					expectedReceivedSharesLength,
					len(verifyingMember.receivedValidSharesT),
				)
			}
			if len(verifyingMember.receivedValidOtherMemberCommitments) != expectedReceivedSharesLength {
				t.Fatalf("\nexpected: %v received commitments\nactual:   %v\n",
					expectedReceivedSharesLength,
					len(verifyingMember.receivedValidOtherMemberCommitments),
				)
			}
		})
	}
}
func TestGeneratePolynomial(t *testing.T) {
	degree := 3
	config := &DKG{P: big.NewInt(100), Q: big.NewInt(9)}

	coefficients, err := generatePolynomial(degree, config)
	if err != nil {
		t.Fatalf("unexpected error [%s]", err)
	}

	if len(coefficients) != degree+1 {
		t.Fatalf("\nexpected: %d coefficients\nactual:   %d\n",
			degree+1,
			len(coefficients),
		)
	}
	for _, c := range coefficients {
		if c.Sign() <= 0 || c.Cmp(config.Q) >= 0 {
			t.Fatalf("coefficient out of range\nexpected: 0 < value < %d\nactual:   %v\n",
				config.Q,
				c,
			)
		}
	}
}

func initializeCommittingMembersGroup(threshold, groupSize int, dkg *DKG) ([]*CommittingMember, error) {
	var err error
	if dkg == nil {
		dkg, err = predefinedDKG()
		if err != nil {
			return nil, fmt.Errorf("DKG Config initialization failed [%v]", err)
		}
	}

	symmetricKeyMembers, err := initializeSymmetricKeyMembersGroup(
		threshold,
		groupSize,
		dkg,
	)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%v]", err)
	}

	vss, err := pedersen.NewVSS(crand.Reader, dkg.P, dkg.Q)
	if err != nil {
		return nil, fmt.Errorf("VSS initialization failed [%v]", err)
	}

	var members []*CommittingMember
	for _, member := range symmetricKeyMembers {
		members = append(members,
			&CommittingMember{
				SymmetricKeyGeneratingMember: member,
				vss:                          vss,
			})
	}

	return members, nil
}

func initializeCommitmentsVerifiyingMembersGroup(threshold, groupSize int, dkg *DKG) ([]*CommitmentsVerifyingMember, error) {
	committingMembers, err := initializeCommittingMembersGroup(threshold, groupSize, dkg)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%v]", err)
	}

	var members []*CommitmentsVerifyingMember
	for _, member := range committingMembers {
		//TODO Handle by Next function
		members = append(members,
			&CommitmentsVerifyingMember{
				CommittingMember:                    member,
				receivedValidSharesS:                make(map[MemberID]*big.Int),
				receivedValidSharesT:                make(map[MemberID]*big.Int),
				receivedValidOtherMemberCommitments: make(map[MemberID][]*big.Int),
			})
	}

	return members, nil
}

// predefinedDKGconfig initializez DKG configuration with predefined 512-bit
// p and q values.
func predefinedDKG() (*DKG, error) {
	// `p` is 512-bit safe prime.
	pStr := "0xde41693a1522be3f2c14113e26ec7bea81f19d15095fbc0d0aca845ce086537535c6f9bdf4c4e3ac0526f3cf8c064c11483beddbc29464a9baaf6bb7ae5a024b"
	// `q` is 511-bit Sophie Germain prime.
	qStr := "0x6f20b49d0a915f1f960a089f13763df540f8ce8a84afde068565422e704329ba9ae37cdefa6271d6029379e7c6032608a41df6ede14a3254dd57b5dbd72d0125"

	var result bool

	p, result := new(big.Int).SetString(pStr, 0)
	if !result {
		return nil, fmt.Errorf("failed to initialize p")
	}

	q, result := new(big.Int).SetString(qStr, 0)
	if !result {
		return nil, fmt.Errorf("failed to initialize q")
	}
	return &DKG{p, q}, nil
}

func filterOtherMemberSharesMessage(
	messages []*OtherMemberSharesMessage, receiverID MemberID,
) []*OtherMemberSharesMessage {
	var result []*OtherMemberSharesMessage
	for _, msg := range messages {
		if msg.senderID != receiverID &&
			msg.receiverID == receiverID {
			result = append(result, msg)
		}
	}
	return result
}

func filterMemberCommitmentsMessages(
	messages []*MemberCommitmentsMessage,
	receiverID MemberID,
) []*MemberCommitmentsMessage {
	var result []*MemberCommitmentsMessage
	for _, msg := range messages {
		if msg.senderID != receiverID {
			result = append(result, msg)
		}
	}
	return result
}
