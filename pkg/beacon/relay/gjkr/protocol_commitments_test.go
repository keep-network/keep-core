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
	sharesMessage, commitmentsMessage, err := member.CalculateMembersSharesAndCommitments()
	if err != nil {
		t.Fatalf("shares and commitments calculation failed [%s]", err)
	}

	if len(member.secretCoefficients) != (threshold + 1) {
		t.Fatalf("\nexpected: %v secret coefficients\nactual:   %v\n",
			threshold+1,
			len(member.secretCoefficients),
		)
	}
	if len(sharesMessage.shares) != (groupSize - 1) {
		t.Fatalf("\nexpected: %v shares in message\nactual:   %v\n",
			groupSize-1,
			len(sharesMessage.shares),
		)
	}

	if len(commitmentsMessage.commitments) != (threshold + 1) {
		t.Fatalf("\nexpected: %v calculated commitments\nactual:   %v\n",
			threshold+1,
			len(commitmentsMessage.commitments),
		)
	}
}

func TestStoreSharesMessageForEvidence(t *testing.T) {
	groupSize := 2

	config, err := predefinedDKG()
	if err != nil {
		t.Fatalf("predefined config initialization failed [%s]", err)
	}

	members, err := initializeCommittingMembersGroup(
		groupSize, // threshold = group size
		groupSize,
		config,
	)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}

	member1 := members[0]
	member2 := members[1]

	sharesMsg1, commitmentsMsg1, err := member1.CalculateMembersSharesAndCommitments()
	if err != nil {
		t.Fatal(err)
	}

	if _, _, err := member2.CalculateMembersSharesAndCommitments(); err != nil {
		t.Fatal(err)
	}

	verifyingMember2 := member2.InitializeCommitmentsVerification()

	if _, err := verifyingMember2.VerifyReceivedSharesAndCommitmentsMessages(
		[]*PeerSharesMessage{sharesMsg1},
		[]*MemberCommitmentsMessage{commitmentsMsg1},
	); err != nil {
		t.Fatal(err)
	}

	evidenceMsg := verifyingMember2.evidenceLog.peerSharesMessage(member1.ID)

	if !reflect.DeepEqual(sharesMsg1, evidenceMsg) {
		t.Fatalf(
			"unexpected message in evidence log\nexpected: %v\n actual:   %v",
			sharesMsg1,
			evidenceMsg,
		)
	}
}

func TestSharesAndCommitmentsCalculationAndVerification(t *testing.T) {
	threshold := 2
	groupSize := 3

	config, err := predefinedDKG()
	if err != nil {
		t.Fatalf("predefined config initialization failed [%s]", err)
	}

	members, err := initializeCommittingMembersGroup(threshold, groupSize, config)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}

	member1 := members[0]
	member2 := members[1]
	member3 := members[2]

	verifyingMemberID := member3.ID
	verifyingMemberKeys := member3.symmetricKeys

	var tests = map[string]struct {
		modifyPeerSharesMessage  func(messages map[MemberID]*PeerSharesMessage) error
		modifyCommitmentsMessage func(messages map[MemberID]*MemberCommitmentsMessage)
		expectedAccusedIDs       []MemberID
	}{
		"no accusations": {
			expectedAccusedIDs: []MemberID{},
		},
		"invalid S share": {
			modifyPeerSharesMessage: func(messages map[MemberID]*PeerSharesMessage) error {
				return alterPeerSharesMessage(
					messages[member2.ID],
					verifyingMemberID,
					verifyingMemberKeys[member2.ID],
					true,
					false,
					config,
				)
			},
			expectedAccusedIDs: []MemberID{member2.ID},
		},
		"invalid T share": {
			modifyPeerSharesMessage: func(messages map[MemberID]*PeerSharesMessage) error {
				return alterPeerSharesMessage(
					messages[member1.ID],
					verifyingMemberID,
					verifyingMemberKeys[member1.ID],
					false,
					true,
					config,
				)
			},
			expectedAccusedIDs: []MemberID{member1.ID},
		},
		"invalid commitment": {
			modifyCommitmentsMessage: func(messages map[MemberID]*MemberCommitmentsMessage) {
				message := messages[member2.ID]
				message.commitments[0] = testutils.NewRandInt(
					message.commitments[0],
					config.Q,
				)
			},
			expectedAccusedIDs: []MemberID{member2.ID},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			shareMessages := make(map[MemberID]*PeerSharesMessage)
			commitmentMessages := make(map[MemberID]*MemberCommitmentsMessage)

			for _, member := range members {
				shares, commitments, err := member.CalculateMembersSharesAndCommitments()
				if err != nil {
					t.Fatal(err)
				}

				shareMessages[member.ID] = shares
				commitmentMessages[member.ID] = commitments
			}

			if test.modifyPeerSharesMessage != nil {
				if err = test.modifyPeerSharesMessage(shareMessages); err != nil {
					t.Fatal(err)
				}

			}

			if test.modifyCommitmentsMessage != nil {
				test.modifyCommitmentsMessage(commitmentMessages)
			}

			verifyingMember := member3.InitializeCommitmentsVerification()

			accusationMessage, err := verifyingMember.VerifyReceivedSharesAndCommitmentsMessages(
				[]*PeerSharesMessage{
					shareMessages[member1.ID],
					shareMessages[member2.ID],
				},
				[]*MemberCommitmentsMessage{
					commitmentMessages[member1.ID],
					commitmentMessages[member2.ID],
				},
			)
			if err != nil {
				t.Fatal(err)
			}

			assertAccusedMembers(
				test.expectedAccusedIDs,
				verifyingMember,
				accusationMessage,
				t,
			)

			assertValidSharesAndCommitments(
				test.expectedAccusedIDs,
				verifyingMember,
				groupSize,
				t,
			)
		})
	}
}

func alterPeerSharesMessage(
	message *PeerSharesMessage,
	receiverID MemberID,
	symmetricKey ephemeral.SymmetricKey,
	alterS bool,
	alterT bool,
	config *DKG,
) error {
	oldShareS, err := message.decryptShareS(receiverID, symmetricKey)
	if err != nil {
		return err
	}

	oldShareT, err := message.decryptShareT(receiverID, symmetricKey)
	if err != nil {
		return err
	}

	var newShareS = oldShareS
	var newShareT = oldShareT

	if alterS {
		newShareS = testutils.NewRandInt(oldShareS, config.Q)
	}
	if alterT {
		newShareT = testutils.NewRandInt(oldShareT, config.Q)
	}

	err = message.addShares(receiverID, newShareS, newShareT, symmetricKey)
	if err != nil {
		return err
	}

	return nil
}

func assertAccusedMembers(
	expectedAccusedIDs []MemberID,
	verifyingMember *CommitmentsVerifyingMember,
	accusationMessage *SecretSharesAccusationsMessage,
	t *testing.T,
) {
	expectedAccusedMembersKeys := make(map[MemberID]*ephemeral.PrivateKey)
	for _, id := range expectedAccusedIDs {
		expectedAccusedMembersKeys[id] = verifyingMember.ephemeralKeyPairs[id].PrivateKey
	}

	if !reflect.DeepEqual(accusationMessage.accusedMembersKeys, expectedAccusedMembersKeys) {
		t.Errorf("incorrect accused members IDs\nexpected: %v\nactual:   %v\n",
			expectedAccusedMembersKeys,
			accusationMessage.accusedMembersKeys,
		)
	}
}

func assertValidSharesAndCommitments(
	expectedAccusedIDs []MemberID,
	verifyingMember *CommitmentsVerifyingMember,
	groupSize int,
	t *testing.T,
) {
	expectedReceivedSharesLength := groupSize - 1 - len(expectedAccusedIDs)
	if len(verifyingMember.receivedValidSharesS) != expectedReceivedSharesLength {
		t.Errorf("\nexpected: %v received shares S\nactual:   %v\n",
			expectedReceivedSharesLength,
			len(verifyingMember.receivedValidSharesS),
		)
	}
	if len(verifyingMember.receivedValidSharesT) != expectedReceivedSharesLength {
		t.Errorf("\nexpected: %v received shares T\nactual:   %v\n",
			expectedReceivedSharesLength,
			len(verifyingMember.receivedValidSharesT),
		)
	}
	if len(verifyingMember.receivedValidPeerCommitments) != expectedReceivedSharesLength {
		t.Errorf("\nexpected: %v received commitments\nactual:   %v\n",
			expectedReceivedSharesLength,
			len(verifyingMember.receivedValidPeerCommitments),
		)
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

	symmetricKeyMembers, err := generateGroupWithEphemeralKeys(
		threshold,
		groupSize,
		dkg,
	)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%v]", err)
	}

	var members []*CommittingMember
	for _, member := range symmetricKeyMembers {
		committingMember := member.InitializeCommitting()
		members = append(members, committingMember)
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
		members = append(members, member.InitializeCommitmentsVerification())
	}

	return members, nil
}

// predefinedDKGconfig initializes DKG configuration with predefined 512-bit
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

	vss, err := pedersen.GenerateVSS(crand.Reader, p, q)
	if err != nil {
		return nil, fmt.Errorf("could not generate DKG paramters [%v]", err)
	}

	return &DKG{p, q, vss}, nil
}

func filterPeerSharesMessage(
	messages []*PeerSharesMessage,
	receiverID MemberID,
) []*PeerSharesMessage {
	var result []*PeerSharesMessage
	for _, msg := range messages {
		if msg.senderID != receiverID {
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
