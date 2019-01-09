package gjkr

import (
	"math/big"
	"strconv"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

// MemberID is a unique-in-group identifier of a member.
type MemberID uint32

type memberCore struct {
	// ID of this group member.
	ID MemberID

	// Group to which this member belongs.
	group *Group

	// Evidence log provides access to messages from earlier protocol phases
	// for the sake of compliant resolution.
	evidenceLog evidenceLog

	// Cryptographic protocol parameters, the same for all members in the group.
	protocolParameters *protocolParameters
}

// LocalMember represents one member in a threshold group, prior to the
// initiation of distributed key generation process
type LocalMember struct {
	*memberCore
}

// EphemeralKeyPairGeneratingMember represents one member in a distributed key
// generating group performing ephemeral key pair generation. It has a full list
// of `memberIDs` that belong to its threshold group.
//
// Executes Phase 1 of the protocol.
type EphemeralKeyPairGeneratingMember struct {
	*LocalMember

	// Ephemeral key pairs used to create symmetric keys,
	// generated individually for each other group member.
	ephemeralKeyPairs map[MemberID]*ephemeral.KeyPair
}

// SymmetricKeyGeneratingMember represents one member in a distributed key
// generating group performing ephemeral symmetric key generation.
//
// Executes Phase 2 of the protocol.
type SymmetricKeyGeneratingMember struct {
	*EphemeralKeyPairGeneratingMember

	// Symmetric keys used to encrypt confidential information,
	// generated individually for each other group member by ECDH'ing the
	// broadcasted ephemeral public key intended for this member and the
	// ephemeral private key generated for the other member.
	symmetricKeys map[MemberID]ephemeral.SymmetricKey
}

// CommittingMember represents one member in a distributed key generation group,
// after it has fully initialized ephemeral symmetric keys with all other group
// members.
//
// Executes Phase 3 of the protocol.
type CommittingMember struct {
	*SymmetricKeyGeneratingMember

	// Polynomial `a` coefficients generated by the member. Polynomial is of
	// degree `dishonestThreshold`, so the number of coefficients equals
	// `dishonestThreshold + 1`
	//
	// This is a private value and should not be exposed.
	secretCoefficients []*big.Int
	// Shares calculated by the current member for themself. They are defined as
	// `s_ii` and `t_ii` respectively across the protocol specification.
	//
	// These are private values and should not be exposed.
	selfSecretShareS, selfSecretShareT *big.Int
}

// CommitmentsVerifyingMember represents one member in a distributed key generation
// group, after it has received secret shares and commitments from other group
// members and it performs verification of received values.
//
// Executes Phase 4 of the protocol.
type CommitmentsVerifyingMember struct {
	*CommittingMember

	// Shares calculated for the current member by peer group members which passed
	// the validation.
	//
	// receivedValidSharesS are defined as `s_ji` and receivedValidSharesT are
	// defined as `t_ji` across the protocol specification.
	receivedValidSharesS, receivedValidSharesT map[MemberID]*big.Int
	// Valid commitments to secret shares polynomial coefficients received from
	// other group members.
	receivedValidPeerCommitments map[MemberID][]*bn256.G1
}

// SharesJustifyingMember represents one member in a threshold key sharing group,
// after it completed secret shares and commitments verification and enters
// justification phase where it resolves invalid share accusations.
//
// Executes Phase 5 of the protocol.
type SharesJustifyingMember struct {
	*CommitmentsVerifyingMember
}

// QualifiedMember represents one member in a threshold key sharing group, after
// it completed secret shares justification. The member holds a share of group
// group private key.
//
// Executes Phase 6 of the protocol.
type QualifiedMember struct {
	*SharesJustifyingMember

	// Member's share of the secret group private key. It is denoted as `z_ik`
	// in protocol specification.
	groupPrivateKeyShare *big.Int
}

// SharingMember represents one member in a threshold key sharing group, after it
// has been qualified to the group private key sharing. A member shares
// public values of it's polynomial coefficients with peer members.
//
// Executes Phase 7 and Phase 8 of the protocol.
type SharingMember struct {
	*QualifiedMember

	// Public values of each polynomial `a` coefficient defined in secretCoefficients
	// field. It is denoted as `A_ik` in protocol specification. The zeroth
	// public key share point `A_i0` is a member's public key share.
	publicKeySharePoints []*bn256.G1
	// Public key share points received from other group members which passed
	// the validation. Defined as `A_jk` across the protocol documentation.
	receivedValidPeerPublicKeySharePoints map[MemberID][]*bn256.G1
}

// PointsJustifyingMember represents one member in a threshold key sharing group,
// after it completed public key share points verification and enters justification
// phase where it resolves public key share points accusations.
//
// Executes Phase 9 of the protocol.
type PointsJustifyingMember struct {
	*SharingMember
}

// RevealingMember represents one member in a threshold sharing group who is
// revealing ephemeral private keys used to create ephemeral symmetric key
// to communicate with other members disqualified in Phase 9.
//
// Executes Phase 10 of the protocol.
type RevealingMember struct {
	*PointsJustifyingMember
}

// ReconstructingMember represents one member in a threshold sharing group who
// is reconstructing individual private and public keys of disqualified group members.
//
// Executes Phase 11 of the protocol.
type ReconstructingMember struct {
	*RevealingMember

	// Disqualified members' individual private keys reconstructed from shares
	// revealed by other group members.
	// Stored as `<m, z_m>`, where:
	// - `m` is disqualified member's ID
	// - `z_m` is reconstructed individual private key of member `m`
	reconstructedIndividualPrivateKeys map[MemberID]*big.Int
	// Individual public keys calculated from reconstructed individual private keys.
	// Stored as `<m, y_m>`, where:
	// - `m` is disqualified member's ID
	// - `y_m` is reconstructed individual public key of member `m`
	reconstructedIndividualPublicKeys map[MemberID]*bn256.G1
}

// CombiningMember represents one member in a threshold sharing group who is
// combining individual public keys of group members to receive group public key.
//
// Executes Phase 12 of the protocol.
type CombiningMember struct {
	*ReconstructingMember

	// Group public key calculated from individual public keys of all group members.
	// Denoted as `Y` across the protocol specification.
	groupPublicKey *bn256.G1
}

// InitializeFinalization returns a member to perform next protocol operations.
func (cm *CombiningMember) InitializeFinalization() *FinalizingMember {
	return &FinalizingMember{CombiningMember: cm}
}

// FinalizingMember represents one member in a threshold key sharing group,
// after it completed distributed key generation.
//
// Prepares a result to publish in Phase 13 of the protocol.
type FinalizingMember struct {
	*CombiningMember
}

// NewMember creates a new member in an initial state, ready to execute DKG
// protocol.
func NewMember(
	memberID MemberID,
	groupMembers []MemberID,
	dishonestThreshold int,
	seed *big.Int,
) *LocalMember {
	return &LocalMember{
		memberCore: &memberCore{
			memberID,
			&Group{
				dishonestThreshold,
				groupMembers,
				[]MemberID{},
				[]MemberID{},
			},
			newDkgEvidenceLog(),
			newProtocolParameters(seed),
		},
	}
}

// PublishingIndex returns sequence number of the current member in a publishing
// group. Counting starts with `0`.
func (fm *FinalizingMember) PublishingIndex() int {
	for index, memberID := range fm.group.MemberIDs() {
		if fm.ID == memberID {
			return index
		}
	}
	return -1 // should never happen
}

// Result can be either the successful computation of a round of distributed key
// generation, or a notification of failure.
//
// If the number of disqualified and inactive members is greater than half of the
// configured dishonest threshold, the group is deemed too weak, and the result
// is set to failure. Otherwise, it returns the generated group public key along
// with the disqualified and inactive members.
func (fm *FinalizingMember) Result() *Result {
	return &Result{
		Success:        fm.group.isThresholdSatisfied(),
		GroupPublicKey: fm.groupPublicKey,              // nil if threshold not satisfied
		Disqualified:   fm.group.disqualifiedMemberIDs, // DQ
		Inactive:       fm.group.inactiveMemberIDs,     // IA
	}
}

// Int converts `MemberID` to `big.Int`.
func (id MemberID) Int() *big.Int {
	return new(big.Int).SetUint64(uint64(id))
}

// AddToGroup adds the provided MemberID to the group
func (mc *memberCore) AddToGroup(memberID MemberID) {
	mc.group.RegisterMemberID(memberID)
}

// InitializeEphemeralKeysGeneration performs a transition of a member state
// from the local state to phase 1 of the protocol.
func (lm *LocalMember) InitializeEphemeralKeysGeneration() *EphemeralKeyPairGeneratingMember {
	return &EphemeralKeyPairGeneratingMember{
		LocalMember:       lm,
		ephemeralKeyPairs: make(map[MemberID]*ephemeral.KeyPair),
	}
}

// InitializeSymmetricKeyGeneration performs a transition of the member state
// from phase 1 to phase 2. It returns a member instance ready to execute the
// next phase of the protocol.
func (ekgm *EphemeralKeyPairGeneratingMember) InitializeSymmetricKeyGeneration() *SymmetricKeyGeneratingMember {
	return &SymmetricKeyGeneratingMember{
		EphemeralKeyPairGeneratingMember: ekgm,
		symmetricKeys:                    make(map[MemberID]ephemeral.SymmetricKey),
	}
}

// InitializeCommitting returns a member to perform next protocol operations.
func (skgm *SymmetricKeyGeneratingMember) InitializeCommitting() *CommittingMember {
	return &CommittingMember{
		SymmetricKeyGeneratingMember: skgm,
	}
}

// InitializeCommitmentsVerification returns a member to perform next protocol operations.
func (cm *CommittingMember) InitializeCommitmentsVerification() *CommitmentsVerifyingMember {
	return &CommitmentsVerifyingMember{
		CommittingMember:             cm,
		receivedValidSharesS:         make(map[MemberID]*big.Int),
		receivedValidSharesT:         make(map[MemberID]*big.Int),
		receivedValidPeerCommitments: make(map[MemberID][]*bn256.G1),
	}
}

// InitializeSharesJustification returns a member to perform next protocol operations.
func (cvm *CommitmentsVerifyingMember) InitializeSharesJustification() *SharesJustifyingMember {
	return &SharesJustifyingMember{cvm}
}

// InitializeQualified returns a member to perform next protocol operations.
func (sjm *SharesJustifyingMember) InitializeQualified() *QualifiedMember {
	return &QualifiedMember{SharesJustifyingMember: sjm}
}

// InitializeSharing returns a member to perform next protocol operations.
func (qm *QualifiedMember) InitializeSharing() *SharingMember {
	return &SharingMember{
		QualifiedMember:                       qm,
		receivedValidPeerPublicKeySharePoints: make(map[MemberID][]*bn256.G1),
	}
}

// InitializePointsJustification returns a member to perform next protocol operations.
func (sm *SharingMember) InitializePointsJustification() *PointsJustifyingMember {
	return &PointsJustifyingMember{sm}
}

// InitializeRevealing returns a member to perform next protocol operations.
func (sm *PointsJustifyingMember) InitializeRevealing() *RevealingMember {
	return &RevealingMember{sm}
}

// InitializeReconstruction returns a member to perform next protocol operations.
func (rm *RevealingMember) InitializeReconstruction() *ReconstructingMember {
	return &ReconstructingMember{
		RevealingMember:                    rm,
		reconstructedIndividualPrivateKeys: make(map[MemberID]*big.Int),
		reconstructedIndividualPublicKeys:  make(map[MemberID]*bn256.G1),
	}
}

// InitializeCombining returns a member to perform next protocol operations.
func (rm *ReconstructingMember) InitializeCombining() *CombiningMember {
	return &CombiningMember{ReconstructingMember: rm}
}

// individualPrivateKey returns current member's individual private key.
// Individual private key is zeroth polynomial coefficient `a_i0`.
func (rm *ReconstructingMember) individualPrivateKey() *big.Int {
	return rm.secretCoefficients[0]
}

// individualPublicKey returns current member's individual public key.
// Individual public key is zeroth public key share point `A_i0`.
func (rm *ReconstructingMember) individualPublicKey() *bn256.G1 {
	return rm.publicKeySharePoints[0]
}

// receivedValidPeerIndividualPublicKeys returns individual public keys received
// from other members which passed the validation. Individual public key is zeroth
// public key share point `A_j0`.
func (sm *SharingMember) receivedValidPeerIndividualPublicKeys() []*bn256.G1 {
	var receivedValidPeerIndividualPublicKeys []*bn256.G1

	for _, peerPublicKeySharePoints := range sm.receivedValidPeerPublicKeySharePoints {
		receivedValidPeerIndividualPublicKeys = append(
			receivedValidPeerIndividualPublicKeys,
			peerPublicKeySharePoints[0],
		)
	}
	return receivedValidPeerIndividualPublicKeys
}

// HexString converts `MemberID` to hex `string` representation.
func (id MemberID) HexString() string {
	return strconv.FormatInt(int64(id), 16)
}

// MemberIDFromHex returns a `MemberID` created from the hex `string`
// representation.
func MemberIDFromHex(hex string) (MemberID, error) {
	id, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0, err
	}

	return MemberID(id), nil
}
