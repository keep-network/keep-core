package tbtc

import (
	"github.com/keep-network/keep-core/pkg/ecdsa"
	"github.com/keep-network/keep-core/pkg/subscription"
)

// Chain represents the interface that the TBTC module expects to interact
// with the anchoring blockchain on.
type Chain interface {
	// OnNewWalletRegistered registers a callback that is invoked when an
	// on-chain notification of the new wallet registration is seen.
	OnNewWalletRegistered(
		func(event *NewWalletRegisteredEvent),
	) subscription.EventSubscription

	ecdsa.Chain
}

// NewWalletRegisteredEvent represents an event informing about a new wallet
// registration.
type NewWalletRegisteredEvent struct {
	WalletID         ecdsa.WalletID
	WalletPubKeyHash [20]byte
}
