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
		modifyPeerShareMessages   func(messages []*PeerSharesMessage) []*PeerSharesMessage
		modifyCommitmentsMessages func(messages []*MemberCommitmentsMessage) []*MemberCommitmentsMessage
		expectedError             error
		expectedAccusations       int
	}{
		"positive validation - no accusations": {
			expectedError:       nil,
			expectedAccusations: 0,
		},
		"negative validation - changed random share": {
			modifyPeerShareMessages: func(messages []*PeerSharesMessage) []*PeerSharesMessage {
				messages[0].shareT = big.NewInt(13)
				return messages
			},
			expectedError:       nil,
			expectedAccusations: 1,
		},
		"negative validation - changed commitment": {
			modifyCommitmentsMessages: func(messages []*MemberCommitmentsMessage) []*MemberCommitmentsMessage {
				messages[0].commitments[0] = big.NewInt(13)
				return messages
			},
			expectedError:       nil,
			expectedAccusations: 1,
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
				filteredSharesMessages = test.modifyPeerShareMessages(filteredSharesMessages)
			}
			if test.modifyCommitmentsMessages != nil {
				filteredCommitmentsMessages = test.modifyCommitmentsMessages(filteredCommitmentsMessages)
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

			if len(accusedMessage.accusedIDs) != test.expectedAccusations {
				t.Fatalf("\nexpected: %v accusations\nactual:   %v\n",
					test.expectedAccusations,
					len(accusedMessage.accusedIDs),
				)
			}

			expectedReceivedSharesLength := groupSize - 1 - test.expectedAccusations
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
		})
	}
}

func TestCombineReceivedShares(t *testing.T) {
	receivedShareS := make(map[int]*big.Int)
	receivedShareT := make(map[int]*big.Int)
	for i := 0; i <= 5; i++ {
		receivedShareS[100+i] = big.NewInt(int64(i))
		receivedShareT[100+i] = big.NewInt(int64(10 + i))
	}

	expectedShareS := big.NewInt(15)
	expectedShareT := big.NewInt(75)

	config, err := predefinedDKGconfig()
	if err != nil {
		t.Fatalf("DKG Config initialization failed [%s]", err)
	}
	member := &SharingMember{
		CommittingMember: &CommittingMember{
			memberCore: &memberCore{
				protocolConfig: config,
			},
			receivedSharesS: receivedShareS,
			receivedSharesT: receivedShareT,
		},
	}

	member.CombineReceivedShares()

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
		big.NewInt(216),
		big.NewInt(148),
		big.NewInt(36),
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
		expectedAccusations              int
	}{
		"positive validation - no accusations": {
			expectedError:       nil,
			expectedAccusations: 0,
		},
		"negative validation - changed public key share": {
			modifyPublicCoefficientsMessages: func(messages []*MemberPublicCoefficientsMessage) {
				messages[1].publicCoefficients[1] = big.NewInt(13)
			},
			expectedError:       nil,
			expectedAccusations: 1,
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

			if len(accusedMessage.accusedIDs) != test.expectedAccusations {
				t.Fatalf("\nexpected: %v accusations\nactual:   %v\n",
					test.expectedAccusations,
					len(accusedMessage.accusedIDs),
				)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	threshold := 3
	groupSize := 5

	committingMembers, err := initializeCommittingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}

	var sharesMessages []*PeerSharesMessage
	var commitmentsMessages []*MemberCommitmentsMessage
	for _, member := range committingMembers {
		sharesMessage, commitmentsMessage, err := member.CalculateMembersSharesAndCommitments()
		if err != nil {
			t.Fatalf("shares and commitments calculation failed [%s]", err)
		}
		sharesMessages = append(sharesMessages, sharesMessage...)
		commitmentsMessages = append(commitmentsMessages, commitmentsMessage)
	}

	for i := range committingMembers {
		committingMember := committingMembers[i]

		accusedSecretSharesMessage, err := committingMember.VerifyReceivedSharesAndCommitmentsMessages(
			filterPeerSharesMessage(sharesMessages, committingMember.ID),
			filterMemberCommitmentsMessages(commitmentsMessages, committingMember.ID),
		)
		if err != nil {
			t.Fatalf("shares and commitments verification failed [%s]", err)
		}

		if len(accusedSecretSharesMessage.accusedIDs) > 0 {
			t.Fatalf("\nexpected: 0 accusations\nactual:   %d\n",
				accusedSecretSharesMessage.accusedIDs,
			)
		}
	}

	var sharingMembers []*SharingMember
	for _, cm := range committingMembers {
		sharingMembers = append(sharingMembers, &SharingMember{
			CommittingMember: cm,
		})
	}

	sharingMember := sharingMembers[0]
	if len(sharingMember.receivedSharesS) != groupSize-1 {
		t.Fatalf("\nexpected: %d received shares\nactual:   %d\n",
			groupSize-1,
			len(sharingMember.receivedSharesS),
		)
	}

	for _, member := range sharingMembers {
		member.CombineReceivedShares()
	}

	publicCoefficientsMessages := make([]*MemberPublicCoefficientsMessage, groupSize)
	for i, member := range sharingMembers {
		publicCoefficientsMessages[i] = member.CalculatePublicCoefficients()
	}

	for i := range sharingMembers {
		member := sharingMembers[i]

		accusedCoefficientsMessage, err := member.VerifyPublicCoefficients(
			filterMemberPublicCoefficientsMessages(publicCoefficientsMessages, member.ID),
		)
		if err != nil {
			t.Fatalf("public coefficients verification failed [%s]", err)
		}
		if len(accusedCoefficientsMessage.accusedIDs) > 0 {
			t.Fatalf("\nexpected: 0 accusations\nactual:   %d\n",
				accusedCoefficientsMessage.accusedIDs,
			)
		}
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
	config, err := predefinedDKGconfig()
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
			vss:             vss,
			receivedSharesS: make(map[int]*big.Int),
			receivedSharesT: make(map[int]*big.Int),
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

// predefinedDKGconfig initializes DKG configuration with predefined values.
func predefinedDKGconfig() (*DKG, error) {
	// `p` is 4096-bit safe prime.
	pStr := "0xc8526644a9c4739683742b7003640b2023ca42cc018a42b02a551bb825c6828f86e2e216ea5d31004c433582a3fa720459efb42e091d73fb281810e1825691f0799811be62ae57f62ab00670edd35426d108d3b9c4fd008eddc67275a0489fe132e4c31bd7069ea7884cbb8f8f9255fe7b87fc0099f246776c340912df48f7945bc2bc0bc6814978d27b7af2ebc41f458ae795186db0fd7e6151bb8a7fe2b41370f7a2848ef75d3ec88f3439022c10e78b434c2f24b2f40bd02930e6c8aadef87b0dc87cdba07dcfa86884a168bd1381a4f48be12e5d98e41f954c37aec011cc683570e8890418756ed98ace8c8e59ae1df50962c1622fe66b5409f330cad6b7c68f2e884786d9807190b89ac4a3b3507e49b2dd3f33d765ad29e2015180c8cd0258dd8bdaab17be5d74871fec04c492240c6a2692b2c9a62c9adbaac34a333f135801ff948e8dfb6bbd6212a67950fb8edd628d05d19d1b94e9be7c52ed484831d50adaa29e71de197e351878f1c40ec67ee809e824124529e27bd5ecf3054f6784153f7db27ff0c87420bb2b2754ed363fc2ba8399d49d291f342173e7619183467a9694efa243e1d41b26c13b38ca0f43bb7c9050eb966461f28436583a9d13d2c1465b78184eae360f009505ccea288a053d111988d55c12befd882a857a530efac2c0592987cd83c39844a10e058739ab1c39006a3123e7fc887845675f"
	// `q` is 4095-bit Sophie Germain prime.
	qStr := "0x6429332254e239cb41ba15b801b2059011e5216600c52158152a8ddc12e34147c371710b752e988026219ac151fd39022cf7da17048eb9fd940c0870c12b48f83ccc08df31572bfb1558033876e9aa13688469dce27e80476ee3393ad0244ff09972618deb834f53c4265dc7c7c92aff3dc3fe004cf9233bb61a04896fa47bca2de15e05e340a4bc693dbd7975e20fa2c573ca8c36d87ebf30a8ddc53ff15a09b87bd142477bae9f64479a1c81160873c5a1a61792597a05e814987364556f7c3d86e43e6dd03ee7d4344250b45e89c0d27a45f0972ecc720fcaa61bd76008e6341ab87444820c3ab76cc56746472cd70efa84b160b117f335aa04f998656b5be347974423c36cc038c85c4d6251d9a83f24d96e9f99ebb2d694f100a8c06466812c6ec5ed558bdf2eba438ff602624912063513495964d3164d6dd561a5199f89ac00ffca4746fdb5deb109533ca87dc76eb14682e8ce8dca74df3e2976a42418ea856d514f38ef0cbf1a8c3c78e207633f7404f412092294f13deaf67982a7b3c20a9fbed93ff8643a105d9593aa769b1fe15d41ccea4e948f9a10b9f3b0c8c1a33d4b4a77d121f0ea0d93609d9c6507a1ddbe482875cb3230f9421b2c1d4e89e960a32dbc0c27571b07804a82e6751445029e888cc46aae095f7ec41542bd29877d61602c94c3e6c1e1cc22508702c39cd58e1c80351891f3fe443c22b3af"

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
