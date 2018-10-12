package dkg

import (
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/pedersen"
)

type memberCore struct {
	// ID of this group member.
	ID *big.Int
	// Group to which this member belongs.
	group *Group
	// DKG Protocol configuration parameters.
	protocolConfig *config.DKG
}

func (m *memberCore) ProtocolConfig() *config.DKG {
	return m.protocolConfig
}

// CommittingMember represents one member in a threshold key sharing group, after
// it has a full list of `memberIDs` that belong to its threshold group. A
// member in this state has two maps of `memberShares` for each member of the
// group.
type CommittingMember struct {
	*memberCore

	// Pedersen VSS scheme used to calculate commitments.
	vss *pedersen.VSS
	// Shares calculated for current member by themself and peer group member.
	// secretShares are defined as `s_ij` and randomShares are defined as `t_ij`
	// across [documentation].
	//
	// [documentation]: http://docs.keep.network/cryptography/beacon_dkg.html#_protocol
	secretShares, randomShares map[*big.Int]*big.Int
}
