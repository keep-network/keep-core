package ethereum

import (
	"context"
	"math/big"
	"time"

	"github.com/keep-network/keep-core/pkg/chain"

	"github.com/ethereum/go-ethereum/common"
)

// BalanceSource provides a balance info for the given address.
type BalanceSource func(address common.Address) (*big.Int, error)

// BalanceMonitor provides the possibility to monitor balances for given
// accounts.
type BalanceMonitor struct {
	balanceSource BalanceSource
}

// Observe starts a process which checks the address balance with the given
// tick and triggers an alert in case the balance falls below the
// alert threshold value.
func (bm *BalanceMonitor) Observe(
	ctx context.Context,
	address string,
	alertThreshold *big.Int,
	tick time.Duration,
) {
	check := func() {
		balance, err := bm.balanceSource(common.HexToAddress(address))
		if err != nil {
			logger.Errorf("ethereum balance monitor error: [%v]", err)
			return
		}

		if balance.Cmp(alertThreshold) == -1 {
			logger.Errorf(
				"ethereum balance for account [%v] is below [%v] wei; "+
					"account should be funded",
				address,
				alertThreshold.Text(10),
			)
		}
	}

	go func() {
		ticker := time.NewTicker(tick)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				check()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (ec *ethereumChain) BalanceMonitor() (chain.BalanceMonitor, error) {
	return &BalanceMonitor{ec.WeiBalanceOf}, nil
}
