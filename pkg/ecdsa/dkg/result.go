package dkg

import (
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// Result of distributed key generation protocol.
type Result struct {
	// TODO: Temporary result. Add real items.
	SymmetricKeys map[group.MemberIndex]ephemeral.SymmetricKey
}
