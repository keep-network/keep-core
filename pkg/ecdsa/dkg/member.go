package dkg

import (
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"math/big"
	"strconv"
)

// Member represents a DKG protocol member.
type member struct {
	// id of this group member.
	id group.MemberIndex
	// Group to which this member belongs.
	group *group.Group
	// Validator allowing to check public key and member index against
	// group members
	membershipValidator group.MembershipValidator
	// Identifier of the particular DKG session this member is part of.
	sessionID string
}

// newMember creates a new member in an initial state
func newMember(
	memberID group.MemberIndex,
	groupSize,
	dishonestThreshold int,
	membershipValidator group.MembershipValidator,
	sessionID string,
) *member {
	return &member{
		id:                  memberID,
		group:               group.NewGroup(dishonestThreshold, groupSize),
		membershipValidator: membershipValidator,
		sessionID:           sessionID,
	}
}

// messageFilter returns a new instance of the message filter.
func (m *member) messageFilter() *group.InactiveMemberFilter {
	return group.NewInactiveMemberFilter(m.id, m.group)
}

// IsSenderAccepted returns true if sender with the given index is accepted
// as an operating group member.
func (m *member) IsSenderAccepted(senderID group.MemberIndex) bool {
	return m.group.IsOperating(senderID)
}

// IsSenderValid returns true if sender with the given index is considered
// a valid member of the given group.
func (m *member) IsSenderValid(
	senderID group.MemberIndex,
	senderPublicKey []byte,
) bool {
	return m.membershipValidator.IsValidMembership(senderID, senderPublicKey)
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
	filter := skgm.messageFilter()
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
		tss.EC(),
		tss.NewPeerContext(tss.SortPartyIDs(groupTssPartiesIDs)),
		tssPartyID,
		len(groupTssPartiesIDs),
		skgm.group.DishonestThreshold(),
	)

	tssOutgoingMessageChan := make(chan tss.Message, len(groupTssPartiesIDs))
	tssResultChan := make(chan keygen.LocalPartySaveData)

	// TODO: Use pre-computed pre-params.
	tssParty := keygen.NewLocalParty(
		tssParameters,
		tssOutgoingMessageChan,
		tssResultChan,
	)

	return &tssRoundOneMember{
		symmetricKeyGeneratingMember: skgm,
		tssParty:                     tssParty,
	}
}

// tssRoundOneMember represents one member in a distributed key generating
// group performing the first round of the TSS keygen.
type tssRoundOneMember struct {
	*symmetricKeyGeneratingMember

	tssParty               tss.Party
	tssOutgoingMessageChan <-chan tss.Message
	tssResultChan          <-chan keygen.LocalPartySaveData
}

// finalizingMember represents one member of the given group, after it
// completed the distributed key generation process.
//
// Prepares a result to publish in the last phase of the protocol.
type finalizingMember struct {
	*symmetricKeyGeneratingMember
}

// Result can be either the successful computation of the distributed key
// generation process or a notification of failure.
func (fm *finalizingMember) Result() *Result {
	return &Result{
		// TODO: Temporary result. Add real items.
		SymmetricKeys: fm.symmetricKeys,
	}
}

func generateTssPartiesIDs(
	memberID group.MemberIndex,
	groupMembersIDs []group.MemberIndex,
) (*tss.PartyID, []*tss.PartyID) {
	var partyID *tss.PartyID
	var groupPartiesIDs []*tss.PartyID

	for _, groupMemberID := range groupMembersIDs {
		newPartyID := tss.NewPartyID(
			strconv.Itoa(int(groupMemberID)),
			"",
			big.NewInt(int64(groupMemberID)),
		)

		if memberID == groupMemberID {
			partyID = newPartyID
		}

		groupPartiesIDs = append(groupPartiesIDs, newPartyID)
	}

	return partyID, groupPartiesIDs
}
