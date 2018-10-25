package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/pedersen"
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
		modifyPeerShareMessages   func(messages []*PeerSharesMessage) []int
		modifyCommitmentsMessages func(messages []*MemberCommitmentsMessage) []int
		expectedError             error
		expectedAccusations       int
	}{
		"positive validation - no accusations": {
			expectedError:       nil,
			expectedAccusations: 0,
		},
		"negative validation - changed share S": {
			modifyPeerShareMessages: func(messages []*PeerSharesMessage) []int {
				messages[0].shareS = big.NewInt(13)
				return []int{messages[0].senderID}
			},
			expectedError:       nil,
			expectedAccusations: 1,
		},
		"negative validation - changed two shares T": {
			modifyPeerShareMessages: func(messages []*PeerSharesMessage) []int {
				messages[1].shareT = big.NewInt(13)
				messages[2].shareT = big.NewInt(23)
				return []int{messages[1].senderID, messages[2].senderID}
			},
			expectedError:       nil,
			expectedAccusations: 2,
		},
		"negative validation - changed commitment": {
			modifyCommitmentsMessages: func(messages []*MemberCommitmentsMessage) []int {
				messages[3].commitments[1] = big.NewInt(33)
				return []int{messages[3].senderID}
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
			var expectedAccusedIDs []int

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
				expectedAccusedIDs = append(
					expectedAccusedIDs,
					test.modifyPeerShareMessages(filteredSharesMessages)...,
				)
			}
			if test.modifyCommitmentsMessages != nil {
				expectedAccusedIDs = append(
					expectedAccusedIDs,
					test.modifyCommitmentsMessages(filteredCommitmentsMessages)...,
				)
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
			if !reflect.DeepEqual(accusedMessage.accusedIDs, expectedAccusedIDs) {
				t.Fatalf("incorrect accused members IDs\nexpected: %v\nactual:   %v\n",
					expectedAccusedIDs,
					accusedMessage.accusedIDs,
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
			if len(currentMember.receivedCommitments) != expectedReceivedSharesLength {
				t.Fatalf("\nexpected: %v received commitments\nactual:   %v\n",
					expectedReceivedSharesLength,
					len(currentMember.receivedCommitments),
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

	vss, err := pedersen.NewVSS(config.P, config.Q)
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

// predefinedDKGconfig initializez DKG configuration with predefined 2048-bit
// p and q values.
func predefinedDKG() (*DKG, error) {
	// `p` is 2048-bit safe prime.
	pStr := "0x93cef9a05e49e4701ab80ec2be6fa7b77524520f4bdad03b8b1a4424c0329588ace3f597cf1e99d8c54486cf2970bd9833b1d83a80ae3315459f9d6ca55dd4ab73e6e84d98d6e0b8f06a409374c646c79aaad075ea4685c6d91b1b2a034044dcfed7b7d5d628e939a63fa03185a71570819c830cb2f8d8d5a8a5b757f4966c362317e96a181d213afff464783bc31b196b5971d8988a98e1c81db6e7ad06c151ca6e4801fe566ae212a8bdbf56c971bc9bb8e64b61ec5bb36a2eb6d5842e4b95e6175d862fbfd8b71ae9912c0a94df6c77c5feeb1c82fb05976d07cad53f012f6910d55d8617ecf166c0856da0932c7d0e6ca858367642295113a1d72ca2408b"
	// `q` is 2047-bit Sophie Germain prime.
	qStr := "0x49e77cd02f24f2380d5c07615f37d3dbba922907a5ed681dc58d221260194ac45671facbe78f4cec62a2436794b85ecc19d8ec1d4057198aa2cfceb652aeea55b9f37426cc6b705c78352049ba632363cd55683af52342e36c8d8d9501a0226e7f6bdbeaeb14749cd31fd018c2d38ab840ce4186597c6c6ad452dbabfa4b361b118bf4b50c0e909d7ffa323c1de18d8cb5acb8ec4c454c70e40edb73d68360a8e5372400ff2b357109545edfab64b8de4ddc7325b0f62dd9b5175b6ac21725caf30baec317dfec5b8d74c896054a6fb63be2ff758e417d82cbb683e56a9f8097b4886aaec30bf678b36042b6d049963e8736542c1b3b2114a889d0eb96512045"

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
