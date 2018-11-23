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

	members, err := initializeCommittingMembersGroup(threshold, groupSize, nil)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}

	member := members[0]
	err = member.CalculateMembersSharesAndCommitments()
	if err != nil {
		t.Fatalf("shares and commitments calculation failed [%s]", err)
	}

	if len(member.secretCoefficients) != (threshold + 1) {
		t.Fatalf("\nexpected: %v secret coefficients\nactual:   %v\n",
			threshold+1,
			len(member.secretCoefficients),
		)
	}
	if len(member.evaluatedSecretSharesS) != (groupSize) {
		t.Fatalf("\nexpected: %v evaluated shares S\nactual:   %v\n",
			groupSize,
			len(member.evaluatedSecretSharesS),
		)
	}

	if len(member.evaluatedSecretSharesT) != (groupSize) {
		t.Fatalf("\nexpected: %v evaluated shares T\nactual:   %v\n",
			groupSize,
			len(member.evaluatedSecretSharesT),
		)
	}

	if len(member.commitments) != (threshold + 1) {
		t.Fatalf("\nexpected: %v calculated commitments\nactual:   %v\n",
			threshold+1,
			len(member.commitments),
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

	var tests = map[string]struct {
		modifySharesAndCommitments func(s map[int]*SharesAndCommitments)
		expectedError              error
		expectedAccusedIDs         []int
	}{
		"positive validation - no accusations": {
			expectedError: nil,
		},
		"negative validation - changed share S": {
			modifySharesAndCommitments: func(s map[int]*SharesAndCommitments) {
				s[2].shareS = testutils.NewRandInt(s[2].shareS, config.Q)
			},
			expectedError:      nil,
			expectedAccusedIDs: []int{2},
		},
		"negative validation - changed two shares T": {
			modifySharesAndCommitments: func(s map[int]*SharesAndCommitments) {
				s[3].shareT = testutils.NewRandInt(s[3].shareT, config.Q)
				s[4].shareT = testutils.NewRandInt(s[4].shareT, config.Q)
			},
			expectedError:      nil,
			expectedAccusedIDs: []int{3, 4},
		},
		"negative validation - changed commitment": {
			modifySharesAndCommitments: func(s map[int]*SharesAndCommitments) {
				s[5].commitments[1] = testutils.NewRandInt(
					s[5].commitments[1], config.Q,
				)
			},
			expectedError:      nil,
			expectedAccusedIDs: []int{5},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			members, err := initializeCommittingMembersGroup(threshold, groupSize, nil)
			if err != nil {
				t.Fatalf("group initialization failed [%s]", err)
			}
			currentMember := members[0]

			receivedPeerSharesAndCommitments := make(map[int]*SharesAndCommitments)
			for _, peer := range members {
				err = peer.CalculateMembersSharesAndCommitments()
				if err != nil {
					t.Fatalf("shares and commitments calculation failed [%s]", err)
				}

				if peer.ID != currentMember.ID {
					receivedPeerSharesAndCommitments[peer.ID] = &SharesAndCommitments{
						peerID:      peer.ID,
						shareS:      peer.evaluatedSecretSharesS[currentMember.ID],
						shareT:      peer.evaluatedSecretSharesT[currentMember.ID],
						commitments: peer.commitments,
					}
				}
			}

			if test.modifySharesAndCommitments != nil {
				test.modifySharesAndCommitments(receivedPeerSharesAndCommitments)
			}

			for peerID, sac := range receivedPeerSharesAndCommitments {
				currentMember.VerifyReceivedSharesAndCommitments(
					peerID,
					sac.shareS,
					sac.shareT,
					sac.commitments,
				)

				if !reflect.DeepEqual(test.expectedError, err) {
					t.Fatalf(
						"\nexpected: %v\nactual:   %v\n",
						test.expectedError,
						err,
					)
				}
			}

			if len(currentMember.accusedMembersIDs) != len(test.expectedAccusedIDs) {
				t.Fatalf("\nexpected: %v accusations\nactual:   %v\n",
					len(test.expectedAccusedIDs),
					len(currentMember.accusedMembersIDs),
				)
			}
			if !reflect.DeepEqual(currentMember.accusedMembersIDs, test.expectedAccusedIDs) {
				t.Fatalf("incorrect accused members IDs\nexpected: %v\nactual:   %v\n",
					test.expectedAccusedIDs,
					currentMember.accusedMembersIDs,
				)
			}

			expectedReceivedSharesLength := groupSize - 1 - len(test.expectedAccusedIDs)
			if len(currentMember.receivedValidSharesS) != expectedReceivedSharesLength {
				t.Fatalf("\nexpected: %v received shares S\nactual:   %v\n",
					expectedReceivedSharesLength,
					len(currentMember.receivedValidSharesS),
				)
			}
			if len(currentMember.receivedValidSharesT) != expectedReceivedSharesLength {
				t.Fatalf("\nexpected: %v received shares T\nactual:   %v\n",
					expectedReceivedSharesLength,
					len(currentMember.receivedValidSharesT),
				)
			}
			if len(currentMember.receivedValidPeerCommitments) != expectedReceivedSharesLength {
				t.Fatalf("\nexpected: %v received commitments\nactual:   %v\n",
					expectedReceivedSharesLength,
					len(currentMember.receivedValidPeerCommitments),
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
			return nil, fmt.Errorf("DKG Config initialization failed [%s]", err)
		}
	}

	vss, err := pedersen.NewVSS(crand.Reader, dkg.P, dkg.Q)
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
			LocalMember: &LocalMember{
				memberCore: &memberCore{
					ID:             id,
					group:          group,
					protocolConfig: dkg,
				},
			},
			vss:                          vss,
			receivedValidSharesS:         make(map[int]*big.Int),
			receivedValidSharesT:         make(map[int]*big.Int),
			receivedValidPeerCommitments: make(map[int][]*big.Int),
		})
		group.RegisterMemberID(id)
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
