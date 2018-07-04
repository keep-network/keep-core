package relay

import (
	"sync"

	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

// Node represents the current state of a relay node.
type Node struct {
	mutex sync.Mutex

	// StakeID is the ID this node is using to prove its stake in the system.
	StakeID string

	// External interactors.
	netProvider  net.Provider
	blockCounter chain.BlockCounter
	chainConfig  config.Chain

	// The IDs of the known stakes in the system, including this node's StakeID.
	stakeIDs []string

	// lastSeenEntry is the last relay entry this node is aware of.
	lastSeenEntry event.Entry
}
