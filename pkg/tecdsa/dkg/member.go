package dkg

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/keep-network/keep-core/pkg/tecdsa"

	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/tss"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// Member represents a DKG protocol member.
type member struct {
	// Logger used to produce log messages.
	logger log.StandardLogger
	// id of this group member.
	id group.MemberIndex
	// Group to which this member belongs.
	group *group.Group
	// Validator allowing to check public key and member index against
	// group members
	membershipValidator *group.MembershipValidator
	// Identifier of the particular DKG session this member is part of.
	sessionID string
	// TSS pre-parameters.
	tssPreParams *keygen.LocalPreParams
}

// newMember creates a new member in an initial state
func newMember(
	logger log.StandardLogger,
	memberID group.MemberIndex,
	groupSize,
	dishonestThreshold int,
	membershipValidator *group.MembershipValidator,
	sessionID string,
	tssPreParams *keygen.LocalPreParams,
) *member {
	return &member{
		logger:              logger,
		id:                  memberID,
		group:               group.NewGroup(dishonestThreshold, groupSize),
		membershipValidator: membershipValidator,
		sessionID:           sessionID,
		tssPreParams:        tssPreParams,
	}
}

// inactiveMemberFilter returns a new instance of the inactive member filter.
func (m *member) inactiveMemberFilter() *group.InactiveMemberFilter {
	return group.NewInactiveMemberFilter(m.logger, m.id, m.group)
}

// shouldAcceptMessage indicates whether the given member should accept
// a message from the given sender.
func (m *member) shouldAcceptMessage(
	senderID group.MemberIndex,
	senderPublicKey []byte,
) bool {
	isMessageFromSelf := senderID == m.id
	isSenderValid := m.membershipValidator.IsValidMembership(
		senderID,
		senderPublicKey,
	)
	isSenderAccepted := m.group.IsOperating(senderID)

	return !isMessageFromSelf && isSenderValid && isSenderAccepted
}

// initializeEphemeralKeysGeneration performs a transition of a member state
// from the initial state to the first phase of the protocol.
func (m *member) initializeEphemeralKeysGeneration() *ephemeralKeyPairGeneratingMember {
	return &ephemeralKeyPairGeneratingMember{
		member:            m,
		ephemeralKeyPairs: make(map[group.MemberIndex]*ephemeral.KeyPair),
	}
}

// ephemeralKeyPairGeneratingMember represents one member in a distributed key
// generating group performing ephemeral key pair generation. It has a full list
// of `memberIDs` that belong to its threshold group.
type ephemeralKeyPairGeneratingMember struct {
	*member

	// Ephemeral key pairs used to create symmetric keys,
	// generated individually for each other group member.
	ephemeralKeyPairs map[group.MemberIndex]*ephemeral.KeyPair
}

// initializeSymmetricKeyGeneration performs a transition of the member state
// to the next phase. It returns a member instance ready to execute the
// next phase of the protocol.
func (ekpgm *ephemeralKeyPairGeneratingMember) initializeSymmetricKeyGeneration() *symmetricKeyGeneratingMember {
	return &symmetricKeyGeneratingMember{
		ephemeralKeyPairGeneratingMember: ekpgm,
		symmetricKeys:                    make(map[group.MemberIndex]ephemeral.SymmetricKey),
	}
}

// symmetricKeyGeneratingMember represents one member in a distributed key
// generating group performing ephemeral symmetric key generation.
type symmetricKeyGeneratingMember struct {
	*ephemeralKeyPairGeneratingMember

	// Symmetric keys used to encrypt confidential information,
	// generated individually for each other group member by ECDH'ing the
	// broadcasted ephemeral public key intended for this member and the
	// ephemeral private key generated for the other member.
	symmetricKeys map[group.MemberIndex]ephemeral.SymmetricKey
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (skgm *symmetricKeyGeneratingMember) MarkInactiveMembers(
	ephemeralPubKeyMessages []*ephemeralPublicKeyMessage,
) {
	filter := skgm.inactiveMemberFilter()
	for _, message := range ephemeralPubKeyMessages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

// initializeTssRoundOne returns a member to perform next protocol operations.
func (skgm *symmetricKeyGeneratingMember) initializeTssRoundOne() *tssRoundOneMember {
	// Set up the local TSS party using only operating members. This effectively
	// removes all excluded members who were marked as disqualified at the
	// beginning of the protocol.
	tssPartyID, groupTssPartiesIDs := generateTssPartiesIDs(
		skgm.id,
		skgm.group.OperatingMemberIDs(),
	)

	tssParameters := tss.NewParameters(
		tecdsa.Curve,
		tss.NewPeerContext(tss.SortPartyIDs(groupTssPartiesIDs)),
		tssPartyID,
		len(groupTssPartiesIDs),
		skgm.group.HonestThreshold()-1,
	)

	tssOutgoingMessagesChan := make(chan tss.Message, len(groupTssPartiesIDs))
	tssResultChan := make(chan keygen.LocalPartySaveData, 1)

	tssParty := keygen.NewLocalParty(
		tssParameters,
		tssOutgoingMessagesChan,
		tssResultChan,
		*skgm.tssPreParams,
	)

	return &tssRoundOneMember{
		symmetricKeyGeneratingMember: skgm,
		tssParty:                     tssParty,
		tssParameters:                tssParameters,
		tssOutgoingMessagesChan:      tssOutgoingMessagesChan,
		tssResultChan:                tssResultChan,
	}
}

// tssRoundOneMember represents one member in a distributed key generating
// group performing the first round of the TSS keygen.
type tssRoundOneMember struct {
	*symmetricKeyGeneratingMember

	tssParty                tss.Party
	tssParameters           *tss.Parameters
	tssOutgoingMessagesChan <-chan tss.Message
	tssResultChan           <-chan keygen.LocalPartySaveData
}

// initializeTssRoundTwo returns a member to perform next protocol operations.
func (trom *tssRoundOneMember) initializeTssRoundTwo() *tssRoundTwoMember {
	return &tssRoundTwoMember{
		tssRoundOneMember: trom,
	}
}

// tssRoundTwoMember represents one member in a distributed key generating
// group performing the second round of the TSS keygen.
type tssRoundTwoMember struct {
	*tssRoundOneMember
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as inactive.
func (trtm *tssRoundTwoMember) MarkInactiveMembers(
	tssRoundOneMessages []*tssRoundOneMessage,
) {
	filter := trtm.inactiveMemberFilter()
	for _, message := range tssRoundOneMessages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

// initializeTssRoundThree returns a member to perform next protocol operations.
func (trtm *tssRoundTwoMember) initializeTssRoundThree() *tssRoundThreeMember {
	return &tssRoundThreeMember{
		tssRoundTwoMember: trtm,
	}
}

// tssRoundThreeMember represents one member in a distributed key generating
// group performing the third round of the TSS keygen.
type tssRoundThreeMember struct {
	*tssRoundTwoMember
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (trtm *tssRoundThreeMember) MarkInactiveMembers(
	tssRoundTwoMessages []*tssRoundTwoMessage,
) {
	filter := trtm.inactiveMemberFilter()
	for _, message := range tssRoundTwoMessages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

// initializeFinalization returns a member to perform next protocol operations.
func (trtm *tssRoundThreeMember) initializeFinalization() *finalizingMember {
	return &finalizingMember{
		tssRoundThreeMember: trtm,
	}
}

// finalizingMember represents one member of the given group, after it
// completed the distributed key generation process.
//
// Prepares a result to publish in the last phase of the protocol.
type finalizingMember struct {
	*tssRoundThreeMember

	tssResult keygen.LocalPartySaveData
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (fm *finalizingMember) MarkInactiveMembers(
	tssRoundThreeMessages []*tssRoundThreeMessage,
) {
	filter := fm.inactiveMemberFilter()
	for _, message := range tssRoundThreeMessages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

// Result can be either the successful computation of the distributed key
// generation process or a notification of failure.
func (fm *finalizingMember) Result() *Result {
	return &Result{
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(fm.tssResult),
	}
}

// generateTssPartiesIDs converts group member ID to parties ID suitable for
// the TSS protocol execution.
func generateTssPartiesIDs(
	memberID group.MemberIndex,
	groupMembersIDs []group.MemberIndex,
) (*tss.PartyID, []*tss.PartyID) {
	var partyID *tss.PartyID
	groupPartiesIDs := make([]*tss.PartyID, len(groupMembersIDs))

	for i, groupMemberID := range groupMembersIDs {
		newPartyID := newTssPartyIDFromMemberID(groupMemberID)

		if memberID == groupMemberID {
			partyID = newPartyID
		}

		groupPartiesIDs[i] = newPartyID
	}

	return partyID, groupPartiesIDs
}

// newTssPartyIDFromMemberID creates a new instance of a TSS party ID using
// the given member ID. Such a created party ID has an unset index since it
// does not yet belong to a sorted parties IDs set.
func newTssPartyIDFromMemberID(memberID group.MemberIndex) *tss.PartyID {
	return tss.NewPartyID(
		strconv.Itoa(int(memberID)),
		fmt.Sprintf("member-%v", memberID),
		memberIDToTssPartyIDKey(memberID),
	)
}

// memberIDToTssPartyIDKey converts a single group member ID to a key that
// can be used to create a TSS party ID.
func memberIDToTssPartyIDKey(memberID group.MemberIndex) *big.Int {
	return big.NewInt(int64(memberID))
}

// tssPartyIDToMemberID converts a single TSS party ID to a group member ID.
func tssPartyIDToMemberID(partyID *tss.PartyID) group.MemberIndex {
	return group.MemberIndex(partyID.KeyInt().Int64())
}

// resolveSortedTssPartyID resolves the TSS party ID for the given member ID
// based on the sorted parties IDs stored in the given TSS parameters set. Such
// a resolved party ID has an index which indicates its position in the parties
// IDs set.
func resolveSortedTssPartyID(
	tssParameters *tss.Parameters,
	memberID group.MemberIndex,
) *tss.PartyID {
	sortedPartiesIDs := tssParameters.Parties().IDs()
	partyIDKey := memberIDToTssPartyIDKey(memberID)
	return sortedPartiesIDs.FindByKey(partyIDKey)
}
