package gjkr

import (
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay/pedersen"
	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
)

type memberCore struct {
	// ID of this group member.
	ID int
	// Group to which this member belongs.
	group *Group
	// DKG Protocol configuration parameters.
	protocolConfig *DKG
}

// CommittingMember represents one member in a threshold key sharing group, after
// it has a full list of `memberIDs` that belong to its threshold group. A
// member in this state has two maps of member shares for each member of the
// group.
//
// Executes Phase 3 and Phase 4 of the protocol.
type CommittingMember struct {
	*memberCore

	// Pedersen VSS scheme used to calculate commitments.
	vss *pedersen.VSS
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
	// Shares calculated for the current member by peer group members which passed
	// the validation.
	//
	// receivedValidSharesS are defined as `s_ji` and receivedValidSharesT are
	// defined as `t_ji` across the protocol specification.
	receivedValidSharesS, receivedValidSharesT map[int]*big.Int
	// Commitments to coefficients received from peer group members which passed
	// the validation.
	receivedValidPeerCommitments map[int][]*big.Int
}

// SharesJustifyingMember represents one member in a threshold key sharing group,
// after it completed secret shares and commitments verification and enters
// justification phase where it resolves invalid share accusations.
//
// Executes Phase 5 of the protocol.
type SharesJustifyingMember struct {
	*CommittingMember
}

// QualifiedMember represents one member in a threshold key sharing group, after
// it completed secret shares justification. The member holds a share of group
// master private key.
//
// Executes Phase 6 of the protocol.
type QualifiedMember struct {
	*SharesJustifyingMember

	// Member's share of the secret master private key. It is denoted as `z_ik`
	// in protocol specification.
	// TODO: unsure if we need shareT `x'_i` field, it should be removed if not used in further steps
	masterPrivateKeyShare, shareT *big.Int
}

// SharingMember represents one member in a threshold key sharing group, after it
// has been qualified to the master private key sharing group. A member shares
// public values of it's polynomial coefficients with peer members.
//
// Executes Phase 7 and Phase 8 of the protocol.
type SharingMember struct {
	*QualifiedMember

	// Public values of each polynomial `a` coefficient defined in secretCoefficients
	// field. It is denoted as `A_ik` in protocol specification. The zeroth
	// public key share point `A_i0` is a member's public key share.
	publicKeySharePoints []*big.Int
	// Public key share points received from peer group members which passed the
	// validation. Defined as `A_jk` across the protocol documentation.
	receivedValidPeerPublicKeySharePoints map[int][]*big.Int
}

// PointsJustifyingMember represents one member in a threshold key sharing group,
// after it completed public key share points verification and enters justification
// phase where it resolves public key share points accusations.
//
// Executes Phase 9 of the protocol.
type PointsJustifyingMember struct {
	*SharingMember
}

// PublishingMember represents one member in a threshold key sharing group,
// after it completed public key share points justification and proceeds to
// result publication phase.
//
// Executes Phase 13 of the protocol.
type PublishingMember struct {
	*PointsJustifyingMember

	result *result.Result
}

//type PublishingMember struct {
//	*SharingMember
//
//	result *result.Result
//}
