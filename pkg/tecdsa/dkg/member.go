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
	tssPartyID, groupTssPartiesIDs := generateTssPartiesIDs(
		skgm.id,
		skgm.group.MemberIDs(),
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
		Group:           fm.group,
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(fm.tssResult),
	}
}

// signingMember represents a group member sharing their preferred DKG result hash
// and signature (over this hash) with other peer members.
type signingMember struct {
	logger      log.StandardLogger
	memberIndex group.MemberIndex
	// Group to which this member belongs.
	group *group.Group
	// Validator allowing to check public key and member index
	// against group members
	membershipValidator *group.MembershipValidator
	// Identifier of the particular DKG session this member is part of.
	sessionID string
	// Hash of DKG result preferred by the current participant.
	preferredDKGResultHash ResultHash
	// Signature over preferredDKGResultHash calculated by the member.
	selfDKGResultSignature []byte
}

// newSigningMember creates a new signingMember in the initial state.
func newSigningMember(
	logger log.StandardLogger,
	memberIndex group.MemberIndex,
	group *group.Group,
	membershipValidator *group.MembershipValidator,
	sessionID string,
) *signingMember {
	return &signingMember{
		logger:              logger,
		memberIndex:         memberIndex,
		group:               group,
		membershipValidator: membershipValidator,
		sessionID:           sessionID,
	}
}

// shouldAcceptMessage indicates whether the given member should accept
// a message from the given sender.
func (sm *signingMember) shouldAcceptMessage(
	senderID group.MemberIndex,
	senderPublicKey []byte,
) bool {
	isMessageFromSelf := senderID == sm.memberIndex
	isSenderValid := sm.membershipValidator.IsValidMembership(
		senderID,
		senderPublicKey,
	)
	isSenderAccepted := sm.group.IsOperating(senderID)

	return !isMessageFromSelf && isSenderValid && isSenderAccepted
}

// initializeSubmittingMember performs a transition of a member state to the
// next phase of the protocol.
func (sm *signingMember) initializeSubmittingMember() *submittingMember {
	return &submittingMember{
		signingMember: sm,
	}
}

// submittingMember represents a member submitting a DKG result to the
// blockchain along with signatures received from other group members supporting
// the result.
type submittingMember struct {
	*signingMember
}

// SubmitDKGResult sends a result, which contains the group public key and
// signatures, to the chain.
// It checks if the result has already been published to the blockchain by
// checking if a group with the given public key is already registered. If not,
// it determines if the current member is eligible to submit a result.
// If allowed, it submits the result to the chain.
// A user's turn to publish is determined based on the user's index and block
// step.
// If a result is submitted by another member and it's accepted by the chain,
// the current member finishes the phase immediately, without submitting
// their own result.
// It returns the on-chain block height of the moment when the result was
// successfully submitted on chain by the member. In case of failure or result
// already submitted by another member it returns `0`.
func (sm *submittingMember) SubmitDKGResult(
	result *Result,
	signatures map[group.MemberIndex][]byte,
	startBlockNumber uint64,
	resultSubmitter ResultSubmitter,
) error {
	return resultSubmitter.SubmitResult(
		result,
		signatures,
		startBlockNumber,
		sm.memberIndex,
	)
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
