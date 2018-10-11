package dkg

import (
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
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
