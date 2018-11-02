package gjkr

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/pedersen"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func TestCalculateSharesAndCommitments(t *testing.T) {
	threshold := 3
	groupSize := 5

	members, err := initializeCommittingMembersGroup(threshold, groupSize)
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
		t.Fatalf("\nexpected: %v peer shares messages\nactual:   %v\n",
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

	var tests = map[string]struct {
		modifyPeerShareMessages   func(messages []*PeerSharesMessage)
		modifyCommitmentsMessages func(messages []*MemberCommitmentsMessage)
		expectedError             error
		expectedAccusedIDs        []int
	}{
		"positive validation - no accusations": {
			expectedError: nil,
		},
		"negative validation - changed share S": {
			modifyPeerShareMessages: func(messages []*PeerSharesMessage) {
				messages[0].shareS = big.NewInt(13)
			},
			expectedError:      nil,
			expectedAccusedIDs: []int{2},
		},
		"negative validation - changed two shares T": {
			modifyPeerShareMessages: func(messages []*PeerSharesMessage) {
				messages[1].shareT = big.NewInt(13)
				messages[2].shareT = big.NewInt(23)
			},
			expectedError:      nil,
			expectedAccusedIDs: []int{3, 4},
		},
		"negative validation - changed commitment": {
			modifyCommitmentsMessages: func(messages []*MemberCommitmentsMessage) {
				messages[3].commitments[1] = big.NewInt(33)
			},
			expectedError:      nil,
			expectedAccusedIDs: []int{5},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			members, err := initializeCommittingMembersGroup(threshold, groupSize)
			if err != nil {
				t.Fatalf("group initialization failed [%s]", err)
			}
			currentMember := members[0]

			var sharesMessages []*PeerSharesMessage
			var commitmentsMessages []*MemberCommitmentsMessage

			for _, member := range members {
				sharesMessage, commitmentsMessage, err := member.CalculateMembersSharesAndCommitments()
				if err != nil {
					t.Fatalf("shares and commitments calculation failed [%s]", err)
				}
				sharesMessages = append(sharesMessages, sharesMessage...)
				commitmentsMessages = append(commitmentsMessages, commitmentsMessage)
			}

			filteredSharesMessages := filterPeerSharesMessage(sharesMessages, currentMember.ID)
			filteredCommitmentsMessages := filterMemberCommitmentsMessages(commitmentsMessages, currentMember.ID)

			if test.modifyPeerShareMessages != nil {
				test.modifyPeerShareMessages(filteredSharesMessages)
			}
			if test.modifyCommitmentsMessages != nil {
				test.modifyCommitmentsMessages(filteredCommitmentsMessages)
			}

			accusedMessage, err := currentMember.VerifyReceivedSharesAndCommitmentsMessages(
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

			if len(accusedMessage.accusedIDs) != len(test.expectedAccusedIDs) {
				t.Fatalf("\nexpected: %v accusations\nactual:   %v\n",
					len(test.expectedAccusedIDs),
					len(accusedMessage.accusedIDs),
				)
			}
			if !reflect.DeepEqual(accusedMessage.accusedIDs, test.expectedAccusedIDs) {
				t.Fatalf("incorrect accused members IDs\nexpected: %v\nactual:   %v\n",
					test.expectedAccusedIDs,
					accusedMessage.accusedIDs,
				)
			}

			expectedReceivedSharesLength := groupSize - 1 - len(test.expectedAccusedIDs)
			if len(currentMember.receivedSharesS) != expectedReceivedSharesLength {
				t.Fatalf("\nexpected: %v received shares S\nactual:   %v\n",
					expectedReceivedSharesLength,
					len(currentMember.receivedSharesS),
				)
			}
			if len(currentMember.receivedSharesT) != expectedReceivedSharesLength {
				t.Fatalf("\nexpected: %v received shares T\nactual:   %v\n",
					expectedReceivedSharesLength,
					len(currentMember.receivedSharesT),
				)
			}
			if len(currentMember.receivedCommitments) != expectedReceivedSharesLength {
				t.Fatalf("\nexpected: %v received commitments\nactual:   %v\n",
					expectedReceivedSharesLength,
					len(currentMember.receivedCommitments),
				)
			}
		})
	}
}

func TestCombineReceivedShares(t *testing.T) {
	selfShareS := big.NewInt(9)
	selfShareT := big.NewInt(19)
	q := big.NewInt(59)

	receivedShareS := make(map[int]*big.Int)
	receivedShareT := make(map[int]*big.Int)
	for i := 0; i <= 5; i++ {
		receivedShareS[100+i] = big.NewInt(int64(10 + i))
		receivedShareT[100+i] = big.NewInt(int64(20 + i))
	}

	// 9 + 10 + 11 + 12 + 13 + 14 + 15 = 84 mod 59 = 25
	expectedShareS := big.NewInt(25)
	// 19 + 20 + 21 + 22 + 23 + 24 + 25 = 154 mod 59 = 36
	expectedShareT := big.NewInt(36)

	config := &DKG{Q: q}
	member := &SharingMember{
		CommittingMember: &CommittingMember{
			memberCore: &memberCore{
				protocolConfig: config,
			},
			selfSecretShareS: selfShareS,
			selfSecretShareT: selfShareT,
			receivedSharesS:  receivedShareS,
			receivedSharesT:  receivedShareT,
		},
	}

	member.CombineMemberShares()

	if member.shareS.Cmp(expectedShareS) != 0 {
		t.Errorf("incorrect combined shares S value\nexpected: %v\nactual:   %v\n",
			expectedShareS,
			member.shareS,
		)
	}
	if member.shareT.Cmp(expectedShareT) != 0 {
		t.Errorf("incorrect combined shares T value\nexpected: %v\nactual:   %v\n",
			expectedShareT,
			member.shareT,
		)
	}
}

func TestCalculatePublicCoefficients(t *testing.T) {
	secretCoefficients := []*big.Int{
		big.NewInt(3),
		big.NewInt(5),
		big.NewInt(2),
	}
	expectedPublicCoefficients := []*big.Int{
		big.NewInt(216), // 6^3 mod 1907 = 216
		big.NewInt(148), // 6^5 mod 1907 = 148
		big.NewInt(36),  // 6^2 mod 1907 = 36
	}

	config := &DKG{P: big.NewInt(1907), Q: big.NewInt(953)}

	// This test uses rand.Reader mock to get specific `g` value in `NewVSS`
	// initialization.
	mockRandomReader := testutils.NewMockRandReader(big.NewInt(6))
	vss, err := pedersen.NewVSS(mockRandomReader, config.P, config.Q)
	if err != nil {
		t.Fatalf("VSS initialization failed [%s]", err)
	}

	member := &SharingMember{
		CommittingMember: &CommittingMember{
			memberCore: &memberCore{
				protocolConfig: config,
			},
			vss:                vss,
			secretCoefficients: secretCoefficients,
		},
	}

	message := member.CalculatePublicCoefficients()

	if !reflect.DeepEqual(member.publicCoefficients, expectedPublicCoefficients) {
		t.Errorf("incorrect member's public shares\nexpected: %v\nactual:   %v\n",
			expectedPublicCoefficients,
			member.publicCoefficients,
		)
	}

	if !reflect.DeepEqual(message.publicCoefficients, expectedPublicCoefficients) {
		t.Errorf("incorrect public shares in message\nexpected: %v\nactual:   %v\n",
			expectedPublicCoefficients,
			message.publicCoefficients,
		)
	}
}

func TestCalculateAndVerifyPublicCoefficients(t *testing.T) {
	threshold := 3
	groupSize := 5

	sharingMembers, err := initializeSharingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("Group initialization failed [%s]", err)
	}

	sharingMember := sharingMembers[0]

	var tests = map[string]struct {
		modifyPublicCoefficientsMessages func(messages []*MemberPublicCoefficientsMessage)
		expectedError                    error
		expectedAccusedIDs               []int
	}{
		"positive validation - no accusations": {
			expectedError: nil,
		},
		"negative validation - changed public key share - one accused member": {
			modifyPublicCoefficientsMessages: func(messages []*MemberPublicCoefficientsMessage) {
				messages[1].publicCoefficients[1] = big.NewInt(13)
			},
			expectedError:      nil,
			expectedAccusedIDs: []int{3},
		},
		"negative validation - changed public key share - two accused members": {
			modifyPublicCoefficientsMessages: func(messages []*MemberPublicCoefficientsMessage) {
				messages[0].publicCoefficients[1] = big.NewInt(13)
				messages[3].publicCoefficients[1] = big.NewInt(18)
			},
			expectedError:      nil,
			expectedAccusedIDs: []int{2, 5},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			messages := make([]*MemberPublicCoefficientsMessage, groupSize)

			for i, m := range sharingMembers {
				messages[i] = m.CalculatePublicCoefficients()
			}

			filteredMessages := filterMemberPublicCoefficientsMessages(
				messages,
				sharingMember.ID,
			)

			if test.modifyPublicCoefficientsMessages != nil {
				test.modifyPublicCoefficientsMessages(filteredMessages)
			}

			accusedMessage, err := sharingMember.VerifyPublicCoefficients(filteredMessages)

			if !reflect.DeepEqual(test.expectedError, err) {
				t.Fatalf(
					"\nexpected: %s\nactual:   %s\n",
					test.expectedError,
					err,
				)
			}

			if !reflect.DeepEqual(accusedMessage.accusedIDs, test.expectedAccusedIDs) {
				t.Fatalf("incorrect accused IDs\nexpected: %v\nactual:   %v\n",
					test.expectedAccusedIDs,
					accusedMessage.accusedIDs,
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

func initializeCommittingMembersGroup(threshold, groupSize int) ([]*CommittingMember, error) {
	config, err := predefinedDKG()
	if err != nil {
		return nil, fmt.Errorf("DKG Config initialization failed [%s]", err)
	}

	vss, err := pedersen.NewVSS(crand.Reader, config.P, config.Q)
	if err != nil {
		return nil, fmt.Errorf("VSS initialization failed [%s]", err)
	}

	group := &Group{
		groupSize:          groupSize,
		dishonestThreshold: threshold,
	}

	var members []*CommittingMember

	for i := 1; i <= groupSize; i++ {
		id := i
		members = append(members, &CommittingMember{
			memberCore: &memberCore{
				ID:             id,
				group:          group,
				protocolConfig: config,
			},
			vss:                 vss,
			receivedSharesS:     make(map[int]*big.Int),
			receivedSharesT:     make(map[int]*big.Int),
			receivedCommitments: make(map[int][]*big.Int),
		})
		group.RegisterMemberID(id)
	}
	return members, nil
}

func initializeSharingMembersGroup(threshold, groupSize int) ([]*SharingMember, error) {
	committingMembers, err := initializeCommittingMembersGroup(threshold, groupSize)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%s]", err)
	}

	var sharingMembers []*SharingMember
	for _, cm := range committingMembers {
		cm.secretCoefficients = make([]*big.Int, threshold+1)
		for i := 0; i < threshold+1; i++ {
			cm.secretCoefficients[i], err = crand.Int(crand.Reader, cm.protocolConfig.Q)
			if err != nil {
				return nil, fmt.Errorf("secret share generation failed [%s]", err)
			}
		}
		sharingMembers = append(sharingMembers, &SharingMember{
			CommittingMember: cm,
		})
	}

	for _, sm := range sharingMembers {
		for _, cm := range committingMembers {
			sm.receivedSharesS[cm.ID] = evaluateMemberShare(sm.ID, cm.secretCoefficients)
		}
	}

	return sharingMembers, nil
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

func filterPeerSharesMessage(
	messages []*PeerSharesMessage,
	receiverID int,
) []*PeerSharesMessage {
	var result []*PeerSharesMessage
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
	receiverID int,
) []*MemberCommitmentsMessage {
	var result []*MemberCommitmentsMessage
	for _, msg := range messages {
		if msg.senderID != receiverID {
			result = append(result, msg)
		}
	}
	return result
}

func filterMemberPublicCoefficientsMessages(
	messages []*MemberPublicCoefficientsMessage, receiverID int,
) []*MemberPublicCoefficientsMessage {
	var result []*MemberPublicCoefficientsMessage
	for _, msg := range messages {
		if msg.senderID != receiverID {
			result = append(result, msg)
		}
	}
	return result
}
