package coordinator

import (
	"github.com/keep-network/keep-core/pkg/tbtc"
)

// Chain represents the interface that the coordinator module expects to interact
// with the anchoring blockchain on.
type Chain interface {
	// TODO: Change to something more specific once https://github.com/keep-network/keep-core/issues/3632
	// is handled.
	tbtc.Chain
}
