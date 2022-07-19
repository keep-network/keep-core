package tbtc

import (
	"github.com/keep-network/keep-core/pkg/sortition"
	"github.com/keep-network/keep-core/pkg/subscription"
	"math/big"
)

// Chain represents the interface that the TBTC module expects to interact
// with the anchoring blockchain on.
type Chain interface {
	// OnDKGStarted registers a callback that is invoked when an on-chain
	// notification of the DKG process start is seen.
	OnDKGStarted(func(event *DKGStartedEvent)) subscription.EventSubscription

	sortition.Chain
}

// DKGStartedEvent represents a DKG start event.
type DKGStartedEvent struct {
	Seed        *big.Int
	BlockNumber uint64
}
