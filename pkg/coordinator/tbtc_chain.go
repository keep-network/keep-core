package coordinator

import (
	"github.com/keep-network/keep-core/pkg/tbtc"
)

// TbtcChain represents the interface that the Wallet Coordinator module expects to
// interact with the anchoring blockchain on.
type TbtcChain interface {
	tbtc.BridgeChain
	tbtc.WalletCoordinatorChain
}
