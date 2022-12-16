package maintainer

import (
	"context"
	"time"

	"github.com/ipfs/go-log"
)

// TODO: Store the logger inside the wallet maintainer.
var walletLogger = log.Logger("keep-maintainer-wallet")

func newWalletMaintainer(
	ctx context.Context,
	chain WalletChain,
) *WalletMaintainer {
	walletMaintainer := &WalletMaintainer{
		chain: chain,
	}

	go walletMaintainer.startControlLoop(ctx)

	return walletMaintainer
}

// TODO: Description
type WalletMaintainer struct {
	chain WalletChain
}

// startControlLoop launches the loop responsible for controlling the wallet
// maintainer.
func (r *WalletMaintainer) startControlLoop(ctx context.Context) {
	walletLogger.Info("starting wallet maintainer")

	defer func() {
		walletLogger.Info("stopping wallet maintainer")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// TODO: Implement the maintainer loop. For now just print a message.
			walletLogger.Info("Wallet maintainer is working")
		}

		time.Sleep(1 * time.Second)
	}
}
